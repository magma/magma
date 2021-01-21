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
	io_prometheus_client "github.com/prometheus/client_model/go"
)

type MetricLabel struct {
	Name, Value string
}

// MakeSingleGaugeFamily returns a MetricFamily with a single gauge value
// as specified by the function arguments.
// label is nil-able - a nil input will return a gauge metric without a label.
func MakeSingleGaugeFamily(
	name string, help string,
	label *MetricLabel,
	value float64,
) *io_prometheus_client.MetricFamily {
	mtype := io_prometheus_client.MetricType_GAUGE
	return &io_prometheus_client.MetricFamily{
		Name:   &name,
		Help:   &help,
		Type:   &mtype,
		Metric: []*io_prometheus_client.Metric{MakeSingleGaugeMetric(label, value)},
	}
}

// MakeSingleGaugeMetric returns a Metric with a single gauge value as
// specified by the function argument.
// label is nil-able - a nil input will return a gauge metric without a label.
func MakeSingleGaugeMetric(
	label *MetricLabel,
	value float64,
) *io_prometheus_client.Metric {
	if label == nil {
		return &io_prometheus_client.Metric{
			Gauge: &io_prometheus_client.Gauge{Value: &value},
		}
	} else {
		return &io_prometheus_client.Metric{
			Label: []*io_prometheus_client.LabelPair{
				{Name: &label.Name, Value: &label.Value},
			},
			Gauge: &io_prometheus_client.Gauge{Value: &value},
		}
	}
}
