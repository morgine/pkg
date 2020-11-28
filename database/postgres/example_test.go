package postgres

import "github.com/morgine/pkg/config"

func ExampleNewPostgres() {
	var data = `
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
`

	configs, err := config.UnmarshalMemory([]byte(data))
	if err != nil {
		panic(err)
	}
	db, err := NewPostgres("postgres", configs)
	if err != nil {
		panic(err)
	}
	db.Query("select * from TABLE_NAME ")
}
