package model

import "uuid"

type ParticipantRole string

const (
	ParticipantRoleOwner  ParticipantRole = "owner"  // 所有者
	ParticipantRoleEditor ParticipantRole = "editor" // 可以编辑
	ParticipantRoleViewer ParticipantRole = "viewer" // 仅查看
)

// 用户
type Participant struct {
	ID       uuid.UUID       `json:"id"`       // 共享时临时UUID
	Nickname string          `json:"nickname"` // 昵称, 用于显示
	Role     ParticipantRole `json:"role"`     // 用户身份
}

// 用户认证
type ParticipantCredential struct {
	ParticipantID uuid.UUID
	TokenHash     [32]byte
}
