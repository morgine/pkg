package admin

import "errors"

var (
	ErrUsernameAlreadyExist         = errors.New("用户名已存在")
	ErrMismatchedUsernameOrPassword = errors.New("用户名或密码错误")
)
