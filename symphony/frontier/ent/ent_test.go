// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ent

import (
	"context"
	"testing"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/cloud/testdb"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T) *Client {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	client := NewClient(Driver(sql.OpenDB(name, db)))
	err = client.Schema.Create(context.Background())
	require.NoError(t, err)
	return client
}
