/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package helpers

import (
	"context"
	"go.opencensus.io/stats"
)

func IncrementCounter(ctx context.Context, stat *stats.Int64Measure) {
	stats.Record(ctx, stat.M(1))
}

func ResetCounter(ctx context.Context, stat *stats.Int64Measure) {
	stats.Record(ctx, stat.M(0))
}

func SetCounter(ctx context.Context, stat *stats.Int64Measure, quantity int64) {
	stats.Record(ctx, stat.M(quantity))
}
