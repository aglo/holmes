package main

import (
	"fmt"
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

		filterResult = DoFilter(accesslog)
		if filterResult == YES {
			redisConn.ListLeftPush("accesslog_yes", accesslogLine)
		} else if filterResult == NO {
			redisConn.ListLeftPush("accesslog_no", accesslogLine)
		} else {
			redisConn.ListLeftPush("accesslog_unkown", accesslogLine)
		}
	}
	c <- 1
}

func DoFilter(accesslog AccessLog) int {
	//return IPFilter(accesslog)
	return YES
}

/*func GUIDFilter(accesslog AccessLog) int {
	if accesslog.GUID == "-" {
		return NO
	} else {
		ListLeftPush("guid", accesslog.GUID)
		ListLeftPush(accesslog.GUID, "----"+accesslog.Referer)
		uri := accesslog.LogTimeString() + "==>" + accesslog.RequestURI
		ListLeftPush(accesslog.GUID, uri)
		return YES
	}
}

func IPFilter(accesslog AccessLog) int {
	SetAdd("ip", accesslog.RemoteAddr)
	ListLeftPush(accesslog.RemoteAddr, "----"+accesslog.Referer)
	uri := accesslog.LogTimeString() + "==>" + accesslog.RequestURI
	ListLeftPush(accesslog.RemoteAddr, uri)
	return UNKNOWN
}*/
