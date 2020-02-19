// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql

import (
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// The following server GraphQL measures are supported for use in custom views.
var (
	ServerRequestCount = stats.Int64(
		"graphql/server/request_count",
		"Number of GraphQL requests started",
		stats.UnitDimensionless,
	)
	ServerResponseCount = stats.Int64(
		"graphql/server/response_count",
		"Number of GraphQL requests ended",
		stats.UnitDimensionless,
	)
	ServerRequestLatency = stats.Float64(
		"graphql/server/request_latency",
		"Latency of GraphQL requests",
		stats.UnitMilliseconds,
	)
	ServerResolveCount = stats.Int64(
		"graphql/server/resolve_count",
		"Number of GraphQL resolves",
		stats.UnitDimensionless,
	)
	ServerResolveLatency = stats.Float64(
		"graphql/server/resolve_latency",
		"Latency of GraphQL resolves",
		stats.UnitMilliseconds,
	)
)

// The following tags are applied to stats recorded by this package.
var (
	// Object is the GraphQL object being resolved.
	Object = tag.MustNewKey("graphql.object")

	// Field is the GraphQL object field being resolved.
	Field = tag.MustNewKey("graphql.field")

	// Error is the GraphQL exit error.
	Error = tag.MustNewKey("graphql.error")
)

// Default distributions used by views in this package.
var (
	DefaultLatencyDistribution = view.Distribution(0, 0.01, 0.05, 0.1, 0.3, 0.6, 0.8, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000)
)

// Package ocgql provides some convenience views for server measures.
// You still need to register these views for data to actually be collected.
var (
	ServerRequestCountView = &view.View{
		Name:        "graphql/server/request_count",
		Description: "Count of GraphQL requests started",
		TagKeys:     []tag.Key{Error},
		Measure:     ServerRequestCount,
		Aggregation: view.Count(),
	}

	ServerResponseCountView = &view.View{
		Name:        "graphql/server/response_count",
		Description: "Count of GraphQL requests ended",
		TagKeys:     []tag.Key{Error},
		Measure:     ServerResponseCount,
		Aggregation: view.Count(),
	}

	ServerRequestLatencyView = &view.View{
		Name:        "graphql/server/request_latency",
		Description: "Latency distribution of GraphQL requests",
		TagKeys:     []tag.Key{Error},
		Measure:     ServerRequestLatency,
		Aggregation: DefaultLatencyDistribution,
	}

	ServerResponseCountByError = &view.View{
		Name:        "graphql/server/response_count_by_error",
		Description: "Count of GraphQL responses by error",
		TagKeys:     []tag.Key{Error},
		Measure:     ServerResponseCount,
		Aggregation: view.Count(),
	}

	ServerResolveCountView = &view.View{
		Name:        "graphql/server/resolve_count",
		Description: "Count of GraphQL resolves",
		TagKeys:     []tag.Key{Error},
		Measure:     ServerResolveCount,
		Aggregation: view.Count(),
	}

	ServerResolveLatencyView = &view.View{
		Name:        "graphql/server/resolve_latency",
		Description: "Latency distribution of GraphQL resolves",
		TagKeys:     []tag.Key{Error},
		Measure:     ServerResolveLatency,
		Aggregation: DefaultLatencyDistribution,
	}

	ServerResolveCountByError = &view.View{
		Name:        "graphql/server/resolve_count_by_error",
		Description: "Count of GraphQL resolves by error",
		TagKeys:     []tag.Key{Error},
		Measure:     ServerResolveCount,
		Aggregation: view.Count(),
	}

	ServerResolveCountByObjectField = &view.View{
		Name:        "graphql/server/resolve_count_by_object_field",
		Description: "Count of GraphQL resolves by object and field",
		TagKeys:     []tag.Key{Object, Field},
		Measure:     ServerResolveCount,
		Aggregation: view.Count(),
	}
)

// DefaultServerViews are the default server views provided by this package.
var DefaultServerViews = []*view.View{
	ServerRequestCountView,
	ServerResponseCountView,
	ServerRequestLatencyView,
	ServerResponseCountByError,
	ServerResolveCountView,
	ServerResolveLatencyView,
	ServerResolveCountByError,
	ServerResolveCountByObjectField,
}
