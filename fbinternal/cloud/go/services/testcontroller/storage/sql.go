/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/Masterminds/squirrel"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

const (
	testCaseTable = "testcontroller_tests"

	pkCol            = "pk"
	typeCol          = "type"
	configCol        = "config"
	isExecutingCol   = "is_executing"
	lastExecutedCol  = "last_executed_sec"
	nextScheduledCol = "next_scheduled_sec"
	stateCol         = "state"
	errCol           = "error"
)

const testCaseTimeout = 1 * time.Hour

var (
	// selectedNextTestCase is an unexported hook for testing
	selectedNextTestCase = func() {}
)

func NewSQLTestcontrollerStorage(db *sql.DB, sqlBuilder sqorc.StatementBuilder) TestControllerStorage {
	return &sqlTestControllerStorage{db: db, builder: sqlBuilder}
}

type sqlTestControllerStorage struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func (s *sqlTestControllerStorage) Init() (err error) {
	tx, err := s.db.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		} else {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = fmt.Errorf("%s; rollback error: %s", err, rollbackErr)
			}
		}
	}()

	_, err = s.builder.CreateTable(testCaseTable).
		IfNotExists().
		Column(pkCol).Type(sqorc.ColumnTypeInt).PrimaryKey().EndColumn().
		Column(typeCol).Type(sqorc.ColumnTypeText).NotNull().EndColumn().
		Column(configCol).Type(sqorc.ColumnTypeBytes).NotNull().EndColumn().
		Column(isExecutingCol).Type(sqorc.ColumnTypeBool).NotNull().Default(false).EndColumn().
		Column(stateCol).Type(sqorc.ColumnTypeText).NotNull().Default(fmt.Sprintf("'%s'", CommonStartState)).EndColumn().
		Column(errCol).Type(sqorc.ColumnTypeText).EndColumn().
		Column(lastExecutedCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		Column(nextScheduledCol).Type(sqorc.ColumnTypeInt).NotNull().Default(0).EndColumn().
		RunWith(tx).
		Exec()
	if err != nil {
		err = errors.Wrap(err, "failed to create test case table")
	}
	return
}

func (s *sqlTestControllerStorage) GetTestCases(pks []int64) (map[int64]*TestCase, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		builder := s.builder.Select(pkCol, typeCol, configCol, isExecutingCol, stateCol, errCol, lastExecutedCol, nextScheduledCol).
			From(testCaseTable).
			RunWith(tx)
		if !funk.IsEmpty(pks) {
			builder = builder.Where(squirrel.Eq{pkCol: pks})
		}

		rows, err := builder.Query()
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve test cases")
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetTestCases")

		return scanTestCases(rows)
	}

	ret, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	if err != nil {
		return map[int64]*TestCase{}, err
	}
	return ret.(map[int64]*TestCase), nil
}

func (s *sqlTestControllerStorage) CreateOrUpdateTestCase(testCase *MutableTestCase) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := s.builder.Insert(testCaseTable).
			Columns(pkCol, typeCol, configCol).
			Values(testCase.Pk, testCase.TestCaseType, testCase.TestConfig).
			OnConflict(
				[]sqorc.UpsertValue{
					{Column: typeCol, Value: testCase.TestCaseType},
					{Column: configCol, Value: testCase.TestConfig},
				},
				pkCol,
			).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "failed to write test case")
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	return err
}

func (s *sqlTestControllerStorage) DeleteTestCase(pk int64) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		_, err := s.builder.Delete(testCaseTable).
			Where(squirrel.Eq{pkCol: pk}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "failed to delete test case")
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	return err
}

func (s *sqlTestControllerStorage) GetNextTestForExecution() (*TestCase, error) {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		now := clock.Now()
		timeoutThreshold := now.Add(-testCaseTimeout)

		// SELECT pk, type, config, is_executing, state, error, last_executed_sec, next_scheduled_sec
		// FROM testcontroller_tests
		// WHERE
		//   (NOT is_executing AND next_scheduled_sec < 42)
		//   OR
		//   (is_executing AND last_executed_sec < 24)
		// LIMIT 1
		// FOR UPDATE SKIP LOCKED
		// Query instead of QueryRow because we don't want to error out if
		// there's no results (this is a valid response)
		rows, err := s.builder.Select(pkCol, typeCol, configCol, isExecutingCol, stateCol, errCol, lastExecutedCol, nextScheduledCol).
			From(testCaseTable).
			Where(
				squirrel.Or{
					// normal case: not executing and timer is up
					squirrel.And{
						squirrel.Eq{isExecutingCol: false},
						squirrel.Lt{nextScheduledCol: now.Unix()},
					},
					// timeout case: is executing but has not been released after 2 hours
					squirrel.And{
						squirrel.Eq{isExecutingCol: true},
						squirrel.Lt{lastExecutedCol: timeoutThreshold.Unix()},
					},
				},
			).
			Limit(1).
			Suffix("FOR UPDATE SKIP LOCKED").
			RunWith(tx).
			Query()
		if err != nil {
			return nil, errors.Wrap(err, "failed to load next test case")
		}
		defer sqorc.CloseRowsLogOnError(rows, "GetNextTestForExecution")

		testsByPk, err := scanTestCases(rows)
		if err != nil {
			return nil, err
		}
		if funk.IsEmpty(testsByPk) {
			return nil, nil
		}
		selectedTest := testsByPk[funk.Head(funk.Keys(testsByPk)).(int64)]

		// Call this callback in case we want to pause here during a test
		selectedNextTestCase()

		// Update the test case, mark it as executing
		_, err = s.builder.Update(testCaseTable).
			Set(isExecutingCol, true).
			Set(lastExecutedCol, now.Unix()).
			Where(squirrel.Eq{pkCol: selectedTest.Pk}).
			RunWith(tx).
			Exec()
		if err != nil {
			return nil, errors.Wrap(err, "failed to mark test case as executing")
		}
		return selectedTest, nil
	}

	ret, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	switch {
	case err != nil:
		return nil, err
	case ret == nil:
		return nil, nil
	default:
		return ret.(*TestCase), nil
	}
}

func (s *sqlTestControllerStorage) ReleaseTest(pk int64, newState string, errorString *string, nextSchedule time.Duration) error {
	txFn := func(tx *sql.Tx) (interface{}, error) {
		builder := s.builder.Update(testCaseTable).
			Set(stateCol, newState).
			Set(isExecutingCol, false).
			Set(lastExecutedCol, clock.Now().Unix()).
			Set(nextScheduledCol, clock.Now().Add(nextSchedule).Unix()).
			Where(squirrel.Eq{pkCol: pk}).
			RunWith(tx)
		if errorString != nil {
			builder = builder.Set(errCol, sql.NullString{String: *errorString, Valid: true})
		} else {
			builder = builder.Set(errCol, sql.NullString{Valid: false})
		}

		_, err := builder.Exec()
		if err != nil {
			return nil, errors.Wrap(err, "failed to release test case")
		}
		return nil, nil
	}

	_, err := sqorc.ExecInTx(s.db, nil, nil, txFn)
	return err
}

func scanTestCases(rows *sql.Rows) (map[int64]*TestCase, error) {
	ret := map[int64]*TestCase{}
	for rows.Next() {
		var pk, lastExecd, nextSched int64
		var typeVal, state string
		var errMsg sql.NullString
		var config []byte
		var isExecuting bool

		err := rows.Scan(&pk, &typeVal, &config, &isExecuting, &state, &errMsg, &lastExecd, &nextSched)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan test case row")
		}

		lastExecTs, err := ptypes.TimestampProto(time.Unix(lastExecd, 0))
		if err != nil {
			return nil, errors.Wrapf(err, "could not validate last exec time for test %d", pk)
		}
		nextSchedTs, err := ptypes.TimestampProto(time.Unix(nextSched, 0))
		if err != nil {
			return nil, errors.Wrapf(err, "could not validate next scheduled time for test %d", pk)
		}

		ret[pk] = &TestCase{
			Pk:                   pk,
			TestCaseType:         typeVal,
			TestConfig:           config,
			State:                state,
			IsCurrentlyExecuting: isExecuting,
			LastExecutionTime:    lastExecTs,
			NextScheduledTime:    nextSchedTs,
			Error:                errMsg.String,
		}
	}
	return ret, nil
}
