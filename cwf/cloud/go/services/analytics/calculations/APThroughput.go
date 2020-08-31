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
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"magma/cwf/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/metrics"
	"time"
)

type APThroughputCalculation struct {
	CalculationParams
	QueryStepSize time.Duration
	Direction     ConsumptionDirection
}

func (x *APThroughputCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]Result, error) {
	// Get datapoints for throughput when the value is not 0 segmented by apn
	avgRateQuery := fmt.Sprintf(`avg(rate(octets_%s[3m]) > 0) by (%s, %s)`, x.Direction, APNLabel, metrics.NetworkLabelName)

	timeRange := v1.Range{End: time.Now(), Start: time.Now().Add(-time.Duration(x.Days * int(time.Hour) * 24)), Step: x.QueryStepSize}
	avgRateMatrix, err := query_api.QueryPrometheusMatrix(prometheusClient, avgRateQuery, timeRange)
	if err != nil {
		return nil, fmt.Errorf("AP Throughput query error: %s", err)
	}

	results := make([]Result, 0)
	for _, apnAverages := range avgRateMatrix {
		apn := string(apnAverages.Metric[APNLabel])
		nID := string(apnAverages.Metric[metrics.NetworkLabelName])
		avgThroughputOverTime := averageDatapoints(apnAverages.Values)
		if apn == "" || nID == "" {
			glog.Errorf("Missing tags from AP Throughput Calculation: APN: %s, NetworkID: %s", apn, nID)
			continue
		}
		results = append(results, Result{
			value:      avgThroughputOverTime,
			metricName: x.Name,
			labels:     combineLabels(x.Labels, map[string]string{APNLabel: apn, metrics.NetworkLabelName: nID, DirectionLabel: string(x.Direction)}),
		})
	}
	for _, res := range results {
		x.RegisteredGauge.With(res.labels).Set(res.value)
	}
	return results, nil
}
