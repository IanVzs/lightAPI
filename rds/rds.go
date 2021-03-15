package rds

import (
	"context"
	"flag"
	"strconv"

	"github.com/IanVzs/lightAPI/log"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func getRds() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     *log.RedisIP,
		Password: *log.RedisPsw, // no password set
		DB:       *log.RedisDB,  // use default DB
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Logger.Error("redis connecnt error")
		// error.Error(err)
		return rdb
	}
	log.Logger.Info("redis connect: ", *log.RedisIP)
	return rdb
}

// GetValueByKey: 从redis中获取数值 入参key
func GetValueByKey(key string) int {
	intRst := 0

	rst, err := Client.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Logger.Warn("rds.GetValueByKey no get key:", key)
		rst = "0"
	} else if err != nil {
		panic(err)
	} else {
		intRst, err = strconv.Atoi(rst)
		if err != nil {
			panic(err)
		}
		log.Logger.Info("rds.GetValueByKey value: ", intRst)
	}
	return intRst
}

// 从redis中获取数据并删除
func GetDelByKey(key string) int {
	intRst := GetValueByKey(key)
	_, err := Client.Del(ctx, key).Result()
	if err == redis.Nil {
		log.Logger.Errorf("rds.GetDelByKey not found key: %s", key)
	} else if err != nil {
		panic(err)
	}
	return intRst
}

var Client *redis.Client

func init() {
	flag.Parse()
	Client = getRds()
}
