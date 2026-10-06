package dto

import "encoding/json"

// 客户端上报消息
type WSClientMessage struct {
	Type WSClientMessageType `json:"type"`
	Data json.RawMessage     `json:"data,omitempty"`
}

// 客户端消息类型枚举
type WSClientMessageType string

const (
	EventCursorMoved   string = "cursor.move"    // 鼠标移动
	EventObjectCreate  string = "object.create"  // 创建新节点
	EventObjectUpdated string = "object.updated" // 更新节点属性
	EventObjectDeleted string = "object.deleted" // 删除节点
)

// 服务端回传消息
type WSServerMessage struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

// 服务器消息类型枚举
type WSServerMessageType string

const (
	EventMemberJoined string = "member.join"  // 新用户加入
	EventMemberLeft   string = "member.leave" // 用户离开
	EventMembersSync  string = "sync.members" // 房间成员同步
)
