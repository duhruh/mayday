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
	createTypesTable = `CREATE TABLE IF NOT EXISTS types(
    id UUID,
    name VARCHAR,
    schema jsonb,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    PRIMARY KEY (id)
);`
	createTypeQuery = `INSERT INTO types (id,name,schema,created_at,updated_at) VALUES($1,$2,$3,$4,$5)`
	selectTypeQuery = `SELECT * FROM types WHERE id=$1`
	listTypes       = `SELECT * FROM types LIMIT $1 OFFSET $2`
)

// Type -
type Type interface {
	Create(context.Context, *proto.Type) error
	List(context.Context, int, int) ([]*proto.Type, error)
	FindByID(context.Context, string) (*proto.Type, error)
	Init(context.Context) error
}

type typeRepo struct {
	connection db.Connection
}

// NewType -
func NewType(connection db.Connection) Type {
	return typeRepo{
		connection: connection,
	}
}

func (t typeRepo) Init(ctx context.Context) error {
	_, err := t.connection.ExecContext(ctx, createTypesTable)

	if err != nil {
		return err
	}

	return nil
}
func (t typeRepo) Create(ctx context.Context, typ *proto.Type) error {
	stmt, err := t.connection.PrepareContext(ctx, createTypeQuery)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(typ.GetSchema())
	if err != nil {
		return err
	}

	created, _ := ptypes.Timestamp(typ.GetCreated())
	updated, _ := ptypes.Timestamp(typ.GetUpdated())
	_, err = stmt.Exec(
		typ.GetId().GetValue(),
		typ.GetName(),
		bytes,
		created,
		updated,
	)
	if err != nil {
		return err
	}

	return nil
}

func (t typeRepo) List(ctx context.Context, limit, page int) ([]*proto.Type, error) {
	page = page * limit
	println(page, limit)
	res, err := t.connection.QueryContext(ctx, listTypes, limit, page)
	if err != nil {
		return nil, err
	}

	return t.fromSQLRows(res)
}

func (t typeRepo) FindByID(ctx context.Context, id string) (*proto.Type, error) {
	res := t.connection.QueryRowContext(ctx, selectTypeQuery, id)
	return t.fromSQLRow(res)
}

func (t typeRepo) fromSQLRow(row db.Scannable) (*proto.Type, error) {
	var (
		obs         proto.Type
		id          string
		name        string
		schemaBytes []byte
		createdAt   time.Time
		updatedAt   time.Time
	)

	err := row.Scan(&id, &name, &schemaBytes, &createdAt, &updatedAt)
	if err == sql.ErrNoRows {
		return &obs, nil
	}
	if err != nil {
		return &obs, err
	}

	schema := make(map[string]string)
	err = json.Unmarshal(schemaBytes, &schema)
	if err != nil {
		return &obs, err
	}
	created, _ := ptypes.TimestampProto(createdAt)
	updated, _ := ptypes.TimestampProto(updatedAt)
	return &proto.Type{
		Id:      &proto.UUID{Value: id},
		Name:    name,
		Schema:  schema,
		Created: created,
		Updated: updated,
	}, err
}

func (t typeRepo) fromSQLRows(rows db.ScannableList) ([]*proto.Type, error) {
	var types []*proto.Type
	defer rows.Close()
	for rows.Next() {
		obs, err := t.fromSQLRow(rows)
		if err == nil {
			types = append(types, obs)
		}
	}

	return types, rows.Err()
}
