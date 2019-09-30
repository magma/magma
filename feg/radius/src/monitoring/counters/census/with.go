package census

import (
	ocprom "contrib.go.opencensus.io/exporter/prometheus"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// WithNamespace sets handler namespace.
func WithNamespace(ns string) Option {
	return func(opts *ocprom.Options) error {
		opts.Namespace = ns
		return nil
	}
}

// WithLogger will emit warnings on errors.
func WithLogger(logger *zap.Logger) Option {
	return func(opts *ocprom.Options) error {
		opts.OnError = func(err error) {
			logger.Warn("prometheus exporter error", zap.Error(err))
		}
		return nil
	}
}

// WithProcessCollector registers prometheus process collector.
func WithProcessCollector() Option {
	return func(opts *ocprom.Options) error {
		if err := opts.Registry.Register(prometheus.NewProcessCollector(
			prometheus.ProcessCollectorOpts{Namespace: opts.Namespace},
		)); err != nil {
			return errors.Wrap(err, "registering process collector")
		}
		return nil
	}
}

// WithGoCollector registers prometheus go collector.
func WithGoCollector() Option {
	return func(opts *ocprom.Options) error {
		if err := opts.Registry.Register(prometheus.NewGoCollector()); err != nil {
			return errors.Wrap(err, "registering go collector")
		}
		return nil
	}
}
