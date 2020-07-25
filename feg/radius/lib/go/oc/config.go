/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package oc

import (
	"encoding/json"
	"fbc/lib/go/oc/ocstats"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

type (
	// Config offers a declarative way to construct a census.
	Config struct {
		DisableStats        bool      `env:"NO_STATS" long:"no-stats" description:"Disables statistics gathering and exporting" json:"disable_stats"`
		StatViews           StatViews `env:"VIEWS" long:"view" default:"proc" description:"Set of metric types to expose" json:"stat_views"`
		SamplingProbability float64   `env:"SAMPLING_PROBABILITY" long:"sampling-probability" default:"1.0" description:"Trace sampling probability" json:"sampling_probability"`
		Jaeger              *Jaeger   `env:"JAEGER" long:"jaeger" description:"Jaeger exporter options as json" json:"jaeger"`
	}

	// StatViews attaches flags methods to []string.
	StatViews []string

	// Jaeger attaches flags methods to jaeger.Options.
	Jaeger struct{ jaeger.Options }

	// An Option configured census.
	Option func(*options)

	options struct {
		service   string
		namespace string
		logger    *zap.Logger
	}

	// Census defines opencensus runtime settings.
	Census struct {
		// handler for /metrics endpoint.
		StatsHandler http.Handler

		// Set of closers to run on Close.
		closers []func()
	}
)

// WithService sets census service name.
func WithService(name string) Option {
	return func(opts *options) {
		opts.service = name
	}
}

// WithNamespace sets census namespace.
func WithNamespace(ns string) Option {
	return func(opts *options) {
		opts.namespace = ns
	}
}

// WithLogger sets census logger.
func WithLogger(logger *zap.Logger) Option {
	return func(opts *options) {
		opts.logger = logger
	}
}

// Build constructs a census from the Config and Options.
func (cfg Config) Build(opt ...Option) (census *Census, err error) {
	var closers []func()
	defer func() {
		if err != nil {
			for _, closer := range closers {
				closer()
			}
		}
	}()

	var (
		opts    = cfg.buildOptions(opt)
		handler http.Handler
	)
	if handler, closers, err = cfg.buildStats(opts); err != nil {
		return nil, err
	}

	var tc []func()
	if tc, err = cfg.buildTraces(opts); err != nil {
		return nil, err
	}
	closers = append(closers, tc...)

	return &Census{
		StatsHandler: handler,
		closers:      closers,
	}, nil
}

func (Config) buildOptions(opts []Option) options {
	var o options
	for _, opt := range opts {
		opt(&o)
	}
	// fill defaults
	if o.service == "" {
		if exec, err := os.Executable(); err == nil {
			o.service = filepath.Base(exec)
		}
	}
	if o.logger == nil {
		o.logger = zap.L()
	}
	return o
}

func (cfg *Config) buildStats(opt options) (handler http.Handler, closers []func(), err error) {
	// nothing to do if not enabled
	if cfg.DisableStats {
		handler = http.NotFoundHandler()
		return
	}

	// run accumulated closers on error
	defer func() {
		if err == nil {
			return
		}
		for _, closer := range closers {
			closer()
		}
	}()

	// start with basic stats options
	opts := []ocstats.Option{
		ocstats.WithNamespace(opt.namespace),
		ocstats.WithLogger(opt.logger),
	}

	// track previously processed views
	views := map[string]struct{}{}
	for _, name := range cfg.StatViews {
		if _, ok := views[name]; ok {
			continue
		}
		views[name] = struct{}{}
		if name == "proc" {
			opts = append(opts,
				ocstats.WithProcessCollector(),
				ocstats.WithGoCollector(),
			)
			continue
		}
		viewer := GetViewer(name)
		if viewer == nil {
			err = errors.Errorf("unknown view name %q", name)
			return
		}
		views := viewer.Views()
		if err = view.Register(views...); err != nil {
			err = errors.WithMessagef(err, "registering %s views", name)
			return
		}
		closers = append(closers, func() {
			view.Unregister(views...)
		})
	}

	var closer func()
	if handler, closer, err = ocstats.NewHandler(opts...); err != nil {
		err = errors.WithMessage(err, "creating stats handler")
		return
	}
	closers = append(closers, closer)

	return handler, closers, nil
}

func (cfg *Config) buildTraces(opt options) (closers []func(), err error) {
	// nothing to do if not enabled
	if cfg.SamplingProbability <= 0 {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.NeverSample()})
		return
	}

	// run accumulated closers on error
	defer func() {
		if err == nil {
			return
		}
		for _, closer := range closers {
			closer()
		}
	}()

	// build jaeger exporter
	if cfg.Jaeger != nil {
		var closer func()
		if closer, err = cfg.buildJaeger(opt); err != nil {
			return
		}
		closers = append(closers, closer)
	}

	// configure sampling rate / sampling rate undoer
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(cfg.SamplingProbability)})
	closers = append(closers, func() {
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.NeverSample()})
	})

	return closers, nil
}

func (cfg *Config) buildJaeger(opt options) (func(), error) {
	if cfg.Jaeger.Process.ServiceName == "" {
		cfg.Jaeger.Process.ServiceName = opt.service
	}
	cfg.Jaeger.OnError = func(err error) {
		opt.logger.Warn("jaeger export failure", zap.Error(err))
	}
	exporter, err := jaeger.NewExporter(cfg.Jaeger.Options)
	if err != nil {
		return nil, errors.Wrap(err, "creating jaeger exporter")
	}
	trace.RegisterExporter(exporter)
	return func() {
		exporter.Flush()
		trace.UnregisterExporter(exporter)
	}, nil
}

// Close implements io.Closer interface.
func (c *Census) Close() error {
	for _, closer := range c.closers {
		closer()
	}
	return nil
}

// UnmarshalFlag implements flags.Unmarshaler interface.
func (s *StatViews) UnmarshalFlag(value string) error {
	*s = append(*s, strings.Split(value, ",")...)
	return nil
}

// IsValidValue implements flags.ValueValidator interface.
func (StatViews) IsValidValue(value string) error {
	for _, name := range strings.Split(value, ",") {
		if GetViewer(name) == nil {
			return &flags.Error{
				Type:    flags.ErrInvalidChoice,
				Message: "Invalid viewer name " + name,
			}
		}
	}
	return nil
}

// UnmarshalFlag implements flags.Unmarshaler interface.
func (j *Jaeger) UnmarshalFlag(value string) error {
	if err := json.Unmarshal([]byte(value), &j.Options); err != nil {
		return &flags.Error{
			Type:    flags.ErrMarshal,
			Message: err.Error(),
		}
	}
	return nil
}
