package service

import (
	"github.com/google/uuid"

	"github.com/KernyrMindDev/core/internal/model"
)

type Room struct {
	// 白板ID, 从哪个白板加载和保存数据
	boardID string
	// 运行时存储
	participants map[uuid.UUID]*model.Participant // 临时身份
	clients      map[uuid.UUID]*Client            // 实际连接
	// 消息通道
	join  chan *Client         // 新用户加入
	leave chan *Client         // 用户离开
	event chan ClientOperation // 用户操作事件
	done  chan struct{}        // 关闭房间
}

// 新建Room实例
func NewRoom(boardID string) *Room {
	return &Room{
		boardID: boardID,

		participants: make(map[uuid.UUID]*model.Participant),
		clients:      make(map[uuid.UUID]*Client),

		join:  make(chan *Client),
		leave: make(chan *Client),
		event: make(chan ClientOperation, 256), // 可能存在大量操作,提供缓冲

		done: make(chan struct{}),
	}
}

// Room实例事件循环协程
func (room *Room) Run() {
	for {
		select {
		case client := <-room.join:
			// 新加入用户
			room.handleJoin(client)
		case client := <-room.leave:
			// 用户离开
			room.handleLeave(client)
		case event := <-room.event:
			// 操作事件
			room.handleEvent(event)
		case <-room.done:
			// 房间关闭, 退出当前协程
			return
		}
	}
}

// 处理用户加入
func (room *Room) handleJoin(client *Client) {
	// 获取用户房间内UUID
	id := client.participant.ID
	// 通道确保了一次只有一个, 无需加锁
	room.participants[id] = client.participant
	room.clients[id] = client
}

// 处理用户离开
func (room *Room) handleLeave(client *Client) {
	id := client.participant.ID

	// Client存在且一致
	if current, ok := room.clients[id]; ok && current == client {
		delete(room.clients, id)
		delete(room.participants, id)
		// 关闭send通道, 让WriteMessage函数退出
		close(client.send)
	}
}

// 处理用户操作消息
func (room *Room) handleEvent(event ClientOperation) {
	// TODO: 后续实现
}

// 新用户加入
func (room *Room) Join(client *Client) {
	room.join <- client
}

// 用户离开
func (room *Room) Leave(client *Client) {
	room.leave <- client
}

// 提交新事件
func (room *Room) Submit(event ClientOperation) {
	room.event <- event
}
