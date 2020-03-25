// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql

import (
	"math"
	"testing"

	"github.com/99designs/gqlgen/client"

	"github.com/99designs/gqlgen/handler"
	"github.com/facebookincubator/symphony/pkg/oc/ocgql/internal/todo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats/view"
)

func TestStatsCollection(t *testing.T) {
	err := view.Register(DefaultServerViews...)
	require.NoError(t, err)

	c := client.New(handler.GraphQL(
		todo.NewExecutableSchema(todo.Config{
			Resolvers: &todo.Resolver{},
		}),
		RequestMiddleware(),
		ResolverMiddleware(),
	))

	var rsp struct {
		Todos []struct {
			ID   string
			Text string
		}
	}
	c.MustPost(`query { todos { id text } }`, &rsp)

	for _, v := range DefaultServerViews {
		v := view.Find(v.Name)
		assert.NotNil(t, v)

		rows, err := view.RetrieveData(v.Name)
		require.NoError(t, err)
		require.NotZero(t, len(rows))

		var (
			count int
			sum   = math.NaN()
		)
		switch data := rows[0].Data.(type) {
		case *view.CountData:
			count = int(data.Value)
		case *view.DistributionData:
			count = int(data.Count)
			sum = data.Sum()
		default:
			assert.Failf(t, "unknown data type", "value=%v", data)
		}

		assert.NotZero(t, count)
		assert.True(t, math.IsNaN(sum) || sum > 0)
	}
}
