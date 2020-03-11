// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xserver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/facebookincubator/symphony/pkg/log"
	"github.com/facebookincubator/symphony/pkg/oc"
	"github.com/facebookincubator/symphony/pkg/server"
	"github.com/facebookincubator/symphony/pkg/server/recovery"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/google/wire"
	promclient "github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"gocloud.dev/server/requestlog"
)

// ServiceSet is a wire provider set for services.
var ServiceSet = wire.NewSet(
	Set,
	oc.Set,
	wire.Value(server.ProfilingEnabler(true)),
)

// Set is a wire provider set that provides the diagnostic hooks for *server.Server.
var Set = wire.NewSet(
	server.Set,
	NewRequestLogger,
	wire.Bind(new(requestlog.Logger), new(*ZapLogger)),
	NewRecoveryHandler,
	NewJaegerExporter,
	NewPrometheusExporter,
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

// NewJaegerExporter returns a new jaeger trace exporter.
func NewJaegerExporter(logger log.Logger, opts jaeger.Options) (trace.Exporter, func(), error) {
	if opts.AgentEndpoint == "" && opts.CollectorEndpoint == "" {
		return nil, func() {}, nil
	}

	if opts.Process.ServiceName == "" {
		if exec, err := os.Executable(); err == nil {
			opts.Process.ServiceName = filepath.Base(exec)
		}
	}
	if opts.OnError == nil {
		opts.OnError = func(err error) {
			logger.Background().Warn("jaeger exporter error", zap.Error(err))
		}
	}

	exporter, err := jaeger.NewExporter(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("creating jaeger exporter: %w", err)
	}
	return exporter, exporter.Flush, nil
}

// NewPrometheusExporter returns a new prometheus view exporter.
func NewPrometheusExporter(logger log.Logger) (view.Exporter, error) {
	return prometheus.NewExporter(prometheus.Options{
		Registry: promclient.DefaultRegisterer.(*promclient.Registry),
		OnError: func(err error) {
			logger.Background().Warn("prometheus exporter error", zap.Error(err))
		},
	})
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
