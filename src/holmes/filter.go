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

		filterResult = DoFilter(redisConn, accesslog)
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

func DoFilter(redisConn RedisConn, accesslog AccessLog) int {
	FilterFlag := UNKNOWN
	switch {
	case (UNKNOWN == GUIDFilter(redisConn, accesslog)):
		return FilterFlag
	case (UNKNOWN == IPFilter(redisConn, accesslog)):
		return FilterFlag
	}
	////FilterFlag = GUIDFilter(accesslog)
	//FilterFlag = IPFilter(accesslog)
	return FilterFlag
}

func GUIDFilter(redisConn RedisConn, accesslog AccessLog) int {
	if accesslog.GUID == "-" {
		return NO
	} else {
		redisConn.ListLeftPush("guid", accesslog.GUID)
		redisConn.ListLeftPush(accesslog.GUID, "----"+accesslog.Referer)
		uri := accesslog.LogTimeString() + "==>" + accesslog.RequestURI
		redisConn.ListLeftPush(accesslog.GUID, uri)
		return YES
	}
}

func IPFilter(redisConn RedisConn, accesslog AccessLog) int {
	redisConn.SetAdd("ip", accesslog.RemoteAddr)
	redisConn.ListLeftPush(accesslog.RemoteAddr, "----"+accesslog.Referer)
	uri := accesslog.LogTimeString() + "==>" + accesslog.RequestURI
	redisConn.ListLeftPush(accesslog.RemoteAddr, uri)
	return UNKNOWN
}
