package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisConf struct {
	Network        string
	Address        string
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	BlockTimeout   int64
}

var redisConf RedisConf

func InitRedisConf(holmesConfig *HolmesConfig) {
	redisConf.Network = holmesConfig.RedisNet
	redisConf.Address = holmesConfig.RedisIP + ":" + holmesConfig.RedisPort
	redisConf.ConnectTimeout = time.Duration(holmesConfig.ConnectTimeout)
	redisConf.ReadTimeout = time.Duration(holmesConfig.ReadTimeout)
	redisConf.WriteTimeout = time.Duration(holmesConfig.WriteTimeout)
	redisConf.BlockTimeout = holmesConfig.BlockTimeout
}

func ConnectRedisTimeout() redis.Conn {
	c, err := redis.DialTimeout(redisConf.Network, redisConf.Address, redisConf.ConnectTimeout, redisConf.ReadTimeout, redisConf.WriteTimeout)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return c
}

func ConnectRedis() redis.Conn {
	c, err := redis.Dial(redisConf.Network, redisConf.Address)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return c
}

func CloseConn(c redis.Conn) {
	c.Close()
}

/*
*ListPush push an item into a list at the end of the list
 */
func ListPush(list, item string) {
	c := ConnectRedis()
	if c != nil {
		defer CloseConn(c)
		r, err := c.Do("LPUSH", list, item)
		if err != nil {
			panic(err)
		}
		_ = r
	}
}

/*
*ListPop return the first element of a list,if the list have no element,
*block forever
 */
func ListPop(list string) string {
	var result string
	c := ConnectRedis()
	if c != nil {
		defer CloseConn(c)
		r, err := c.Do("LPOP", list)
		if err != nil {
			panic(err)
		}
		if r == nil {
			result = ""
		} else {
			result = string(r.([]uint8))
		}
	}
	return result
}

/*
*BlockListPop return the first element of a list,when the list we want to
*pop have no element,block at most timeout seconds
*input:	1)list name type of string
*		2)timeout second type of int64
*output:1)list name type of string
		2)element type of string
		if failed,return null string
*/
func BlockListPop(list string, timeout int64) (string, string) {
	c := ConnectRedis()
	if c != nil {
		defer CloseConn(c)
		r, err := c.Do("BLPOP", list, timeout)
		if err != nil {
			panic(err)
		}
		v, err := redis.Values(r, err)
		if err != nil {
			panic(err)
		}
		listname := string(v[0].([]uint8))
		item := string(v[1].([]uint8))
		return listname, item
	}
	return "", ""
}

/*
*HashSet set a field to value if the field is not exist,or update the
*value of the field
*input:	1)key which represent the hash table
*		2)field
*		3)value
*output:if the field is not exist,return 1,else return 0
 */
func HashSet(ht string, field string, value string) int64 {
	var result int64
	c := ConnectRedis()
	if c != nil {
		defer CloseConn(c)
		r, err := c.Do("HSET", ht, field, value)
		if err != nil {
			panic(err)
		}
		result = r.(int64)
	}
	return result
}

func HashGet(ht string, field string) string {
	var result string
	c := ConnectRedis()
	if c != nil {
		defer CloseConn(c)
		r, err := c.Do("HGET", ht, field)
		if err != nil {
			panic(err)
		}
		result = string(r.([]uint8))
	}
	return result
}

func HashIncrby(ht string, field string, increment int) int64 {
	var result int64
	c := ConnectRedis()
	if c != nil {
		defer CloseConn(c)
		r, err := c.Do("HINCRBY", ht, field, increment)
		if err != nil {
			panic(err)
		}
		result = r.(int64)
	}
	return result
}
