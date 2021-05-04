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
	"flag"
	"fmt"
	"testing"
	"time"

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

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	queueTableName   = "reindex_job_queue"
	versionTableName = "indexer_versions"

	twoAttempts = 2

	defaultJobTimeout  = 5 * time.Minute // copied from queue_sql.go
	defaultTestTimeout = 5 * time.Second

	// Cause 3 batches per network
	// 200 directory records + 1 gw status => ceil(201 / 100) = 3 batches per network
	nStatesToReindexPerCall    = 100 // copied from reindex.go
	directoryRecordsPerNetwork = 2 * nStatesToReindexPerCall
	nNetworks                  = 3
	nBatches                   = 9

	nid0 = "some_networkid_0"
	nid1 = "some_networkid_1"
	nid2 = "some_networkid_2"

	hwid0 = "some_hwid_0"
	hwid1 = "some_hwid_1"
	hwid2 = "some_hwid_2"

	id0 = "some_indexerid_0"
	id1 = "some_indexerid_1"
	id2 = "some_indexerid_2"
	id3 = "some_indexerid_3"
	id4 = "some_indexerid_4"

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

func TestRun(t *testing.T) {
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

	r, q := initReindexTest(t, dbName)
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
	go r.Run(ctx)
	defer cancel()

	// Job exists but indexer's broken
	// Populate
	broken := getBasicIndexer(id0, version0)
	broken.On("GetTypes").Return(allTypes).Once()
	broken.On("PrepareReindex", zero, version0, true).Return(nil).Once()
	broken.On("Index", mock.Anything, mock.Anything).Return(nil, someErr).Once()
	registerAndPopulate(t, q, broken)
	// Check
	recvCh(t, ch)
	recvCh(t, ch) // twice to go through full loop at least once with indexer available
	broken.AssertExpectations(t)
	assertErrored(t, q, id0, reindex.ErrReindex, someErr)
}

func TestRunMissingIndexer(t *testing.T) {
	dbName := "state___reindex_test___run_missing_indexer"

	// Writes to channel after completing a job
	ch := make(chan interface{})
	reindex.TestHookReindexDone = func() { ch <- nil }
	defer func() { reindex.TestHookReindexDone = func() {} }()

	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	// Job exists but indexer doesn't exist
	r, q := initReindexTest(t, dbName)
	ctx, cancel := context.WithCancel(context.Background())
	missing := getIndexer(id0, zero, version0, true)
	missing.On("GetTypes").Return(allTypes).Once()
	registerAndPopulate(t, q, missing)
	indexer.DeregisterAllForTest(t)
	go r.Run(ctx)
	defer cancel()

	// Check
	recvCh(t, ch)
	recvCh(t, ch) // twice to go through full loop at least once with indexer available
	assertErrored(t, q, id0, "", errors.New("indexer not found"))
}

func TestRunUnsafe(t *testing.T) {
	dbName := "state___reindex_test___run_unsafe"
	r, q := initReindexTest(t, dbName)
	ctx := context.Background()

	// New indexer => reindex
	idx0 := getIndexer(id0, zero, version0, true)
	idx0.On("GetTypes").Return(allTypes).Once()
	register(t, idx0)

	updates := func(m string) {
		assert.Contains(t, m, id0)
	}
	err := r.RunUnsafe(ctx, id0, updates) // this run gets updates to ensure it doesn't break
	assert.NoError(t, err)
	assertVersions(t, q, id0, version0, version0)

	// Old version => reindex
	idx0a := getIndexer(id0, version0, version0a, false)
	idx0a.On("GetTypes").Return(allTypes).Once()
	register(t, idx0a)
	err = r.RunUnsafe(ctx, id0, nil)
	assert.NoError(t, err)
	assertVersions(t, q, id0, version0a, version0a)

	// Up-to-date version => no reindex
	idx0b := getIndexer(id0, version0a, version0a, false)
	idx0b.On("GetTypes").Return(allTypes).Once()
	register(t, idx0b)
	err = r.RunUnsafe(ctx, id0, nil)
	assert.NoError(t, err)
	assertVersions(t, q, id0, version0a, version0a)

	// Reindex: filter all but one state
	idx1 := getIndexerNoIndex(id1, zero, version1, true)
	idx1.On("GetTypes").Return(gwStateType).Once()
	idx1.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(nNetworks)
	register(t, idx1)
	err = r.RunUnsafe(ctx, id1, nil)
	assert.NoError(t, err)
	assertVersions(t, q, id1, version1, version1)

	// Reindex: filter all
	idx2 := getIndexerNoIndex(id2, zero, version2, true)
	idx2.On("GetTypes").Return(noTypes).Once()
	register(t, idx2)
	err = r.RunUnsafe(ctx, id2, nil)
	assert.NoError(t, err)
	assertVersions(t, q, id2, version2, version2)

	// Two new indexers => reindex
	idx3 := getIndexer(id3, zero, version3, true)
	idx4 := getIndexer(id4, zero, version4, true)
	idx3.On("GetTypes").Return(allTypes).Once()
	idx4.On("GetTypes").Return(allTypes).Once()
	register(t, idx3, idx4)
	err = r.RunUnsafe(ctx, "", nil)
	assert.NoError(t, err)
	assertVersions(t, q, id3, version3, version3)
	assertVersions(t, q, id4, version4, version4)
}

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
