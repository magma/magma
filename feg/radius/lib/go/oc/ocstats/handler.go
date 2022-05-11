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

package ocstats

import (
	"fmt"
	"net/http"

	ocprom "contrib.go.opencensus.io/exporter/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
)

// An Option configures a Handler.
type Option func(*ocprom.Options) error

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
			return fmt.Errorf("registering process collector: %w", err)
		}
		return nil
	}
}

// WithGoCollector registers prometheus go collector.
func WithGoCollector() Option {
	return func(opts *ocprom.Options) error {
		if err := opts.Registry.Register(prometheus.NewGoCollector()); err != nil {
			return fmt.Errorf("registering go collector: %w", err)
		}
		return nil
	}
}

// NewHandler creates a stats http handler.
func NewHandler(opt ...Option) (http.Handler, func(), error) {
	opts := ocprom.Options{Registry: prometheus.NewRegistry()}
	for i := range opt {
		if err := opt[i](&opts); err != nil {
			return nil, nil, fmt.Errorf("applying option: %w", err)
		}
	}
	exporter, err := ocprom.NewExporter(opts)
	if err != nil {
		return nil, nil, fmt.Errorf("creating prometheus exporter: %w", err)
	}
	view.RegisterExporter(exporter)
	closer := func() { view.UnregisterExporter(exporter) }
	return exporter, closer, nil
}
