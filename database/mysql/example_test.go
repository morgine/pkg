package mysql

import "github.com/morgine/pkg/config"

func ExampleNewMysql() {
	var data = `
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
`
	configs, err := config.UnmarshalMemory([]byte(data))
	if err != nil {
		panic(err)
	}
	db, err := NewMysql("mysql", configs)
	if err != nil {
		panic(err)
	}
	db.Query("select * from TABLE_NAME ")
}
