/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package counters

import (
	"context"
	"fmt"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// Operation counterset to track an operation
type Operation struct {
	name           string
	startTime      int64
	ctx            context.Context
	start          *stats.Int64Measure
	success        *stats.Int64Measure
	successLatency *stats.Int64Measure
	failed         *stats.Int64Measure
	failedLatency  *stats.Int64Measure
}

// NewOperation creates a new operation counter set
func NewOperation(name string, tagKeys ...tag.Key) Operation {
	operation := Operation{
		name:      name,
		startTime: 0,
		ctx:       context.Background(),
		start: stats.Int64(
			fmt.Sprintf("%s/start", name),
			fmt.Sprintf("Operation '%s' started", name),
			stats.UnitDimensionless,
		),
		success: stats.Int64(
			fmt.Sprintf("%s/success", name),
			fmt.Sprintf("Operation '%s' succeeded", name),
			stats.UnitDimensionless,
		),
		successLatency: stats.Int64(
			fmt.Sprintf("%s/success_latency", name),
			fmt.Sprintf("Operation '%s' success latency", name),
			stats.UnitMilliseconds,
		),
		failed: stats.Int64(
			fmt.Sprintf("%s/failed", name),
			fmt.Sprintf("Operation '%s' failed", name),
			stats.UnitDimensionless,
		),
		failedLatency: stats.Int64(
			fmt.Sprintf("%s/failed_latency", name),
			fmt.Sprintf("Operation '%s' failure latency", name),
			stats.UnitMilliseconds,
		),
	}
	views := []*view.View{
		{
			Name:        fmt.Sprintf("%s/start/count", name),
			Measure:     operation.start,
			Description: fmt.Sprintf("The number of time '%s' was started", name),
			Aggregation: view.Count(),
			TagKeys:     tagKeys,
		},
		{
			Name:        fmt.Sprintf("%s/failure/count", name),
			Measure:     operation.failed,
			Description: fmt.Sprintf("The number of time '%s' has failed", name),
			Aggregation: view.Count(),
			TagKeys:     tagKeys,
		},
		{
			Name:        fmt.Sprintf("%s/failure/latency", name),
			Measure:     operation.failedLatency,
			Description: fmt.Sprintf("The latency of failed '%s' operations", name),
			Aggregation: view.LastValue(),
			TagKeys:     tagKeys,
		},
		{
			Name:        fmt.Sprintf("%s/success/latency", name),
			Measure:     operation.successLatency,
			Description: fmt.Sprintf("The latency of successful '%s' operations", name),
			Aggregation: view.LastValue(),
			TagKeys:     tagKeys,
		},
		{
			Name:        fmt.Sprintf("%s/success/count", name),
			Measure:     operation.success,
			Description: fmt.Sprintf("The number of time '%s' has succeeded", name),
			Aggregation: view.Count(),
			TagKeys:     tagKeys,
		},
	}

	view.Register(views...)

	return operation
}

// SetTag sets a sticky tag which will be emitted with every counter log
func (o Operation) SetTag(key tag.Key, value string) Operation {
	o.ctx, _ = tag.New(o.ctx, tag.Upsert(key, value))
	return o
}

// DeleteTag delets a sticky tag so it is no longer emitted with every counter log
func (o Operation) DeleteTag(key tag.Key) Operation {
	o.ctx, _ = tag.New(o.ctx, tag.Delete(key))
	return o
}

// Start indicates the operation has started
func (o Operation) Start() Operation {
	o.startTime = time.Now().UTC().Unix()
	stats.Record(o.ctx, o.start.M(1))
	return o
}

// Success indicates the operation has completed successfully
func (o Operation) Success() {
	n := time.Now().UTC().Unix()
	stats.Record(
		o.ctx,
		o.success.M(1),
		o.successLatency.M(n-o.startTime),
	)
	o.startTime = 0
}

// Failure indicates the operation has completed successfully
func (o Operation) Failure(errorCode string) {
	n := time.Now().UTC().Unix()
	stats.RecordWithTags(
		o.ctx,
		[]tag.Mutator{tag.Upsert(ErrorCodeTag, errorCode)},
		o.failed.M(1),
		o.failedLatency.M(n-o.startTime),
	)
	o.startTime = 0
}
