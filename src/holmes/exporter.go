package main

func Export(c chan int, holmesConfig *HolmesConfig) {
	for {
		//BlockListPop("accesslog_yes", 1)
		//BlockListPop("accesslog_no", 1)
		BlockListRightPop("accesslog", 0)
		//guid:=BlockListPop("guid",accesslog.GUID)
	}
	c <- 1
}
