package handler

import (
	"net/http"

	"github.com/KernyrMindDev/core/internal/dto"
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

// GetBoardDetail 获取白板详情
// @Summary      获取白板详情
// @Description  根据白板ID查询白板的基础信息、所有节点及节点间的连线
// @Tags         Boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "白板 ID" example("board_abc123")
// @Success      200  {object}  dto.BoardDetailResponse "获取成功"
// @Failure      404  {object}  dto.ErrorResponse       "白板不存在"
// @Router       /boards/{id} [get]
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

// CreateBoard 创建新白板
// @Summary      创建新白板
// @Description  输入标题和所属用户ID创建新白板
// @Tags         Boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateBoardRequest  true  "白板创建参数"
// @Success      200      {object}  model.Board             "创建成功，返回新白板信息"
// @Failure      400      {object}  dto.ErrorResponse       "参数验证失败"
// @Failure      500      {object}  dto.ErrorResponse       "数据库创建失败"
// @Router       /boards [post]
func (h *BoardHandler) CreateBoard(c *gin.Context) {
	var req dto.CreateBoardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	newBoard := model.Board{
		Title:  req.Title,
		UserID: req.UserID,
	}

	if err := h.DB.Create(&newBoard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "创建失败"})
		return
	}

	c.JSON(http.StatusOK, newBoard)
}
