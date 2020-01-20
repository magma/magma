// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graphactions

import (
	"context"
	"testing"

	"github.com/facebookincubator/ent/dialect/sql"
	"github.com/facebookincubator/ent/dialect/sql/schema"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/pkg/actions/core"
	"github.com/facebookincubator/symphony/pkg/testdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newClient(t *testing.T) *ent.Client {
	db, name, err := testdb.Open()
	require.NoError(t, err)
	db.SetMaxOpenConns(1)
	drv := sql.OpenDB(name, db)
	client := ent.NewClient(ent.Driver(drv))
	require.NoError(t, client.Schema.Create(context.Background(), schema.WithGlobalUniqueID(true)))
	return client
}

func TestQueryRules(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	dataLoader := EntDataLoader{client}

	_, err := client.
		ActionsRule.Create().
		SetName("testInput").
		SetTriggerID("trigger1").
		SetRuleActions([]*core.ActionsRuleAction{}).
		SetRuleFilters([]*core.ActionsRuleFilter{}).
		Save(ctx)
	assert.NoError(t, err)

	rules, err := dataLoader.QueryRules(ctx, "trigger1")
	assert.NoError(t, err)
	assert.Equal(t, rules[0].TriggerID, core.TriggerID("trigger1"))

	rules, err = dataLoader.QueryRules(ctx, "trigger2")
	assert.NoError(t, err)
	assert.Equal(t, len(rules), 0)
}
