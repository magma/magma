// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xserver

import (
	"context"

	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/pkg/server/recovery"
	"github.com/facebookincubator/symphony/pkg/telemetry"

	"github.com/google/wire"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.uber.org/zap"
	"gocloud.dev/server/requestlog"
)

// ServiceSet is a wire provider set for services.
var ServiceSet = wire.NewSet(
	Set,
	telemetry.Provider,
	wire.Value(server.ProfilingEnabler(true)),
)

// Set is a wire provider set that provides the diagnostic hooks for *server.Server.
var Set = wire.NewSet(
	server.Set,
	NewRequestLogger,
	wire.Bind(new(requestlog.Logger), new(*ZapLogger)),
	NewRecoveryHandler,
)

// NewRequestLogger returns a request logger that sends entries to background logger.
func NewRequestLogger(logger log.Logger) *ZapLogger {
	return NewZapLogger(logger.Background())
}

// NewRecoveryHandler returns a panic recovery handler that logs to background logger.
func NewRecoveryHandler(logger log.Logger) recovery.HandlerFunc {
	return func(ctx context.Context, p interface{}) error {
		err, _ := p.(error)
		logger.For(ctx).Error("panic recovery", zap.Error(err), zap.Stack("stacktrace"))
		return nil
	}
}

// DefaultViews are predefined views for OpenCensus metrics.
func DefaultViews() []*view.View {
	return []*view.View{
		func() *view.View {
			v := ochttp.ServerResponseCountByStatusCode.
				WithName("http_request_total")
			v.Description = "Total number of HTTP requests"
			v.TagKeys = []tag.Key{
				ochttp.Method,
				ochttp.Path,
				ochttp.StatusCode,
			}
			return v
		}(),
		ochttp.ServerLatencyView.
			WithName("http_request_duration_milliseconds"),
		ochttp.ServerRequestBytesView.
			WithName("http_request_size_bytes"),
		ochttp.ServerResponseBytesView.
			WithName("http_response_size_bytes"),
	}
}
