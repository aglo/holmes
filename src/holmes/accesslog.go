package main

import (
	"io/ioutil"
	"strings"
)

type AccessLog struct {
	Hour                 string
	Year                 string
	Day                  string
	Min                  string
	RequestTime          string
	UpstreamResponseTime string
	RemoteAddr           string
	UpstreamAddr         string
	Hostname             string
	Method               string
	RequestURI           string
	HttpCode             string
	BytesSent            string
	Referer              string
	UserAgent            string
	GzipRatio            string
	HttpXForwardedFor    string
	ServerAddr           string
	GUID                 string
	Sec                  string
	Month                string
	RequestLen           string
	ServerPort           string
}

func ReadLogLines(filename string) []string {
	lines := []string{}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	} else {
		lines = strings.Split(string(data), "\n")
	}
	return lines
}

func GetLog(line string) AccessLog {
	var accessLog AccessLog
	if line != "" {
		fields := strings.Split(line, "\t")
		accessLog.Hour = fields[0]
		accessLog.Year = fields[1]
		accessLog.Day = fields[2]
		accessLog.Min = fields[3]
		accessLog.RequestTime = fields[4]
		accessLog.UpstreamResponseTime = fields[5]
		accessLog.RemoteAddr = fields[6]
		accessLog.UpstreamAddr = fields[7]
		accessLog.Hostname = fields[8]
		accessLog.Method = fields[9]
		accessLog.RequestURI = fields[10]
		accessLog.HttpCode = fields[11]
		accessLog.BytesSent = fields[12]
		accessLog.Referer = fields[13]
		accessLog.UserAgent = fields[14]
		accessLog.GzipRatio = fields[15]
		accessLog.HttpXForwardedFor = fields[16]
		accessLog.ServerAddr = fields[17]
		accessLog.GUID = fields[18]
		accessLog.Sec = fields[19]
		accessLog.Month = fields[20]
		accessLog.RequestLen = fields[21]
		accessLog.ServerPort = fields[22]
	}
	return accessLog
}
