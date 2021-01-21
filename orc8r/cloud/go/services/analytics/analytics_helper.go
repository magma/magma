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
package analytics

import (
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
	"github.com/olivere/elastic/v7"
	promAPI "github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
)

func GetElasticClient() *elastic.Client {
	elasticConfig, err := config.GetServiceConfig(orc8r.ModuleName, "elastic")
	if err != nil {
		glog.Errorf("Error %v reading elastic service configuration", err)
		return nil
	}
	elasticHost := elasticConfig.MustGetString("elasticHost")
	elasticPort := elasticConfig.MustGetInt("elasticPort")

	client, err := elastic.NewSimpleClient(elastic.SetURL(fmt.Sprintf("http://%s:%d", elasticHost, elasticPort)))
	if err != nil {
		glog.Errorf("Error %v getting client handle to elastic service", err)
		return nil
	}
	return client
}

func GetPrometheusClient() v1.API {
	metricsConfig, err := config.GetServiceConfig(orc8r.ModuleName, metricsd.ServiceName)
	if err != nil {
		glog.Fatalf("Could not retrieve metricsd configuration: %s", err)
	}
	promClient, err := promAPI.NewClient(promAPI.Config{Address: metricsConfig.MustGetString(metricsd.PrometheusQueryAddress)})
	if err != nil {
		glog.Fatalf("Error creating prometheus client: %s", promClient)
	}
	return v1.NewAPI(promClient)
}

// GetLogMetricsCalculations gets all the log calculations in the component
func GetLogMetricsCalculations(cfg *calculations.AnalyticsConfig) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	elasticClient := GetElasticClient()
	for metricName, metricConfig := range cfg.Metrics {
		if metricConfig.LogConfig == nil {
			continue
		}

		labels := []string{}
		for k := range metricConfig.Labels {
			labels = append(labels, k)
		}
		glog.V(1).Infof("Adding Log Calculation for %s", metricName)
		gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metricName}, labels)
		prometheus.MustRegister(gauge)
		calcs = append(calcs, &calculations.LogsMetricCalculation{
			BaseCalculation: calculations.BaseCalculation{
				CalculationParams: calculations.CalculationParams{
					Name:                metricName,
					Hours:               GetServiceConfig().AnalysisSchedule,
					AnalyticsConfig:     cfg,
					Labels:              metricConfig.Labels,
					RegisteredGauge:     gauge,
					ExpectedGaugeLabels: labels,
				},
			},
			LogConfig:     metricConfig.LogConfig,
			ElasticClient: elasticClient,
		})
	}
	return calcs
}

// GetRawMetricsCalculations gets all raw metrics calculations configured for the service
func GetRawMetricsCalculations(cfg *calculations.AnalyticsConfig) []calculations.Calculation {
	calcs := make([]calculations.Calculation, 0)
	for metricName, metricConfig := range cfg.Metrics {
		if metricConfig.Expr == "" {
			continue
		}
		glog.V(10).Infof("Adding RawMetrics Calculation for %s", metricName)
		calcs = append(calcs, &calculations.RawMetricsCalculation{
			BaseCalculation: calculations.BaseCalculation{
				CalculationParams: calculations.CalculationParams{
					Name:            metricName,
					AnalyticsConfig: cfg,
					Hours:           GetServiceConfig().AnalysisSchedule,
				},
			},
			MetricExpr: metricConfig.Expr,
		})
	}
	return calcs
}
