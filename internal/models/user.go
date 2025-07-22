package models

import (
	"encoding/json"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Base
	Username       string          `gorm:"size:50;not null;unique" json:"username"`
	Email          string          `gorm:"size:100;not null;unique" json:"email"`
	Password       string          `gorm:"size:100;not null" json:"-"`
	FullName       string          `gorm:"size:100" json:"full_name"`
	Active         bool            `gorm:"default:true" json:"active"`
	CustomAttributes json.RawMessage `gorm:"type:json" json:"custom_attributes"`
	ThemeID        *uint           `json:"theme_id"`
	Theme          *Theme          `gorm:"foreignKey:ThemeID" json:"theme,omitempty"`
	UserRoles      []UserRole      `gorm:"foreignKey:UserID" json:"user_roles,omitempty"`
}

// BeforeSave - 保存前对密码进行哈希处理
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 仅当密码被修改时才进行哈希
	if tx.Statement.Changed("Password") {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword - 验证用户密码是否匹配
func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return errors.New("密码不匹配")
	}
	return nil
}

// GetUserAttributes - 获取用户自定义属性
func (u *User) GetUserAttributes() (map[string]interface{}, error) {
	if u.CustomAttributes == nil || len(u.CustomAttributes) == 0 {
		return make(map[string]interface{}), nil
	}

	var attributes map[string]interface{}
	err := json.Unmarshal(u.CustomAttributes, &attributes)
	return attributes, err
}

// SetUserAttributes - 设置用户自定义属性
func (u *User) SetUserAttributes(attributes map[string]interface{}) error {
	jsonData, err := json.Marshal(attributes)
	if err != nil {
		return err
	}
	u.CustomAttributes = jsonData
	return nil
}