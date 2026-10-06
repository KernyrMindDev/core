package objects

import (
	"uuid"

	"github.com/KernyrMindDev/core/internal/apperr"
)

// 组件构造器
type ObjectFactory struct {
	creators map[string]ObjectCreator
}

func NewObjectFactory() *ObjectFactory {
	factory := &ObjectFactory{
		creators: make(map[string]ObjectCreator),
	}

	// 注册各个组件
	factory.Register("rect", NewRectObject)
	// factory.Register("text", NewTextObject)

	return factory
}

// 一个节点的工厂函数
type ObjectCreator func(
	id uuid.UUID,
	data []byte,
) (BoardObject, error)

// 根据传入的Type, ID及data实例化一个BoardObject
func (o *ObjectFactory) Create(typ string, id uuid.UUID, data []byte) (BoardObject, error) {
	// 尝试获取工厂函数
	creator, ok := o.creators[typ]
	if !ok {
		// 没获取到
		return nil, apperr.UnregisteredCreator{
			Type: typ,
		}
	}
	// 调用工厂函数
	return creator(id, data)
}

func (o *ObjectFactory) Register(typ string, creator ObjectCreator) {
	o.creators[typ] = creator
}
