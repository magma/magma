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

package metricsd

import (
	"magma/orc8r/cloud/go/services/metricsd/collection"
	"magma/orc8r/cloud/go/services/metricsd/exporters"
)

// MetricsProfile is a configuration for the metricsd servicer which specifies
// which collectors and exporters it should run.
//
// Collectors intake metrics, metricsd augments each set of metrics,
// and Exporters output metrics.
// Deprecated: collectors are now static, exporters via service mesh
type MetricsProfile struct {
	// Name is a unique name to assign to this profile. This is how you
	// will tell metricsd which profile to run with.
	Name string

	// Collectors defines the set of functionalities for collecting metrics.
	// Many-inputs.
	Collectors []collection.MetricCollector

	// Exporters defines the set of functionalities for exporting metrics.
	// Many-outputs.
	Exporters []exporters.Exporter
}
