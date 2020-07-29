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

	"contrib.go.opencensus.io/exporter/aws"
	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	"fbc/lib/go/http/server"
	"fbc/lib/go/oc/helpers"
	"github.com/pkg/errors"
	prom_client "github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

type (
	xrayOptions struct {
		Service, Region, Version string
	}

	// CensusConfig defines opencensus exporters config.
	CensusConfig struct {
		XRay       *xrayOptions
		Jaeger     *jaeger.Options
		Prometheus *prometheus.Options
	}
)

// ServerResponseCountByStatusAndPath is an additional view for server response status code and path.
var ServerResponseCountByStatusAndPath = &view.View{
	Name:        "opencensus.io/http/server/response_count_by_status_code_path",
	Description: "Server response count by status code and path",
	TagKeys:     []tag.Key{ochttp.StatusCode, ochttp.KeyServerRoute},
	Measure:     ochttp.ServerLatency,
	Aggregation: view.Count(),
}

// NewConfig creates census server config.
func NewConfig(config string) (*CensusConfig, error) {
	var cc CensusConfig
	if err := json.Unmarshal([]byte(config), &cc); err != nil {
		return nil, errors.Wrapf(err, "parsing census config: %q", config)
	}
	return &cc, nil
}

// WithService sets service name on underlying exporter configs.
func (cc *CensusConfig) WithService(service string) *CensusConfig {
	if cc.XRay != nil && cc.XRay.Service == "" {
		cc.XRay.Service = service
	}
	if cc.Jaeger != nil && cc.Jaeger.Process.ServiceName == "" {
		cc.Jaeger.Process.ServiceName = service
	}
	if cc.Prometheus != nil && cc.Prometheus.Namespace == "" {
		cc.Prometheus.Namespace = service
	}
	return cc
}

// ServerOptions builds server options from internal state.
func (cc *CensusConfig) ServerOptions() (opts []server.Option) {
	if cc.XRay != nil {
		opts = append(opts, xrayOption(*cc.XRay))
	}
	if cc.Jaeger != nil {
		opts = append(opts, jaegerOption(*cc.Jaeger))
	}
	if cc.Prometheus != nil {
		opts = append(opts, prometheusOption(*cc.Prometheus))
	}
	if len(opts) > 0 {
		opts = append(opts, server.OptionFunc(func(*server.Server) error {
			trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
			return nil
		}))
	}
	return
}

func xrayOption(options xrayOptions) server.Option {
	var opts []aws.Option
	if options.Service != "" {
		opts = append(opts, aws.WithServiceName(options.Service))
	}
	if options.Region != "" {
		opts = append(opts, aws.WithRegion(options.Region))
	}
	if options.Version == "" {
		options.Version = "latest"
	}
	opts = append(opts, aws.WithVersion(options.Version))

	return server.OptionFunc(func(srv *server.Server) error {
		opts = append(opts, aws.WithOutput(writerFunc(func(p []byte) (int, error) {
			srv.Logger.Bg().Warn("xray exporter failure", zap.ByteString("error", p))
			return len(p), nil
		})))
		exporter, err := aws.NewExporter(opts...)
		if err != nil {
			return errors.Wrap(err, "creating xray exporter")
		}
		trace.RegisterExporter(exporter)
		err = srv.Apply(server.Closer(closerFunc(func() error { exporter.Flush(); return nil })))
		if err != nil {
			return errors.WithMessage(err, "registering xray flusher")
		}
		return nil
	})
}

func jaegerOption(options jaeger.Options) server.Option {
	return server.OptionFunc(func(srv *server.Server) error {
		options.OnError = func(err error) {
			srv.Logger.Bg().Warn("jaeger exporter failure", zap.Error(err))
		}
		exporter, err := jaeger.NewExporter(options)
		if err != nil {
			return errors.Wrap(err, "creating jaeger exporter")
		}
		trace.RegisterExporter(exporter)
		err = srv.Apply(server.Closer(closerFunc(func() error { exporter.Flush(); return nil })))
		if err != nil {
			return errors.WithMessage(err, "registering jaeger flusher")
		}
		return nil
	})
}

func prometheusOption(options prometheus.Options) server.Option {
	return server.OptionFunc(func(srv *server.Server) error {
		options.OnError = func(err error) {
			srv.Logger.Bg().Warn("prometheus exporter failure", zap.Error(err))
		}
		options.Registry = prom_client.NewRegistry()
		exporter, err := prometheus.NewExporter(options)
		if err != nil {
			return errors.Wrap(err, "creating prometheus exporter")
		}

		// Adding process collector
		if err := options.Registry.Register(prom_client.NewProcessCollector(
			prom_client.ProcessCollectorOpts{Namespace: options.Namespace},
		)); err != nil {
			return errors.Wrap(err, "registering process collector")
		}

		// Adding GO collector
		if err := options.Registry.Register(prom_client.NewGoCollector()); err != nil {
			return errors.Wrap(err, "registering go collector")
		}
		if err := view.Register(
			ochttp.ServerRequestCountView,
			ochttp.ServerRequestBytesView,
			ochttp.ServerResponseBytesView,
			ochttp.ServerLatencyView,
			ochttp.ServerRequestCountByMethod,
			ochttp.ServerResponseCountByStatusCode,
			ServerResponseCountByStatusAndPath,
		); err != nil {
			return errors.Wrap(err, "registering http server views")
		}
		if err := view.Register(
			ochttp.ClientCompletedCount,
			ochttp.ClientSentBytesDistribution,
			ochttp.ClientReceivedBytesDistribution,
			ochttp.ClientRoundtripLatencyDistribution,
			ochttp.ClientCompletedCount,
		); err != nil {
			return errors.Wrap(err, "registering http client views")
		}
		if err := view.Register(
			helpers.LatencyView,
			helpers.ErrorCountView,
			helpers.SuccessCountView,
			helpers.CountView,
		); err != nil {
			return errors.Wrap(err, "registering customized KPI views")
		}

		view.RegisterExporter(exporter)
		srv.Mux.Handle("/metrics", exporter)
		return nil
	})
}

type closerFunc func() error

func (f closerFunc) Close() error {
	return f()
}

type writerFunc func([]byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) {
	return f(p)
}
