package handler

import (
	"errors"
	"net/http"

	"github.com/KernyrMindDev/core/internal/dto"
	"github.com/KernyrMindDev/core/internal/model"
	"github.com/KernyrMindDev/core/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

// Register 用户注册
// @Summary      用户注册
// @Description  通过用户名、邮箱和密码注册新账号
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request  body      dto.RegisterRequest  true  "注册信息"
// @Success      200      {object}  model.User       "注册成功，返回用户信息"
// @Failure      400      {object}  dto.ErrorResponse    "参数验证失败或邮箱已被注册"
// @Failure      500      {object}  dto.ErrorResponse    "密码加密异常等服务器错误"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
	}

	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:  "Password encode failed",
			Detail: err.Error(),
		})
		return
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:  "Write Database failed",
			Detail: err.Error(),
		})
		return
	}

	// 密码不会被转换为json字段
	c.JSON(http.StatusOK, user)
}

// Login 用户登录
// @Summary      用户登录
// @Description  通过邮箱和密码登录已有账号
// @Tags         Auth
// @Accept		 json
// @Produce 	 json
// @Param        request  body      dto.LoginRequest  true  "登录信息"
// @Success      200      {object}  dto.LoginSuccess       "登录成功，返回用户JWT Token"
// @Failure      400      {object}  dto.ErrorResponse    "请求体参数错误或不足"
// @Failure      401      {object}  dto.ErrorResponse    "账户或密码错误"
// @Failure		 403 	  {object}  dto.ErrorResponse	   "账户或IP被封禁"
// @Failure		 500	  {object}  dto.ErrorResponse  "服务器内部错误"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 解析失败
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}
	// 查库
	var user model.User
	err := h.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error:  "邮箱或密码错误",
				Detail: "账户不存在",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:  "服务器内部错误",
				Detail: err.Error(),
			})
			return
		}
	}
	// 校验密码
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "邮箱或密码错误",
		})
		return
	}
	// 生成Token
	token, err := service.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:  "服务器内部错误",
			Detail: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.LoginSuccess{
		Token: token,
	})
}
