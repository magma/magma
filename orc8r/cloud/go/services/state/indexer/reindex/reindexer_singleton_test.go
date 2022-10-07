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

package reindex_test

import (
	"fmt"
	"testing"

	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	directoryd_types "magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	state_test "magma/orc8r/cloud/go/services/state/test_utils"
)

// initSingletonReindexTest reports enough directory records to cause 3 batches per network
// (with the +1 gateway status per network). It creates 3 networks,
// so numBatches following this method will be 3 * 3 = 9
func initSingletonReindexTest(t *testing.T) reindex.Reindexer {
	indexer.DeregisterAllForTest(t)

	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)

	reindexer := state_test_init.StartTestSingletonServiceInternal(t)

	// Report enough directory records to cause 3 batches per network (with the +1 gateway status per network)
	reportAdditionalState(t, nid0, hwid0, 0)
	reportAdditionalState(t, nid1, hwid1, 1)
	reportAdditionalState(t, nid2, hwid2, 2)
	return reindexer
}

// reportMoreState reports enough directory records to cause 3 batches per network
// (with the +1 gateway status per network). It adds an extra network from 3 -> 4,
// so numBatches following this method will be 3 * 4 = 12
func reportAdditionalState(t *testing.T, nid string, hwid string, networkNumber int) {
	configurator_test.RegisterNetwork(t, nid, fmt.Sprintf("Network %v for reindex test", networkNumber))
	configurator_test.RegisterGateway(t, nid, hwid, &models.GatewayDevice{HardwareID: hwid})

	ctx := state_test.GetContextWithCertificate(t, hwid)

	var records []*directoryd_types.DirectoryRecord
	var deviceIDs []string
	for i := 0; i < directoryRecordsPerNetwork; i++ {
		hwidStr := fmt.Sprintf("hwid%d", i)
		imsiStr := fmt.Sprintf("imsiStr%d", i)
		records = append(records, &directoryd_types.DirectoryRecord{LocationHistory: []string{hwidStr}})
		deviceIDs = append(deviceIDs, imsiStr)
	}
	reportDirectoryRecord(t, ctx, deviceIDs, records)

	// Report one gateway status per network
	gwStatus := &models.GatewayStatus{Meta: map[string]string{"foo": "bar"}}
	reportGatewayStatus(t, ctx, gwStatus)
}
