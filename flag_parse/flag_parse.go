package flag_parse

import (
	"flag"
	"runtime"
)

// main 相关启动配置
var Addr = flag.String("addr", ":8080", "http service address")

// log 相关启动配置
var LogPath = flag.String("log_path", "./logs/ligthAPI.log", "set app log path.")
var LogLevel = flag.String("log_level", "debug", "set app log level. 默认debug")

// Redis 相关启动配置
var RedisIP = flag.String("redis_ip", "192.168.0.1:6379", "redis address")
var RedisPsw = flag.String("redis_passwd", "", "redis passwd")
var RedisDB = flag.Int("redis_db", 0, "redis db")

// Redis 连接池配置
var PoolSize = flag.Int("rpool_size", runtime.NumCPU(), "连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU")
var MinIdleConns = flag.Int("ridle_conns", 5, "在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量")

// 接口加密密钥
var KeyAPI = flag.String("key_api", "key_test", "接口加/解密密钥")
var KeySplit = flag.String("key_split", "_|_", "接口加/解密数据分割符")

// 放在初始化方法中执行日志的初始化

func init() {
	flag.Parse()
}
