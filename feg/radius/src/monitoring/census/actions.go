package census

import (
	"fmt"
	"net/http"

	ocprom "contrib.go.opencensus.io/exporter/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
)

// Build constructs a census from the Config and Options.
func (cfg Config) Build(logger *zap.Logger, opts ...Option) (census *Census, err error) {
	var closers []func()
	defer func() {
		if err != nil {
			for _, closer := range closers {
				closer()
			}
		}
	}()

	opts = append(
		opts,
		WithNamespace("radius"),
		WithLogger(logger),
	)

	var handler http.Handler
	if handler, closers, err = cfg.buildStats(opts...); err != nil {
		return nil, err
	}

	return &Census{
		StatsHandler: handler,
		closers:      closers,
	}, nil
}

func (cfg *Config) buildStats(opts ...Option) (handler http.Handler, closers []func(), err error) {
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

	// track previously processed views
	views := map[string]struct{}{}
	for _, name := range cfg.StatViews {
		if _, ok := views[name]; ok {
			continue
		}
		views[name] = struct{}{}
		if name == "proc" {
			opts = append(opts,
				WithProcessCollector(),
				WithGoCollector(),
			)
			continue
		}
		viewer := GetViewer(name)
		if viewer == nil {
			err = fmt.Errorf("unknown view name %q", name)
			return
		}
		views := viewer.Views()
		if err = view.Register(views...); err != nil {
			return
		}
		closers = append(closers, func() {
			view.Unregister(views...)
		})
	}

	var closer func()
	if handler, closer, err = NewHandler(opts...); err != nil {
		return
	}
	closers = append(closers, closer)

	return handler, closers, nil
}

// NewHandler creates a stats http handler.
func NewHandler(opt ...Option) (http.Handler, func(), error) {
	opts := ocprom.Options{Registry: prometheus.NewRegistry()}
	for i := range opt {
		if err := opt[i](&opts); err != nil {
			return nil, nil, err
		}
	}
	exporter, err := ocprom.NewExporter(opts)
	if err != nil {
		return nil, nil, err
	}
	view.RegisterExporter(exporter)
	closer := func() { view.UnregisterExporter(exporter) }
	return exporter, closer, nil
}
