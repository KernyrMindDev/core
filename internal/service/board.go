package service

import (
	"uuid"

	"github.com/KernyrMindDev/core/internal/objects"
)

// BoardState不需要锁, BoardState的写入由Room的事件循环驱动, 而Room的事件循环则建立在channel上, 天然串行

// 白板状态 内存存储
type BoardState struct {
	Objects map[uuid.UUID]objects.BoardObject // 节点对象
	ID      string                            // 白板ID
}

// 实例化一个新的BoardState
func NewRoomState(boardID string) *BoardState {
	return &BoardState{
		Objects: make(map[uuid.UUID]objects.BoardObject), // TODO: 增加数据库查询
		ID:      boardID,                                 // 白板ID
	}
}

// 新建/替换
func (board *BoardState) Set(object objects.BoardObject) {
	id := object.ID()
	board.Objects[id] = object
}

// 获取
func (board *BoardState) Get(id uuid.UUID) objects.BoardObject {
	object, ok := board.Objects[id]
	if !ok {
		return nil
	}
	return object
}

// 删除
func (board *BoardState) Delete(id uuid.UUID) objects.BoardObject {
	object, ok := board.Objects[id]
	if !ok {
		return nil
	}
	delete(board.Objects, id)
	return object
}
