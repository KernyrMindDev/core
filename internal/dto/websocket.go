package dto

import "encoding/json"

// 客户端上报消息
type WSClientMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// 服务端回传消息
type WSServerMessage struct {
	Type string `json:"type"`
	Data any    `json:"data,omitempty"`
}

// 消息类型枚举
type WSServerMessageType string

const (
	EventParticipantJoined string = "participant.join"  // 新用户加入
	EventParticipantLeft   string = "participant.leave" // 用户离开
	EventParticipantsSync  string = "sync.participants" // 房间成员同步
)
