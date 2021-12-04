package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/morgine/pkg/config"
)

/**
# redis 数据库配置
[redis]
# redis 地址
addr = "localhost:6379"
# redis 密码
password = ""
# db 索引
db = 0
*/

type Config struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

func (e Config) Connect() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     e.Addr,
		Password: e.Password,
		DB:       e.DB,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func NewClient(namespace string, configs config.Configs) (*redis.Client, error) {
	cfg := &Config{}
	err := configs.UnmarshalSub(namespace, cfg)
	if err != nil {
		return nil, err
	}
	return cfg.Connect()
}
