package model

import "gorm.io/gorm"

type MultiFile struct {
	ID     int
	UserID int  `gorm:"index"`
	Kind   Kind `gorm:"index"`
	File   string
	Url    string `gorm:"-"`
}

type MultiFileDB struct {
	db      *gorm.DB
	storage Storage
}

func NewMultiFileDB(db *gorm.DB, s Storage) (*MultiFileDB, error) {
	err := db.AutoMigrate(&MultiFile{})
	if err != nil {
		return nil, err
	}
	return &MultiFileDB{
		db:      db,
		storage: s,
	}, nil
}

type UserKind struct {
	UserID int  // 用户ID限制条件
	Kind   Kind // 文件分类限制条件
}

func (c UserKind) handle(db *gorm.DB) *gorm.DB {
	if c.UserID != 0 {
		db = db.Where("user_id=?", c.UserID)
	}
	if c.Kind != 0 {
		db = db.Where("kind=?", c.Kind)
	}
	return db
}

// Create 创建文件并初始化服务地址
func (db *MultiFileDB) Create(file *MultiFile, data []byte) error {
	err := db.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(file).Error
		if err != nil {
			return err
		}
		return db.storage.CreateFile(file.File, data)
	})
	if err != nil {
		return err
	} else {
		file.Url, err = db.storage.GetServeUrl(file.File)
		if err != nil {
			return err
		}
	}
	return nil
}

// Count 统计数据总量
func (db *MultiFileDB) Count(uk UserKind) (total int64, err error) {
	err = handle(db.db.Model(&MultiFile{}), uk).Count(&total).Error
	if err != nil {
		return 0, err
	} else {
		return total, nil
	}
}

type List struct {
	UserKind
	OrderBy
	Pagination
}

// Find 查询多条数据, 并获得服务地址
func (db *MultiFileDB) Find(l List) (files []*MultiFile, err error) {
	err = handle(db.db, l.UserKind, l.OrderBy, l.Pagination).Find(&files).Error
	if err != nil {
		return nil, err
	} else {
		for _, file := range files {
			file.Url, err = db.storage.GetServeUrl(file.File)
			if err != nil {
				return nil, err
			}
		}
		return files, nil
	}
}

// Delete 删除多条数据并返回剩余数据量
func (db *MultiFileDB) Delete(uk UserKind, ids []int) (total int64, err error) {
	var files []*MultiFile
	err = handle(db.db.Where("id in (?)", ids), uk).Find(&files).Error
	if err != nil {
		return 0, err
	}
	err = db.db.Transaction(func(tx *gorm.DB) error {
		err = handle(db.db.Where("id in (?)", ids), uk).Delete(&MultiFile{}).Error
		if err != nil {
			return err
		}
		for _, file := range files {
			err = db.storage.DeleteFile(file.File)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return db.Count(uk)
}

// GetServeUrl 获得文件服务地址
func (db *MultiFileDB) GetServeUrl(file string) (string, error) {
	return db.storage.GetServeUrl(file)
}

// SetServeUrlGetter 获得文件服务地址
func (db *MultiFileDB) SetServeUrlGetter(getter func(file string) (url string, err error)) error {
	return db.storage.SetServeUrlGetter(getter)
}

// GetFile 获得文件内容
func (db *MultiFileDB) GetFile(file string) (data []byte, err error) {
	return db.storage.GetFile(file)
}
