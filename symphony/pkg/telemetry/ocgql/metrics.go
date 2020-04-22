// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ocgql

import (
	"context"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

// Metrics for opencensus.
type Metrics struct{}

var _ interface {
	graphql.HandlerExtension
	graphql.ResponseInterceptor
	graphql.FieldInterceptor
} = Metrics{}

// ExtensionName returns the metrics extension name.
func (Metrics) ExtensionName() string {
	return "OpenCensusMetrics"
}

// Validate the executable graphql schema.
func (Metrics) Validate(graphql.ExecutableSchema) error {
	return nil
}

// InterceptResponse measures graphql response execution.
func (Metrics) InterceptResponse(ctx context.Context, next graphql.ResponseHandler) *graphql.Response {
	stats.Record(ctx, ServerRequestCount.M(1))

	start := graphql.GetOperationContext(ctx).Stats.OperationStart
	rsp := next(ctx)
	end := graphql.Now()

	latency := float64(end.Sub(start)) / float64(time.Millisecond)
	tags := []tag.Mutator{
		tag.Upsert(Error, strconv.FormatBool(
			graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)),
		)),
	}
	_ = stats.RecordWithTags(
		ctx,
		tags,
		ServerRequestLatency.M(latency),
		ServerResponseCount.M(1),
	)
	return rsp
}

// InterceptField measures graphql field execution.
func (Metrics) InterceptField(ctx context.Context, next graphql.Resolver) (interface{}, error) {
	fc := graphql.GetFieldContext(ctx)
	deprecated := fc.Field.Definition.Directives.ForName("deprecated")
	ctx, _ = tag.New(ctx,
		tag.Upsert(Object, fc.Object),
		tag.Upsert(Field, fc.Field.Name),
		tag.Upsert(Deprecated, strconv.FormatBool(deprecated != nil)),
	)

	start := graphql.Now()
	res, err := next(ctx)
	end := graphql.Now()

	tags := []tag.Mutator{tag.Upsert(Error, strconv.FormatBool(err != nil))}
	latency := float64(end.Sub(start)) / float64(time.Millisecond)
	_ = stats.RecordWithTags(
		ctx,
		tags,
		ServerResolveLatency.M(latency),
		ServerResolveCount.M(1),
	)
	return res, err
}
