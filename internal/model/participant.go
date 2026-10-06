package model

import "uuid"

type MemberRole string

const (
	MemberRoleOwner  MemberRole = "owner"  // 所有者
	MemberRoleEditor MemberRole = "editor" // 可以编辑
	MemberRoleViewer MemberRole = "viewer" // 仅查看
)

// 用户
type Member struct {
	ID       uuid.UUID  `json:"id"`       // 共享时临时UUID
	Nickname string     `json:"nickname"` // 昵称, 用于显示
	Role     MemberRole `json:"role"`     // 用户身份
}

// 用户认证
type MemberCredential struct {
	MemberID  uuid.UUID
	TokenHash [32]byte
}
