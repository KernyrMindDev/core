package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PublicHandler struct{}

func NewPublicHandler(_ *gorm.DB) *PublicHandler {
	return &PublicHandler{}
}

// PingHandler 	 Ping 接口
// @Summary      心跳检测
// @Description  用于检测服务健康状态
// @Tags         System
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string  "成功返回 pong"
// @Router       /ping [get]
func (h *PublicHandler) PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
