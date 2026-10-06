package objects

// 二维坐标
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// 客户端上传数据
type PointPatch struct {
	X *float64 `json:"x"`
	Y *float64 `json:"y"`
}

// 传入PointPatch指针,完成数据修改
func (o *Point) Patch(patch *PointPatch) {
	if patch.X != nil {
		o.X = *patch.X
	}
	if patch.Y != nil {
		o.Y = *patch.Y
	}
}

// 方形宽高
type Size struct {
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// 客户端上传数据
type SizePatch struct {
	Width  *float64 `json:"width"`
	Height *float64 `json:"height"`
}

// 传入SizePatch指针,完成数据修改
func (o *Size) Patch(patch *SizePatch) {
	if patch.Width != nil {
		o.Width = *patch.Width
	}
	if patch.Height != nil {
		o.Height = *patch.Height
	}
}
