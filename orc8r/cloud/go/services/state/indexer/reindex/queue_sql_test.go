/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package reindex_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state/indexer"
	"magma/orc8r/cloud/go/services/state/indexer/mocks"
	"magma/orc8r/cloud/go/services/state/indexer/reindex"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/DATA-DOG/go-sqlmock"
	assert "github.com/stretchr/testify/require"
)

var (
	queueCols         = []string{"indexer_id", "from_version", "to_version", "status", "attempts", "error", "last_status_change"}
	versionCols       = []string{"indexer_id", "version_actual", "version_desired"}
	queueColsJoined   = strings.Join(queueCols, ", ")
	versionColsJoined = strings.Join(versionCols, ", ")

	sqlIndexer0 = mocks.NewMockIndexer(id0, version0, nil, nil, nil, nil)
	sqlIndexer1 = mocks.NewMockIndexer(id1, version1a, nil, nil, nil, nil)
	sqlIndexer2 = mocks.NewMockIndexer(id2, version2a, nil, nil, nil, nil)
	sqlIndexer3 = mocks.NewMockIndexer(id3, version3, nil, nil, nil, nil)
	sqlIndexer4 = mocks.NewMockIndexer(id4, version4, nil, nil, nil, nil)
)

func init() {
	//_ = flag.Set("alsologtostderr", "true") // uncomment to view logs during test
}

func TestSqlJobQueue_PopulateJobs(t *testing.T) {
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(4*time.Hour))
	defer clock.UnfreezeClock(t)
	now := clock.Now()
	recent := now.Add(-time.Minute)
	pastTimeout := recent.Add(-defaultJobTimeout)

	// Five indexers
	//	- id0 tracked, needs upgrade, previous reindex job doesn't exist
	//	- id1 tracked, needs upgrade, previous reindex job completed
	//	- id2 tracked, needs upgrade, previous reindex job unfinished
	//	- id3 tracked, no upgrade
	//	- id4 untracked, needs upgrade
	happyPath := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			// Get tracked versions
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", versionColsJoined, versionTableName)).
				WillReturnRows(
					sqlmock.NewRows(versionCols).
						AddRow(id0, zero, version0).
						AddRow(id1, version1, version1).
						AddRow(id2, zero, version2).
						AddRow(id3, version3, version3),
				)
			// Clear tracked versions
			m.ExpectExec(fmt.Sprintf("DELETE FROM %s", versionTableName)).
				WillReturnResult(sqlmock.NewResult(1, 1))
			// Insert new versions
			m.ExpectExec(fmt.Sprintf("INSERT INTO %s", versionTableName)).
				WithArgs(
					id0, zero, version0,
					id1, version1, version1a,
					id2, zero, version2a,
					id3, version3, version3,
					id4, zero, version4,
				).
				WillReturnResult(sqlmock.NewResult(1, 1))
			// Get existing jobs
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnRows(
					sqlmock.NewRows(queueCols).
						AddRow(id2, zero, version2, reindex.StatusInProgress, reindex.DefaultMaxAttempts, "", pastTimeout.Unix()),
				)
			// Clear existing jobs
			m.ExpectExec(fmt.Sprintf("DELETE FROM %s", queueTableName)).
				WillReturnResult(sqlmock.NewResult(1, 1))
			// Insert new jobs
			m.ExpectExec(fmt.Sprintf("INSERT INTO %s", queueTableName)).
				WithArgs(
					id0, zero, version0, now.Unix(),
					id1, version1, version1a, now.Unix(),
					id2, zero, version2a, now.Unix(),
					id4, zero, version4, now.Unix(),
				).
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},

		registered: []indexer.Indexer{sqlIndexer0, sqlIndexer1, sqlIndexer2, sqlIndexer3, sqlIndexer4},
		run:        func(queue reindex.JobQueue) (interface{}, error) { return queue.PopulateJobs() },
		result:     true,
	}

	noNewJobs := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			// Get tracked versions
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", versionColsJoined, versionTableName)).
				WillReturnRows(
					sqlmock.NewRows(versionCols).
						AddRow(id0, version0, version0).
						AddRow(id1, version1, version1),
				)
			m.ExpectCommit()
		},

		run:    func(queue reindex.JobQueue) (interface{}, error) { return queue.PopulateJobs() },
		result: false,
	}

	runCase(t, happyPath)
	runCase(t, noNewJobs)
}

func TestSqlJobQueue_ClaimAvailableJob(t *testing.T) {
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(4*time.Hour))
	defer clock.UnfreezeClock(t)
	now := clock.Now()

	oneAvailable := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnRows(
					sqlmock.NewRows(queueCols).
						AddRow(id0, version0, version0a, reindex.StatusAvailable, 0, "", 42),
				)
			m.ExpectExec(fmt.Sprintf("UPDATE %s", queueTableName)).
				WithArgs(reindex.StatusInProgress, 1, now.Unix(), id0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},

		registered: []indexer.Indexer{sqlIndexer0, sqlIndexer1, sqlIndexer2},
		run:        func(queue reindex.JobQueue) (interface{}, error) { return queue.ClaimAvailableJob() },
		result:     &reindex.Job{Idx: sqlIndexer0, From: version0, To: version0a},
	}

	selectEmpty := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnRows(sqlmock.NewRows(queueCols))
			m.ExpectRollback()
		},
		registered: []indexer.Indexer{sqlIndexer0, sqlIndexer1, sqlIndexer2},
		run:        func(queue reindex.JobQueue) (interface{}, error) { return queue.ClaimAvailableJob() },
		result:     nil,
	}

	selectErr := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnError(someErr)
			m.ExpectRollback()
		},
		registered: []indexer.Indexer{sqlIndexer0, sqlIndexer1, sqlIndexer2},
		run:        func(queue reindex.JobQueue) (interface{}, error) { return queue.ClaimAvailableJob() },
		err:        someErr,
	}

	updateErr := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnRows(
					sqlmock.NewRows(queueCols).
						AddRow(id0, version0, version0a, reindex.StatusAvailable, 0, "", 42),
				)
			m.ExpectExec(fmt.Sprintf("UPDATE %s", queueTableName)).
				WithArgs(reindex.StatusInProgress, 1, now.Unix(), id0).
				WillReturnError(someErr)
			m.ExpectRollback()
		},
		registered: []indexer.Indexer{sqlIndexer0, sqlIndexer1, sqlIndexer2},
		run:        func(queue reindex.JobQueue) (interface{}, error) { return queue.ClaimAvailableJob() },
		err:        someErr,
	}

	runCase(t, oneAvailable)
	runCase(t, selectEmpty)
	runCase(t, selectErr)
	runCase(t, updateErr)
}

func TestSqlJobQueue_CompleteJob(t *testing.T) {
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(4*time.Hour))
	defer clock.UnfreezeClock(t)
	now := clock.Now()

	completeWithSuccess := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("UPDATE %s", versionTableName)).
				WithArgs(version0a, id0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectExec(fmt.Sprintf("UPDATE %s", queueTableName)).
				WithArgs(reindex.StatusComplete, "", now.Unix(), id0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},

		run: func(queue reindex.JobQueue) (interface{}, error) {
			return nil, queue.CompleteJob(&reindex.Job{Idx: sqlIndexer0, From: version0, To: version0a}, nil)
		},
	}

	completeWithErr := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("UPDATE %s", queueTableName)).
				WithArgs(reindex.StatusAvailable, someErr.Error(), now.Unix(), id0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},

		run: func(queue reindex.JobQueue) (interface{}, error) {
			return nil, queue.CompleteJob(&reindex.Job{Idx: sqlIndexer0, From: version0, To: version0a}, someErr)
		},
	}

	versionsUpdateErr := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("UPDATE %s", versionTableName)).
				WithArgs(version0a, id0).
				WillReturnError(someErr)
			m.ExpectRollback()
		},

		run: func(queue reindex.JobQueue) (interface{}, error) {
			return nil, queue.CompleteJob(&reindex.Job{Idx: sqlIndexer0, From: version0, To: version0a}, nil)
		},
		err: someErr,
	}

	queueUpdateErr := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("UPDATE %s", versionTableName)).
				WithArgs(version0a, id0).
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectExec(fmt.Sprintf("UPDATE %s", queueTableName)).
				WithArgs(reindex.StatusComplete, "", now.Unix(), id0).
				WillReturnError(someErr)
			m.ExpectRollback()
		},

		run: func(queue reindex.JobQueue) (interface{}, error) {
			return nil, queue.CompleteJob(&reindex.Job{Idx: sqlIndexer0, From: version0, To: version0a}, nil)
		},
		err: someErr,
	}

	runCase(t, completeWithSuccess)
	runCase(t, completeWithErr)
	runCase(t, versionsUpdateErr)
	runCase(t, queueUpdateErr)
}

func TestSqlJobQueue_GetAllErrors(t *testing.T) {
	twoFound := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnRows(
					sqlmock.NewRows(queueCols).
						AddRow(id0, version0, version0a, reindex.StatusAvailable, 42, someErr.Error(), 42).
						AddRow(id1, version1, version1a, reindex.StatusAvailable, 43, someErr1.Error(), 43),
				)
			m.ExpectCommit()
		},

		run:    func(queue reindex.JobQueue) (interface{}, error) { return reindex.GetErrors(queue) },
		result: map[string]string{id0: someErr.Error(), id1: someErr1.Error()},
	}

	zeroFound := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnRows(sqlmock.NewRows(queueCols))
			m.ExpectCommit()
		},

		run:    func(queue reindex.JobQueue) (interface{}, error) { return reindex.GetErrors(queue) },
		result: map[string]string{},
	}

	selectErr := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnError(someErr)
			m.ExpectRollback()
		},

		run: func(queue reindex.JobQueue) (interface{}, error) { return reindex.GetErrors(queue) },
		err: someErr,
	}

	runCase(t, twoFound)
	runCase(t, zeroFound)
	runCase(t, selectErr)
}

func TestSqlJobQueue_GetAllJobInfo(t *testing.T) {
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(4*time.Hour))
	defer clock.UnfreezeClock(t)
	now := clock.Now()
	timedOut := clock.Now().Add(-10 * time.Minute)
	maxAttempts := reindex.DefaultMaxAttempts

	// Jobs:
	//	- max attempts, available => err
	//	- max attempts, in progress + timed out => err
	//	- max attempts, in progress + not timed out => no err
	//	- min attempts, available => no err
	//	- complete => no err
	happyPath := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnRows(
					sqlmock.NewRows(queueCols).
						AddRow(id0, version0, version0a, reindex.StatusAvailable, maxAttempts, someErr.Error(), now.Unix()).
						AddRow(id1, version1, version1a, reindex.StatusInProgress, maxAttempts, someErr1.Error(), timedOut.Unix()).
						AddRow(id2, version2, version2a, reindex.StatusInProgress, maxAttempts, someErr2.Error(), now.Unix()).
						AddRow(id3, version3, version3a, reindex.StatusAvailable, 1, someErr3.Error(), timedOut.Unix()).
						AddRow(id4, version4, version4a, reindex.StatusComplete, 1, "", timedOut.Unix()),
				)
			m.ExpectCommit()
		},

		run: func(queue reindex.JobQueue) (interface{}, error) { return queue.GetJobInfos() },
		result: map[string]reindex.JobInfo{
			id0: {IndexerID: id0, Status: reindex.StatusAvailable, Error: someErr.Error(), Attempts: maxAttempts},
			id1: {IndexerID: id1, Status: reindex.StatusInProgress, Error: someErr1.Error(), Attempts: maxAttempts},
			id2: {IndexerID: id2, Status: reindex.StatusInProgress, Error: "", Attempts: maxAttempts},
			id3: {IndexerID: id3, Status: reindex.StatusAvailable, Error: "", Attempts: 1},
			id4: {IndexerID: id4, Status: reindex.StatusComplete, Error: "", Attempts: 1},
		},
	}

	zeroFound := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnRows(sqlmock.NewRows(queueCols))
			m.ExpectCommit()
		},

		run:    func(queue reindex.JobQueue) (interface{}, error) { return queue.GetJobInfos() },
		result: map[string]reindex.JobInfo{},
	}

	selectErr := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", queueColsJoined, queueTableName)).
				WillReturnError(someErr)
			m.ExpectRollback()
		},

		run: func(queue reindex.JobQueue) (interface{}, error) { return queue.GetJobInfos() },
		err: someErr,
	}

	runCase(t, happyPath)
	runCase(t, zeroFound)
	runCase(t, selectErr)
}

type testCase struct {
	setup      func(m sqlmock.Sqlmock)
	registered []indexer.Indexer

	run func(queue reindex.JobQueue) (interface{}, error)

	result interface{}
	err    error
}

func runCase(t *testing.T, test *testCase) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %v", err)
	}

	mock.ExpectBegin()
	test.setup(mock)

	indexer.DeregisterAllForTest(t)
	if test.registered != nil {
		err = indexer.RegisterAll(test.registered...)
		if err != nil {
			t.Fatalf("Error registering indexers: %v", err)
		}
	}

	queue := reindex.NewSQLJobQueue(reindex.DefaultMaxAttempts, db, sqorc.GetSqlBuilder())
	actual, err := test.run(queue)

	if test.err != nil {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), test.err.Error())
	} else {
		assert.NoError(t, err)
	}

	if test.result != nil {
		assert.Equal(t, test.result, actual)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}
