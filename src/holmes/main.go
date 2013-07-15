package main

import (
	"runtime"
)

var holmesConf HolmesConfig

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	confFile := "holmes.conf"
	ua_pattern_file := "../data/user_agent_pattern.json"
	holmesConf = LoadConfig(confFile)
	InitRedisConf(&holmesConf)
	InitUAParsers(ua_pattern_file)

	c := make(chan int, 3)
	go StageLog(c, &holmesConf)
	go Filter(c)
	go Export(c, &holmesConf)
	for i := 0; i < 3; i++ { // waiting all goroutines to finish
		<-c
	}
}
