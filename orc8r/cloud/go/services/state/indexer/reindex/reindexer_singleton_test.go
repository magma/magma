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
	"context"
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	directoryd_types "magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	state_test "magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/mock"
)

func initSingletonReindexTest(t *testing.T, dbName string) (reindex.Reindexer) {
	indexer.DeregisterAllForTest(t)

	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	configurator_test.RegisterNetwork(t, nid0, "Network 0 for reindex test")
	configurator_test.RegisterNetwork(t, nid1, "Network 1 for reindex test")
	configurator_test.RegisterNetwork(t, nid2, "Network 2 for reindex test")
	configurator_test.RegisterGateway(t, nid0, hwid0, &models.GatewayDevice{HardwareID: hwid0})
	configurator_test.RegisterGateway(t, nid1, hwid1, &models.GatewayDevice{HardwareID: hwid1})
	configurator_test.RegisterGateway(t, nid2, hwid2, &models.GatewayDevice{HardwareID: hwid2})

	reindexer := state_test_init.StartTestSingletonServiceInternal(t, dbName, sqorc.PostgresDriver)
	ctxByNetwork := map[string]context.Context{
		nid0: state_test.GetContextWithCertificate(t, hwid0),
		nid1: state_test.GetContextWithCertificate(t, hwid1),
		nid2: state_test.GetContextWithCertificate(t, hwid2),
	}

	// Report enough directory records to cause 3 batches per network (with the +1 gateway status per network)
	for _, nid := range []string{nid0, nid1, nid2} {
		var records []*directoryd_types.DirectoryRecord
		var deviceIDs []string
		for i := 0; i < directoryRecordsPerNetwork; i++ {
			hwid := fmt.Sprintf("hwid%d", i)
			imsi := fmt.Sprintf("imsi%d", i)
			records = append(records, &directoryd_types.DirectoryRecord{LocationHistory: []string{hwid}})
			deviceIDs = append(deviceIDs, imsi)
		}
		reportDirectoryRecord(t, ctxByNetwork[nid], deviceIDs, records)
	}

	// Report one gateway status per network
	gwStatus := &models.GatewayStatus{Meta: map[string]string{"foo": "bar"}}
	for _, nid := range []string{nid0, nid1, nid2} {
		reportGatewayStatus(t, ctxByNetwork[nid], gwStatus)
	}

	return reindexer
}


func TestSingletonRun(t *testing.T) {
	dbName := "state___reindex_test___run"

	// Make nullipotent calls to handle code coverage indeterminacy
	reindex.TestHookReindexSuccess()
	reindex.TestHookReindexDone()

	// Writes to channel after completing a job
	ch := make(chan interface{})
	reindex.TestHookReindexSuccess = func() { ch <- nil }
	defer func() { reindex.TestHookReindexSuccess = func() {} }()

	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	r := initSingletonReindexTest(t, dbName)
	ctx, cancel := context.WithCancel(context.Background())
	go r.Run(ctx)
	defer cancel()

	// Single indexer
	// Populate
	idx0 := getIndexer(id0, zero, version0, true)
	idx0.On("GetTypes").Return(allTypes).Once()
	registerAndPopulate(t, q, idx0)
	// Check
	recvCh(t, ch)
	idx0.AssertExpectations(t)
	assertComplete(t, q, id0)

	// Bump existing indexer version
	// Populate
	idx0a := getIndexerNoIndex(id0, version0, version0a, false)
	idx0a.On("GetTypes").Return(gwStateType).Once()
	idx0a.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(nNetworks)
	registerAndPopulate(t, q, idx0a)
	// Check
	recvCh(t, ch)
	idx0a.AssertExpectations(t)
	assertComplete(t, q, id0)

	// Indexer returns err => reindex jobs fail
	// Populate
	// Fail1 at PrepareReindex
	fail1 := getBasicIndexer(id1, version1)
	fail1.On("GetTypes").Return(allTypes).Once()
	fail1.On("PrepareReindex", zero, version1, true).Return(someErr1).Once()
	// Fail2 at first Reindex
	fail2 := getBasicIndexer(id2, version2)
	fail2.On("GetTypes").Return(allTypes).Once()
	fail2.On("PrepareReindex", zero, version2, true).Return(nil).Once()
	fail2.On("Index", mock.Anything, mock.Anything).Return(nil, someErr2).Once()
	// Fail3 at CompleteReindex
	fail3 := getBasicIndexer(id3, version3)
	fail3.On("GetTypes").Return(allTypes).Once()
	fail3.On("PrepareReindex", zero, version3, true).Return(nil).Once()
	fail3.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(nBatches)
	fail3.On("CompleteReindex", zero, version3).Return(someErr3).Once()
	registerAndPopulate(t, q, fail1, fail2, fail3)
	// Check
	recvCh(t, ch)
	recvCh(t, ch)
	recvCh(t, ch)
	fail1.AssertExpectations(t)
	fail2.AssertExpectations(t)
	fail3.AssertExpectations(t)
	assertErrored(t, q, id1, reindex.ErrPrepare, someErr1)
	assertErrored(t, q, id2, reindex.ErrReindex, someErr2)
	assertErrored(t, q, id3, reindex.ErrComplete, someErr3)
}



