package repository

import (
	"github.com/docker/mayday/proto"
)

type Type interface {
	Create(*proto.Type) *proto.Type
	List() []*proto.Type
}

type inmemoryType struct {
	types []*proto.Type
}

func NewInMemoryType() Type {
	return &inmemoryType{}
}

func (in *inmemoryType) Create(o *proto.Type) *proto.Type {
	in.types = append(in.types, o)
	return o
}
func (in *inmemoryType) List() []*proto.Type {
	return in.types
}
