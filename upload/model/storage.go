package model

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// Storage 数据存储器
type Storage interface {
	// CreateFile 创建/上传文件
	CreateFile(file string, data []byte) error
	// DeleteFile 删除文件，如果文件不存在则视作删除成功，禁止返回错误
	DeleteFile(file string) error
	// GetFile 获得字节文件
	GetFile(file string) (data []byte, err error)
	// GetServeUrl 获得文件服务地址，有可能是远程服务或本地服务
	GetServeUrl(file string) (url string, err error)
	// SetServeUrlGetter 设置文件服务地址提供器
	SetServeUrlGetter(urlGetter func(file string) (url string, err error)) error
}

// FileStorage 本地存储器
type FileStorage struct {
	Dir       string
	UrlGetter func(file string) (url string, err error)
}

func (f *FileStorage) SetServeUrlGetter(urlGetter func(file string) (url string, err error)) error {
	f.UrlGetter = urlGetter
	return nil
}

func NewFileStorage(dir string) (Storage, error) {
	return &FileStorage{
		Dir: dir,
	}, os.MkdirAll(dir, os.ModePerm)
}

func (f *FileStorage) GetFile(file string) (data []byte, err error) {
	return ioutil.ReadFile(filepath.Join(f.Dir, file))
}

func (f *FileStorage) CreateFile(file string, data []byte) error {
	return ioutil.WriteFile(filepath.Join(f.Dir, file), data, os.ModePerm)
}

func (f *FileStorage) DeleteFile(file string) error {
	err := os.Remove(filepath.Join(f.Dir, file))
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	return err
}

func (f *FileStorage) GetServeUrl(file string) (url string, err error) {
	return f.UrlGetter(file)
}
