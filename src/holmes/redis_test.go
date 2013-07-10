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

    len := ListLen(list)
    for i := 0; i < int(len); i++ {
        ListLeftPop(list)
    }

    len = ListLeftPush(list, item)
    if len != 1 {
        t.Errorf("push an item int a list failed")
    }
    len = ListLeftPush(list, item)
    if len != 2 {
        t.Errorf("push an item int a list failed")
    }
    len = ListLeftPush(list, item)
    if len != 3 {
        t.Errorf("push an item int a list failed")
    }
    ListLeftPop(list)
    ListLeftPop(list)
    ListLeftPop(list)
}

func TestBlockListRightPop(t *testing.T) {
    list := "list_test"

    len := ListLen(list)
    for i := 0; i < int(len); i++ {
        ListLeftPop(list)
    }

    ListLeftPush(list, "a")
    ListLeftPush(list, "b")
    ListLeftPush(list, "c")

    resultList, item := BlockListRightPop(list, 1)
    if resultList != list {
        t.Errorf("listname:%s", resultList)
    } else if item != "a" {
        t.Errorf("item:%s", item)
    }

    resultList, item = BlockListRightPop(list, 1)
    if resultList != list {
        t.Errorf("listname:%s", resultList)
    } else if item != "b" {
        t.Errorf("failed")
    }

    resultList, item = BlockListRightPop(list, 1)
    if resultList != list {
        t.Errorf("listname:%s", resultList)
    } else if item != "c" {
        t.Errorf("item:%s", item)
    }
}
