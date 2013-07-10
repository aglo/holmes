package main

import (
    "fmt"
)

const (
    YES = iota // is a human
    NO         // is not a human
    UNKNOWN
)

func Filter(c chan int) {
    for {
        _, accesslogLine := BlockListRightPop("accesslog", 0)

        accesslog := GetLog(accesslogLine)
        fmt.Printf("filter==>%s\n", accesslogLine)

        filterResult := DoFilter(accesslog)
        if filterResult == YES {
            ListLeftPush("accesslog_yes", accesslogLine)
        } else if filterResult == NO {
            ListLeftPush("accesslog_no", accesslogLine)
        } else {
            ListLeftPush("accesslog_unkown", accesslogLine)
        }
    }
    c <- 1

}

func DoFilter(accesslog AccessLog) int {
    return GUIDFilter(accesslog)
}

func GUIDFilter(accesslog AccessLog) int {
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
