package main

import (
	"fmt"
	"strings"
	//"net"
	//"net/http"
	"regexp"
	"time"
	"log"
)

const (
	YES = iota // is a human
	NO         // is not a human
	UNKNOWN
)

func Filter(c chan int) {
	var accesslogLine string
	var accesslog AccessLog
	var filterResult int
	var i int
	redisConn := NewRedisConn()
	defer CloseRedisConn(redisConn)
	for {
		_, accesslogLine = redisConn.BlockListRightPop("accesslog", 5)
		//fmt.Printf("%dfilter==>%s\n", i, accesslogLine)
		if accesslogLine == "" {
			fmt.Printf("%s now list have no log to process,continue to wait others to add log to list\n", time.Now())
			continue
		}

		accesslog = GetLog(accesslogLine)
		i++
		if i%100000 == 0 {
			fmt.Printf("%s holmes have processed %d logs\n", time.Now(), i)
		}

		redisConn.HashIncrby("accesslog_result", "total_request", 1)
		filterResult = DoFilter(redisConn, accesslog)

		//  these should done in filter function
		//
		if filterResult == YES {
			logTimeMin := accesslog.LogTimeMinString()
			//log.Println("Result of DoFilter() is YES,increment accesslog_result_vppv_per_min at ",logTimeMin)
			redisConn.HashIncrby("accesslog_result_vppv_per_min", logTimeMin, 1)
			redisConn.HashIncrby("accesslog_result", "effective_vppv", 1)
			//redisConn.ListLeftPush("accesslog_yes", accesslogLine)
		}
		//else if filterResult == NO {
		//	redisConn.ListLeftPush("accesslog_no", accesslogLine)
		//} else {
		//	redisConn.ListLeftPush("accesslog_unkown", accesslogLine)
		//}
	}
	c <- 1
}

func DoFilter(redisConn RedisConn, accesslog AccessLog) int {
	return UserAgentFilter(redisConn, accesslog)
}

func UserAgentFilter(redisConn RedisConn, accesslog AccessLog) int {
	if accesslog.UserAgent == "-" {
		redisConn.HashIncrby("accesslog_result", "UserAgentFilter_NO_PASS", 1)
		return NO
	} else {
		uaFamily := Parse(accesslog.UserAgent)
		if uaFamily == "" {
			redisConn.HashIncrby("accesslog_result", "UserAgentFilter_NO_PASS", 1)
			return NO
		} else {
			uaFamily = strings.ToLower(uaFamily)
			if strings.Contains(uaFamily, "bot") {
				redisConn.HashIncrby("accesslog_result", "UserAgentFilter_NO_PASS", 1)
				return NO
			} else {
				redisConn.HashIncrby("accesslog_result", "UserAgentFilter_PASS", 1)
				redisConn.SetAdd("ua", uaFamily)
				AddRefererList(redisConn, accesslog)
				return URIFilter(redisConn, accesslog)
			}
		}
	}
}

func URIFilter(redisConn RedisConn, accesslog AccessLog) int {
	//
	//redisConn.SetAdd(accesslog.RemoteAddr, accesslog.RequestURI) // record all logs of each ip
	if matched, err := regexp.MatchString("^/prop/view/", accesslog.RequestURI); err == nil && matched {
		redisConn.HashIncrby("accesslog_result", "total_vppv", 1)
		return HttpCodeFilter(redisConn, accesslog)
	} else {
		Analysis(redisConn, accesslog)
		return UNKNOWN
	}
}

func HttpCodeFilter(redisConn RedisConn, accesslog AccessLog) int {
	redisConn.HashIncrby("accesslog_result", "vppv_"+accesslog.HttpCode, 1)

	if matched, err := regexp.MatchString("^2", accesslog.HttpCode); err == nil && matched {
		return WhiteIpFilter(redisConn, accesslog)
	} else {
		return UNKNOWN
	}
}

func RefererFilter(redisConn RedisConn, accesslog AccessLog) int {
	/*if accesslog.Referer == "-" && accesslog.GUID == "-" {
		return NO
	}
	if accesslog.Referer == "-" && accesslog.GUID != "-" {
		if redisConn.SetIsMember("s.anjuke.com", accesslog.RemoteAddr) == 1 {
			return YES
		} else {
			return NO
		}
	}*/
	if strings.Contains(accesslog.Referer, "my.anjuke.com") == true {
		log.Println(accesslog.String(),"is come from my.anjuke.com")
		return NO
	} else {
		if redisConn.SetIsMember("Referer_"+accesslog.RemoteAddr, accesslog.Referer) == 1 {
			log.Println(accesslog.String(),"is come from ",accesslog.Referer)
			return YES
		} else {
			//log.Println(accesslog.String(),"is not come from ",accesslog.Referer)
			return NO
		}
	}

	//////////// get UA type from website
	//
	//_, err := http.Get("http://www.useragentstring.com/?usa=" + accesslog.UserAgent + "&getText=all")
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	//	fmt.Println("success", res)
	//}

	////////////  DNS reverse lookup
	//
	//if matched,err := regexp.MatchString("[S|s]pider",accesslog.UserAgent) ; err != nil || !matched{
	//    return UNKNOWN
	//} else {
	//    ans , err1 := net.LookupAddr(accesslog.RemoteAddr)
	//    if err1 != nil{
	//        fmt.Println("Failed",accesslog.UserAgent,"+", accesslog.RemoteAddr,err1)
	//    } else {
	//        fmt.Println("Successful",accesslog.UserAgent,"+",accesslog.RemoteAddr,"-->",ans)
	//    }
	//    return UNKNOWN
	//}
}

func WhiteIpFilter(redisConn RedisConn, accesslog AccessLog) int {
	if 1 == redisConn.SetIsMember("WhiteList", accesslog.RemoteAddr) {
		return YES
	} else {
		AddWatchingList(redisConn, accesslog)
		return UNKNOWN
	}
}

func AddRefererList(redisConn RedisConn, accesslog AccessLog) {
	//log.Println("add to Referer_"+accesslog.RemoteAddr, "member:","http://"+accesslog.Hostname+accesslog.RequestURI)
	redisConn.SetAdd("RefererList", accesslog.RemoteAddr)
	redisConn.SetAdd("Referer_"+accesslog.RemoteAddr, "http://"+accesslog.Hostname+accesslog.RequestURI)
}

func DelRefererList(redisConn RedisConn, accesslog AccessLog) {
	//log.Println("DelRefererList delete member ",accesslog.RemoteAddr," and key Referer_" + accesslog.RemoteAddr)
	redisConn.SetRem("RefererList", accesslog.RemoteAddr)
	redisConn.KeyDel("Referer_" + accesslog.RemoteAddr)
}

func AddWatchingList(redisConn RedisConn, accesslog AccessLog) {
	redisConn.SetAdd("WatchingList", accesslog.RemoteAddr)
	redisConn.ListLeftPush("WL_"+accesslog.RemoteAddr, accesslog.String())
}

func DelWatchingList(redisConn RedisConn, accesslog AccessLog) {
	redisConn.SetRem("WatchingList", accesslog.RemoteAddr)
	redisConn.KeyDel("WL_" + accesslog.RemoteAddr)
}

func AddWhiteList(redisConn RedisConn, accesslog AccessLog) {
	redisConn.SetAdd("WhiteList", accesslog.RemoteAddr)
}

func AddIgnoreList(redisConn RedisConn, accesslog AccessLog) {
	redisConn.SetAdd("IgnoreList", accesslog.RemoteAddr)
}

func Analysis(redisConn RedisConn, accesslog AccessLog) {
	if matched, err := regexp.MatchString("^s.anjuke.com", accesslog.Hostname); err == nil && matched {
		//AddWhiteList(redisConn, accesslog)
		//log.Println("call Analysis... Hostname is: ",accesslog.Hostname)
		ProcessWatchingList(redisConn, accesslog)
	}
}

func ProcessWatchingList(redisConn RedisConn, accesslog AccessLog) {
	//log.Println("call ProcessWatchingList... : ",accesslog.String())
	trustFlag:=false
	listLen := redisConn.ListLen("WL_" + accesslog.RemoteAddr)
	for i := 0; i < int(listLen); i++ {
		line := redisConn.ListLeftPop("WL_" + accesslog.RemoteAddr)
		watchAccesslog := GetLog(line)
		if matched, err := regexp.MatchString("^/prop/view/", watchAccesslog.RequestURI); err == nil && matched {
			if matched, err := regexp.MatchString("^2", watchAccesslog.HttpCode); err == nil && matched {
				//if RefererFilter(redisConn, watchAccesslog) == YES {
				if watchAccesslog.Referer != "-" {
					trustFlag=true
					logTimeMin := watchAccesslog.LogTimeMinString()
					//log.Println("Result of RefererFilter() is YES,increment accesslog_result_vppv_per_min at ",logTimeMin)
					redisConn.HashIncrby("accesslog_result_vppv_per_min", logTimeMin, 1)
					redisConn.HashIncrby("accesslog_result", "effective_vppv", 1)
				}
			}
		}
	}
	DelWatchingList(redisConn, accesslog)
	DelRefererList(redisConn, accesslog)
	if trustFlag{
		AddWhiteList(redisConn, accesslog)
	}
}

//func GUIDFilter(redisConn RedisConn, accesslog AccessLog) int {
//	if accesslog.GUID == "-" {
//		return NO
//	} else {
//		redisConn.ListLeftPush("guid", accesslog.GUID)
//		redisConn.ListLeftPush(accesslog.GUID, "----"+accesslog.Referer)
//		uri := accesslog.LogTimeString() + "==>" + accesslog.RequestURI
//		redisConn.ListLeftPush(accesslog.GUID, uri)
//		return YES
//	}
//}
//
//func IPFilter(redisConn RedisConn, accesslog AccessLog) int {
//	redisConn.SetAdd("ip", accesslog.RemoteAddr)
//	redisConn.ListLeftPush(accesslog.RemoteAddr, "----"+accesslog.Referer)
//	uri := accesslog.LogTimeString() + "==>" + accesslog.RequestURI
//	redisConn.ListLeftPush(accesslog.RemoteAddr, uri)
//	return UNKNOWN
//}
