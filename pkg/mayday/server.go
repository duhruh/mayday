package mayday

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/mayday/pkg/repository"
	"github.com/docker/mayday/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

type server struct {
	observationRepo repository.Observation
	typeRepo        repository.Type
}

// NewServer -
func NewServer(g *grpc.Server) proto.MaydayServiceServer {
	m := server{
		observationRepo: repository.NewInMemoryObservation(),
		typeRepo:        repository.NewInMemoryType(),
	}
	proto.RegisterMaydayServiceServer(g, m)
	return m
}

func (m server) CreateObservation(ctx context.Context, req *proto.CreateObservationRequest) (*proto.CreateObservationResponse, error) {
	u := uuid.NewV4()
	id := &proto.UUID{
		Value: fmt.Sprintf("%s", u),
	}
	req.Observation.Id = id
	now := time.Now()
	secs := now.Unix()
	req.Observation.Created = &timestamp.Timestamp{
		Seconds: secs,
	}
	req.Observation.Updated = &timestamp.Timestamp{
		Seconds: secs,
	}
	m.observationRepo.Create(req.Observation)

	return &proto.CreateObservationResponse{
		Observation: req.Observation,
	}, nil
}
func (m server) CreateType(ctx context.Context, req *proto.CreateTypeRequest) (*proto.CreateTypeResponse, error) {
	u := uuid.NewV4()
	id := &proto.UUID{
		Value: fmt.Sprintf("%s", u),
	}
	req.Type.Id = id
	now := time.Now()
	secs := now.Unix()
	req.Type.Created = &timestamp.Timestamp{
		Seconds: secs,
	}
	req.Type.Updated = &timestamp.Timestamp{
		Seconds: secs,
	}
	m.typeRepo.Create(req.Type)

	return &proto.CreateTypeResponse{
		Type: req.Type,
	}, nil
}
func (m server) ListObservations(context.Context, *proto.ListObservationsRequest) (*proto.ListObservationsResponse, error) {
	return &proto.ListObservationsResponse{
		Observations: m.observationRepo.List(),
	}, nil
}
func (m server) ListTypes(context.Context, *proto.ListTypesRequest) (*proto.ListTypesResponse, error) {
	return &proto.ListTypesResponse{
		Types: m.typeRepo.List(),
	}, nil
}
