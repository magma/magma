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
	"fmt"
	"strings"

	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// ParsePrometheusText parses the HTTP response body from common exporters
// that expose prometheus metrics in a common text format.
// Returns metric families keyed by the metric name.
func ParsePrometheusText(prometheusText string) (map[string]*dto.MetricFamily, error) {
	reader := strings.NewReader(prometheusText)

	parser := expfmt.TextParser{}
	metricFamilies, err := parser.TextToMetricFamilies(reader)
	if err != nil {
		return nil, fmt.Errorf("Error parsing metric families from text: %s", err)
	}

	return metricFamilies, nil
}
