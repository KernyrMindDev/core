package objects

import (
	"uuid"

	"github.com/KernyrMindDev/core/internal/dto"
)

// 白板内节点对象的抽象接口
type BoardObject interface {
	ID() uuid.UUID
	ApplyPatch(data []byte) error

	ToDTO() dto.ObjectDTO
}
