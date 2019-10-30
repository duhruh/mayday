package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/docker/mayday/pkg/db"
	"github.com/docker/mayday/proto"
	"github.com/golang/protobuf/ptypes"
)

const (
	createObservationsTable = `CREATE TABLE IF NOT EXISTS observations(
	id UUID,
	name VARCHAR,
	types_id UUID,
	payload jsonb,
	created_at TIMESTAMPTZ,
	updated_at TIMESTAMPTZ,
	PRIMARY KEY (id),
	FOREIGN KEY (types_id) REFERENCES types (id) ON DELETE CASCADE
);`
	createObservationQuery = `INSERT INTO observations (id,name,types_id,payload,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6)`
	selectObservationQuery = `SELECT * FROM observations WHERE id=$1`
	listObservations       = `SELECT * FROM observations LIMIT $1 OFFSET $2`
)

// Observation -
type Observation interface {
	Create(context.Context, *proto.Observation) error
	List(context.Context, int, int) ([]*proto.Observation, error)
	Init(context.Context) error
}

type observation struct {
	connection db.Connection
}

// NewObservation - s
func NewObservation(con db.Connection) Observation {
	return observation{
		connection: con,
	}
}

func (o observation) Init(ctx context.Context) error {
	_, err := o.connection.ExecContext(ctx, createObservationsTable)

	if err != nil {
		return err
	}

	return nil
}

func (o observation) Create(ctx context.Context, obs *proto.Observation) error {
	stmt, err := o.connection.PrepareContext(ctx, createObservationQuery)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(obs.GetPayload())
	if err != nil {
		return err
	}
	created, _ := ptypes.Timestamp(obs.GetCreated())
	updated, _ := ptypes.Timestamp(obs.GetUpdated())
	_, err = stmt.Exec(
		obs.GetId().GetValue(),
		obs.GetName(),
		obs.GetType().GetId().GetValue(),
		bytes,
		created,
		updated,
	)
	if err != nil {
		return err
	}
	println("created observation")

	return nil
}

func (o observation) List(ctx context.Context, limit int, page int) ([]*proto.Observation, error) {
	page = page * limit
	res, err := o.connection.QueryContext(ctx, listObservations, limit, page)
	if err != nil {
		return nil, err
	}

	return o.fromSQLRows(res)
}

func (o observation) fromSQLRow(row db.Scannable) (*proto.Observation, error) {
	var (
		obs          proto.Observation
		id           string
		name         string
		typeID       string
		payloadBytes []byte
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := row.Scan(&id, &name, &typeID, &payloadBytes, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return &obs, nil
	}
	if err != nil {
		return &obs, err
	}

	payload := make(map[string]string)
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return &obs, err
	}

	created, _ := ptypes.TimestampProto(createdAt)
	updated, _ := ptypes.TimestampProto(updatedAt)
	return &proto.Observation{
		Id:   &proto.UUID{Value: id},
		Name: name,
		Type: &proto.Type{
			Id: &proto.UUID{
				Value: typeID,
			},
		},
		Payload: payload,
		Created: created,
		Updated: updated,
	}, err
}

func (o observation) fromSQLRows(rows db.ScannableList) ([]*proto.Observation, error) {
	var observations []*proto.Observation
	defer rows.Close()
	for rows.Next() {

		obs, err := o.fromSQLRow(rows)
		if err == nil {
			observations = append(observations, obs)
		}
	}
	err := rows.Err()
	return observations, err
}
