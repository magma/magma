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
	"sync"
	"testing"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/assert"
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestSQLReindexJobQueue_Integration_PopulateJobs(t *testing.T) {
	dbName := "state___reindex_queue___populate_jobs"
	queue := initSQLTest(t, dbName)

	// Start and register indexer servicers
	mocks.NewMockIndexer(t, id0, version0, nil, nil, nil, nil)
	mocks.NewMockIndexer(t, id1, version1, nil, nil, nil, nil)
	mocks.NewMockIndexer(t, id2, version2, nil, nil, nil, nil)

	ch := make(chan interface{})
	defer close(ch)
	wg := sync.WaitGroup{}

	// tx0 -- will be held up by the test hook, eventually fail to commit
	reindex.TestHookGet = func() {
		ch <- nil
		ch <- nil
	}
	wg.Add(1)
	go func() {
		populated, err := queue.PopulateJobs()
		assert.NoError(t, err)
		assert.False(t, populated)
		wg.Done()
	}()

	select {
	// tx0 has begun
	case <-ch:
	// Prevent test hanging
	case <-time.After(15 * time.Second):
		t.Fatal("PopulateJobs failed to start transaction")
		return
	}

	// tx1 -- will begin after tx0 has begun, but tx1 will move first and commit its update
	reindex.TestHookGet = func() {}
	populated, err := queue.PopulateJobs()
	assert.NoError(t, err)
	assert.True(t, populated)

	<-ch // tx0 can continue, will fail
	wg.Wait()

	// All jobs are available
	statuses, err := reindex.GetStatuses(queue)
	assert.NoError(t, err)
	for _, st := range statuses {
		assert.Equal(t, reindex.StatusAvailable, st)
	}
}

func TestSQLJobQueue_Integration_ClaimAvailableReindexJob(t *testing.T) {
	dbName := "state___reindex_queue___claim_jobs"
	queue := initSQLTest(t, dbName)

	// Start and register indexer servicers
	mocks.NewMockIndexer(t, id0, version0, nil, nil, nil, nil)
	mocks.NewMockIndexer(t, id1, version1, nil, nil, nil, nil)
	mocks.NewMockIndexer(t, id2, version2, nil, nil, nil, nil)

	populated, err := queue.PopulateJobs()
	assert.NoError(t, err)
	assert.True(t, populated)

	// Claim all idxs
	jobX, err := queue.ClaimAvailableJob()
	assertJob(t, jobX, err)
	jobY, err := queue.ClaimAvailableJob()
	assertJob(t, jobY, err)
	jobZ, err := queue.ClaimAvailableJob()
	assertJob(t, jobZ, err)

	// No jobs left
	j, err := queue.ClaimAvailableJob()
	assertNoJob(t, j, err)

	// All jobs are in_progress
	statuses, err := reindex.GetStatuses(queue)
	assert.NoError(t, err)
	for _, st := range statuses {
		assert.Equal(t, reindex.StatusInProgress, st)
	}

	// Extract jobs/indexers to properly keep track by number
	jobs := map[string]*reindex.Job{}
	jobs[jobX.Idx.GetID()] = jobX
	jobs[jobY.Idx.GetID()] = jobY
	jobs[jobZ.Idx.GetID()] = jobZ
	job0, job1, job2 := jobs[id0], jobs[id1], jobs[id2]
	idx0, idx1, idx2 := job0.Idx, job1.Idx, job2.Idx

	// Got correct from/to versions
	assert.Equal(t, zero, job0.From)
	assert.Equal(t, zero, job1.From)
	assert.Equal(t, zero, job2.From)
	assert.Equal(t, version0, job0.To)
	assert.Equal(t, version1, job1.To)
	assert.Equal(t, version2, job2.To)

	// Successfully complete idx0
	err = queue.CompleteJob(job0, nil)
	assert.NoError(t, err)
	status, err := reindex.GetStatus(queue, job0.Idx.GetID())
	assert.NoError(t, err)
	assert.Equal(t, reindex.StatusComplete, status)

	// Fail to complete idx1 => retry=1, no error saved
	err = queue.CompleteJob(job1, someErr)
	assert.NoError(t, err)
	errVal, err := reindex.GetError(queue, idx1.GetID())
	assert.NoError(t, err)
	assert.Empty(t, errVal)

	// Claim new idx -- should be idx1 again
	job1a, err := queue.ClaimAvailableJob()
	assertJob(t, job1a, err)
	idx1a := job1a.Idx
	assert.Equal(t, idx1.GetID(), idx1a.GetID())
	assert.Equal(t, zero, job1a.From)
	assert.Equal(t, version1, job1a.To)

	// Still no errors saved
	errVals, err := reindex.GetErrors(queue)
	assert.NoError(t, err)
	assert.Empty(t, errVals)

	// Fail to complete idx1 (aka idx1a) again => retry=2, error now saved
	err = queue.CompleteJob(job1a, someErr)
	assert.NoError(t, err)
	status, err = reindex.GetStatus(queue, job1a.Idx.GetID())
	assert.NoError(t, err)
	assert.Equal(t, reindex.StatusAvailable, status)
	errVal, err = reindex.GetError(queue, idx1a.GetID())
	assert.NoError(t, err)
	assert.Equal(t, someErr.Error(), errVal)

	// Can't claim idx1 again -- no idxs left
	j, err = queue.ClaimAvailableJob()
	assertNoJob(t, j, err)

	// Get all errors -- should just be for idx1
	errVals, err = reindex.GetErrors(queue)
	assert.NoError(t, err)
	assert.Contains(t, errVals, idx1.GetID())
	assert.Equal(t, someErr.Error(), errVals[idx1.GetID()])

	// Fail idx2, then claim but allow to time out -- should result in an err
	err = queue.CompleteJob(job2, someErr)
	assert.NoError(t, err)
	errVal, err = reindex.GetError(queue, idx0.GetID())
	assert.NoError(t, err)
	assert.Empty(t, errVal)
	job2a, err := queue.ClaimAvailableJob()
	assertJob(t, job2a, err)
	idx2a := job2a.Idx
	assert.Equal(t, idx2.GetID(), idx2a.GetID())

	errVal, err = reindex.GetError(queue, idx2a.GetID())
	assert.NoError(t, err)
	assert.Empty(t, errVal)
	clock.SetAndFreezeClock(t, time.Now().Add(defaultJobTimeout).Add(time.Minute))
	defer clock.UnfreezeClock(t)
	errVal, err = reindex.GetError(queue, idx2a.GetID())
	assert.NoError(t, err)
	assert.Equal(t, someErr.Error(), errVal)

	// Complete idx2 -- unspecified behavior but gracefully handle a job taking longer than default timeout
	err = queue.CompleteJob(job2a, nil)
	assert.NoError(t, err)
	errVal, err = reindex.GetError(queue, idx0.GetID())
	assert.NoError(t, err)
	assert.Empty(t, errVal)
	status, err = reindex.GetStatus(queue, job2a.Idx.GetID())
	assert.NoError(t, err)
	assert.Equal(t, reindex.StatusComplete, status)
}

// Update indexer version, repopulate should add new job
func TestSQLJobQueue_Integration_RepopulateNewJobs(t *testing.T) {
	dbName := "state___reindex_queue___repopulate_jobs"
	queue := initSQLTest(t, dbName)

	// Start and register indexer servicers
	mocks.NewMockIndexer(t, id0, version0, nil, nil, nil, nil)
	mocks.NewMockIndexer(t, id1, version1, nil, nil, nil, nil)
	mocks.NewMockIndexer(t, id2, version2, nil, nil, nil, nil)

	// Empty to start
	j, err := queue.ClaimAvailableJob()
	assertNoJob(t, j, err)

	// Populate indexers
	populated, err := queue.PopulateJobs()
	assert.NoError(t, err)
	assert.True(t, populated)

	// Claim all idxs
	jobX, err := queue.ClaimAvailableJob()
	assertJob(t, jobX, err)
	jobY, err := queue.ClaimAvailableJob()
	assertJob(t, jobY, err)
	jobZ, err := queue.ClaimAvailableJob()
	assertJob(t, jobZ, err)
	// No jobs left
	j, err = queue.ClaimAvailableJob()
	assertNoJob(t, j, err)
	// Extract jobs/indexers to properly keep track by number
	jobs := map[string]*reindex.Job{}
	jobs[jobX.Idx.GetID()] = jobX
	jobs[jobY.Idx.GetID()] = jobY
	jobs[jobZ.Idx.GetID()] = jobZ
	job0, job1, job2 := jobs[id0], jobs[id1], jobs[id2]

	// Complete all idxs
	// Complete with success idx0, idx2
	err = queue.CompleteJob(job0, nil)
	assert.NoError(t, err)
	err = queue.CompleteJob(job2, nil)
	assert.NoError(t, err)
	// Complete with fail idx1
	err = queue.CompleteJob(job1, someErr)
	assert.NoError(t, err)
	_, err = queue.ClaimAvailableJob()
	assert.NoError(t, err)
	err = queue.CompleteJob(job1, someErr)
	assert.NoError(t, err)
	errVal, err := reindex.GetError(queue, job1.Idx.GetID())
	assert.NoError(t, err)
	assert.Equal(t, someErr.Error(), errVal)
	// No jobs left
	j, err = queue.ClaimAvailableJob()
	assertNoJob(t, j, err)
	assert.Nil(t, j)

	// Update version of indexer 0 -- previously succeeded
	indexer0a, _ := mocks.NewMockIndexer(t, id0, version0a, nil, nil, nil, nil)
	updated, err := queue.PopulateJobs()
	assert.NoError(t, err)
	assert.True(t, updated)
	// Update version of indexer 1 -- previously failed
	indexer1a, _ := mocks.NewMockIndexer(t, id1, version1a, nil, nil, nil, nil)
	updated, err = queue.PopulateJobs()
	assert.NoError(t, err)
	assert.True(t, updated)

	// Claim jobs -- idx0 and idx1 should both be present, across re-populations
	jobZ, err = queue.ClaimAvailableJob()
	assertJob(t, jobZ, err)
	jobY, err = queue.ClaimAvailableJob()
	assertJob(t, jobY, err)
	// No jobs remaining
	j, err = queue.ClaimAvailableJob()
	assertNoJob(t, j, err)

	// Extract jobs/indexers to properly keep track by number
	jobs = map[string]*reindex.Job{}
	jobs[jobZ.Idx.GetID()] = jobZ
	jobs[jobY.Idx.GetID()] = jobY
	job0, job1 = jobs[id0], jobs[id1]
	idx0, idx1 := job0.Idx, job1.Idx

	// Check idx0 version -- previously succeeded
	assert.Equal(t, indexer0a.GetID(), idx0.GetID())
	assert.Equal(t, version0, job0.From)
	assert.Equal(t, version0a, job0.To)
	// Check idx1 version -- previously failed
	assert.Equal(t, indexer1a.GetID(), idx1.GetID())
	assert.Equal(t, zero, job1.From)
	assert.Equal(t, version1a, job1.To)

	// Complete job for indexer 0
	err = queue.CompleteJob(job0, nil)
	assert.NoError(t, err)
	// Complete job for indexer 1
	err = queue.CompleteJob(job1, nil)
	assert.NoError(t, err)

	// No jobs remaining
	j, err = queue.ClaimAvailableJob()
	assertNoJob(t, j, err)

	// All jobs succeeded
	statuses, err := reindex.GetStatuses(queue)
	assert.NoError(t, err)
	for _, st := range statuses {
		assert.Equal(t, reindex.StatusComplete, st)
	}
}

func TestSQLJobQueue_Integration_IndexerVersions(t *testing.T) {
	dbName := "state___reindex_queue___indexer_versions"
	q := initSQLTest(t, dbName)

	// Empty initially
	v, err := q.GetIndexerVersions()
	assert.NoError(t, err)
	assert.Empty(t, v)

	// Write some versions, ensure they stuck
	want := []*indexer.Versions{
		{IndexerID: id0, Actual: zero, Desired: version0},
		{IndexerID: id1, Actual: zero, Desired: version1},
		{IndexerID: id2, Actual: zero, Desired: version2},
	}

	// Start and register indexer servicers
	mocks.NewMockIndexer(t, id0, version0, nil, nil, nil, nil)
	mocks.NewMockIndexer(t, id1, version1, nil, nil, nil, nil)
	mocks.NewMockIndexer(t, id2, version2, nil, nil, nil, nil)

	assert.NoError(t, err)
	got, err := q.GetIndexerVersions()
	assert.NoError(t, err)
	assert.Equal(t, want, got)

	// Update one actual version
	err = q.SetIndexerActualVersion(id2, version2)
	assert.NoError(t, err)
	gotv, err := reindex.GetIndexerVersion(q, id2)
	assert.NoError(t, err)
	assert.Equal(t, version2, gotv.Actual)

	// Bump desired version for same indexer
	mocks.NewMockIndexer(t, id2, version2a, nil, nil, nil, nil)
	assert.NoError(t, err)
	got, err = q.GetIndexerVersions()
	assert.NoError(t, err)
	want = []*indexer.Versions{
		{IndexerID: id0, Actual: zero, Desired: version0},
		{IndexerID: id1, Actual: zero, Desired: version1},
		{IndexerID: id2, Actual: version2, Desired: version2a},
	}
	assert.Equal(t, want, got)
}

func initSQLTest(t *testing.T, dbName string) reindex.JobQueue {
	indexer.DeregisterAllForTest(t)
	db := sqorc.OpenCleanForTest(t, dbName, sqorc.PostgresDriver)

	q := reindex.NewSQLJobQueue(twoAttempts, db, sqorc.GetSqlBuilder())
	err := q.Initialize()
	assert.NoError(t, err)
	return q
}

func assertJob(t *testing.T, job *reindex.Job, err error) {
	assert.NoError(t, err)
	assert.NotNil(t, job)
}

func assertNoJob(t *testing.T, job *reindex.Job, err error) {
	assert.NoError(t, err)
	assert.Nil(t, job)
}
