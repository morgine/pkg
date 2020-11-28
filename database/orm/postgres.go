package orm

import (
	"github.com/morgine/pkg/config"
	"github.com/morgine/pkg/database/postgres"
	"gorm.io/gorm"
)

func NewPostgresORM(postgresNamespace, gormNamespace string, configs config.Configs) (*gorm.DB, error) {
	db, err := postgres.NewPostgres(postgresNamespace, configs)
	if err != nil {
		return nil, err
	}
	gormConfig := &Config{}
	err = configs.UnmarshalSub(gormNamespace, gormConfig)
	if err != nil {
		return nil, err
	}
	return gormConfig.Init(NewPostgresDialector(db))
}
