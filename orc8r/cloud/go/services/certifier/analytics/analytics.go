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
	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/lib/go/metrics"

	"github.com/prometheus/client_golang/prometheus"

	certifier_calcs "magma/orc8r/cloud/go/services/certifier/analytics/calculations"
)

// GetAnalyticsCalculations returns all calculations computed by the component
func GetAnalyticsCalculations(serviceConfig *certifier.Config) []calculations.Calculation {
	certLabels := []string{metrics.CertNameLabel, metrics.NetworkLabelName}
	gauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: metrics.CertExpiresInHoursMetric}, certLabels)
	prometheus.MustRegister(gauge)
	calcs := make([]calculations.Calculation, 0)
	calcs = append(calcs, &certifier_calcs.CertLifespanCalculation{
		CertsDirectory: serviceConfig.CertsDirectory,
		Certs:          serviceConfig.Certs,
		BaseCalculation: calculations.BaseCalculation{
			CalculationParams: calculations.CalculationParams{
				AnalyticsConfig:     &serviceConfig.Analytics,
				RegisteredGauge:     gauge,
				ExpectedGaugeLabels: certLabels,
			},
		},
	})
	return calcs
}
