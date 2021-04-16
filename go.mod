module github.com/IanVzs/lightAPI

go 1.15

replace github.com/IanVzs/lightAPI/chat => ./chat

replace github.com/IanVzs/lightAPI/rds => ./rds

replace github.com/IanVzs/lightAPI/log => ./log

replace github.com/IanVzs/lightAPI/flag_parse => ./flag_parse

require (
	github.com/IanVzs/lightAPI/chat v0.0.0-00010101000000-000000000000
	github.com/IanVzs/lightAPI/flag_parse v0.0.0-00010101000000-000000000000
	github.com/IanVzs/lightAPI/log v0.0.0-00010101000000-000000000000
	github.com/IanVzs/lightAPI/rds v0.0.0-00010101000000-000000000000
	github.com/go-redis/redis/v8 v8.6.0 // indirect
	github.com/natefinch/lumberjack v2.0.0+incompatible // indirect
	go.uber.org/zap v1.16.0 // indirect
)
