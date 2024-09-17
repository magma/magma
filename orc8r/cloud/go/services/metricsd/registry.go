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

// File registry.go provides a metrics exporter registry by forwarding calls to
// the service registry.

package metricsd

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/lib/go/registry"
)

// GetMetricsExporters returns all registered metrics exporters.
func GetMetricsExporters() ([]exporters.Exporter, error) {
	services, err := registry.FindServices(orc8r.MetricsExporterLabel)
	if err != nil {
		return []exporters.Exporter{}, err
	}
	var exps []exporters.Exporter
	for _, s := range services {
		exps = append(exps, exporters.NewRemoteExporter(s))
	}

	return exps, nil
}
