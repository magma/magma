// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"database/sql"
	"fmt"

	"github.com/facebookincubator/symphony/graph/viewer"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TenantService is a tenant service
type TenantService struct {
	db *sql.DB
}

// NewTenantService creates a tenant service.
func NewTenantService(db *sql.DB) *TenantService {
	return &TenantService{db}
}

func (s *TenantService) Create(ctx context.Context, name *wrappers.StringValue) (*Tenant, error) {
	if name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "missing tenant name")
	}
	exist, err := s.exist(ctx, name.Value)
	if err != nil {
		return nil, status.FromContextError(err).Err()
	}
	if exist {
		return nil, status.Errorf(codes.AlreadyExists, "tenant %q exists", name.Value)
	}
	if _, err := s.db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_bin", viewer.DBName(name.Value))); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	return &Tenant{Id: name.Value, Name: name.Value}, nil
}

func (s *TenantService) List(ctx context.Context, _ *empty.Empty) (*TenantList, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME LIKE ?", viewer.DBName("%"))
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

func (s *TenantService) Get(ctx context.Context, name *wrappers.StringValue) (*Tenant, error) {
	if name.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "missing tenant name")
	}
	exist, err := s.exist(ctx, name.Value)
	if err != nil {
		return nil, status.FromContextError(err).Err()
	}
	if !exist {
		return nil, status.Errorf(codes.NotFound, "missing tenant %s", name.Value)
	}
	return &Tenant{Id: name.Value, Name: name.Value}, nil
}

func (s *TenantService) Delete(ctx context.Context, id *wrappers.StringValue) (*empty.Empty, error) {
	if id.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "missing tenant id")
	}
	exist, err := s.exist(ctx, id.Value)
	if err != nil {
		return nil, status.FromContextError(err).Err()
	}
	if !exist {
		return nil, status.Errorf(codes.NotFound, "missing tenant %s", id.Value)
	}
	if _, err := s.db.ExecContext(ctx, fmt.Sprintf("DROP DATABASE `%s`", viewer.DBName(id.Value))); err != nil {
		return nil, status.FromContextError(err).Err()
	}
	return &empty.Empty{}, nil
}

func (s *TenantService) exist(ctx context.Context, name string) (bool, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", viewer.DBName(name))
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return false, sql.ErrNoRows
	}
	var n int
	if err := rows.Scan(&n); err != nil {
		return false, errors.WithMessage(err, "scanning count")
	}
	return n > 0, nil
}
