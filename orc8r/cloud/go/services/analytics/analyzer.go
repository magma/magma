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
	"context"
	"fmt"
	"net/http"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/robfig/cron/v3"
)

// Analyzer generic interface to schedule any analysis to run.
type Analyzer interface {
	// Schedule the analyzer to run calculations periodically based on the
	// cron expression format schedule parameter
	Schedule() error

	// Run triggers the analyzer's cronjob to start running. This function
	// blocks.
	Run()
}

// PrometheusAnalyzer accesses prometheus metrics and performs
// queries/aggregations to calculate various metrics
type PrometheusAnalyzer struct {
	Cron             *cron.Cron
	Config           *Config
	PrometheusClient query_api.PrometheusAPI
	Calculations     []calculations.Calculation
	Exporter         Exporter
}

func NewPrometheusAnalyzer(config *Config, prometheusClient v1.API, exporter Exporter) Analyzer {
	cronJob := cron.New()
	return &PrometheusAnalyzer{
		Config:           config,
		Cron:             cronJob,
		PrometheusClient: prometheusClient,
		Exporter:         exporter,
	}
}

// Schedule method takes in a schedule string in cron format and schedules the
// analyze job to be run at that schedule
func (a *PrometheusAnalyzer) Schedule() error {
	glog.V(1).Infof("Analyzer will be run every %d hours", a.Config.AnalysisSchedule)

	a.Cron = cron.New()
	_, err := a.Cron.AddFunc(getAnalysisSchedule(a.Config.AnalysisSchedule), a.Analyze)
	if err != nil {
		glog.Infof("error scheduling the local analytics function %v", err)
		return err
	}

	return nil
}

// Analyze methods runs through collectors and exports their metrics
func (a *PrometheusAnalyzer) Analyze() {
	glog.V(2).Info("Running Analyze")
	collectorClients, err := getRemoteCollectors()
	if err != nil {
		glog.Infof("err %v failed to get remote collectors", err)
		return
	}
	for _, c := range collectorClients {
		collectResp, err := c.Collect(context.Background(), &protos.CollectRequest{})
		if err != nil || collectResp == nil {
			glog.Infof("err %v or empty response when attempting to collect from service", err)
			continue
		}
		for _, res := range collectResp.GetResults() {
			if a.Exporter == nil {
				continue
			}

			err = a.Exporter.Export(res, http.DefaultClient)
			if err != nil {
				glog.Errorf("Error exporting result: %v", err)
			} else {
				glog.V(10).Infof("Exported %s, %s, %f", res.MetricName, res.Labels, res.Value)
			}
		}
	}
}

func (a *PrometheusAnalyzer) Run() {
	a.Cron.Run()
}

func getRemoteCollectors() ([]protos.AnalyticsCollectorClient, error) {
	services, err := registry.FindServices(orc8r.AnalyticsCollectorLabel)
	if err != nil {
		glog.Errorf("Failed finding analytics collectors %v", err)
		return nil, err
	}

	var collectorClientList []protos.AnalyticsCollectorClient
	for _, s := range services {
		conn, err := registry.GetConnection(s)
		if err != nil {
			glog.Errorf("Unable to get a remote connection %s error %v", s, err)
			continue
		}
		collectorClientList = append(collectorClientList, protos.NewAnalyticsCollectorClient(conn))
	}
	return collectorClientList, nil
}

func getAnalysisSchedule(analysisSchedule uint) string {
	return fmt.Sprintf("0 */%d * * *", analysisSchedule)
}
