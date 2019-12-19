// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"net/http"
	"net/http/pprof"
	"path"
	"sync"
	"time"

	"github.com/facebookincubator/symphony/pkg/server/driver"
	"github.com/facebookincubator/symphony/pkg/server/recovery"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"gocloud.dev/server/health"
	"gocloud.dev/server/requestlog"
)

// Set is a Wire provider set that produces a *Server given the fields of Options.
var Set = wire.NewSet(
	New,
	wire.Struct(new(Options), "*"),
	wire.Value(&DefaultDriver{}),
	wire.Bind(new(driver.Server), new(*DefaultDriver)),
)

// Server is a preconfigured HTTP server with diagnostic hooks.
// The zero value is a server with the default options.
type Server struct {
	reqlog  requestlog.Logger
	handler http.Handler
	health  health.Handler
	metrics http.Handler
	views   []*view.View
	ve      view.Exporter
	te      trace.Exporter
	sampler trace.Sampler
	profile ProfilingEnabler
	recover recovery.HandlerFunc
	once    sync.Once
	driver  driver.Server
}

// Options is the set of optional parameters.
type Options struct {
	// RequestLogger specifies the logger that will be used to log requests.
	RequestLogger requestlog.Logger

	// HealthChecks specifies the health checks to be run when the
	// /healthz/readiness endpoint is requested.
	HealthChecks []health.Checker

	// Views construct view data.
	Views []*view.View

	// ViewExporter exports view data.
	ViewExporter view.Exporter

	// TraceExporter exports sampled trace spans.
	TraceExporter trace.Exporter

	// EnableProfiling enables server profiling.
	EnableProfiling ProfilingEnabler

	// DefaultSamplingPolicy is a function that takes a
	// trace.SamplingParameters struct and returns a true or false decision about
	// whether it should be sampled and exported.
	DefaultSamplingPolicy trace.Sampler

	// RecoveryHandler handles panic recovery.
	RecoveryHandler recovery.HandlerFunc

	// Driver serves HTTP requests.
	Driver driver.Server
}

// ProfilingEnabler toggles server profiling.
type ProfilingEnabler bool

// New creates a new server.
func New(h http.Handler, opts *Options) *Server {
	srv := &Server{handler: h}
	if opts != nil {
		srv.reqlog = opts.RequestLogger
		for _, c := range opts.HealthChecks {
			srv.health.Add(c)
		}
		srv.ve = opts.ViewExporter
		srv.metrics, _ = opts.ViewExporter.(http.Handler)
		srv.views = opts.Views
		srv.te = opts.TraceExporter
		srv.sampler = opts.DefaultSamplingPolicy
		srv.profile = opts.EnableProfiling
		srv.recover = opts.RecoveryHandler
		srv.driver = opts.Driver
	}
	return srv
}

func (srv *Server) init() {
	srv.once.Do(func() {
		if srv.ve != nil {
			view.RegisterExporter(srv.ve)
		}
		if srv.te != nil {
			trace.RegisterExporter(srv.te)
		}
		if srv.sampler != nil {
			trace.ApplyConfig(trace.Config{DefaultSampler: srv.sampler})
		}
		if srv.driver == nil {
			srv.driver = NewDefaultDriver()
		}
		if srv.handler == nil {
			srv.handler = http.DefaultServeMux
		}
	})
}

// ListenAndServe is a wrapper to use wherever http.ListenAndServe is used.
// It wraps the passed-in http.Handler with a handler that handles metrics, tracing
// and request logging. If the handler is nil, then http.DefaultServeMux will be used.
func (srv *Server) ListenAndServe(addr string) error {
	srv.init()

	// Setup health checks, /healthz route is taken by health checks by default.
	hr := "/healthz"
	mux := http.NewServeMux()
	mux.HandleFunc(hr, health.HandleLive)
	mux.HandleFunc(path.Join(hr, "liveness"), health.HandleLive)
	mux.Handle(path.Join(hr, "readiness"), &srv.health)

	// Setup metrics endpoint, /metrics route is taken by default.
	if srv.metrics != nil {
		mux.Handle("/metrics", srv.metrics)
	}
	// Register metrics views
	if err := view.Register(srv.views...); err != nil {
		return errors.Wrap(err, "registering views")
	}

	// Setup profiling, /debug/pprof route is taken by default.
	if srv.profile {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	// Setup middleware chain
	var h http.Handler = &recovery.Handler{
		Handler:     srv.handler,
		HandlerFunc: srv.recover,
	}
	if srv.reqlog != nil {
		h = requestlog.NewHandler(srv.reqlog, h)
	}
	h = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if span := trace.FromContext(r.Context()); span != nil {
				if sc := span.SpanContext(); sc.IsSampled() {
					w.Header().Set("X-Correlation-ID", sc.TraceID.String())
				}
			}
			h.ServeHTTP(w, r)
		})
	}(h)
	h = &ochttp.Handler{
		Handler: h,
		FormatSpanName: func(r *http.Request) string {
			return r.URL.Host + r.URL.Path
		},
	}

	mux.Handle("/", h)
	return srv.driver.ListenAndServe(addr, mux)
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (srv *Server) Shutdown(ctx context.Context) error {
	if srv.driver == nil {
		return nil
	}
	return srv.driver.Shutdown(ctx)
}

// DefaultDriver implements the driver.Server interface. The zero value is a valid http.Server.
type DefaultDriver struct {
	Server http.Server
}

// NewDefaultDriver creates a driver with an http.Server with default timeouts.
func NewDefaultDriver() *DefaultDriver {
	return &DefaultDriver{
		Server: http.Server{
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  2 * time.Minute,
		},
	}
}

// ListenAndServe sets the address and handler on DefaultDriver's http.Server,
// then calls ListenAndServe on it.
func (dd *DefaultDriver) ListenAndServe(addr string, h http.Handler) error {
	dd.Server.Addr = addr
	dd.Server.Handler = h
	return dd.Server.ListenAndServe()
}

// Shutdown gracefully shuts down the server without interrupting any active connections,
// by calling Shutdown on DefaultDriver's http.Server
func (dd *DefaultDriver) Shutdown(ctx context.Context) error {
	return dd.Server.Shutdown(ctx)
}
