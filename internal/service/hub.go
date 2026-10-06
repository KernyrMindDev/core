package service

import (
	"sync"

	"github.com/KernyrMindDev/core/internal/objects"
)

type Hub struct {
	// 内存存储
	rooms map[string]*Room
	// 同步锁
	mu sync.Mutex
	// 对象管理器
	objectFactory *objects.ObjectFactory
}

// 构造函数
func NewHub(objectFactory *objects.ObjectFactory) *Hub {
	return &Hub{
		rooms:         make(map[string]*Room),
		mu:            sync.Mutex{},
		objectFactory: objectFactory,
	}
}

// 获取或创建一个Room
func (hub *Hub) CreateOrGetRoom(boardID string) *Room {
	// 获取同步锁
	hub.mu.Lock()
	defer hub.mu.Unlock()
	// 查询
	room, ok := hub.rooms[boardID]
	if ok {
		return room
	}
	// 没查到, 新建
	room = NewRoom(boardID, hub.objectFactory)
	hub.rooms[boardID] = room
	// 启动Room处理协程
	go room.Run()
	return room
}
