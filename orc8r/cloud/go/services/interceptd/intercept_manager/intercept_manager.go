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

package interceptmgr

import (
	"magma/orc8r/cloud/go/services/interceptd"
	"magma/orc8r/cloud/go/services/interceptd/collector"
)

// InterceptManager provides the main functionality for the interceptd
// service. It collects ES events, encode IRIs records and export
// them to LIMS.
type InterceptManager struct {
	Collector               *collector.EventsCollector
	MaxEventsCollectRetries uint32
	MaxRecordsExportRetries uint32
}

// NewInterceptManager creates and returns a new interceptManager
func NewInterceptManager(
	collector *collector.EventsCollector,
	config interceptd.Config,
) *InterceptManager {
	return &InterceptManager{
		Collector:               collector,
		MaxEventsCollectRetries: config.MaxEventsCollectRetries,
		MaxRecordsExportRetries: config.MaxRecordsExportRetries,
	}
}

// CollectAndProcessEvents collects events, process them and
// sends IRIs reports to an external LIMS
func (im *InterceptManager) CollectAndProcessEvents() error {
	// TBD
	// Support retrieving the intercept config from X1 interface
	// Encode IRIs records
	// Export IRIs records to LIMS
	return nil
}
