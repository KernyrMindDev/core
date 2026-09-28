package handler

import (
	"errors"
	"fmt"
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
// @Failure      401  {object}  dto.ErrorResponse       "未登录"
// @Failure      403  {object}  dto.ErrorResponse       "无权限访问"
// @Failure      404  {object}  dto.ErrorResponse       "白板不存在"
// @Router       /boards/{id} [get]
func (h *BoardHandler) GetBoardDetail(c *gin.Context) {
	// 获取登录状态
	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:  "Unauthorized",
			Detail: "Lost login from middleware",
		})
		return
	}
	// 数据库查询
	boardID := c.Param("id")

	var board model.Board
	if err := h.DB.First(&board, "id = ?", boardID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "白板不存在"})
		return
	}
	// 判断权限
	if board.UserID != uid {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "无权限访问此白板"})
		return
	}

	// 查询所有节点与连线
	var nodes []model.BoardNode
	var connections []model.BoardConnection

	h.DB.Where("board_id = ?", boardID).Find(&nodes)
	h.DB.Where("board_id = ?", boardID).Find(&connections)

	c.JSON(http.StatusOK, dto.BoardDetailResponse{
		Board:       board,
		Connections: connections,
		Nodes:       nodes,
	})
}

// DeleteBoard   删除白板
// @Summary      删除指定白板
// @Description  根据白板ID删除白板的基础信息、所有节点及节点间的连线
// @Tags         Boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string                  true  "白板 ID" example("board_abc123")
// @Success      204   "删除成功"
// @Failure      401  {object}  dto.ErrorResponse       "未登录"
// @Failure      403  {object}  dto.ErrorResponse       "无权限访问"
// @Failure      404  {object}  dto.ErrorResponse       "白板不存在"
// @Failure      500  {object}  dto.ErrorResponse       "删除失败"
// @Router       /boards/{id} [delete]
func (h *BoardHandler) DeleteBoard(c *gin.Context) {
	// 获取登录状态
	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:  "Unauthorized",
			Detail: "Lost login from middleware",
		})
		return
	}
	// 提取白板ID
	boardID := c.Param("id")
	// 数据库查询
	var board model.Board
	if err := h.DB.First(&board, "id = ?", boardID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "白板不存在"})
		return
	}
	// 判断权限
	if board.UserID != uid {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: "无权限访问此白板"})
		return
	}
	// 开启数据库事务，确保原子性
	err := h.DB.Transaction(func(tx *gorm.DB) error {
		// 删除关联的节点
		if err := tx.Where("board_id = ?", boardID).Delete(&model.BoardNode{}).Error; err != nil {
			return err
		}

		// 删除关联的连线
		if err := tx.Where("board_id = ?", boardID).Delete(&model.BoardConnection{}).Error; err != nil {
			return err
		}

		// 删除白板本身
		if err := tx.Delete(&board).Error; err != nil {
			return err
		}

		// 返回 nil 自动提交事务，返回 error 自动回滚
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:  "删除白板及关联数据失败",
			Detail: err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetBoards 	 获取白板列表
// @Summary      获取白板列表
// @Description  获取当前账户下的所有白板
// @Tags         Boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.BoardsList "获取成功"
// @Failure      401  {object}  dto.ErrorResponse       "未登录"
// @Router       /boards [get]
func (h *BoardHandler) GetBoards(c *gin.Context) {
	// 获取UID
	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:  "Unauthorized",
			Detail: "Lost login from middleware",
		})
		return
	}

	// 数据库查询
	var boards []model.Board

	err := h.DB.Where("user_id = ?", uid).Find(&boards).Error
	if err != nil {
		// 数据库查询出错
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:  "Query Database failed",
			Detail: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.BoardsList{
		Boards: boards,
	})
}

// CreateBoard   创建新白板
// @Summary      在用户账户下创建新白板
// @Description  输入标题创建新白板
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
	// 获取UID
	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:  "Unauthorized",
			Detail: "Lost login from middleware",
		})
		return
	}

	var req dto.CreateBoardRequest
	// 解析请求体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	newBoard := model.Board{
		Title:  req.Title,
		UserID: uid,
	}
	// 写入数据库, UUID会自动生成
	if err := h.DB.Create(&newBoard).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "创建失败"})
		return
	}

	c.JSON(http.StatusOK, newBoard)
}

// UpdateBoard   更新白板信息
// @Summary      更新指定白板的信息
// @Description  更新指定白板的标题
// @Tags         Boards
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string       true  "白板 ID" example("board_abc123")
// @Param        request  body      dto.ChangeBoardDetailRequest  true  "白板创建参数"
// @Success      200      {object}  model.Board             "修改成功，返回白板信息"
// @Failure      400      {object}  dto.ErrorResponse       "参数不足"
// @Failure      403      {object}  dto.ErrorResponse       "权限不足"
// @Failure      404      {object}  dto.ErrorResponse       "白板不存在"
// @Failure      500      {object}  dto.ErrorResponse       "数据库独学而失败"
// @Router       /boards/{id} [put]
func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	// 获取UID
	uid := c.GetString("uid")
	if uid == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error:  "Unauthorized",
			Detail: "Lost login from middleware",
		})
		return
	}
	// 解析请求体
	var req dto.ChangeBoardDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	// 数据库查询
	var board model.Board
	if err := h.DB.First(&board, c.Param("id")).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:  "Not found",
				Detail: "Board not found",
			})
		} else {
			// 数据库错误
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:  "Not found",
				Detail: fmt.Sprintf("Database error: %v", err),
			})
		}
		return
	}
	// 权限检查
	if board.UserID != uid {
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:  "Permission denied",
			Detail: "Try to change another user's board",
		})
		return
	}
	// 修改
	board.Title = req.Title
	// 写数据库
	if err := h.DB.Save(&board).Error; err != nil {
		// 写入失败
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:  "Internet Server Errror",
			Detail: err.Error(),
		})
	} else {
		// 成功
		c.JSON(http.StatusOK, board)
	}
}
