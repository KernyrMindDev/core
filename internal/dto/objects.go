package dto

import "uuid"

// 创建对象
type ObjectCreateRequest struct {
	Type string `json:"type"`
}

// 返回前端的节点对象的数据结构
type ObjectDTO struct {
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"`
	Data any       `json:"data"`
}
