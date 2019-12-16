// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphgrpc

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/facebookincubator/symphony/graph/ent/migrate"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestTenantServer_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	ts := NewTenantService(func(context.Context) ExecQueryer { return db })

	tenant, err := ts.Create(context.Background(), &wrappers.StringValue{Value: ""})
	require.Nil(t, tenant)
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?")).
		WithArgs("tenant_foo").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	tenant, err = ts.Create(context.Background(), &wrappers.StringValue{Value: "foo"})
	require.Nil(t, tenant)
	require.IsType(t, codes.AlreadyExists, status.Code(err))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?")).
		WithArgs("tenant_foo").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectExec(regexp.QuoteMeta("CREATE DATABASE `tenant_foo`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	tenant, err = ts.Create(context.Background(), &wrappers.StringValue{Value: "foo"})
	require.NoError(t, err)
	require.NotNil(t, tenant)
	require.Equal(t, "foo", tenant.Id)
	require.Equal(t, "foo", tenant.Name)
}

func TestTenantServer_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	ts := NewTenantService(func(context.Context) ExecQueryer { return db })

	tenant, err := ts.Get(context.Background(), &wrappers.StringValue{Value: ""})
	require.Nil(t, tenant)
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?")).
		WithArgs("tenant_foo").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	tenant, err = ts.Get(context.Background(), &wrappers.StringValue{Value: "foo"})
	require.NoError(t, err)
	require.NotNil(t, tenant)
	require.Equal(t, "foo", tenant.Id)
	require.Equal(t, "foo", tenant.Name)
}

func TestTenantServer_Truncate(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	ts := NewTenantService(func(context.Context) ExecQueryer { return db })

	_, err = ts.Truncate(context.Background(), &wrappers.StringValue{Value: ""})
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?",
	)).
		WithArgs("tenant_foo").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	result := sqlmock.NewResult(0, 0)
	mock.ExpectExec("SET FOREIGN_KEY_CHECKS=0").WillReturnResult(result)
	for _, table := range migrate.Tables {
		query := fmt.Sprintf("DELETE FROM `tenant_foo`.`%s`", table.Name)
		mock.ExpectExec(query).WillReturnResult(result)
	}
	mock.ExpectExec("SET FOREIGN_KEY_CHECKS=1").WillReturnResult(result)
	_, err = ts.Truncate(context.Background(), &wrappers.StringValue{Value: "foo"})
	require.NoError(t, err)
}

func TestTenantServer_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	ts := NewTenantService(func(context.Context) ExecQueryer { return db })

	_, err = ts.Delete(context.Background(), &wrappers.StringValue{Value: ""})
	require.IsType(t, codes.InvalidArgument, status.Code(err))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?")).
		WithArgs("tenant_foo").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	_, err = ts.Delete(context.Background(), &wrappers.StringValue{Value: "foo"})
	require.IsType(t, codes.NotFound, status.Code(err))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?")).
		WithArgs("tenant_foo").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta("DROP DATABASE `tenant_foo`")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	_, err = ts.Delete(context.Background(), &wrappers.StringValue{Value: "foo"})
	require.NoError(t, err)
}

func TestTenantServer_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	ts := NewTenantService(func(context.Context) ExecQueryer { return db })

	mock.ExpectQuery(regexp.QuoteMeta("SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME LIKE ?")).
		WithArgs("tenant_%").
		WillReturnRows(sqlmock.NewRows([]string{"SCHEMA_NAME"}).AddRow("tenant_foo").AddRow("tenant_bar"))
	res, err := ts.List(context.Background(), nil)
	require.NoError(t, err)
	require.Len(t, res.Tenants, 2)
	require.Equal(t, "foo", res.Tenants[0].Id)
	require.Equal(t, "foo", res.Tenants[0].Name)
	require.Equal(t, "bar", res.Tenants[1].Id)
	require.Equal(t, "bar", res.Tenants[1].Name)
}
