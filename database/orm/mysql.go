package orm

import (
	"github.com/morgine/pkg/config"
	"github.com/morgine/pkg/database/mysql"
	"gorm.io/gorm"
)

func NewMysqlORM(mysqlNamespace, gormNamespace string, configs config.Configs) (*gorm.DB, error) {
	db, err := mysql.NewMysql(mysqlNamespace, configs)
	if err != nil {
		return nil, err
	}
	gormConfig := &Config{}
	err = configs.UnmarshalSub(gormNamespace, gormConfig)
	if err != nil {
		return nil, err
	}
	return gormConfig.Init(NewMysqlDialector(db))
}
