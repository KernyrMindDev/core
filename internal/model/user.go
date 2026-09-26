package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// 用户数据
type User struct {
	ID           string    `gorm:"primaryKey;size:36" json:"id"` // 用户UUID
	Username     string    `gorm:"size:64;not null" json:"username"`
	Email        string    `gorm:"size:128;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:128;column:password_hash;not null" json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
}

// 将用户的明文密码哈希加密
func (u *User) SetPassword(plainPassword string) error {
	// bcrypt对密码加盐哈希
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(bytes)
	return nil
}

// 校验传入的明文密码是否与哈希匹配
func (u *User) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plainPassword))
	return err == nil
}
