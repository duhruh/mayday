package mayday

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"

	"github.com/docker/mayday/pkg/db"
	"github.com/docker/mayday/pkg/repository"
	"github.com/docker/mayday/proto"
	"github.com/golang/protobuf/ptypes"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
)

type server struct {
	observationRepo repository.Observation
	typeRepo        repository.Type
	config          Config
	logger          logrus.FieldLogger
	databaseCon     db.Connection
}

// Server -
type Server interface {
	Start(context.Context) error
}

// NewServer -
func NewServer(cfg Config, logger logrus.FieldLogger, databaseCon db.Connection) Server {
	return server{
		observationRepo: repository.NewObservation(databaseCon),
		typeRepo:        repository.NewType(databaseCon),
		config:          cfg,
		logger:          logger,
		databaseCon:     databaseCon,
	}
}

func (m server) Start(ctx context.Context) error {
	m.logger.Info("initializing database")
	err := repository.Init(ctx, m.observationRepo, m.typeRepo)
	if err != nil {
		return err
	}

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
	now := ptypes.TimestampNow()
	req.Observation.Created = now
	req.Observation.Updated = now

	foundType, err := m.typeRepo.FindByID(ctx, req.Observation.GetType().GetId().GetValue())
	if err != nil {
		return nil, err
	}
	if foundType == nil {
		return nil, errors.New("unknown type")
	}
	req.Observation.Type = foundType
	err = m.observationRepo.Create(ctx, req.Observation)
	if err != nil {
		return nil, err
	}
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

	now := ptypes.TimestampNow()
	req.Type.Created = now
	req.Type.Updated = now
	err := m.typeRepo.Create(ctx, req.Type)
	if err != nil {
		return nil, err
	}

	return &proto.CreateTypeResponse{
		Type: req.Type,
	}, nil
}
func (m server) ListObservations(ctx context.Context, l *proto.ListObservationsRequest) (*proto.ListObservationsResponse, error) {
	observations, err := m.observationRepo.List(ctx, int(l.GetLimit()), int(l.GetPage()))
	if err != nil {
		return nil, err
	}
	return &proto.ListObservationsResponse{
		Observations: observations,
	}, nil
}
func (m server) ListTypes(ctx context.Context, l *proto.ListTypesRequest) (*proto.ListTypesResponse, error) {
	types, err := m.typeRepo.List(ctx, int(l.GetLimit()), int(l.GetPage()))
	if err != nil {
		return nil, err
	}
	return &proto.ListTypesResponse{
		Types: types,
	}, nil
}
