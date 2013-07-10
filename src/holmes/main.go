package main

var holmesConf HolmesConfig

func main() {
    confFile := "holmes.conf"
    holmesConf = LoadConfig(confFile)
    InitRedisConf(&holmesConf)
    c := make(chan int, 1)
    go StageLog(c, &holmesConf)
    go Filter(c)
    go Export(c, &holmesConf)
    for i := 0; i < 3; i++ { // waiting all goroutines to finish
        <-c
    }
}
