package admin

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	ID       int
	Username string `gorm:"index"`
	Password string
}

type model struct {
	db *gorm.DB
}

// 注册账号，如果账号已存在则返回 ErrUsernameAlreadyExist 错误
func (m *model) RegisterAdmin(username, password string) (err error) {
	admin, err := m.GetAdminByUsername(username)
	if err != nil {
		return err
	}
	if admin != nil {
		return ErrUsernameAlreadyExist
	} else {
		password, err := bcrypt.GenerateFromPassword([]byte(password), 10)
		if err != nil {
			return err
		}
		return m.db.Create(&Admin{Username: username, Password: string(password)}).Error
	}
}

func (m *model) LoginAdmin(username, password string) (*Admin, error) {
	admin, err := m.GetAdminByUsername(username)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, ErrMismatchedUsernameOrPassword
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return nil, ErrMismatchedUsernameOrPassword
			} else {
				return nil, err
			}
		} else {
			return admin, nil
		}
	}
}

func (m *model) GetAdminByUsername(username string) (*Admin, error) {
	admin := &Admin{}
	err := m.db.First(admin, "username=?", username).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if admin.ID > 0 {
		return admin, nil
	} else {
		return nil, nil
	}
}

func (m *model) GetAdminByID(id int) (*Admin, error) {
	admin := &Admin{}
	err := m.db.First(admin, "id=?", id).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if admin.ID > 0 {
		return admin, nil
	} else {
		return nil, nil
	}
}

func (m *model) ResetPassword(authAdminID int, newPassword string) error {
	password, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return err
	}
	return m.db.Where("id=?", authAdminID).Updates(&Admin{Password: string(password)}).Error
}
