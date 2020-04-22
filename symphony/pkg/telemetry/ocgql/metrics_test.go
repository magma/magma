// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql_test

import (
	"math"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler/testserver"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/facebookincubator/symphony/pkg/telemetry/ocgql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats/view"
)

func TestMetrics(t *testing.T) {
	err := view.Register(ocgql.DefaultServerViews...)
	require.NoError(t, err)
	defer view.Unregister(ocgql.DefaultServerViews...)

	h := testserver.New()
	h.AddTransport(transport.POST{})
	h.Use(ocgql.Metrics{})

	c := client.New(h)
	err = c.Post(`query { name }`, &struct{ Name string }{})
	require.NoError(t, err)

	for _, v := range ocgql.DefaultServerViews {
		v := view.Find(v.Name)
		require.NotNil(t, v)

		rows, err := view.RetrieveData(v.Name)
		require.NoError(t, err)
		require.NotEmpty(t, rows)

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
			require.Failf(t, "unknown data type", "value=%v", data)
		}

		assert.NotZero(t, count)
		assert.True(t, math.IsNaN(sum) || sum > 0)
	}
}
