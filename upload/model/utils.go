package model

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type condition interface {
	handle(db *gorm.DB) *gorm.DB
}

func handle(db *gorm.DB, cs ...condition) *gorm.DB {
	for _, c := range cs {
		db = c.handle(db)
	}
	return db
}

type Pagination struct {
	Limit  int
	Offset int
}

func (p Pagination) handle(db *gorm.DB) *gorm.DB {
	if p.Limit <= 0 || p.Limit > 500 {
		p.Limit = 10
	}
	db = db.Limit(p.Limit)
	if p.Offset != 0 {
		db = db.Offset(p.Offset)
	}
	return db
}

type OrderBy struct {
	OrderBy    string // 排序字段
	Descending bool   // 逆序
}

func (p OrderBy) handle(db *gorm.DB) *gorm.DB {
	if p.OrderBy != "" {
		db.Order(clause.OrderByColumn{Column: clause.Column{Name: p.OrderBy}, Desc: p.Descending})
	}
	return db
}
