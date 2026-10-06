package dto

import "uuid"

// 单个成员信息
type MemberInfo struct {
	Nickname string    `json:"nickname"`
	UUID     uuid.UUID `json:"uuid"`
}

// 房间成员列表
type MemberList struct {
	Members []MemberInfo `json:"members"`
}
