package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func StageLog(c chan int, holmesConfig *HolmesConfig) {
	dirname := holmesConfig.InLogDir
	filenames := ReadFilenames(dirname)
	redisConn := NewRedisConn()
	defer CloseRedisConn(redisConn)
	for _, filename := range filenames {
		fmt.Println(time.Now(), " read file:", dirname+"/"+filename)
		file, err := os.Open(dirname + "/" + filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			redisConn.ListLeftPush("accesslog", scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading ", filename, ":", err)
		}
	}
	c <- 1
}
