/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package calculations

import (
	"fmt"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
	"github.com/influxdata/tdigest"
)

type UserThroughputCalculation struct {
	calculations.BaseCalculation
	Direction calculations.ConsumptionDirection
}

// Calculate computes the average and p95 quantile for user upload/download throughput.
func (x *UserThroughputCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Info("Calculate UserThroughputCalculation Metrics")

	q := fmt.Sprintf(`sum(rate(ue_reported_usage{direction="%s"}[5m])) by (%s, IMSI)`, x.Direction, metrics.NetworkLabelName)
	rateVector, err := query_api.QueryPrometheusVector(prometheusClient, q)
	if err != nil {
		return nil, fmt.Errorf("User Throughput query %s error: %s", q, err)
	}
	results := make([]*protos.CalculationResult, 0)

	networkThroughputVector := map[string][]float64{}
	for _, v := range rateVector {
		networkID := ""
		for label, value := range v.Metric {
			if string(label) == metrics.NetworkLabelName {
				networkID = string(value)
				break
			}
		}
		_, ok := networkThroughputVector[networkID]
		if !ok {
			networkThroughputVector[networkID] = make([]float64, 0)
		}
		networkThroughputVector[networkID] = append(networkThroughputVector[networkID], float64(v.Value))
	}

	for networkID, networkThroughputVec := range networkThroughputVector {
		if len(networkThroughputVec) == 0 {
			continue
		}
		td := tdigest.NewWithCompression(1000)
		for _, v := range networkThroughputVec {
			td.Add(v, 1)
		}
		p50thValue := td.Quantile(0.5)
		p95thValue := td.Quantile(0.95)

		labels := map[string]string{
			calculations.DirectionLabel: string(x.Direction),
			metrics.NetworkLabelName:    networkID,
		}
		glog.V(2).Infof("Network ID %s p50th %f p95 %f direction %s", networkID, p50thValue, p95thValue, string(x.Direction))
		results = append(results, calculations.NewResult(
			p50thValue,
			metrics.SubscriberThroughputMetric,
			calculations.CombineLabels(labels, map[string]string{metrics.QuantileLabel: "0.5"})))
		results = append(results, calculations.NewResult(
			p95thValue,
			metrics.SubscriberThroughputMetric,
			calculations.CombineLabels(labels, map[string]string{metrics.QuantileLabel: "0.95"})))
	}

	return results, nil
}
