// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package enttest

import (
	"context"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/symphony/frontier/ent"
	"github.com/facebookincubator/symphony/pkg/testdb"
)

func NewClient() (*ent.Client, error) {
	db, name, err := testdb.Open()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	client := ent.NewClient(ent.Driver(sql.OpenDB(name, db)))
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}
	return client, nil
}
