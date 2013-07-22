package main

import (
	"io/ioutil"
	"regexp"
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

func (accessLog *AccessLog) LogTimeMinString() string {
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

func ReadFilenames(dirname string) []string {
	filenames := []string{}
	fileInfos, err := ioutil.ReadDir(dirname)
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range fileInfos {
		filenames = append(filenames, fileInfo.Name())
	}
	return filenames
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

func GetLogNginx(line string) AccessLog {
	monthMap := map[string]string{"Jan": "1", "Feb": "2", "Mar": "3", "Apr": "4", "May": "5", "Jun": "6", "Jul": "7", "Aug": "8", "Sep": "9", "Oct": "10", "Nov": "11", "Dec": "12"}
	pattern := `(\d+\.\d+|\-)` +
		`\s` +
		`(\d+\.\d+|\-)` +
		`\s` +
		`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})` +
		`\s` +
		`(\d+)` +
		`\s` +
		`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}|\-)(:(\d{1,5}))?` +
		`\s+` +
		`\[(\d{2})\/` +
		`([A-Z][a-z]{2}?)\/` +
		`(\d{4})\:` +
		`(\d{2})\:` +
		`(\d{2})\:` +
		`(\d{2})` +
		`\s+\+0800\]` +
		`\s` +
		`([^\s]+?)` +
		`\s` +
		`"` +
		`([A-Z]+)` +
		`\s` +
		`([^\s]+?)` +
		`\s` +
		`HTTP/[0-9.]+` +
		`"` +
		`\s` +
		`(\d{3})` +
		`\s` +
		`(\d+)` +
		`\s` +
		`"` +
		`([^\"]+|\-)` +
		`"` +
		`\s` +
		`"([^\"]+|\-)"` +
		`\s` +
		`"([^\"]+|\-)"` +
		`\s` +
		`"([^\"]+|\-)"` +
		`\s` +
		`-` +
		`\s` +
		`"` +
		`(` +
		`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})` +
		`(:` +
		`(\d+)?` +
		`|\-)` +
		`\s?` +
		`(.+?))?` +
		`"` +
		`.*`

	myRegexp, _ := regexp.Compile(pattern)

	var accessLog AccessLog
	if line != "" {
		fields := myRegexp.FindSubmatch([]byte(line))
		//for i, d := range fields {
		//fmt.Println(i, ":", d)
		//}
		//fmt.Println(string(fields[0]))
		accessLog.Hour = string(fields[11])
		accessLog.Year = string(fields[10])
		accessLog.Day = string(fields[8])

		accessLog.Min = string(fields[12])
		accessLog.RequestTime = string(fields[1])
		accessLog.UpstreamResponseTime = string(fields[2])
		accessLog.RemoteAddr = string(fields[3])
		accessLog.UpstreamAddr = string(fields[5])
		accessLog.Hostname = string(fields[14])
		accessLog.Method = string(fields[15])
		accessLog.RequestURI = string(fields[16])
		accessLog.HttpCode = string(fields[17])
		accessLog.BytesSent = string(fields[18])
		accessLog.Referer = string(fields[19])
		accessLog.UserAgent = string(fields[20])
		accessLog.GzipRatio = string(fields[21])
		accessLog.HttpXForwardedFor = string(fields[22])
		accessLog.ServerAddr = string(fields[24])
		accessLog.GUID = string(fields[27])
		accessLog.Sec = string(fields[13])
		accessLog.Month = monthMap[string(fields[9])]
		accessLog.RequestLen = string(fields[4])
		accessLog.ServerPort = string(fields[26])
	}
	return accessLog
}
