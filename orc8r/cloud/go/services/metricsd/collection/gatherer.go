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

package collection

import (
	"fmt"
	"time"

	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	prometheus_proto "github.com/prometheus/client_model/go"
)

// MetricsGatherer wraps a set of MetricCollectors, polling each collector
// at the configured interval and putting the results onto an output channel.
type MetricsGatherer struct {
	StaticCollectors []MetricCollector
	CollectInterval  time.Duration
	OutputChan       chan *prometheus_proto.MetricFamily
}

// NewMetricsGatherer validates params and returns a new metrics gatherer.
// The gatherer accepts a static set of collectors, to which it will
// dynamically append a collector per orc8r service.
func NewMetricsGatherer(
	staticCollectors []MetricCollector,
	collectInterval time.Duration,
	outputChan chan *prometheus_proto.MetricFamily,
) (*MetricsGatherer, error) {
	if collectInterval < 0 {
		return nil, fmt.Errorf("collectInterval must be positive")
	}
	return &MetricsGatherer{
		StaticCollectors: staticCollectors,
		CollectInterval:  collectInterval,
		OutputChan:       outputChan,
	}, nil
}

func (g *MetricsGatherer) Run() {
	glog.V(2).Info("Running metrics gatherer")
	// Gather metrics from each collector periodically in separate goroutines
	// so a hanging collector doesn't block other collectors
	for _, collector := range g.getCollectors() {
		go g.gatherEvery(collector)
	}
}

func (g *MetricsGatherer) gatherEvery(collector MetricCollector) {
	for range time.Tick(g.CollectInterval) {
		fams, err := collector.GetMetrics()
		if err != nil {
			glog.Errorf("Metric collector error: %s", err)
		}
		for _, fam := range fams {
			g.OutputChan <- fam
		}
	}
}

// getCollectors returns the set of metrics collectors.
// Returned collectors include disk usage, process statistics, and
// per-service custom metrics.
func (g *MetricsGatherer) getCollectors() []MetricCollector {
	collectors := g.StaticCollectors

	services, err := registry.ListAllServices()
	if err != nil {
		err = errors.Wrap(err, "error getting metrics collectors: list all services")
		glog.Warning(err)
		return collectors
	}

	for _, s := range services {
		collectors = append(collectors, NewCloudServiceMetricCollector(s))
	}

	return collectors
}
