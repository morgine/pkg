package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/morgine/pkg/config"
	"github.com/morgine/pkg/database"
)

/**
# mysql 数据库配置
[mysql]
# 连接地址
host = "127.0.0.1"
# 连接端口
port = 3306
# 用户名
user = "root"
# 密码
password = "123456"
# 数据库
db_name = ""
# 连接参数
parameters = "charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true"
# 最长等待断开时间(单位: 秒), 如果该值为 0, 则不限制时间
max_lifetime = 0
# 最多打开数据库的连接数量, 如果该值为 0, 则不限制连接数量
max_open_conns = 10
# 连接池中最多空闲链接数量, 如果该值为 0, 则不保留空闲链接
max_idle_conns = 10
*/

type Config struct {
	Host       string `toml:"host"`
	Port       int    `toml:"port"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	DBName     string `toml:"db_name"`
	Parameters string `toml:"parameters"`
	database.Config
}

// DSN 数据库连接串
func (e Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		e.User, e.Password, e.Host, e.Port, e.DBName, e.Parameters)
}

func (e Config) Connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", e.DSN())
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

func NewMysql(namespace string, configs config.Configs) (*sql.DB, error) {
	cfg := &Config{}
	err := configs.UnmarshalSub(namespace, cfg)
	if err != nil {
		return nil, err
	}
	return cfg.Connect()
}
