package service

import (
	"github.com/KernyrMindDev/core/internal/objects"
)

// 白板状态 内存存储
type BoardState struct {
	Objects map[string]objects.BoardObject // 节点对象
	ID      string                         // 白板ID

}

// 实例化一个新的BoardState
func NewRoomState(boardID string) *BoardState {
	return &BoardState{
		Objects: make(map[string]objects.BoardObject), // TODO: 增加数据库查询
		ID:      boardID,                              // 白板ID
	}
}

func (board *BoardState) NewObject(id string, objectType string, data []byte)
