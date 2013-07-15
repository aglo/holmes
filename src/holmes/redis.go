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

type RedisConn struct {
	c redis.Conn
}

func NewRedisConn() RedisConn {
	var redisConn RedisConn
	redisConn.c = ConnectRedis()
	return redisConn
}

func CloseRedisConn(redisConn RedisConn) {
	CloseConn(redisConn.c)
}

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

// ListLen return the lenght of a list
// output:the lenght of list
func (redisConn *RedisConn) ListLen(list string) int64 {
	var result int64
	if redisConn.c != nil {
		r, err := redisConn.c.Do("LLEN", list)
		if err != nil {
			panic(err)
		}
		result = r.(int64)
	}
	return result
}

// ListLeftPush push an item into a list at the left side of the list
// output:the lenght of list after push this item
func (redisConn *RedisConn) ListLeftPush(list, item string) int64 {
	var result int64
	if redisConn.c != nil {
		r, err := redisConn.c.Do("LPUSH", list, item)
		if err != nil {
			fmt.Println("error")
			panic(err)
		}
		result = r.(int64)
	}
	return result
}

// ListRightPush push an item into a list at the right side of the list
// output:the lenght of list after push this item
func (redisConn *RedisConn) ListRightPush(list, item string) int64 {
	var result int64
	if redisConn.c != nil {
		r, err := redisConn.c.Do("RPUSH", list, item)
		if err != nil {
			panic(err)
		}
		result = r.(int64)
	}
	return result
}

// ListLeftPop return the most left side element of a list
// output:if list a items return the most left side element,else,return null string
func (redisConn *RedisConn) ListLeftPop(list string) string {
	var result string
	if redisConn.c != nil {
		r, err := redisConn.c.Do("LPOP", list)
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

// ListRightPop return the most right side element of a list
// output:if list a items return the most right side element,else,return null string
func (redisConn *RedisConn) ListRightPop(list string) string {
	var result string
	if redisConn.c != nil {
		r, err := redisConn.c.Do("RPOP", list)
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

// BlockListLeftPop return the most left side element of a list,when the list we want to
// pop have no element,block at most timeout seconds
// input:
//     1)list name type of string;
//     2)timeout second type of int64
// output:
//     if success,return a <list,item> pair;else return a <"",""> pair
func (redisConn *RedisConn) BlockListLeftPop(list string, timeout int64) (string, string) {
	if redisConn.c != nil {
		r, err := redisConn.c.Do("BLPOP", list, timeout)
		if err != nil {
			panic(err)
		}
		if r != nil {
			v, err := redis.Values(r, err)
			if err != nil {
				panic(err)
			}
			listname := string(v[0].([]uint8))
			item := string(v[1].([]uint8))
			return listname, item
		}
	}
	return "", ""
}

// BlockListRightPop return the most right side element of a list,when the list we want to
// pop have no element,block at most timeout seconds
// input:
//     1)list name type of string;
//     2)timeout second type of int64
// output:
//     if success,return a <list,item> pair;else return a <"",""> pair
func (redisConn *RedisConn) BlockListRightPop(list string, timeout int64) (string, string) {
	if redisConn.c != nil {
		r, err := redisConn.c.Do("BRPOP", list, timeout)
		if err != nil {
			panic(err)
		}
		if r != nil {
			v, err := redis.Values(r, err)
			if err != nil {
				panic(err)
			}
			listname := string(v[0].([]uint8))
			item := string(v[1].([]uint8))
			return listname, item
		}
	}
	return "", ""
}

// HashSet set a field to value if the field is not exist,or update the value of
// the field
// input:
//     1)key which represent the hash table;
//     2)field;
//     3)value
// output:
//     if the field is not exist,return 1,else return 0
func (redisConn *RedisConn) HashSet(ht string, field string, value string) int64 {
	var result int64
	if redisConn.c != nil {
		r, err := redisConn.c.Do("HSET", ht, field, value)
		if err != nil {
			panic(err)
		}
		result = r.(int64)
	}
	return result
}

func (redisConn *RedisConn) HashGet(ht string, field string) string {
	var result string
	if redisConn.c != nil {
		r, err := redisConn.c.Do("HGET", ht, field)
		if err != nil {
			panic(err)
		}
		result = string(r.([]uint8))
	}
	return result
}

func (redisConn *RedisConn) HashIncrby(ht string, field string, increment int) int64 {
	var result int64
	if redisConn.c != nil {
		r, err := redisConn.c.Do("HINCRBY", ht, field, increment)
		if err != nil {
			panic(err)
		}
		result = r.(int64)
	}
	return result
}

func (redisConn *RedisConn) SetAdd(set string, member string) int64 {
	var result int64
	if redisConn.c != nil {
		r, err := redisConn.c.Do("SADD", set, member)
		if err != nil {
			panic(err)
		}
		result = r.(int64)
	}
	return result
}

// SetIsMember determine if a member is belong to a set
// if true,return 1,else return 0
func (redisConn *RedisConn) SetIsMember(set string, member string) int64 {
	var result int64
	if redisConn.c != nil {
		r, err := redisConn.c.Do("SISMEMBER", set, member)
		if err != nil {
			panic(err)
		}
		result = r.(int64)
	}
	return result
}
