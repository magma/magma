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

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"magma/orc8r/cloud/go/clock"
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

func TestSingletonRunSuccess(t *testing.T) {
	// Make nullimpotent calls to handle code coverage indeterminacy
	reindex.TestHookReindexSuccess()
	reindex.TestHookReindexDone()

	// Writes to channel after completing a job
	reindexSuccessNum, reindexDoneNum := 0, 0
	ch := make(chan interface{})

	reindex.TestHookReindexSuccess = func() {
		reindexSuccessNum += 1
	}
	defer func() { reindex.TestHookReindexSuccess = func() {} }()

	reindex.TestHookReindexDone = func() {
		reindexDoneNum += 1
		ch <- nil
	}
	defer func() { reindex.TestHookReindexDone = func() {} }()

	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	r := initSingletonReindexTest(t)

	ctx, cancel := context.WithCancel(context.Background())
	go r.Run(ctx)

	// Single indexer
	idx0 := getIndexer(id0, zero, version0, true)
	idx0.On("GetTypes").Return(allTypes).Once()
	// Register indexers
	register(t, idx0)

	// Check
	recvCh(t, ch)
	recvNoCh(t, ch)

	idx0.AssertExpectations(t)
	require.Equal(t, 1, reindexSuccessNum)
	require.Equal(t, 1, reindexDoneNum)

	// Bump existing indexer version
	idx0a := getIndexerNoIndex(id0, version0, version0a, false)
	idx0a.On("GetTypes").Return(gwStateType).Once()
	idx0a.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(nNetworks)
	// Register indexers
	register(t, idx0a)

	// Check
	recvCh(t, ch)
	recvNoCh(t, ch)

	idx0a.AssertExpectations(t)
	require.Equal(t, 2, reindexSuccessNum)
	require.Equal(t, 2, reindexDoneNum)

	// Test that a network/hardware pair that has been added after Run
	// will have its states reindexed as well

	// reportAdditionalState reports enough directory records to cause 3 batches per network
	// (with the +1 gateway status per network). It adds an extra network from 3 -> 4,
	// so numBatches following this method will be 3 * 4 = 12
	reportAdditionalState(t, nid3, hwid3, 3)

	idx5 := getIndexerNoIndex(id5, zero, version5, true)
	idx5.On("GetTypes").Return(allTypes).Once()
	idx5.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(newNBatches)
	// Register indexers
	register(t, idx5)

	// Check
	recvCh(t, ch)
	recvNoCh(t, ch)

	idx5.AssertExpectations(t)
	require.Equal(t, 3, reindexSuccessNum)
	require.Equal(t, 3, reindexDoneNum)
	cancel()
}

func TestSingletonRunFail(t *testing.T) {
	// Make nullimpotent calls to handle code coverage indeterminacy
	reindex.TestHookReindexSuccess()
	reindex.TestHookReindexDone()

	// Writes to channel after completing a job
	reindexSuccessNum, reindexDoneNum := 0, 0
	ch := make(chan interface{})

	reindex.TestHookReindexSuccess = func() {
		reindexSuccessNum += 1
	}
	defer func() { reindex.TestHookReindexSuccess = func() {} }()

	reindex.TestHookReindexDone = func() {
		if reindexDoneNum < 6 {
			ch <- nil
		}
		reindexDoneNum += 1
	}
	defer func() { reindex.TestHookReindexDone = func() {} }()

	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	r := initSingletonReindexTest(t)
	reportAdditionalState(t, nid3, hwid3, 3)

	ctx, cancel := context.WithCancel(context.Background())
	go r.Run(ctx)

	// Indexer returns err => reindex jobs fail

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
	fail3.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(newNBatches)
	fail3.On("CompleteReindex", zero, version3).Return(someErr3).Once()

	// Register indexers
	register(t, fail1, fail2, fail3)

	// Check
	recvCh(t, ch)
	recvCh(t, ch)
	recvCh(t, ch)
	cancel()
	recvNoCh(t, ch)

	fail1.AssertExpectations(t)
	fail2.AssertExpectations(t)
	fail3.AssertExpectations(t)
	require.Equal(t, 0, reindexSuccessNum)
	require.Equal(t, 0, reindexDoneNum%3)
}

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
