package main

import (
	"testing"
)

func init() {
	confFile := "../../conf/holmes.conf"
	holmesConf = LoadConfig(confFile)
	InitRedisConf(&holmesConf)
}

func TestListLeftPush(t *testing.T) {
	list := "list_test"
	item := "item_test"

	redisConn := NewRedisConn()
	defer CloseRedisConn(redisConn)
	len := redisConn.ListLen(list)
	for i := 0; i < int(len); i++ {
		redisConn.ListLeftPop(list)
	}

	len = redisConn.ListLeftPush(list, item)
	if len != 1 {
		t.Errorf("push an item int a list failed")
	}
	len = redisConn.ListLeftPush(list, item)
	if len != 2 {
		t.Errorf("push an item int a list failed")
	}
	len = redisConn.ListLeftPush(list, item)
	if len != 3 {
		t.Errorf("push an item int a list failed")
	}
	redisConn.ListLeftPop(list)
	redisConn.ListLeftPop(list)
	redisConn.ListLeftPop(list)
}

func TestBlockListRightPop(t *testing.T) {
	list := "list_test"

	redisConn := NewRedisConn()
	defer CloseRedisConn(redisConn)

	len := redisConn.ListLen(list)
	for i := 0; i < int(len); i++ {
		redisConn.ListLeftPop(list)
	}

	redisConn.ListLeftPush(list, "a")
	redisConn.ListLeftPush(list, "b")
	redisConn.ListLeftPush(list, "c")

	resultList, item := redisConn.BlockListRightPop(list, 1)
	if resultList != list {
		t.Errorf("listname:%s", resultList)
	} else if item != "a" {
		t.Errorf("item:%s", item)
	}

	resultList, item = redisConn.BlockListRightPop(list, 1)
	if resultList != list {
		t.Errorf("listname:%s", resultList)
	} else if item != "b" {
		t.Errorf("failed")
	}

	resultList, item = redisConn.BlockListRightPop(list, 1)
	if resultList != list {
		t.Errorf("listname:%s", resultList)
	} else if item != "c" {
		t.Errorf("item:%s", item)
	}
}
