package model

import (
	"gorm.io/gorm"
)

type Kind int

// SingleFile 单一文件模型，同一用户的同种类型只对应一个文件
type SingleFile struct {
	ID     int
	UserID int  `gorm:"index"`
	Kind   Kind `gorm:"index"`
	File   string
	Url    string `gorm:"-"`
}

type SingleFileDB struct {
	db      *gorm.DB
	storage Storage
}

func NewSingleFileDB(db *gorm.DB, s Storage) (*SingleFileDB, error) {
	err := db.AutoMigrate(&SingleFile{})
	if err != nil {
		return nil, err
	}
	return &SingleFileDB{
		db:      db,
		storage: s,
	}, nil
}

// First 获得单种文件，文件不存在并不返回错误
func (sfm *SingleFileDB) First(userID int, kind Kind) (*SingleFile, error) {
	file := &SingleFile{}
	err := sfm.db.Where("user_id=? AND kind=?", userID, kind).First(file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else {
		return file, nil
	}
}

// GetServeUrl 获得文件服务地址
func (sfm *SingleFileDB) GetServeUrl(file string) (string, error) {
	return sfm.storage.GetServeUrl(file)
}

// GetFile 获得文件内容
func (sfm *SingleFileDB) GetFile(file string) (data []byte, err error) {
	return sfm.storage.GetFile(file)
}

// Create 创建单种文件，如果文件已存在则自动覆盖已有文件, 创建完成之后将会自动初始化 url 地址
func (sfm *SingleFileDB) Create(file *SingleFile, data []byte) error {
	err := sfm.Delete(file.UserID, file.Kind)
	if err != nil {
		return err
	}
	return sfm.db.Transaction(func(tx *gorm.DB) error {
		err = sfm.db.Create(file).Error
		if err != nil {
			return err
		}
		err = sfm.storage.CreateFile(file.File, data)
		if err != nil {
			return err
		}
		file.Url, err = sfm.storage.GetServeUrl(file.File)
		if err != nil {
			return err
		}
		return nil
	})
}

// Delete 删除单种文件, 如果文件不存在则不做处理
func (sfm *SingleFileDB) Delete(userID int, k Kind) error {
	exist, err := sfm.First(userID, k)
	if err != nil {
		return err
	}
	if exist != nil {
		return sfm.db.Transaction(func(tx *gorm.DB) error {
			err = tx.Where("user_id=? AND kind=?", userID, k).Delete(&SingleFile{}).Error
			if err != nil {
				return nil
			}
			return sfm.storage.DeleteFile(exist.File)
		})
	}
	return nil
}
