package main

import (
	"fmt"
)

const (
	YES = iota
	NO
	UNKNOWN
)

func Filter(c chan int) {
	for {
		_, accesslogLine := BlockListPop("accesslog", 0)
		fmt.Println("filter==> ", accesslogLine)
		accesslog := GetLog(accesslogLine)
		filterResult := DoFilter(accesslog)
		if filterResult == YES {
			ListPush("accesslog_yes", accesslogLine)
		} else if filterResult == NO {
			ListPush("accesslog_no", accesslogLine)
		} else {
			fmt.Println("UNKNOWN")
		}
	}
	c <- 1

}

func DoFilter(accesslog AccessLog) int {
	return 0
}
