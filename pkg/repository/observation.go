package repository

import (
	"github.com/docker/mayday/proto"
)

type Observation interface {
	Create(*proto.Observation) *proto.Observation
	List() []*proto.Observation
}

type inmemoryObservation struct {
	observations []*proto.Observation
}

func NewInMemoryObservation() Observation {
	return &inmemoryObservation{}
}

func (in *inmemoryObservation) Create(o *proto.Observation) *proto.Observation {
	in.observations = append(in.observations, o)
	return o
}
func (in *inmemoryObservation) List() []*proto.Observation {
	return in.observations
}
