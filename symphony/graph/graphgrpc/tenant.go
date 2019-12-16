// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/facebookincubator/symphony/graph/viewer"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	// TenantService is a tenant service.
	TenantService struct{ DB Provider }

	// ExecQueryer wraps QueryContext and ExecContext methods.
	ExecQueryer interface {
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	}

	// Provider provides a db from context.
	Provider func(context.Context) ExecQueryer
)

// NewTenantService create a new tenant service.
func NewTenantService(provider Provider) TenantService {
	return TenantService{provider}
}

// Create a tenant by name.
func (s TenantService) Create(ctx context.Context, name *wrappers.StringValue) (*Tenant, error) {
	if name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "missing tenant name")
	}
	switch exist, err := s.exist(ctx, name.Value); {
	case err != nil:
		return nil, status.FromContextError(err).Err()
	case exist:
		return nil, status.Errorf(codes.AlreadyExists, "tenant %q exists", name.Value)
	}
	if _, err := s.DB(ctx).ExecContext(ctx, fmt.Sprintf("CREATE DATABASE `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_bin", viewer.DBName(name.Value))); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	return &Tenant{Id: name.Value, Name: name.Value}, nil
}

// List all tenants.
func (s TenantService) List(ctx context.Context, _ *empty.Empty) (*TenantList, error) {
	rows, err := s.DB(ctx).QueryContext(ctx,
		"SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME LIKE ?", viewer.DBName("%"),
	)
	if err != nil {
		return nil, status.FromContextError(err).Err()
	}
	defer rows.Close()
	var tenants []*Tenant
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, status.FromContextError(err).Err()
		}
		name = viewer.FromDBName(name)
		tenants = append(tenants, &Tenant{Id: name, Name: name})
	}
	return &TenantList{Tenants: tenants}, nil
}

// Get tenant by name.
func (s TenantService) Get(ctx context.Context, name *wrappers.StringValue) (*Tenant, error) {
	if name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "missing tenant name")
	}
	switch exist, err := s.exist(ctx, name.Value); {
	case err != nil:
		return nil, status.FromContextError(err).Err()
	case !exist:
		return nil, status.Errorf(codes.NotFound, "missing tenant %s", name.Value)
	default:
		return &Tenant{Id: name.Value, Name: name.Value}, nil
	}
}

// Truncate tenant data by name.
func (s TenantService) Truncate(ctx context.Context, name *wrappers.StringValue) (_ *empty.Empty, err error) {
	if name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "missing tenant name")
	}
	switch exist, err := s.exist(ctx, name.Value); {
	case err != nil:
		return nil, status.FromContextError(err).Err()
	case !exist:
		return nil, status.Errorf(codes.NotFound, "missing tenant %s", name.Value)
	}
	db, dbname := s.DB(ctx), viewer.DBName(name.Value)
	if _, err := db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=0"); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	for _, table := range migrate.Tables {
		query := fmt.Sprintf("DELETE FROM `%s`.`%s`", dbname, table.Name) // nolint:gosec
		if _, err = db.ExecContext(ctx, query); err != nil {
			err = status.FromContextError(err).Err()
			break
		}
	}
	if _, err := db.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS=1"); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	return &empty.Empty{}, err
}

// Delete tenant by name.
func (s TenantService) Delete(ctx context.Context, name *wrappers.StringValue) (*empty.Empty, error) {
	if name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "missing tenant name")
	}
	switch exist, err := s.exist(ctx, name.Value); {
	case err != nil:
		return nil, status.FromContextError(err).Err()
	case !exist:
		return nil, status.Errorf(codes.NotFound, "missing tenant %s", name.Value)
	}
	if _, err := s.DB(ctx).ExecContext(ctx,
		fmt.Sprintf("DROP DATABASE `%s`", viewer.DBName(name.Value)),
	); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	return &empty.Empty{}, nil
}

func (s TenantService) exist(ctx context.Context, name string) (bool, error) {
	rows, err := s.DB(ctx).QueryContext(ctx,
		"SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", viewer.DBName(name),
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return false, sql.ErrNoRows
	}
	var n int
	if err := rows.Scan(&n); err != nil {
		return false, fmt.Errorf("scanning count: %w", err)
	}
	return n > 0, nil
}
