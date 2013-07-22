package main

import (
	"fmt"
	"strings"
	//"net"
	//"net/http"
	"regexp"
	"time"
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
		if accesslogLine == "" {
			fmt.Printf("%s now list have no log to process,continue to wait others to add log to list\n", time.Now())
			continue
		}

		accesslog = GetLogNginx(accesslogLine)
		// fmt.Printf("%dfilter==>%s\n", i, accesslogLine)
		if i%100000 == 0 {
			fmt.Printf("%s %d\n", time.Now(), i)
		}
		i++

		filterResult = DoFilter(redisConn, accesslog)

		//  these should done in filter function
		//
		if filterResult == YES {
			logTimeMin := accesslog.LogTimeMinString()
			redisConn.HashIncrby("accesslog_result_time", logTimeMin, 1)
			redisConn.HashIncrby("accesslog_result", "effective_pv", 1)
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
		return NO
	} else {
		uaFamily := Parse(accesslog.UserAgent)
		if uaFamily == "" {
			return NO
		} else {
			uaFamily = strings.ToLower(uaFamily)
			if strings.Contains(uaFamily, "bot") {
				return NO
			} else {
				redisConn.SetAdd("ua", uaFamily)
				return URIFilter(redisConn, accesslog)
			}
		}
	}
}

func URIFilter(redisConn RedisConn, accesslog AccessLog) int {
	redisConn.SetAdd(accesslog.RemoteAddr, accesslog.RequestURI) // record all logs of each ip
	if strings.Contains(accesslog.Hostname, "s.anjuke.com") {
		redisConn.SetAdd("s.anjuke.com", accesslog.RemoteAddr)
	}
	if matched, err := regexp.MatchString("^/prop/view/", accesslog.RequestURI); err == nil && matched {
		redisConn.HashIncrby("accesslog_result", "total_pv", 1)
		return HttpCodeFilter(redisConn, accesslog)
	} else {
		//TODO Analysis(redisConn,accesslog)
		return UNKNOWN
	}
}

func HttpCodeFilter(redisConn RedisConn, accesslog AccessLog) int {
	if matched, err := regexp.MatchString("^2", accesslog.HttpCode); err == nil && matched {
		redisConn.HashIncrby("accesslog_result", accesslog.HttpCode, 1)
		return RefererFilter(redisConn, accesslog)
		//return YES
	} else {
		redisConn.HashIncrby("accesslog_result", accesslog.HttpCode, 1)
		//log.Printf("http code: %s", accesslog.HttpCode)
		return UNKNOWN
	}
}

func RefererFilter(redisConn RedisConn, accesslog AccessLog) int {
	if accesslog.Referer == "-" && accesslog.GUID == "-" {
		return NO
	}
	if accesslog.Referer == "-" && accesslog.GUID != "-" {
		if redisConn.SetIsMember("s.anjuke.com", accesslog.RemoteAddr) == 1 {
			return YES
		} else {
			return NO
		}
	}
	if redisConn.SetIsMember(accesslog.RemoteAddr, accesslog.Referer) == 1 {
		return YES
	} else {
		return NO
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
	if true {
		//TODO if accesslog . RemoteAddr == (check whether in the whitelist)

		return YES
	} else {
		return UNKNOWN
	}
}

func AddWatchingList(redisConn RedisConn, accesslog AccessLog) {
	redisConn.ListLeftPush("WatchingList", accesslog.String())
}

func DelWatchingList(redisConn RedisConn, accesslog AccessLog) {
	//redisConn.ListRigthPop("WatchingList", accesslog.String())
	//
	//   TODO pop specific record
}

func AddWhiteList(redisConn RedisConn, accesslog AccessLog) {
	redisConn.ListLeftPush("WhiteList", accesslog.String())
}

func AddIgnoreList(redisConn RedisConn, accesslog AccessLog) {
	redisConn.ListLeftPush("IgnoreList", accesslog.String())
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
