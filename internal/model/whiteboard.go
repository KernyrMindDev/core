package model

import "time"

// 白板基础信息
type Board struct {
	ID        string    `gorm:"primaryKey;size:64" json:"id"`
	UserID    string    `gorm:"index;size:36;not null" json:"userId"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// 节点表
type BoardNode struct {
	ID      string  `gorm:"primaryKey;size:64" json:"id"`
	BoardID string  `gorm:"primaryKey;index;size:64" json:"boardId"`
	Type    string  `gorm:"size:32;not null" json:"type"` // 如 "rect", "circle"
	X       float64 `gorm:"not null" json:"x"`
	Y       float64 `gorm:"not null" json:"y"`
	Width   float64 `gorm:"not null" json:"width"`
	Height  float64 `gorm:"not null" json:"height"`
	Data    string  `gorm:"type:text" json:"data"` // 存储文本内容、颜色等自定义元数据
}

// 连线表
type BoardConnection struct {
	ID         string `gorm:"primaryKey;size:64" json:"id"`
	BoardID    string `gorm:"primaryKey;index;size:64" json:"boardId"`
	SourceID   string `gorm:"size:64;not null" json:"sourceId"`
	SourcePort string `gorm:"size:16;not null" json:"sourcePort"` // top, bottom, left, right, center
	TargetID   string `gorm:"size:64;not null" json:"targetId"`
	TargetPort string `gorm:"size:16;not null" json:"targetPort"`
	Style      string `gorm:"type:text" json:"style"` // 箭头样式、粗细等 JSON 字符串
}
