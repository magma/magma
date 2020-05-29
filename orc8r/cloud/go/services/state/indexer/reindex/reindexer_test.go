/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

// NOTE: to run these tests outside the testing environment, e.g. from IntelliJ,
// ensure postgres_test and maria_test containers are running, and use the
// following environment variables to point to the relevant DB endpoints:
//	- TEST_DATABASE_HOST=localhost
//	- TEST_DATABASE_PORT_POSTGRES=5433
//	- TEST_DATABASE_PORT_MARIA=3307

// reindex_test.go also contains the consts and vars shared by reindex testing code.

package reindex_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/serde"
	configurator_test_init "magma/orc8r/cloud/go/services/configurator/test_init"
	configurator_test "magma/orc8r/cloud/go/services/configurator/test_utils"
	device_test_init "magma/orc8r/cloud/go/services/device/test_init"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	state_test_init "magma/orc8r/cloud/go/services/state/test_init"
	state_test "magma/orc8r/cloud/go/services/state/test_utils"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	assert "github.com/stretchr/testify/require"
)

const (
	queueTableName   = "reindex_job_queue"
	versionTableName = "indexer_versions"

	twoAttempts = 2

	defaultJobTimeout  = 5 * time.Minute // copied from queue_sql.go
	defaultTestTimeout = 5 * time.Second

	// Cause 3 batches per network
	numStatesToReindexPerCall = 100 // copied from reindex.go
	numBatches                = numNetworks * 3
	numNetworks               = 3
	statesPerNetwork          = 2*numStatesToReindexPerCall + 1

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

	matchAll  = []indexer.Subscription{{Type: orc8r.DirectoryRecordType, KeyMatcher: indexer.MatchAll}}
	matchOne  = []indexer.Subscription{{Type: orc8r.DirectoryRecordType, KeyMatcher: indexer.NewMatchExact("imsi0")}}
	matchNone = []indexer.Subscription{{Type: orc8r.DirectoryRecordType, KeyMatcher: indexer.NewMatchExact("0xdeadbeef")}}
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestRun(t *testing.T) {
	dbName := "state___reindex_test___run"

	// Writes to channel after completing a job
	ch := make(chan interface{})
	reindex.TestHookReindexComplete = func() { ch <- nil }
	defer func() { reindex.TestHookReindexComplete = func() {} }()

	clock.SkipSleeps(t)
	defer clock.ResumeSleeps(t)

	r, q := initReindexTest(t, dbName)
	ctx, cancel := context.WithCancel(context.Background())
	go r.Run(ctx)
	defer cancel()

	// Single indexer
	// Populate
	idx0 := getIndexer(id0, zero, version0, true)
	idx0.On("GetSubscriptions").Return(matchAll).Once()
	registerAndPopulate(t, q, idx0)
	// Check
	recvCh(t, ch)
	idx0.AssertExpectations(t)
	assertComplete(t, q, id0)

	// Bump existing indexer version
	// Populate
	idx0a := getIndexerNoIndex(id0, version0, version0a, false)
	idx0a.On("GetSubscriptions").Return(matchOne).Once()
	idx0a.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(numNetworks)
	registerAndPopulate(t, q, idx0a)
	// Check
	recvCh(t, ch)
	idx0a.AssertExpectations(t)
	assertComplete(t, q, id0)

	// Indexer returns err => reindex jobs fail
	// Populate
	// Fail1 at PrepareReindex
	fail1 := getBasicIndexer(id1, version1)
	fail1.On("GetSubscriptions").Return(matchAll).Once()
	fail1.On("PrepareReindex", zero, version1, true).Return(someErr1).Once()
	// Fail2 at first Reindex
	fail2 := getBasicIndexer(id2, version2)
	fail2.On("GetSubscriptions").Return(matchAll).Once()
	fail2.On("PrepareReindex", zero, version2, true).Return(nil).Once()
	fail2.On("Index", mock.Anything, mock.Anything).Return(nil, someErr2).Once()
	// Fail3 at CompleteReindex
	fail3 := getBasicIndexer(id3, version3)
	fail3.On("GetSubscriptions").Return(matchAll).Once()
	fail3.On("PrepareReindex", zero, version3, true).Return(nil).Once()
	fail3.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(numBatches)
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

func TestRunUnsafe(t *testing.T) {
	dbName := "state___reindex_test___run_unsafe"
	r, q := initReindexTest(t, dbName)
	ctx := context.Background()

	// New indexer => reindex
	idx0 := getIndexer(id0, zero, version0, true)
	idx0.On("GetSubscriptions").Return(matchAll).Once()
	register(t, idx0)

	updates := func(m string) {
		assert.Contains(t, m, id0)
	}
	err := r.RunUnsafe(ctx, id0, updates) // this run gets a ctx+updates to ensure it doesn't break
	assert.NoError(t, err)
	assertVersions(t, q, id0, version0, version0)

	// Old version => reindex
	idx0a := getIndexer(id0, version0, version0a, false)
	idx0a.On("GetSubscriptions").Return(matchAll).Once()
	register(t, idx0a)
	err = r.RunUnsafe(nil, id0, nil)
	assert.NoError(t, err)
	assertVersions(t, q, id0, version0a, version0a)

	// Up-to-date version => no reindex
	idx0b := getIndexer(id0, version0a, version0a, false)
	idx0b.On("GetSubscriptions").Return(matchAll).Once()
	register(t, idx0b)
	err = r.RunUnsafe(nil, id0, nil)
	assert.NoError(t, err)
	assertVersions(t, q, id0, version0a, version0a)

	// Reindex: filter all but one state
	idx1 := getIndexerNoIndex(id1, zero, version1, true)
	idx1.On("GetSubscriptions").Return(matchOne).Once()
	idx1.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(numNetworks)
	register(t, idx1)
	err = r.RunUnsafe(nil, id1, nil)
	assert.NoError(t, err)
	assertVersions(t, q, id1, version1, version1)

	// Reindex: filter all
	idx2 := getIndexerNoIndex(id2, zero, version2, true)
	idx2.On("GetSubscriptions").Return(matchNone).Once()
	register(t, idx2)
	err = r.RunUnsafe(nil, id2, nil)
	assert.NoError(t, err)
	assertVersions(t, q, id2, version2, version2)

	// Two new indexers => reindex
	idx3 := getIndexer(id3, zero, version3, true)
	idx4 := getIndexer(id4, zero, version4, true)
	idx3.On("GetSubscriptions").Return(matchAll).Once()
	idx4.On("GetSubscriptions").Return(matchAll).Once()
	register(t, idx3, idx4)
	err = r.RunUnsafe(nil, "", nil)
	assert.NoError(t, err)
	assertVersions(t, q, id3, version3, version3)
	assertVersions(t, q, id4, version4, version4)
}

func initReindexTest(t *testing.T, dbName string) (reindex.Reindexer, reindex.JobQueue) {
	assert.NoError(t, plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{}))
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
	for _, nid := range []string{nid0, nid1, nid2} {
		var records []*directoryd.DirectoryRecord
		var deviceIDs []string
		for i := 0; i < statesPerNetwork; i++ {
			hwid := fmt.Sprintf("hwid%d", i)
			imsi := fmt.Sprintf("imsi%d", i)
			records = append(records, &directoryd.DirectoryRecord{LocationHistory: []string{hwid}})
			deviceIDs = append(deviceIDs, imsi)
		}
		reportStates(t, ctxByNetwork[nid], deviceIDs, records)
	}

	return reindexer, q
}

func reportStates(t *testing.T, ctx context.Context, deviceIDs []string, records []*directoryd.DirectoryRecord) {
	client, err := state.GetStateClient()
	assert.NoError(t, err)

	var states []*protos.State
	for i, st := range records {
		serialized, err := serde.Serialize(state.SerdeDomain, orc8r.DirectoryRecordType, st)
		assert.NoError(t, err)
		pState := &protos.State{
			Type:     orc8r.DirectoryRecordType,
			DeviceID: deviceIDs[i],
			Value:    serialized,
		}
		states = append(states, pState)
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
	idx.On("Index", mock.Anything, mock.Anything).Return(nil, nil).Times(numBatches)
	idx.On("CompleteReindex", from, to).Return(nil).Once()
	return idx
}

func register(t *testing.T, idx ...indexer.Indexer) {
	indexer.DeregisterAllForTest(t)
	err := indexer.RegisterAll(idx...)
	assert.NoError(t, err)
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
	st, err := reindex.GetStatus(q, id)
	assert.NoError(t, err)
	assert.Equal(t, reindex.StatusComplete, st)
	e, err := reindex.GetError(q, id)
	assert.NoError(t, err)
	assert.Empty(t, e)
}

func assertErrored(t *testing.T, q reindex.JobQueue, indexerID string, sentinel reindex.Error, rootErr error) {
	st, err := reindex.GetStatus(q, indexerID)
	assert.NoError(t, err)
	assert.Equal(t, reindex.StatusAvailable, st)
	e, err := reindex.GetError(q, indexerID)
	assert.NoError(t, err)
	// Job err contains relevant info
	assert.Contains(t, e, indexerID)
	assert.Contains(t, e, sentinel)
	assert.Contains(t, e, rootErr.Error())
}

func assertVersions(t *testing.T, queue reindex.JobQueue, indexerID string, actual, desired indexer.Version) {
	v, err := reindex.GetIndexerVersion(queue, indexerID)
	assert.NoError(t, err)
	assert.Equal(t, actual, v.Actual)
	assert.Equal(t, desired, v.Desired)
}
