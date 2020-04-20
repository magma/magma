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

package storage_test

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"magma/fbinternal/cloud/go/services/testcontroller/storage"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

var allExpectedCols = []string{
	"pk", "type", "config", "is_executing", "state", "error", "last_executed_sec", "next_scheduled_sec",
}
var expectedColsJoined = strings.Join(allExpectedCols, ", ")

const expectedTableName = "testcontroller_tests"

func TestSqlTestControllerStorage_GetTestCases(t *testing.T) {
	happyPath := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", expectedColsJoined, expectedTableName)).
				WithArgs(1, 2).
				WillReturnRows(
					sqlmock.NewRows(allExpectedCols).
						AddRow(1, "foo", []byte("fooconfig"), true, "foostate", nil, 0, 100).
						AddRow(2, "bar", []byte("barconfig"), false, "barstate", "bar error", 42, 69),
				)
			m.ExpectCommit()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return store.GetTestCases([]int64{1, 2})
		},

		expectedResult: map[int64]*storage.TestCase{
			1: {
				Pk:                   1,
				TestCaseType:         "foo",
				TestConfig:           []byte("fooconfig"),
				IsCurrentlyExecuting: true,
				LastExecutionTime:    timestampProto(t, 0),
				State:                "foostate",
				Error:                "",
				NextScheduledTime:    timestampProto(t, 100),
			},
			2: {
				Pk:                   2,
				TestCaseType:         "bar",
				TestConfig:           []byte("barconfig"),
				IsCurrentlyExecuting: false,
				LastExecutionTime:    timestampProto(t, 42),
				State:                "barstate",
				Error:                "bar error",
				NextScheduledTime:    timestampProto(t, 69),
			},
		},
	}

	// load without specified PKs
	loadAll := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", expectedColsJoined, expectedTableName)).
				WillReturnRows(
					sqlmock.NewRows(allExpectedCols).
						AddRow(1, "foo", []byte("fooconfig"), true, "foostate", nil, 0, 100).
						AddRow(2, "bar", []byte("barconfig"), false, "barstate", "bar error", 42, 69),
				)
			m.ExpectCommit()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return store.GetTestCases(nil)
		},

		expectedResult: map[int64]*storage.TestCase{
			1: {
				Pk:                   1,
				TestCaseType:         "foo",
				TestConfig:           []byte("fooconfig"),
				IsCurrentlyExecuting: true,
				LastExecutionTime:    timestampProto(t, 0),
				State:                "foostate",
				Error:                "",
				NextScheduledTime:    timestampProto(t, 100),
			},
			2: {
				Pk:                   2,
				TestCaseType:         "bar",
				TestConfig:           []byte("barconfig"),
				IsCurrentlyExecuting: false,
				LastExecutionTime:    timestampProto(t, 42),
				State:                "barstate",
				Error:                "bar error",
				NextScheduledTime:    timestampProto(t, 69),
			},
		},
	}

	badTimestamp := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", expectedColsJoined, expectedTableName)).
				WillReturnRows(sqlmock.NewRows(allExpectedCols).AddRow(1, "foo", []byte("fooconfig"), true, "foostate", nil, 0, math.MinInt64))
			m.ExpectRollback()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return store.GetTestCases(nil)
		},

		expectedResult: map[int64]*storage.TestCase{},
		expectedError:  errors.New("could not validate next scheduled time for test 1: timestamp: seconds:-9223372036854775808  before 0001-01-01"),
	}

	runCase(t, happyPath)
	runCase(t, loadAll)
	runCase(t, badTimestamp)
}

func TestSqlTestControllerStorage_CreateOrUpdateTestCase(t *testing.T) {
	happyPath := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("INSERT INTO %s", expectedTableName)).
				WithArgs(1, "foo", []byte("fooconfig"), "foo", []byte("fooconfig")).
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return nil, store.CreateOrUpdateTestCase(&storage.MutableTestCase{
				Pk:           1,
				TestCaseType: "foo",
				TestConfig:   []byte("fooconfig"),
			})
		},
	}

	errorCase := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("INSERT INTO %s", expectedTableName)).
				WithArgs(1, "foo", []byte("fooconfig"), "foo", []byte("fooconfig")).
				WillReturnError(errors.New("blah error"))
			m.ExpectRollback()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return nil, store.CreateOrUpdateTestCase(&storage.MutableTestCase{
				Pk:           1,
				TestCaseType: "foo",
				TestConfig:   []byte("fooconfig"),
			})
		},

		expectedError: errors.New("failed to write test case: blah error"),
	}

	runCase(t, happyPath)
	runCase(t, errorCase)
}

func TestSqlTestControllerStorage_DeleteTestCase(t *testing.T) {
	happyPath := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("DELETE FROM %s", expectedTableName)).
				WithArgs(1).
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return nil, store.DeleteTestCase(1)
		},
	}

	errorCase := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("DELETE FROM %s", expectedTableName)).
				WithArgs(1).
				WillReturnError(errors.New("blah error"))
			m.ExpectRollback()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return nil, store.DeleteTestCase(1)
		},

		expectedError: errors.New("failed to delete test case: blah error"),
	}

	runCase(t, happyPath)
	runCase(t, errorCase)
}

func TestSqlTestControllerStorage_GetNextTestForExecution(t *testing.T) {
	frozenTime := 4 * time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenTime))
	defer clock.UnfreezeClock(t)

	happyPath := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", expectedColsJoined, expectedTableName)).
				WithArgs(false, frozenTime/time.Second, true, (frozenTime-time.Hour)/time.Second).
				WillReturnRows(sqlmock.NewRows(allExpectedCols).AddRow(1, "foo", []byte("fooconfig"), true, "foostate", nil, 42, 69))
			m.ExpectExec(fmt.Sprintf("UPDATE %s", expectedTableName)).
				WithArgs(true, frozenTime/time.Second, 1).
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return store.GetNextTestForExecution()
		},

		expectedResult: &storage.TestCase{
			Pk:                   1,
			TestCaseType:         "foo",
			TestConfig:           []byte("fooconfig"),
			IsCurrentlyExecuting: true,
			LastExecutionTime:    timestampProto(t, 42),
			State:                "foostate",
			Error:                "",
			NextScheduledTime:    timestampProto(t, 69),
		},
	}

	emptySelect := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", expectedColsJoined, expectedTableName)).
				WithArgs(false, frozenTime/time.Second, true, (frozenTime-time.Hour)/time.Second).
				WillReturnRows(sqlmock.NewRows(allExpectedCols))
			m.ExpectCommit()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return store.GetNextTestForExecution()
		},
	}

	selectError := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", expectedColsJoined, expectedTableName)).
				WithArgs(false, frozenTime/time.Second, true, (frozenTime-time.Hour)/time.Second).
				WillReturnError(errors.New("blah error"))
			m.ExpectRollback()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return store.GetNextTestForExecution()
		},

		expectedError: errors.New("failed to load next test case: blah error"),
	}

	updateError := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", expectedColsJoined, expectedTableName)).
				WithArgs(false, frozenTime/time.Second, true, (frozenTime-time.Hour)/time.Second).
				WillReturnRows(sqlmock.NewRows(allExpectedCols).AddRow(1, "foo", []byte("fooconfig"), true, "foostate", nil, 42, 69))
			m.ExpectExec(fmt.Sprintf("UPDATE %s", expectedTableName)).
				WithArgs(true, frozenTime/time.Second, 1).
				WillReturnError(errors.New("blah error"))
			m.ExpectRollback()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return store.GetNextTestForExecution()
		},

		expectedError: errors.New("failed to mark test case as executing: blah error"),
	}

	runCase(t, happyPath)
	runCase(t, emptySelect)
	runCase(t, selectError)
	runCase(t, updateError)
}

func TestSqlTestControllerStorage_ReleaseTest(t *testing.T) {
	frozenTime := 4 * time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenTime))
	defer clock.UnfreezeClock(t)

	happyPath := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("UPDATE %s", expectedTableName)).
				WithArgs("newstate", false, frozenTime/time.Second, (frozenTime+5*time.Minute)/time.Second, "errmsg", 1).
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return nil, store.ReleaseTest(1, "newstate", strPtr("errmsg"), 5*time.Minute)
		},
	}

	errorCase := &testCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("UPDATE %s", expectedTableName)).
				WithArgs("newstate", false, frozenTime/time.Second, (frozenTime+5*time.Minute)/time.Second, nil, 1).
				WillReturnError(errors.New("blah error"))
			m.ExpectRollback()
		},

		run: func(store storage.TestControllerStorage) (interface{}, error) {
			return nil, store.ReleaseTest(1, "newstate", nil, 5*time.Minute)
		},

		expectedError: errors.New("failed to release test case: blah error"),
	}

	runCase(t, happyPath)
	runCase(t, errorCase)

}

type testCase struct {
	setup func(m sqlmock.Sqlmock)

	run func(store storage.TestControllerStorage) (interface{}, error)

	expectedError  error
	expectedResult interface{}
}

func runCase(t *testing.T, test *testCase) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}

	mock.ExpectBegin()
	test.setup(mock)

	store := storage.NewSQLTestcontrollerStorage(db, sqorc.GetSqlBuilder())
	actual, err := test.run(store)

	if test.expectedError != nil {
		assert.EqualError(t, err, test.expectedError.Error())
	} else {
		assert.NoError(t, err)
	}

	if test.expectedResult != nil {
		assert.Equal(t, test.expectedResult, actual)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func timestampProto(t *testing.T, unix int64) *timestamp.Timestamp {
	ret, err := ptypes.TimestampProto(time.Unix(unix, 0))
	assert.NoError(t, err)
	return ret
}

func strPtr(s string) *string {
	return &s
}
