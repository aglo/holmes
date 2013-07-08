package main

func StageLog(c chan int, holmesConfig *HolmesConfig) {
	filename := holmesConfig.Infilepath
	lines := ReadLogLines(filename)
	for _, line := range lines {
		ListPush("accesslog", line)
	}
	c <- 1
}
