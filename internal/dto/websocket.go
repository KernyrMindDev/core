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
