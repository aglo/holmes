package main

import (
	"io/ioutil"
	"strings"
)

type AccessLog struct {
	Hour                 string // 1
	Year                 string // 2
	Day                  string // 3
	Min                  string // 4
	RequestTime          string // 5
	UpstreamResponseTime string // 6
	RemoteAddr           string // 7
	UpstreamAddr         string // 8
	Hostname             string // 9
	Method               string // 10
	RequestURI           string // 11
	HttpCode             string // 12
	BytesSent            string // 13
	Referer              string // 14
	UserAgent            string // 15
	GzipRatio            string // 16
	HttpXForwardedFor    string // 17
	ServerAddr           string // 18
	GUID                 string // 19
	Sec                  string // 20
	Month                string // 21
	RequestLen           string // 22
	ServerPort           string // 23
}

func (accessLog *AccessLog) String() string {
	fields := make([]string, 23)
	fields[0] = accessLog.Hour
	fields[1] = accessLog.Year
	fields[2] = accessLog.Day
	fields[3] = accessLog.Min
	fields[4] = accessLog.RequestTime
	fields[5] = accessLog.UpstreamResponseTime
	fields[6] = accessLog.RemoteAddr
	fields[7] = accessLog.UpstreamAddr
	fields[8] = accessLog.Hostname
	fields[9] = accessLog.Method
	fields[10] = accessLog.RequestURI
	fields[11] = accessLog.HttpCode
	fields[12] = accessLog.BytesSent
	fields[13] = accessLog.Referer
	fields[14] = accessLog.UserAgent
	fields[15] = accessLog.GzipRatio
	fields[16] = accessLog.HttpXForwardedFor
	fields[17] = accessLog.ServerAddr
	fields[18] = accessLog.GUID
	fields[19] = accessLog.Sec
	fields[20] = accessLog.Month
	fields[21] = accessLog.RequestLen
	fields[22] = accessLog.ServerPort
	accessLogString := strings.Join(fields, "\t")

	return accessLogString
}

func (accessLog *AccessLog) LogTimeString() string {
	strings := []uint8{}
	strings = append(strings, accessLog.Year...)
	strings = append(strings, "-"...)
	strings = append(strings, accessLog.Month...)
	strings = append(strings, "-"...)
	strings = append(strings, accessLog.Day...)
	strings = append(strings, " "...)
	strings = append(strings, accessLog.Hour...)
	strings = append(strings, ":"...)
	strings = append(strings, accessLog.Min...)
	strings = append(strings, ":"...)
	strings = append(strings, accessLog.Sec...)

	return string(strings)
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
