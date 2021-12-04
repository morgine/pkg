package orm

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

/**
# gorm 配置
[gorm]
# 日志等级 1-Silent, 2-Error, 3-Warn, 4-Info
log_level = 4
# 数据库表名前缀
table_prefix = ""
# 使用单数表名
singular_table = false
*/

type Config struct {
	LogLevel      logger.LogLevel `toml:"log_level"`
	TablePrefix   string          `toml:"table_prefix"`
	SingularTable bool            `toml:"singular_table"`
}

func (e *Config) Init(dialector gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(e.LogLevel),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   e.TablePrefix,
			SingularTable: e.SingularTable,
		},
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func NewMysqlDialector(db *sql.DB) gorm.Dialector {
	return mysql.New(mysql.Config{Conn: db})
}

func NewPostgresDialector(db *sql.DB) gorm.Dialector {
	return postgres.New(postgres.Config{Conn: db})
}
