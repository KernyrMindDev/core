package dto

import "uuid"

// 单个成员信息
type ParticipantInfo struct {
	Nickname string    `json:"nickname"`
	UUID     uuid.UUID `json:"uuid"`
}

// 房间成员列表
type ParticipantList struct {
	Participants []ParticipantInfo `json:"participants"`
}
