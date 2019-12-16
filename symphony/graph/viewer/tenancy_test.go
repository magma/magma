// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viewer

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/graph/ent"

	"github.com/cenkalti/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gocloud.dev/server/health"
)

func TestFixedTenancy(t *testing.T) {
	want := &ent.Client{}
	tenancy := NewFixedTenancy(want)
	assert.Implements(t, (*Tenancy)(nil), tenancy)
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

func createMySQLDatabase(db *sql.DB) (string, func() error, error) {
	name := "testdb_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if _, err := db.Exec("create database " + DBName(name)); err != nil {
		return "", nil, err
	}
	return name, func() error {
		_, err := db.Exec("drop database " + DBName(name))
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
	tenancy, err := NewMySQLTenancy(dsn)
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
	c2, err := tenancy.ClientFor(context.Background(), n1)
	assert.NoError(t, err)
	assert.True(t, c1 == c2)
	c2, err = tenancy.ClientFor(context.Background(), n2)
	assert.NoError(t, err)
	assert.False(t, c1 == c2)
}
