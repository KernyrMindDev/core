package handler

import (
	"net/http"

	"github.com/KernyrMindDev/core/internal/dto"
	"github.com/KernyrMindDev/core/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BoardHandler 只负责解析请求 / 调用 Service / 组装响应，
// 不直接操作数据库，业务逻辑(权限校验、事务)由 BoardService 承担。
type BoardHandler struct {
	svc *service.BoardService
}

// 构造函数
func NewBoardHandler(db *gorm.DB) *BoardHandler {
	return &BoardHandler{svc: service.NewBoardService(db)}
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
	uid := c.GetString("uid")
	boardID := c.Param("id")

	detail, err := h.svc.GetBoardDetail(uid, boardID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, dto.BoardDetailResponse{
		Board:       detail.Board,
		Nodes:       detail.Nodes,
		Connections: detail.Connections,
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
	uid := c.GetString("uid")
	boardID := c.Param("id")

	if err := h.svc.DeleteBoard(uid, boardID); err != nil {
		c.Error(err)
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
	uid := c.GetString("uid")

	boards, err := h.svc.ListBoards(uid)
	if err != nil {
		c.Error(err)
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
	uid := c.GetString("uid")

	var req dto.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	board, err := h.svc.CreateBoard(uid, req.Title)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, board)
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
// @Failure      500      {object}  dto.ErrorResponse       "数据库读写失败"
// @Router       /boards/{id} [put]
func (h *BoardHandler) UpdateBoard(c *gin.Context) {
	uid := c.GetString("uid")

	var req dto.ChangeBoardDetailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	board, err := h.svc.UpdateBoard(uid, c.Param("id"), req.Title)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, board)
}
