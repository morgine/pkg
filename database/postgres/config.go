package postgres

import (
	"database/sql"
	"fmt"
	"github.com/morgine/pkg/config"
	"github.com/morgine/pkg/database"
)

/**
# postgres 数据库配置
[postgres]
# 连接地址
host = "127.0.0.1"
# 连接端口
port= 5432
# 用户名
user = "root"
# 密码
password = "123456"
# 数据库
db_name = "ginadmin"
# SSL模式
ssl_mode = "disable"
# 最长等待断开时间(单位: 秒), 如果该值为 0, 则不限制时间
max_lifetime = 0
# 最多打开数据库的连接数量, 如果该值为 0, 则不限制连接数量
max_open_conns = 10
# 连接池中最多空闲链接数量, 如果该值为 0, 则不保留空闲链接
max_idle_conns = 10
*/
type Config struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"db_name"`
	SSLMode  string `toml:"ssl_mode"`
	database.Config
}

// DSN 数据库连接串
func (e Config) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		e.Host, e.Port, e.User, e.DBName, e.Password, e.SSLMode)
}

func (e Config) Connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", e.DSN())
	if err != nil {
		return nil, err
	}
	err = e.Config.Init(db)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func NewPostgres(namespace string, configs config.Configs) (*sql.DB, error) {
	cfg := &Config{}
	err := configs.UnmarshalSub(namespace, cfg)
	if err != nil {
		return nil, err
	}
	return cfg.Connect()
}
