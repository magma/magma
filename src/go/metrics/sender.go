// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"context"
	"sync"

	io_prometheus_client "github.com/prometheus/client_model/go"

	"github.com/magma/magma/src/go/protos/magma/metricsd"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
)

// MetricSender is only responsible for sending collected metrics to a
// configured destination.
// MetricSender is not intended to cache any metric values, so if a Send
// fails for any reason, the caller should attempt to call Send again later.
type MetricSender interface {
	// Send sends all metric samples to the configured destination.
	// If Send fails, the caller should attempt to call Send again later with
	// the same values. Including newer values as well is permissible.
	Send(gauges []Gauge) error
}

//go:generate go run github.com/golang/mock/mockgen -destination mock_sender/mock_sender.go . MetricSender

type prometheusMetricSender struct {
	metricsContainer *metricsd.MetricsContainer
	gatewayId        string
	client           metricsd.MetricsControllerClient
	// Used to gate access to metricsContainer
	sync.Mutex
}

func NewPrometheusMetricSender(client metricsd.MetricsControllerClient, gatewayId string) MetricSender {
	container := &metricsd.MetricsContainer{
		GatewayId: gatewayId,
		Family:    []*io_prometheus_client.MetricFamily{},
	}
	return &prometheusMetricSender{metricsContainer: container, gatewayId: gatewayId, client: client}
}

func (s *prometheusMetricSender) Send(gauges []Gauge) error {
	s.processGauges(gauges)
	container := s.popMetrics()
	_, err := s.client.Collect(context.Background(), container)
	if err != nil {
		s.Lock()
		defer s.Unlock()
		s.metricsContainer.Family = append(container.Family, s.metricsContainer.Family...)
		return errors.Wrap(err, "failed to send Prometheus metrics")
	}
	return nil
}

// popMetrics returns the metricsContainer, and the local field is emptied.
func (s *prometheusMetricSender) popMetrics() *metricsd.MetricsContainer {
	s.Lock()
	defer s.Unlock()
	metricsContainer := s.metricsContainer

	s.metricsContainer = &metricsd.MetricsContainer{
		GatewayId: s.gatewayId,
		Family:    []*io_prometheus_client.MetricFamily{},
	}
	return metricsContainer
}

func (s *prometheusMetricSender) processGauges(gauges []Gauge) {
	var wg sync.WaitGroup
	for _, gauge := range gauges {
		wg.Add(1)
		go func(gauge Gauge) {
			defer wg.Done()
			s.processGauge(gauge)
		}(gauge)
	}
	wg.Wait()

}

// processGauge pops all recent samples from a Gauge and adds a MetricFamily to be sent
func (s *prometheusMetricSender) processGauge(gauge Gauge) {
	metricValues := []*io_prometheus_client.Metric{}

	samples := gauge.GetSamples()
	for _, sample := range samples {
		promMetric := s.processGaugeSample(sample)
		metricValues = append(metricValues, promMetric)
	}
	store := &io_prometheus_client.MetricFamily{
		Name:   proto.String(gauge.Name()),
		Help:   proto.String(gauge.Description()),
		Type:   io_prometheus_client.MetricType_GAUGE.Enum(),
		Metric: metricValues,
	}
	s.Lock()
	defer s.Unlock()
	s.metricsContainer.Family = append(s.metricsContainer.Family, store)
}

// processGaugeSample processes a GaugeSample into a prometheus Metric value
func (s *prometheusMetricSender) processGaugeSample(sample GaugeSample) *io_prometheus_client.Metric {
	labelPairs := []*io_prometheus_client.LabelPair{}
	for labelName, labelValue := range sample.LabelValues {
		labelPair := &io_prometheus_client.LabelPair{
			Name:  proto.String(labelName),
			Value: proto.String(labelValue),
		}
		labelPairs = append(labelPairs, labelPair)
	}
	return &io_prometheus_client.Metric{
		Label: labelPairs,
		Gauge: &io_prometheus_client.Gauge{
			Value: proto.Float64(sample.Value),
		},
		TimestampMs: proto.Int64(sample.TimestampMs),
	}
}
