package rds

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/IanVzs/lightAPI/flag_parse"
	"github.com/IanVzs/lightAPI/log"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func getRds() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     *flag_parse.RedisIP,
		Password: *flag_parse.RedisPsw, // no password set
		DB:       *flag_parse.RedisDB,  // use default DB
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Logger.Error("redis connecnt error")
		// error.Error(err)
		return rdb
	}
	log.Logger.Info("redis connect: ", *flag_parse.RedisIP)
	return rdb
}

func getRdsFromPool() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		// 连接信息
		Addr:     *flag_parse.RedisIP,
		Password: *flag_parse.RedisPsw,
		DB:       *flag_parse.RedisDB,
		// 连接池容量及闲置连接数量
		PoolSize:     *flag_parse.PoolSize,
		MinIdleConns: *flag_parse.MinIdleConns,

		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔

		//可自定义连接函数
		// Dialer: func() (net.Conn, error) {
		// 	netDialer := &net.Dialer{
		// 		Timeout:   5 * time.Second,
		// 		KeepAlive: 5 * time.Minute,
		// 	}
		// 	return netDialer.Dial("tcp", "127.0.0.1:6379")
		// },

		// //钩子函数
		// OnConnect: func(conn *redis.Conn) error { //仅当客户端执行命令时需要从连接池获取连接时，如果连接池需要新建连接时则会调用此钩子函数
		// 	fmt.Printf("conn=%v\n", conn)
		// 	return nil
		// },
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Logger.Error("redis connecnt error")
		// error.Error(err)
		return rdb
	}
	log.Logger.Info("redis connect: ", *flag_parse.RedisIP)
	return rdb
}

// GetValueByKey: 从redis中获取数值 入参key
func GetValueByKey(key string) int {
	intRst := 0

	// rst, err := Client.Get(ctx, key).Result()
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

func printRedisPool(stats *redis.PoolStats) {
	fmt.Printf("Hits=%d Misses=%d Timeouts=%d TotalConns=%d IdleConns=%d StaleConns=%d\n",
		stats.Hits, stats.Misses, stats.Timeouts, stats.TotalConns, stats.IdleConns, stats.StaleConns)
}

func init() {
	// Client = getRds()
	Client = getRdsFromPool()
	printRedisPool(Client.PoolStats())
}
