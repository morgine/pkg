package admin

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/morgine/pkg/admin/pkg/crypt"
	"github.com/morgine/pkg/admin/pkg/xtime"
	"github.com/morgine/pkg/session"
	"gorm.io/gorm"
	"strconv"
)

var (
	ErrUsernameAlreadyExist         = errors.New("用户名已存在")
	ErrMismatchedUsernameOrPassword = errors.New("用户名或密码错误")
	ErrUnauthorized                 = errors.New("未授权")
)

type Handler struct {
	m           *model
	session     session.Storage
	authExpires int64  // 会话过期时间
	aesCryptKey []byte // 16 位字符串
	sender      MessageSender
}

type Options struct {
	DB      *gorm.DB
	Session session.Storage

	AuthExpires int64  // 会话过期时间
	AesCryptKey []byte // 16 位字符串
	Sender      MessageSender
}

func NewHandler(opts *Options) (*Handler, error) {
	err := opts.DB.AutoMigrate(&Admin{})
	if err != nil {
		return nil, err
	}
	return &Handler{
		m:           &model{opts.DB},
		session:     opts.Session,
		authExpires: opts.AuthExpires,
		aesCryptKey: opts.AesCryptKey,
		sender:      opts.Sender,
	}, nil
}

// 注册账号，如果账号已存在，则返回 ErrUsernameAlreadyExist 错误
func (h *Handler) RegisterAdmin(username, password string) error {
	return h.m.RegisterAdmin(username, password)
}

func (h *Handler) Auth(authorizationKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get(authorizationKey)
		if token == "" {
			h.sender.SendError(ctx, ErrUnauthorized)
		} else {
			adminID, err := h.decryptToken(token)
			if err != nil {
				h.sender.SendError(ctx, err)
			} else {
				ok, err := h.session.CheckAndRefreshToken(strconv.Itoa(adminID), token, h.authExpires)
				if err != nil {
					h.sender.SendError(ctx, err)
				} else {
					if !ok {
						h.sender.SendError(ctx, ErrUnauthorized)
					} else {
						ctx.Set("auth_admin", adminID)
					}
				}
			}
		}
	}
}

func (h *Handler) GetAuthAdmin(ctx *gin.Context) (adminID int, ok bool) {
	v, ok := ctx.Get("auth_admin")
	if ok {
		return v.(int), true
	} else {
		return 0, false
	}
}

func (h *Handler) GetLoginAdmin(ctx *gin.Context) {
	adminID, ok := h.GetAuthAdmin(ctx)
	if !ok {
		h.sender.SendError(ctx, ErrUnauthorized)
	} else {
		admin, err := h.m.GetAdminByID(adminID)
		if err != nil {
			h.sender.SendError(ctx, err)
		} else {
			h.sender.SendData(ctx, admin)
		}
	}
}

func (h *Handler) Login() gin.HandlerFunc {
	type params struct {
		Username string
		Password string
	}
	return func(ctx *gin.Context) {
		ps := &params{}
		err := ctx.Bind(ps)
		if err != nil {
			h.sender.SendError(ctx, err)
		} else {
			admin, err := h.m.LoginAdmin(ps.Username, ps.Password)
			if err != nil {
				h.sender.SendError(ctx, err)
			} else {
				token, err := h.encryptToken(admin.ID)
				if err != nil {
					h.sender.SendError(ctx, err)
				} else {
					err = h.session.SaveToken(strconv.Itoa(admin.ID), token, h.authExpires)
					if err != nil {
						h.sender.SendError(ctx, err)
					} else {
						h.sender.SendData(ctx, token)
					}
				}
			}
		}
	}
}

func (h *Handler) ResetPassword() gin.HandlerFunc {

	type params struct {
		NewPassword string
	}
	return func(ctx *gin.Context) {
		adminID, ok := h.GetAuthAdmin(ctx)
		if !ok {
			h.sender.SendError(ctx, ErrUnauthorized)
			return
		}
		ps := &params{}
		err := ctx.Bind(ps)
		if err != nil {
			h.sender.SendError(ctx, err)
		} else {
			err = h.m.ResetPassword(adminID, ps.NewPassword)
			if err != nil {
				h.sender.SendError(ctx, err)
			} else {
				token, err := h.encryptToken(adminID)
				if err != nil {
					h.sender.SendError(ctx, err)
				} else {
					err = h.session.SaveToken(strconv.Itoa(adminID), token, h.authExpires)
					if err != nil {
						h.sender.SendError(ctx, err)
					} else {
						h.sender.SendData(ctx, token)
					}
				}
			}
		}
	}
}

func (h *Handler) Logout(ctx *gin.Context) {
	admin, ok := h.GetAuthAdmin(ctx)
	if !ok {
		h.sender.SendMsgSuccess(ctx, "已退出")
	} else {
		err := h.session.DelToken(strconv.Itoa(admin))
		if err != nil {
			h.sender.SendError(ctx, err)
		} else {
			h.sender.SendMsgSuccess(ctx, "已退出")
		}
	}
}

// token 加密
func (h *Handler) encryptToken(adminID int) (token string, err error) {
	return crypt.AesCBCEncrypt([]byte(fmt.Sprintf("%d:%10d", adminID, xtime.Now().UnixNano())), h.aesCryptKey)
}

// token 解密
func (h *Handler) decryptToken(token string) (adminID int, err error) {
	data, err := crypt.AesCBCDecrypt(token, h.aesCryptKey)
	if err != nil {
		return 0, err
	} else {
		sepIdx := bytes.Index(data, []byte(":"))
		return strconv.Atoi(string(data[:sepIdx]))
	}
}
