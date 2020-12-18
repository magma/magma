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

//Config represents the configuration provided to analytics service
type Config struct {
	AnalysisSchedule string `yaml:"analysisSchedule"`
	ExportMetrics    bool   `yaml:"exportMetrics"`
	MetricsPrefix    string `yaml:"metricsPrefix"`
	AppSecret        string `yaml:"appSecret"`
	AppID            string `yaml:"appID"`
	MetricExportURL  string `yaml:"metricExportURL"`
	CategoryName     string `yaml:"categoryName"`
}
