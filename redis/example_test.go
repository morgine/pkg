package redis

import (
	"context"
	"github.com/morgine/pkg/config"
)

func ExampleNewClient() {
	var data = `
# redis 数据库配置
[redis]
# redis 地址
addr = "localhost:6379"
# redis 密码
password = ""
# db 索引
db = 1
`

	configs, err := config.UnmarshalMemory([]byte(data))
	if err != nil {
		panic(err)
	}
	client, err := NewClient("redis", configs)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}
}
