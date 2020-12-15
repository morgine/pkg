package admin

import (
	"bytes"
	"fmt"
	"github.com/morgine/pkg/admin/pkg/crypt"
	"github.com/morgine/pkg/session"
	"gorm.io/gorm"
	"strconv"
	"time"
)

var Now = time.Now

type Handler struct {
	m    *model
	opts *Options
}

type Options struct {
	DB          *gorm.DB        // 数据库 ORM
	Session     session.Storage // token 存储器
	AuthExpires int64           // 会话过期时间
	AesCryptKey []byte          // 16 位字符串
}

func NewHandler(opts *Options) (*Handler, error) {
	err := opts.DB.AutoMigrate(&Admin{})
	if err != nil {
		return nil, err
	}
	return &Handler{
		m:    &model{opts.DB},
		opts: opts,
	}, nil
}

// 注册账号，如果账号已存在则返回 ErrUsernameAlreadyExist 错误
func (h *Handler) RegisterAdmin(username, password string) error {
	return h.m.RegisterAdmin(username, password)
}

// CheckAndRefreshToken 验证并刷新 token 过期时间
func (h *Handler) CheckAndRefreshToken(token string) (adminID int, err error) {
	if len(token) == 0 {
		return 0, nil
	}
	admin, err := h.decryptToken(token)
	if err != nil {
		return 0, err
	} else {
		ok, err := h.opts.Session.CheckAndRefreshToken(admin, token, h.opts.AuthExpires)
		if err != nil {
			return 0, err
		} else {
			if !ok {
				return 0, nil
			} else {
				return strconv.Atoi(admin)
			}
		}
	}
}

// 获得账户信息
func (h *Handler) GetAdmin(adminID int) (admin *Admin, err error) {
	return h.m.GetAdminByID(adminID)
}

// Login 登陆账号
func (h *Handler) Login(username, password string) (token string, err error) {
	admin, err := h.m.LoginAdmin(username, password)
	if err != nil {
		return "", err
	} else {
		token, err := h.encryptToken(strconv.Itoa(admin.ID))
		if err != nil {
			return "", err
		} else {
			err = h.opts.Session.SaveToken(strconv.Itoa(admin.ID), token, h.opts.AuthExpires)
			if err != nil {
				return "", err
			} else {
				return token, nil
			}
		}
	}
}

// ResetPassword 重置密码
func (h *Handler) ResetPassword(adminID int, newPassword string) error {
	err := h.m.ResetPassword(adminID, newPassword)
	if err != nil {
		return err
	}
	return h.opts.Session.RemoveUser(strconv.Itoa(adminID))
}

// Logout 退出登陆
func (h *Handler) Logout(adminID int, token string) error {
	return h.opts.Session.RemoveToken(strconv.Itoa(adminID), token)
}

// token 加密
func (h *Handler) encryptToken(adminID string) (token string, err error) {
	return crypt.AesCBCEncrypt([]byte(fmt.Sprintf("%s:%10d", adminID, Now().UnixNano())), h.opts.AesCryptKey)
}

// token 解密
func (h *Handler) decryptToken(token string) (adminID string, err error) {
	data, err := crypt.AesCBCDecrypt(token, h.opts.AesCryptKey)
	if err != nil {
		return "", err
	} else {
		sepIdx := bytes.Index(data, []byte(":"))
		return string(data[:sepIdx]), nil
	}
}
