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

		accesslog = GetLog(accesslogLine)
		// fmt.Printf("%dfilter==>%s\n", i, accesslogLine)
		if i%10000 == 0 {
			fmt.Printf("%s %d\n", time.Now(), i)
		}
		i++

		filterResult = DoFilter(redisConn, accesslog)

		//  these should done in filter function
		//
		if filterResult == YES {
			redisConn.ListLeftPush("accesslog_yes", accesslogLine)
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
	return URIFilter(redisConn, accesslog)
}

func URIFilter(redisConn RedisConn, accesslog AccessLog) int {
	if matched, err := regexp.MatchString("^/prop/view", accesslog.RequestURI); err == nil && matched {
		return HttpCodeFilter(redisConn, accesslog)
	} else {
		Analysis(redisConn, accesslog)
		return UNKNOWN
	}
}

func HttpCodeFilter(redisConn RedisConn, accesslog AccessLog) int {
	if matched, err := regexp.MatchString("^2", accesslog.HttpCode); err == nil && matched {
		return UserAgentFilter(redisConn, accesslog)
	} else {
		return UNKNOWN
	}
}

func UserAgentFilter(redisConn RedisConn, accesslog AccessLog) int {
	if accesslog.UserAgent == "-" {
		return NO
	} else {
		uaFamily := Parse(accesslog.UserAgent)
		if uaFamily == "" {
			return NO
		}
		uaFamily = strings.ToLower(uaFamily)
		if strings.Contains(uaFamily, "bot") {
			return NO
		}
		redisConn.SetAdd("ua", uaFamily)
		return YES
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
