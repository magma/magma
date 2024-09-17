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

// Package metrics provides utility functions for service_health
// services to more easily obtain their metrics
package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

func GetInt64(metricName string) (int64, error) {
	families, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		return 0, fmt.Errorf("Error gathering metrics from registry; %s", err)
	}
	for _, family := range families {
		metric := family.Metric
		if family.Name == nil || metricName != *family.Name {
			continue
		}
		if len(metric) == 0 {
			continue
		}
		value := metric[0].GetCounter().GetValue()
		return int64(value), nil
	}
	return 0, fmt.Errorf("Could not find metric with name: %s in default registry", metricName)
}
