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

package main

import (
	"flag"
	"magma/orc8r/cloud/go/services/analytics"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
)

const (
	defaultAnalysisSchedule = "*/15 * * * *" // Every 15 minutes
)

//ServiceMain - entrypoint for analytics service
func main() {
	flag.Parse()

	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, analytics.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating Analytics service: %s", err)
	}

	var serviceConfig analytics.Config
	_, _, err = config.GetStructuredServiceConfig(orc8r.ModuleName, analytics.ServiceName, &serviceConfig)
	if err != nil {
		glog.Infof("err %v failed parsing the analytics config file ", err)
		return
	}

	if serviceConfig.AnalysisSchedule == "" {
		serviceConfig.AnalysisSchedule = defaultAnalysisSchedule
	}
	glog.Infof("Analytics Service Config %v", serviceConfig)
	promAPIClient := analytics.GetPrometheusClient()
	exporter := getExporterIfRequired(&serviceConfig)
	analyzer := analytics.NewPrometheusAnalyzer(&serviceConfig, promAPIClient, exporter)
	err = analyzer.Schedule()
	if err != nil {
		glog.Fatalf("Error scheduling analyzer: %s", err)
	}

	go analyzer.Run()

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}

func getExporterIfRequired(config *analytics.Config) analytics.Exporter {
	if config.ExportMetrics {
		return analytics.NewWWWExporter(
			config.MetricsPrefix,
			config.AppID,
			config.AppSecret,
			config.MetricExportURL,
			config.CategoryName,
		)
	}
	return nil
}
