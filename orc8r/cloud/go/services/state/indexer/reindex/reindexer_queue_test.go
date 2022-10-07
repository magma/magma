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

// NOTE: to run these tests outside the testing environment, e.g. from IntelliJ,
// ensure postgres_test container is running, and use the following environment
// variables to point to the relevant DB endpoints:
//	- TEST_DATABASE_HOST=localhost
//	- TEST_DATABASE_PORT_POSTGRES=5433

package reindex_test

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/serdes"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	directoryd_types "magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	state_test "magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"
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
	_ = flag.Set("logtostderr", "true") // uncomment to view logs during test
}

func TestRunBrokenIndexer00(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer01(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer02(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer03(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer04(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer05(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer06(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer07(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer08(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer09(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer10(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer11(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer12(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer13(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer14(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer15(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer16(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer17(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer18(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer19(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer20(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer21(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer22(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer23(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer24(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer25(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer26(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer27(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer28(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer29(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer30(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer31(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer32(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer33(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer34(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer35(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer36(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer37(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer38(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer39(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer40(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer41(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer42(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer43(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer44(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer45(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer46(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer47(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer48(t *testing.T) { TestRunBrokenIndexer(t) }
func TestRunBrokenIndexer49(t *testing.T) { TestRunBrokenIndexer(t) }

func TestRunBrokenIndexer(t *testing.T) {
	dbName := "state___reindex_test___run_broken_indexer"

	// Writes to channel after completing a job
	ch := make(chan interface{})
	reindex.TestHookReindexDone = func() { ch <- nil }
	defer func() { reindex.TestHookReindexDone = func() {} }()

	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	r, q := initReindexTest(t, dbName)
	ctx, cancel := context.WithCancel(context.Background())
	// Job exists but indexer's broken
	// Populate
	go r.Run(ctx)
	defer cancel()

	broken := getBasicIndexer(id0, version0)
	broken.On("GetTypes").Return(allTypes).Once()
	broken.On("PrepareReindex", zero, version0, true).Return(nil).Once()
	broken.On("Index", mock.Anything, mock.Anything).Return(nil, someErr).Once()
	registerAndPopulate(t, q, broken)
	// Check
	recvCh(t, ch)
	broken.AssertExpectations(t)
	assertErrored(t, q, id0, reindex.ErrReindex, someErr)
}

// initReindexTest reports enough directory records to cause 3 batches per network
// (with the +1 gateway status per network). It creates 3 networks,
// so numBatches following this method will be 3 * 3 = 9
func initReindexTest(t *testing.T, dbName string) (reindex.Reindexer, reindex.JobQueue) {
	indexer.DeregisterAllForTest(t)

	configurator_test_init.StartTestService(t)
	device_test_init.StartTestService(t)
	configurator_test.RegisterNetwork(t, nid0, "Network 0 for reindex test")
	configurator_test.RegisterNetwork(t, nid1, "Network 1 for reindex test")
	configurator_test.RegisterNetwork(t, nid2, "Network 2 for reindex test")
	configurator_test.RegisterGateway(t, nid0, hwid0, &models.GatewayDevice{HardwareID: hwid0})
	configurator_test.RegisterGateway(t, nid1, hwid1, &models.GatewayDevice{HardwareID: hwid1})
	configurator_test.RegisterGateway(t, nid2, hwid2, &models.GatewayDevice{HardwareID: hwid2})

	reindexer, q := state_test_init.StartTestServiceInternal(t, dbName, sqorc.PostgresDriver)
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

	return reindexer, q
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

func registerAndPopulate(t *testing.T, q reindex.JobQueue, idx ...indexer.Indexer) {
	register(t, idx...)
	populated, err := q.PopulateJobs()
	assert.True(t, populated)
	assert.NoError(t, err)
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

func assertComplete(t *testing.T, q reindex.JobQueue, id string) {
	e, err := reindex.GetError(q, id)
	assert.NoError(t, err)
	assert.Empty(t, e)
	st, err := reindex.GetStatus(q, id)
	assert.NoError(t, err)
	assert.Equal(t, reindex.StatusComplete, st)
}

func assertErrored(t *testing.T, q reindex.JobQueue, indexerID string, sentinel reindex.Error, rootErr error) {
	st, err := reindex.GetStatus(q, indexerID)
	assert.NoError(t, err)
	assert.Equal(t, reindex.StatusAvailable, st)
	e, err := reindex.GetError(q, indexerID)
	assert.NoError(t, err)
	// Job err contains relevant info
	assert.Contains(t, e, indexerID)
	if sentinel != "" {
		assert.Contains(t, e, sentinel)
	}
	assert.Contains(t, e, rootErr.Error())
}

func assertVersions(t *testing.T, queue reindex.JobQueue, indexerID string, actual, desired indexer.Version) {
	v, err := reindex.GetIndexerVersion(queue, indexerID)
	assert.NoError(t, err)
	assert.Equal(t, actual, v.Actual)
	assert.Equal(t, desired, v.Desired)
}
