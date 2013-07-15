package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
	"regexp"
)

type UAParserPattern struct { // UA is stand for User Agent
	RegexpString      string
	FamilyReplacement string
}

type UAParser struct {
	uAParserPattern UAParserPattern
	regexp          *regexp.Regexp
}

var UAParsers = []UAParser{}

func InitUAParsers(pattern_file string) {
	userAgentParserPatterns := LoadPattern(pattern_file)
	for _, pattern := range userAgentParserPatterns {
		regexp, _ := regexp.Compile(pattern.RegexpString)
		uaParser := UAParser{uAParserPattern: pattern, regexp: regexp}
		UAParsers = append(UAParsers, uaParser)
	}
}

func Parse(ua string) string {
	var uaFamily string
	for _, uaParser := range UAParsers {
		uaFamily = uaParser.Parse(ua)
		if uaFamily != "" {
			break
		}
	}
	return uaFamily
}

func (uaParser *UAParser) Parse(ua string) string {
	matchs := uaParser.regexp.FindStringSubmatch(ua)
	if matchs != nil {
		if uaParser.uAParserPattern.FamilyReplacement != "None" {
			return uaParser.uAParserPattern.FamilyReplacement
		}
		return matchs[1]
	}
	return "" // no matchs
}

func LoadPattern(filename string) []UAParserPattern {
	var userAgentParserPatterns []UAParserPattern
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	patternReader := bufio.NewReader(file)
	var content string
	for {
		line, err := patternReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				content = content + line
				break
			}
			log.Fatal(err)
		} else {
			content = content + line
		}
	}
	temp := []byte(content)
	err = json.Unmarshal(temp, &userAgentParserPatterns)
	if err != nil {
		log.Fatal(err)
	}

	return userAgentParserPatterns
}
