package upload

import (
	"github.com/gin-gonic/gin"
	"github.com/morgine/pkg/upload/model"
	"io/ioutil"
	"path/filepath"
)

// 具有文件管理功能的处理方法集合
type SingleFileHandlers struct {
	singleDB *model.SingleFileDB
	chs      CommonHandlers
}

// NewMultiFileHandlers 创建文件管理器
//  singleDB 单文件管理数据库，如用户头像的管理等
//  multiDB 多文件管理数据库，如管理员在后台创建消息图片列表，文章图片列表等
//  chs 常用处理器，包含用户授权，错误处理
func NewSingleFileHandlers(singleDB *model.SingleFileDB, chs CommonHandlers) *SingleFileHandlers {
	return &SingleFileHandlers{
		singleDB: singleDB,
		chs:      chs,
	}
}

// 上传文件配置项
type CreateSingleFileOpts struct {
	Kind    model.Kind                                // 文件种类, 同一用户同一种类文件同时只能存在一个文件，该操作会覆盖已有的文件
	PostKey string                                    // 上传文件在 POST 表单中的 KEY 值
	Success func(f *model.SingleFile, c *gin.Context) // 处理成功则返回文件, 其中 url 已初始化
}

// 上传文件（同一种类限制只能上传一个文件，如果该种类的文件已存在，则将被覆盖）
func (hs *SingleFileHandlers) CreateSingleFile(opts CreateSingleFileOpts) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if userID, ok := hs.chs.GetAuthUser(ctx); ok {
			header, err := ctx.FormFile(opts.PostKey)
			if err != nil {
				hs.chs.HandleError(ctx, err)
				return
			}
			f, err := header.Open()
			if err != nil {
				hs.chs.HandleError(ctx, err)
				return
			}
			defer f.Close()
			data, err := ioutil.ReadAll(f)
			if err != nil {
				hs.chs.HandleError(ctx, err)
				return
			} else {
				file := &model.SingleFile{
					ID:     0,
					UserID: userID,
					Kind:   opts.Kind,
					File:   randStr(18) + filepath.Ext(header.Filename),
				}
				// 创建数据模型并返回文件服务URL地址
				err = hs.singleDB.Create(file, data)
				if err != nil {
					hs.chs.HandleError(ctx, err)
					return
				} else {
					opts.Success(file, ctx)
				}
			}
		}
	}
}

// 获得文件服务地址
func (hs *SingleFileHandlers) GetServeUrl(file string) (string, error) {
	return hs.singleDB.GetServeUrl(file)
}

// 获得文件服务地址
func (hs *SingleFileHandlers) GetFile(file string) ([]byte, error) {
	return hs.singleDB.GetFile(file)
}

// 具有文件管理功能的处理方法集合
type MultiFileHandlers struct {
	multiDB *model.MultiFileDB
	chs     CommonHandlers
}

// NewMultiFileHandlers 创建文件管理器
//  singleDB 单文件管理数据库，如用户头像的管理等
//  multiDB 多文件管理数据库，如管理员在后台创建消息图片列表，文章图片列表等
//  chs 常用处理器，包含用户授权，错误处理
func NewMultiFileHandlers(multiDB *model.MultiFileDB, chs CommonHandlers) *MultiFileHandlers {
	return &MultiFileHandlers{
		multiDB: multiDB,
		chs:     chs,
	}
}

// 常用处理器接口
type CommonHandlers interface {
	GetAuthUser(ctx *gin.Context) (userID int, ok bool) // 获得用户授权
	HandleError(ctx *gin.Context, err error)            // 处理错误信息
}

// 上传多文件配置项
type CreateMultiFileOpts struct {
	Kind    model.Kind                                  // 文件种类
	PostKey string                                      // 上传文件在 POST 表单中的 KEY 值
	Success func(fs []*model.MultiFile, c *gin.Context) // 处理成功则返回文件, 文件 Url 已初始化
	// TotalLimit int64                                       // 文件总量限制, 超过限制禁止上传
}

// 上传多个文件（同一种类允许包含多个文件）
func (hs *MultiFileHandlers) CreateMultiFile(opts CreateMultiFileOpts) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if userID, ok := hs.chs.GetAuthUser(ctx); ok {
			form, err := ctx.MultipartForm()
			if err != nil {
				hs.chs.HandleError(ctx, err)
				return
			}
			var files []*model.MultiFile
			for _, header := range form.File[opts.PostKey] {
				f, err := header.Open()
				if err != nil {
					hs.chs.HandleError(ctx, err)
					return
				}
				err = func() error {
					defer f.Close()
					data, err := ioutil.ReadAll(f)
					if err != nil {
						return err
					} else {
						file := &model.MultiFile{
							ID:     0,
							UserID: userID,
							Kind:   opts.Kind,
							File:   randStr(18) + filepath.Ext(header.Filename),
						}
						// 创建数据模型并返回文件服务URL地址
						err = hs.multiDB.Create(file, data)
						if err != nil {
							return err
						} else {
							files = append(files, file)
						}
					}
					return nil
				}()
				if err != nil {
					hs.chs.HandleError(ctx, err)
					return
				}
			}
			opts.Success(files, ctx)
		}
	}
}

type CountMultiFilesOpts struct {
	Kind    model.Kind                          // 文件类型
	Success func(total int64, ctx *gin.Context) // 处理成功返回文件总量
}

// 统计文件数量
func (hs *MultiFileHandlers) CountMultiFiles(opts CountMultiFilesOpts) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if userID, ok := hs.chs.GetAuthUser(ctx); ok {
			total, err := hs.multiDB.Count(model.UserKind{UserID: userID, Kind: opts.Kind})
			if err != nil {
				hs.chs.HandleError(ctx, err)
			} else {
				opts.Success(total, ctx)
			}
		}
	}
}

type GetMultiFilesOpts struct {
	Kind    model.Kind                                                      // 文件类型
	Params  func(ctx *gin.Context) (model.OrderBy, model.Pagination, error) // 提供限制条件参数
	Success func(fs []*model.MultiFile, ctx *gin.Context)                   // 处理成功返回带 url 地址的文件列表
}

// 获得文件列表
func (hs *MultiFileHandlers) GetMultiFiles(opts GetMultiFilesOpts) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if userID, ok := hs.chs.GetAuthUser(ctx); ok {
			orderBy, pagination, err := opts.Params(ctx)
			if err != nil {
				hs.chs.HandleError(ctx, err)
			} else {
				params := model.List{
					UserKind:   model.UserKind{UserID: userID, Kind: opts.Kind},
					Pagination: pagination,
					OrderBy:    orderBy,
				}
				files, err := hs.multiDB.Find(params)
				if err != nil {
					hs.chs.HandleError(ctx, err)
				} else {
					opts.Success(files, ctx)
				}
			}
		}
	}
}

type DelMultiFilesOpts struct {
	Kind    model.Kind                                        // 文件类型
	Params  func(ctx *gin.Context) (fileIDs []int, err error) // 参数提供需要删除的文件 ID
	Success func(leftTotal int64, ctx *gin.Context)           // 删除成功返回剩余文件总量
}

// 删除多个文件
func (hs *MultiFileHandlers) DelMultiFiles(opts DelMultiFilesOpts) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if userID, ok := hs.chs.GetAuthUser(ctx); ok {
			ids, err := opts.Params(ctx)
			if err != nil {
				hs.chs.HandleError(ctx, err)
			} else {
				total, err := hs.multiDB.Delete(model.UserKind{UserID: userID, Kind: opts.Kind}, ids)
				if err != nil {
					hs.chs.HandleError(ctx, err)
				} else {
					opts.Success(total, ctx)
				}
			}
		}
	}
}

// 获得文件服务地址
func (hs *MultiFileHandlers) GetServeUrl(file string) (string, error) {
	return hs.multiDB.GetServeUrl(file)
}

// 获得文件服务地址
func (hs *MultiFileHandlers) SetServeUrlGetters(getter func(file string) (url string, err error)) error {
	return hs.multiDB.SetServeUrlGetter(getter)
}

// 获得文件服务地址
func (hs *MultiFileHandlers) GetFile(file string) ([]byte, error) {
	return hs.multiDB.GetFile(file)
}
