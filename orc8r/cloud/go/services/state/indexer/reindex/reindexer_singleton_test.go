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

func TestSingletonRun1(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun2(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun3(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun4(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun5(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun6(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun7(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun8(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun9(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun10(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun11(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun12(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun13(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun14(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun15(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun16(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun17(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun18(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun19(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun20(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun21(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun22(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun23(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun24(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun25(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun26(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun27(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun28(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun29(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun30(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun31(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun32(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun33(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun34(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun35(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun36(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun37(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun38(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun39(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun40(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun41(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun42(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun43(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun44(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun45(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun46(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun47(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun48(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun49(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun50(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun51(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun52(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun53(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun54(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun55(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun56(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun57(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun58(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun59(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun60(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun61(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun62(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun63(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun64(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun65(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun66(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun67(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun68(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun69(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun70(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun71(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun72(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun73(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun74(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun75(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun76(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun77(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun78(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun79(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun80(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun81(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun82(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun83(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun84(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun85(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun86(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun87(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun88(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun89(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun90(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun91(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun92(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun93(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun94(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun95(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun96(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun97(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun98(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRun99(t *testing.T) {
	TestSingletonRun(t)
}

func TestSingletonRunInf(t *testing.T) {
	for {
		TestSingletonRun(t)
	}
}

func TestSingletonRun(t *testing.T) {
	// Make nullimpotent calls to handle code coverage indeterminacy
	reindex.TestHookReindexSuccess()
	reindex.TestHookReindexDone()
	// Change to trigger tests
	fmt.Print("Test")

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
	require.Equal(t, reindexSuccessNum, 1)
	require.Equal(t, reindexDoneNum, 1)

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
	require.Equal(t, reindexSuccessNum, 2)
	require.Equal(t, reindexDoneNum, 2)

	// Test that a network/hardware pair that has been added after Run
	// will have its states reindexed as well

	reportMoreState(t)

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
	require.Equal(t, 3, reindexSuccessNum)
	require.Equal(t, 6, reindexDoneNum)
}

// initSingletonReindexTest reports enough directory records to cause 3 batches per network
// (with the +1 gateway status per network). It creates 3 networks,
// so numBatches following this method will be 3 * 3 = 9
func initSingletonReindexTest(t *testing.T) reindex.Reindexer {
	indexer.DeregisterAllForTest(t)

	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)

	configurator_test.RegisterNetwork(t, nid0, "Network 0 for reindex test")
	configurator_test.RegisterNetwork(t, nid1, "Network 1 for reindex test")
	configurator_test.RegisterNetwork(t, nid2, "Network 2 for reindex test")
	configurator_test.RegisterGateway(t, nid0, hwid0, &models.GatewayDevice{HardwareID: hwid0})
	configurator_test.RegisterGateway(t, nid1, hwid1, &models.GatewayDevice{HardwareID: hwid1})
	configurator_test.RegisterGateway(t, nid2, hwid2, &models.GatewayDevice{HardwareID: hwid2})

	reindexer := state_test_init.StartTestSingletonServiceInternal(t)
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

// reportMoreState reports enough directory records to cause 3 batches per network
// (with the +1 gateway status per network). It adds an extra network from 3 -> 4,
// so numBatches following this method will be 3 * 4 = 12
func reportMoreState(t *testing.T) {
	configurator_test.RegisterNetwork(t, nid3, "Network 3 for reindex test")
	configurator_test.RegisterGateway(t, nid3, hwid3, &models.GatewayDevice{HardwareID: hwid3})

	ctxByNetwork := map[string]context.Context{
		nid3: state_test.GetContextWithCertificate(t, hwid3),
	}

	for _, nid := range []string{nid3} {
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
	for _, nid := range []string{nid3} {
		reportGatewayStatus(t, ctxByNetwork[nid], gwStatus)
	}
}
