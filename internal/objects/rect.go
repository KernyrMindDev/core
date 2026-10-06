package objects

import (
	"encoding/json"

	"github.com/KernyrMindDev/core/internal/dto"
)

type RectObject struct {
	id string

	Point
	Size
}

// 客户端上报数据
type RectPatch struct {
	Size  *SizePatch  `json:"size"`
	Point *PointPatch `json:"point"`
}

// 服务端返回数据
type RectDTO struct {
	Point Point `json:"point"`
	Size  Size  `json:"size"`
}

func (o *RectObject) ID() string {
	return o.id
}

func (o *RectObject) ApplyPatch(data []byte) error {
	var patch RectPatch
	if err := json.Unmarshal(data, &patch); err != nil {
		// 先简单透传
		return err
	}
	if patch.Point != nil {
		o.Point.Patch(patch.Point)
	}
	if patch.Size != nil {
		o.Size.Patch(patch.Size)
	}
	return nil
}

func (o *RectObject) ToDTO() dto.ObjectDTO {
	return dto.ObjectDTO{
		ID:   o.id,
		Type: "rect",
		Data: RectDTO{
			Point: o.Point,
			Size:  o.Size,
		},
	}
}

// 接口实现检查
var _ BoardObject = &RectObject{}

// 工厂函数
func NewRectObject(id string, data []byte) (BoardObject, error) {
	return nil, nil
}
