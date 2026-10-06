package service

import (
	"github.com/KernyrMindDev/core/internal/dto"
	"github.com/KernyrMindDev/core/internal/model"
	"github.com/gorilla/websocket"
)

// 每个连接的对象
type Client struct {
	conn   *websocket.Conn          // 实际连接
	member *model.Member            // 用户对象
	room   *Room                    // 房间对象
	send   chan dto.WSServerMessage // 向客户端下发的消息通道
}

// 实例化一个新的Clinet
func NewClient(
	conn *websocket.Conn,
	member *model.Member,
	room *Room,
) *Client {
	return &Client{
		conn:   conn,
		member: member,
		room:   room,
		send:   make(chan dto.WSServerMessage, 64),
	}
}

// 循环读取客户端上报消息并封装后传递的协程
func (c *Client) ReadLoop() {
	defer func() {
		// 用户离开后向room上报以清理连接和协程
		c.room.Leave(c)
		_ = c.conn.Close() // 关闭WebSocket连接
	}()

	for {
		// 准备解析
		var msg dto.WSClientMessage
		// 从WebSocket读取并进行解析
		if err := c.conn.ReadJSON(&msg); err != nil {
			break
		}
		// 提交事件
		c.room.Submit(ClientOperation{
			Client:  c,
			Message: msg,
		})
	}
}

// 向客户端发送消息
func (c *Client) WriteLoop() {
	for msg := range c.send {
		if err := c.conn.WriteJSON(msg); err != nil {
			break
		}
	}
	// 断开WebSocket连接,让ReadLoop退出
	c.conn.Close()
}

// 客户端上报消息封装
type ClientOperation struct {
	Client  *Client
	Message dto.WSClientMessage
}
