// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gocloud.dev/server/health"
)

func TestFixedTenancy(t *testing.T) {
	want := &ent.Client{}
	tenancy := viewer.NewFixedTenancy(want)
	assert.Implements(t, (*viewer.Tenancy)(nil), tenancy)
	t.Run("ClientFor", func(t *testing.T) {
		got, err := tenancy.ClientFor(context.Background(), "")
		assert.NoError(t, err)
		assert.True(t, want == got)
	})
	t.Run("Client", func(t *testing.T) {
		got := tenancy.Client()
		assert.True(t, want == got)
	})
}

type testTenancy struct {
	mock.Mock
}

func (t *testTenancy) ClientFor(ctx context.Context, name string) (*ent.Client, error) {
	args := t.Called(ctx, name)
	client, _ := args.Get(0).(*ent.Client)
	return client, args.Error(1)
}

func TestCacheTenancy(t *testing.T) {
	var m testTenancy
	m.On("ClientFor", mock.Anything, "bar").
		Return(&ent.Client{}, nil).
		Once()
	m.On("ClientFor", mock.Anything, "baz").
		Return(nil, errors.New("try again")).
		Once()
	m.On("ClientFor", mock.Anything, "baz").
		Return(&ent.Client{}, nil).
		Once()
	defer m.AssertExpectations(t)

	var count int
	tenancy := viewer.NewCacheTenancy(&m, func(*ent.Client) { count++ })
	assert.Implements(t, (*health.Checker)(nil), tenancy)

	client, err := tenancy.ClientFor(context.Background(), "bar")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	cached, err := tenancy.ClientFor(context.Background(), "bar")
	assert.NoError(t, err)
	assert.True(t, client == cached)
	client, err = tenancy.ClientFor(context.Background(), "baz")
	assert.Error(t, err)
	assert.Nil(t, client)
	client, err = tenancy.ClientFor(context.Background(), "baz")
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Equal(t, 2, count)
}

func createMySQLDatabase(db *sql.DB) (string, func() error, error) {
	name := "testdb_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if _, err := db.Exec("create database " + viewer.DBName(name)); err != nil {
		return "", nil, err
	}
	return name, func() error {
		_, err := db.Exec("drop database " + viewer.DBName(name))
		return err
	}, nil
}

func TestMySQLTenancy(t *testing.T) {
	dsn, ok := os.LookupEnv("MYSQL_DSN")
	if !ok {
		t.Skip("MYSQL_DSN not provided")
	}

	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err)
	tenancy, err := viewer.NewMySQLTenancy(dsn)
	require.NoError(t, err)

	assert.Implements(t, (*health.Checker)(nil), tenancy)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	b := backoff.WithContext(backoff.NewConstantBackOff(10*time.Millisecond), ctx)
	err = backoff.Retry(tenancy.CheckHealth, b)
	assert.NoError(t, err)

	n1, cleaner, err := createMySQLDatabase(db)
	require.NoError(t, err)
	defer func(cleaner func() error) {
		assert.NoError(t, cleaner())
	}(cleaner)
	n2, cleaner, err := createMySQLDatabase(db)
	require.NoError(t, err)
	defer func(cleaner func() error) {
		assert.NoError(t, cleaner())
	}(cleaner)

	c1, err := tenancy.ClientFor(context.Background(), n1)
	assert.NotNil(t, c1)
	assert.NoError(t, err)
	c2, err := tenancy.ClientFor(context.Background(), n2)
	assert.NoError(t, err)
	assert.True(t, c1 != c2)
}
