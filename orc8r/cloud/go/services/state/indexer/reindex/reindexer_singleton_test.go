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
	"flag"
	"fmt"
	"testing"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

const (
	queueTableName   = "reindex_job_queue"
	versionTableName = "indexer_versions"

	twoAttempts = 2

	defaultJobTimeout  = 5 * time.Minute // copied from queue_sql.go
	defaultTestTimeout = 5 * time.Second
	shortTestTimeout   = 1 * time.Second

	// Cause 3 batches per network
	// 200 directory records + 1 gw status => ceil(201 / 100) = 3 batches per network
	nStatesToReindexPerCall    = 100 // copied from reindex.go
	directoryRecordsPerNetwork = 2 * nStatesToReindexPerCall
	nNetworks                  = 3
	// 3 networks and 3 batches per network = 3 * 3 = 9
	nBatches = 9
	// 4 networks and 3 batches per network = 4 * 3 = 12
	newNBatches = 12

	nid0 = "some_networkid_0"
	nid1 = "some_networkid_1"
	nid2 = "some_networkid_2"
	nid3 = "some_networkid_3"

	hwid0 = "some_hwid_0"
	hwid1 = "some_hwid_1"
	hwid2 = "some_hwid_2"
	hwid3 = "some_hwid_3"

	id0 = "some_indexerid_0"
	id1 = "some_indexerid_1"
	id2 = "some_indexerid_2"
	id3 = "some_indexerid_3"
	id4 = "some_indexerid_4"
	id5 = "some_indexerid_5"

	zero      indexer.Version = 0
	version0  indexer.Version = 10
	version0a indexer.Version = 100
	version1  indexer.Version = 20
	version1a indexer.Version = 200
	version2  indexer.Version = 30
	version2a indexer.Version = 300
	version3  indexer.Version = 40
	version3a indexer.Version = 400
	version4  indexer.Version = 50
	version4a indexer.Version = 500
	version5  indexer.Version = 60
)

var (
	someErr  = errors.New("some_error")
	someErr1 = errors.New("some_error_1")
	someErr2 = errors.New("some_error_2")
	someErr3 = errors.New("some_error_3")

	allTypes    = []string{orc8r.DirectoryRecordType, orc8r.GatewayStateType}
	gwStateType = []string{orc8r.GatewayStateType}
	noTypes     []string
)

func init() {
	// TODO(hcgatewood) after resolving racy CI issue, revert most changes from #6329
	_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestSingletonRun(t *testing.T) {
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

	reindexer := state_test_init.StartTestServiceInternal(t)
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

func reportDirectoryRecord(t *testing.T, ctx context.Context, deviceIDs []string, records []*directoryd_types.DirectoryRecord) {
	client, err := state.GetStateClient()
	assert.NoError(t, err)

	var states []*protos.State
	for i, st := range records {
		serialized, err := serde.Serialize(st, orc8r.DirectoryRecordType, serdes.State)
		assert.NoError(t, err)
		pState := &protos.State{Type: orc8r.DirectoryRecordType, DeviceID: deviceIDs[i], Value: serialized}
		states = append(states, pState)
	}
	_, err = client.ReportStates(ctx, &protos.ReportStatesRequest{States: states})
	assert.NoError(t, err)
}

func reportGatewayStatus(t *testing.T, ctx context.Context, gwStatus *models.GatewayStatus) {
	client, err := state.GetStateClient()
	assert.NoError(t, err)

	serialized, err := serde.Serialize(gwStatus, orc8r.GatewayStateType, serdes.State)
	assert.NoError(t, err)
	states := []*protos.State{
		{
			Type:     orc8r.GatewayStateType,
			DeviceID: hwid0,
			Value:    serialized,
		},
	}
	_, err = client.ReportStates(ctx, &protos.ReportStatesRequest{States: states})
	assert.NoError(t, err)
}

func getBasicIndexer(id string, v indexer.Version) *mocks.Indexer {
	idx := &mocks.Indexer{}
	idx.On("GetID").Return(id)
	idx.On("GetVersion").Return(v)
	return idx
}

func getIndexerNoIndex(id string, from, to indexer.Version, isFirstReindex bool) *mocks.Indexer {
	idx := &mocks.Indexer{}
	idx.On("GetID").Return(id)
	idx.On("GetVersion").Return(to)
	idx.On("PrepareReindex", from, to, isFirstReindex).Return(nil).Once()
	idx.On("CompleteReindex", from, to).Return(nil).Once()
	return idx
}

func getIndexer(id string, from, to indexer.Version, isFirstReindex bool) *mocks.Indexer {
	idx := &mocks.Indexer{}
	idx.On("GetID").Return(id)
	idx.On("GetVersion").Return(to)
	idx.On("PrepareReindex", from, to, isFirstReindex).Return(nil).Once()
	idx.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(nBatches)
	idx.On("CompleteReindex", from, to).Return(nil).Once()
	return idx
}

func register(t *testing.T, indexers ...indexer.Indexer) {
	indexer.DeregisterAllForTest(t)
	for _, x := range indexers {
		state_test_init.StartNewTestIndexer(t, x)
	}
}

func recvCh(t *testing.T, ch chan interface{}) {
	select {
	case <-ch:
		return
	case <-time.After(defaultTestTimeout):
		t.Fatal("receive on hook channel timed out")
	}
}

func recvNoCh(t *testing.T, ch chan interface{}) {
	select {
	case <-ch:
		t.Fatal("should not receive anything from hook channel")
	case <-time.After(shortTestTimeout):
		return
	}
}
