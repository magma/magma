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
	"bytes"
	"fmt"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"text/template"

	"github.com/golang/glog"
)

// RawMetricsCalculation params for querying existing prometheus metrics.
type RawMetricsCalculation struct {
	BaseCalculation
	MetricExpr string
}

// Calculate queries for input promql expression and returns the result unchanged.
func (x *RawMetricsCalculation) Calculate(prometheusClient query_api.PrometheusAPI) ([]*protos.CalculationResult, error) {
	expr := x.getPromExpr()
	glog.V(1).Infof("Calculating Raw Metrics for %s expr %s", x.Name, expr)
	vec, err := query_api.QueryPrometheusVector(prometheusClient, expr)
	if err != nil {
		return nil, fmt.Errorf("query error: %s", err)
	}
	results := MakeVectorResults(vec, x.Labels, x.Name)
	return results, nil
}

// getPromExpr gets the template substituted string for the prom expression
func (x *RawMetricsCalculation) getPromExpr() string {
	t, _ := template.New("").Parse(x.MetricExpr)
	data := struct {
		Duration string
	}{
		Duration: fmt.Sprintf("%dh", x.Hours),
	}
	var tpl bytes.Buffer
	t.Execute(&tpl, data)
	return tpl.String()
}
