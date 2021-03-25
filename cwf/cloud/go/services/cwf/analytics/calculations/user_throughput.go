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
	"time"

	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"

	"github.com/golang/glog"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

// UserThroughputCalculation params for computing average user throughput
type UserThroughputCalculation struct {
	calculations.BaseCalculation
	QueryStepSize time.Duration
	Direction     calculations.ConsumptionDirection
}

//Calculate method calculations average user throughput across the specified direction
func (x *UserThroughputCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	glog.V(1).Infof("Calculating User Throughput. Days: %d, Direction: %s", x.Days, x.Direction)
	// Get datapoints for throughput when the value is not 0 segmented
	avgRateQuery := fmt.Sprintf(`avg(rate(octets_%s[3m]) > 0) by (%s)`, x.Direction, metrics.NetworkLabelName)

	timeRange := v1.Range{End: time.Now(), Start: time.Now().Add(-time.Duration(int(x.Days) * int(time.Hour) * 24)), Step: x.QueryStepSize}
	avgRateMatrix, err := query_api.QueryPrometheusMatrix(prometheusClient, avgRateQuery, timeRange)
	if err != nil {
		return nil, fmt.Errorf("user Throughput query error: %s", err)
	}

	results := make([]*protos.CalculationResult, 0)
	for _, apnAverages := range avgRateMatrix {
		nID := string(apnAverages.Metric[metrics.NetworkLabelName])
		avgThroughputOverTime := calculations.AverageDatapoints(apnAverages.Values)
		if nID == "" {
			glog.Error("Missing NetworkID from Throughput Calculation")
			continue
		}
		results = append(results, &protos.CalculationResult{
			Value:      avgThroughputOverTime,
			MetricName: x.Name,
			Labels:     calculations.CombineLabels(x.Labels, map[string]string{metrics.NetworkLabelName: nID, calculations.DirectionLabel: string(x.Direction)}),
		})
	}
	return results, nil
}
