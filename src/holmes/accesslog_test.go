package main

import (
	"testing"
)

func TestString(t *testing.T) {
	var accessLog AccessLog
	accessLog.Hour = "1"
	accessLog.Year = "2"
	accessLog.Day = "3"
	accessLog.Min = "4"
	accessLog.RequestTime = "5"
	accessLog.UpstreamResponseTime = "6"
	accessLog.RemoteAddr = "7"
	accessLog.UpstreamAddr = "8"
	accessLog.Hostname = "9"
	accessLog.Method = "10"
	accessLog.RequestURI = "11"
	accessLog.HttpCode = "12"
	accessLog.BytesSent = "13"
	accessLog.Referer = "14"
	accessLog.UserAgent = "15"
	accessLog.GzipRatio = "16"
	accessLog.HttpXForwardedFor = "17"
	accessLog.ServerAddr = "18"
	accessLog.GUID = "19"
	accessLog.Sec = "20"
	accessLog.Month = "21"
	accessLog.RequestLen = "22"
	accessLog.ServerPort = "23"
	accessLogString := accessLog.String()
	newAccessLog := GetLog(accessLogString)
	if newAccessLog != accessLog {
		t.Errorf("String() is not a correct method")
		t.Errorf("%s len: %d", accessLogString, len(accessLogString))
	}
}

func TestRequestTimeString(t *testing.T) {
	var accesslog AccessLog
	accesslog.Year = "2013"
	accesslog.Month = "7"
	accesslog.Day = "9"
	accesslog.Hour = "15"
	accesslog.Min = "20"
	accesslog.Sec = "12"
	strings := accesslog.LogTimeString()
	if strings != "2013-7-9 15:20:12" {
		t.Errorf("%s len: %d", strings, len(strings))
	}
}