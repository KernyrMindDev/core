package service

import (
	"errors"
	"net/http"

	"github.com/KernyrMindDev/core/internal/apperr"
	"github.com/KernyrMindDev/core/internal/model"

	"gorm.io/gorm"
)

// AuthService 封装用户注册/登录的业务逻辑。
type AuthService struct {
	DB *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{DB: db}
}

// Register 注册新用户
func (s *AuthService) Register(username, email, password string) (*model.User, error) {
	user := model.User{
		Username: username,
		Email:    email,
	}

	if err := user.SetPassword(password); err != nil {
		return nil, apperr.Internal("密码加密失败", err)
	}

	if err := s.DB.Create(&user).Error; err != nil {
		return nil, &apperr.AppError{
			HTTPStatus: http.StatusBadRequest,
			Message:    "邮箱已被注册",
			Err:        err,
		}
	}

	return &user, nil
}

// Login 校验邮箱密码并返回登录成功后的 JWT Token
func (s *AuthService) Login(email, password string) (string, error) {
	var user model.User
	err := s.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperr.Unauthorized("邮箱或密码错误")
		}
		return "", apperr.Internal("查询用户失败", err)
	}

	if !user.CheckPassword(password) {
		return "", apperr.Unauthorized("邮箱或密码错误")
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		return "", apperr.Internal("生成登录凭证失败", err)
	}

	return token, nil
}
