package handler

import (
	"net/http"

	"github.com/KernyrMindDev/core/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BoardHandler 持有数据库连接
type BoardHandler struct {
	DB *gorm.DB
}

// 构造函数
func NewBoardHandler(db *gorm.DB) *BoardHandler {
	return &BoardHandler{DB: db}
}

// 获取白板的全量节点和连线数据
func (h *BoardHandler) GetBoardDetail(c *gin.Context) {
	boardID := c.Param("id")

	var board model.Board
	if err := h.DB.First(&board, "id = ?", boardID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "白板不存在"})
		return
	}

	// 查询所有节点与连线
	var nodes []model.BoardNode
	var connections []model.BoardConnection

	h.DB.Where("board_id = ?", boardID).Find(&nodes)
	h.DB.Where("board_id = ?", boardID).Find(&connections)

	c.JSON(http.StatusOK, gin.H{
		"board":       board,
		"nodes":       nodes,
		"connections": connections,
	})
}

// 创建新白板
func (h *BoardHandler) CreateBoard(c *gin.Context) {
	var req struct {
		Title  string `json:"title" binding:"required"`
		UserID string `json:"userId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBoard := model.Board{
		Title:  req.Title,
		UserID: req.UserID,
	}

	if err := h.DB.Create(&newBoard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, newBoard)
}
