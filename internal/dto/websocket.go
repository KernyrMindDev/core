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
	EventObjectCreated string = "object.created" // 创建新节点
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
	EventParticipantJoined string = "participant.join"  // 新用户加入
	EventParticipantLeft   string = "participant.leave" // 用户离开
	EventParticipantsSync  string = "sync.participants" // 房间成员同步
)
