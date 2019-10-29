package mayday

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/docker/mayday/pkg/repository"
	"github.com/docker/mayday/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

type server struct {
	observationRepo repository.Observation
	typeRepo        repository.Type
	config          Config
	logger          logrus.FieldLogger
}

// Server -
type Server interface {
	Start() error
}

// NewServer -
func NewServer(cfg Config, logger logrus.FieldLogger) Server {
	return server{
		observationRepo: repository.NewInMemoryObservation(),
		typeRepo:        repository.NewInMemoryType(),
		config:          cfg,
		logger:          logger,
	}
}

func (m server) Start() error {
	lis, err := net.Listen("tcp", m.config.GRPCPort())
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()

	proto.RegisterMaydayServiceServer(grpcServer, m)

	return grpcServer.Serve(lis)
}

func (m server) CreateObservation(ctx context.Context, req *proto.CreateObservationRequest) (*proto.CreateObservationResponse, error) {
	u := uuid.NewV4()
	id := &proto.UUID{
		Value: fmt.Sprintf("%s", u),
	}
	req.Observation.Id = id
	now := time.Now()
	req.Observation.Created = &timestamp.Timestamp{
		Seconds: now.Unix(),
		Nanos:   int32(now.UnixNano()),
	}
	req.Observation.Updated = &timestamp.Timestamp{
		Seconds: now.Unix(),
		Nanos:   int32(now.UnixNano()),
	}

	foundType := m.typeRepo.FindByID(req.Observation.GetType().GetId().GetValue())
	if foundType == nil {
		return nil, errors.New("unknown type")
	}
	req.Observation.Type = foundType
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
	req.Type.Created = &timestamp.Timestamp{
		Seconds: now.Unix(),
		Nanos:   int32(now.UnixNano()),
	}
	req.Type.Updated = &timestamp.Timestamp{
		Seconds: now.Unix(),
		Nanos:   int32(now.UnixNano()),
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
