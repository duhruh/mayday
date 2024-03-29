package mayday

import (
	"context"
	"encoding/json"

	"github.com/docker/mayday/proto"
	"google.golang.org/grpc"
)

// Client -
type Client interface {
	CreateObservation(context.Context, []byte) (*proto.CreateObservationResponse, error)
	CreateType(context.Context, []byte) (*proto.CreateTypeResponse, error)
	ListObservations(context.Context) (*proto.ListObservationsResponse, error)
	ListTypes(context.Context) (*proto.ListTypesResponse, error)
	DeleteObservation(context.Context, []byte) (*proto.DeleteObservationResponse, error)
	DeleteType(context.Context, []byte) (*proto.DeleteTypeResponse, error)
}

type client struct {
	grpcClient proto.MaydayServiceClient
}

// NewClient -
func NewClient(c *grpc.ClientConn) Client {
	cc := client{
		grpcClient: proto.NewMaydayServiceClient(c),
	}
	return cc
}

func (c client) CreateObservation(ctx context.Context, j []byte) (*proto.CreateObservationResponse, error) {
	observation := &proto.Observation{}

	err := json.Unmarshal(j, observation)
	if err != nil {
		return nil, err
	}

	return c.grpcClient.CreateObservation(ctx, &proto.CreateObservationRequest{
		Observation: observation,
	})
}
func (c client) CreateType(ctx context.Context, j []byte) (*proto.CreateTypeResponse, error) {
	protoTypes := &proto.Type{}

	err := json.Unmarshal(j, protoTypes)
	if err != nil {
		return nil, err
	}

	return c.grpcClient.CreateType(ctx, &proto.CreateTypeRequest{
		Type: protoTypes,
	})
}

func (c client) DeleteObservation(ctx context.Context, j []byte) (*proto.DeleteObservationResponse, error) {
	observation := &proto.Observation{}

	err := json.Unmarshal(j, observation)
	if err != nil {
		return nil, err
	}

	return c.grpcClient.DeleteObservation(ctx, &proto.DeleteObservationRequest{
		Observation: observation,
	})
}
func (c client) DeleteType(ctx context.Context, j []byte) (*proto.DeleteTypeResponse, error) {
	protoTypes := &proto.Type{}

	err := json.Unmarshal(j, protoTypes)
	if err != nil {
		return nil, err
	}

	return c.grpcClient.DeleteType(ctx, &proto.DeleteTypeRequest{
		Type: protoTypes,
	})
}
func (c client) ListObservations(ctx context.Context) (*proto.ListObservationsResponse, error) {
	return c.grpcClient.ListObservations(ctx, &proto.ListObservationsRequest{
		Page:  0,
		Limit: 100,
	})
}
func (c client) ListTypes(ctx context.Context) (*proto.ListTypesResponse, error) {
	return c.grpcClient.ListTypes(ctx, &proto.ListTypesRequest{
		Page:  0,
		Limit: 100,
	})
}
