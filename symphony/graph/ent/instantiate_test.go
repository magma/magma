// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ent

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/pkg/testdb"
	"github.com/stretchr/testify/require"
)

func TestInstantiation(t *testing.T) {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	client := NewClient(Driver(sql.OpenDB(name, db)))

	ctx := context.Background()
	err = client.Schema.Create(ctx)
	require.NoError(t, err)

	typ := client.LocationType.
		Create().
		SetName("planet").
		SetMapZoomLevel(5).
		SetSite(true).
		SaveX(ctx)
	_ = client.Location.
		Create().
		SetName("earth").
		SetType(typ).
		SaveX(ctx)

	data, err := json.Marshal(typ)
	require.NoError(t, err)
	typ = nil
	err = json.Unmarshal(data, &typ)
	require.NoError(t, err)

	count := client.LocationType.
		Instantiate(typ).
		QueryLocations().
		CountX(ctx)
	require.Equal(t, 1, count)
}
