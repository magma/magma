// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql

import (
	"context"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/handler"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

// RequestMiddleware is a graphql request statistics emitter.
func RequestMiddleware() handler.Option {
	f := func(ctx context.Context, next func(context.Context) []byte) []byte {
		stats.Record(ctx, ServerRequestCount.M(1))

		start := time.Now()
		res := next(ctx)
		latency := float64(time.Since(start)) / float64(time.Millisecond)

		tags := []tag.Mutator{
			tag.Upsert(Error, strconv.FormatBool(
				graphql.GetRequestContext(ctx).HasError(graphql.GetResolverContext(ctx)),
			)),
		}

		_ = stats.RecordWithTags(
			ctx,
			tags,
			ServerRequestLatency.M(latency),
			ServerResponseCount.M(1),
		)

		return res
	}
	return handler.RequestMiddleware(f)
}

// ResolverMiddleware is a graphql resolve statistics emitter.
func ResolverMiddleware() handler.Option {
	f := func(ctx context.Context, next graphql.Resolver) (interface{}, error) {
		start := time.Now()
		res, err := next(ctx)
		latency := float64(time.Since(start)) / float64(time.Millisecond)

		rc := graphql.GetResolverContext(ctx)
		deprecated := rc.Field.Definition.Directives.ForName("deprecated") != nil
		tags := []tag.Mutator{
			tag.Upsert(Object, rc.Object),
			tag.Upsert(Field, rc.Field.Name),
			tag.Upsert(Error, strconv.FormatBool(err != nil)),
			tag.Upsert(Deprecated, strconv.FormatBool(deprecated)),
		}

		_ = stats.RecordWithTags(
			ctx,
			tags,
			ServerResolveLatency.M(latency),
			ServerResolveCount.M(1),
		)
		return res, err
	}
	return handler.ResolverMiddleware(f)
}

// DefaultServerOptions are the default graphql server
// instrumentation options provided by this package.
var DefaultServerOptions = []handler.Option{
	RequestMiddleware(),
	ResolverMiddleware(),
}
