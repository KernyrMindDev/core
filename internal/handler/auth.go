package handler

import (
	"net/http"

	"github.com/KernyrMindDev/core/internal/dto"
	"github.com/KernyrMindDev/core/internal/model"

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
// @Router       /api/v1/auth/register [post]
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
			Error: "密码加密异常",
		})
		return
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "邮箱已被注册或入库失败",
		})
		return
	}

	// 密码不会被转换为json字段
	c.JSON(http.StatusOK, user)
}
