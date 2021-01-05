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
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
)

const (
	// MaxAnalysisScheduleInHours is the maximum period at which the analyzer
	// job can periodically scheduled to run
	MaxAnalysisScheduleInHours = 12
)

type Config struct {
	// Analysis schedule is time specified in hours at which the analyzer
	// will be periodically run, the typical schedule to run the analyzer
	// will be 3, 6, 9, 12 hours. The max will be 12 hours and any number
	// beyond this will default to max
	AnalysisSchedule uint `yaml:"analysisSchedule"`

	// Export Metrics is a boolean flag which controls the export of metrics
	// by the analytics service
	ExportMetrics   bool   `yaml:"exportMetrics"`
	MetricsPrefix   string `yaml:"metricsPrefix"`
	AppSecret       string `yaml:"appSecret"`
	AppID           string `yaml:"appID"`
	MetricExportURL string `yaml:"metricExportURL"`
	CategoryName    string `yaml:"categoryName"`
}

func GetServiceConfig() Config {
	var serviceConfig Config
	_, _, err := config.GetStructuredServiceConfig(orc8r.ModuleName, ServiceName, &serviceConfig)
	if err != nil {
		glog.Fatalf("Failed parsing the analytics config file: %v ", err)
	}

	if serviceConfig.AnalysisSchedule == 0 ||
		serviceConfig.AnalysisSchedule > MaxAnalysisScheduleInHours {
		serviceConfig.AnalysisSchedule = MaxAnalysisScheduleInHours
	}
	return serviceConfig
}
