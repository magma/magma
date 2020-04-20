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

// package status implements magmad status amd metrics collectors & reporters
package status

import (
	"context"
	"fmt"

	prometheus "github.com/prometheus/client_model/go"

	"magma/gateway/service_registry"
	"magma/gateway/status"
	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/protos"
)

var (
	serviceLabelName = "service"
	metricsQueue     = &MetricsQueue{items: []*prometheus.MetricFamily{}}
)

func reportMetrics(maxQueueLen int) error {
	samples := collectMetrics()
	if len(samples) == 0 {
		return nil
	}
	metricsdConn, err := service_registry.Get().GetSharedCloudConnection(definitions.MetricsdServiceName)
	if err != nil {
		enqueueRetry(samples, maxQueueLen)
		return fmt.Errorf("failed to connect to metricsd service: %v", err)
	}
	_, err = protos.NewMetricsControllerClient(metricsdConn).Collect(
		context.Background(),
		&protos.MetricsContainer{
			GatewayId: status.GetHwId(),
			Family:    samples,
		})
	if err != nil {
		err = fmt.Errorf("metrics reporting error: %v", err)
		enqueueRetry(samples, maxQueueLen)
	}
	return err
}

func enqueueMetrics(service string, serviceMetrics *protos.MetricsContainer) int {
	if len(serviceMetrics.GetFamily()) == 0 {
		return 0
	}
	for _, f := range serviceMetrics.Family {
		if f != nil {
			for _, m := range f.Metric {
				m.Label = append(m.Label, &prometheus.LabelPair{
					Name:  &serviceLabelName,
					Value: &service,
				})
			}
		}
	}
	return metricsQueue.Append(serviceMetrics.Family...)
}

func collectMetrics() (result []*prometheus.MetricFamily) {
	return metricsQueue.Collect()
}

func enqueueRetry(retry []*prometheus.MetricFamily, maxQueueLen int) {
	metricsQueue.Prepend(retry, maxQueueLen)
}

func resetMetricsQueue() {
	metricsQueue.Reset()
}
