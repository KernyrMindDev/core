package dto

import "github.com/KernyrMindDev/core/internal/model"

// CreateBoardRequest 创建白板请求参数
type CreateBoardRequest struct {
	Title string `json:"title" binding:"required" example:"示例白板"`
}

// BoardDetailResponse 白板详情响应数据
type BoardDetailResponse struct {
	Board       model.Board             `json:"board"`
	Nodes       []model.BoardNode       `json:"nodes"`
	Connections []model.BoardConnection `json:"connections"`
}

// BoardsList 白板列表响应数据
type BoardsList struct {
	Boards []model.Board `json:"boards"`
}
