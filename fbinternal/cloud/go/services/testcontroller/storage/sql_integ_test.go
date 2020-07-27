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
	"testing"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/thoas/go-funk"
)

func TestNewSQLTestcontrollerStorage_Integration(t *testing.T) {
	const dbName = "testcontroller__storage__sql_integ_test"
	db1 := sqorc.OpenCleanForTest(t, dbName, sqorc.PostgresDriver)
	db2 := sqorc.OpenForTest(t, dbName, sqorc.PostgresDriver)
	defer db1.Close()
	defer db2.Close()
	_, err := db1.Exec("DROP TABLE IF EXISTS testcontroller_tests")
	assert.NoError(t, err)

	store := NewSQLTestcontrollerStorage(db1, sqorc.GetSqlBuilder())
	err = store.Init()
	assert.NoError(t, err)
	// 2nd store instance using a different DB connection for concurrency
	// testing
	store2 := NewSQLTestcontrollerStorage(db2, sqorc.GetSqlBuilder())

	// Empty case
	actualCases, err := store.GetTestCases([]int64{1, 2, 3})
	assert.NoError(t, err)
	assert.True(t, funk.IsEmpty(actualCases))

	frozenClock := 1000 * time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenClock))
	defer clock.UnfreezeClock(t)

	// Basic CRUD: Create 2, update 1, delete 1
	err = store.CreateOrUpdateTestCase(&MutableTestCase{Pk: 1, TestCaseType: "foo", TestConfig: []byte("fooconfig")})
	assert.NoError(t, err)
	err = store.CreateOrUpdateTestCase(&MutableTestCase{Pk: 2, TestCaseType: "bar", TestConfig: []byte("barconfig")})
	assert.NoError(t, err)
	actualCases, err = store.GetTestCases([]int64{})
	assert.NoError(t, err)
	expectedCases := map[int64]*TestCase{
		1: {
			Pk:                   1,
			TestCaseType:         "foo",
			TestConfig:           []byte("fooconfig"),
			IsCurrentlyExecuting: false,
			LastExecutionTime:    timestampProto(t, 0),
			State:                CommonStartState,
			Error:                "",
			NextScheduledTime:    timestampProto(t, 0),
		},
		2: {
			Pk:                   2,
			TestCaseType:         "bar",
			TestConfig:           []byte("barconfig"),
			IsCurrentlyExecuting: false,
			LastExecutionTime:    timestampProto(t, 0),
			State:                CommonStartState,
			Error:                "",
			NextScheduledTime:    timestampProto(t, 0),
		},
	}
	assert.Equal(t, expectedCases, actualCases)

	err = store.CreateOrUpdateTestCase(&MutableTestCase{Pk: 1, TestCaseType: "foo2", TestConfig: []byte("fooconfig2")})
	assert.NoError(t, err)
	expectedCases[1].TestCaseType, expectedCases[1].TestConfig = "foo2", []byte("fooconfig2")
	actualCases, err = store.GetTestCases([]int64{1, 2})
	assert.NoError(t, err)
	assert.Equal(t, expectedCases, actualCases)

	err = store.DeleteTestCase(2)
	assert.NoError(t, err)
	delete(expectedCases, 2)
	actualCases, err = store.GetTestCases(nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedCases, actualCases)

	// Acquire 1 test case, verify that further requests are empty and the test
	// state is written properly
	actual, err := store.GetNextTestForExecution()
	assert.NoError(t, err)
	expected := expectedCases[1]
	assert.Equal(t, expected, actual)

	actual, err = store.GetNextTestForExecution()
	assert.NoError(t, err)
	assert.Nil(t, actual)

	actualCases, err = store.GetTestCases(nil)
	assert.NoError(t, err)
	expectedCases = map[int64]*TestCase{
		1: {
			Pk:                   1,
			TestCaseType:         "foo2",
			TestConfig:           []byte("fooconfig2"),
			IsCurrentlyExecuting: true,
			LastExecutionTime:    timestampProto(t, int64(frozenClock/time.Second)),
			State:                CommonStartState,
			Error:                "",
			NextScheduledTime:    timestampProto(t, 0),
		},
	}
	assert.Equal(t, expectedCases, actualCases)

	// Release the test case, verify its state is written
	err = store.ReleaseTest(1, "nextstate", nil, 30*time.Minute)
	assert.NoError(t, err)
	actualCases, err = store.GetTestCases(nil)
	assert.NoError(t, err)
	expectedCases = map[int64]*TestCase{
		1: {
			Pk:                   1,
			TestCaseType:         "foo2",
			TestConfig:           []byte("fooconfig2"),
			IsCurrentlyExecuting: false,
			LastExecutionTime:    timestampProto(t, int64(frozenClock/time.Second)),
			State:                "nextstate",
			Error:                "",
			NextScheduledTime:    timestampProto(t, int64((frozenClock+30*time.Minute)/time.Second)),
		},
	}
	assert.Equal(t, expectedCases, actualCases)

	// We should not get any test cases now since this one is not up for
	// execution and it's not timed out
	actual, err = store.GetNextTestForExecution()
	assert.NoError(t, err)
	assert.Nil(t, actual)

	// Advance the clock to 30 minutes past the next scheduled time,
	// we should get this test case now
	frozenClock += time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenClock))
	actual, err = store.GetNextTestForExecution()
	assert.NoError(t, err)
	assert.Equal(t, expectedCases[1], actual)
	actualCases, err = store.GetTestCases(nil)
	assert.NoError(t, err)
	expectedCases = map[int64]*TestCase{
		1: {
			Pk:                   1,
			TestCaseType:         "foo2",
			TestConfig:           []byte("fooconfig2"),
			IsCurrentlyExecuting: true,
			LastExecutionTime:    timestampProto(t, int64(frozenClock/time.Second)),
			State:                "nextstate",
			Error:                "",
			NextScheduledTime:    timestampProto(t, int64((frozenClock-30*time.Minute)/time.Second)),
		},
	}
	assert.Equal(t, expectedCases, actualCases)

	// Release this test case
	err = store.ReleaseTest(1, "anotherstate", strPtr("oh dang"), 30*time.Minute)
	assert.NoError(t, err)
	actualCases, err = store.GetTestCases(nil)
	assert.NoError(t, err)
	expectedCases = map[int64]*TestCase{
		1: {
			Pk:                   1,
			TestCaseType:         "foo2",
			TestConfig:           []byte("fooconfig2"),
			IsCurrentlyExecuting: false,
			LastExecutionTime:    timestampProto(t, int64(frozenClock/time.Second)),
			State:                "anotherstate",
			Error:                "oh dang",
			NextScheduledTime:    timestampProto(t, int64((frozenClock+30*time.Minute)/time.Second)),
		},
	}
	assert.Equal(t, expectedCases, actualCases)

	// Advance the clock again to 30 minutes past scheduled runtime
	frozenClock += time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenClock))

	// Concurrency test: pause one client between the SELECT FOR UPDATE and its
	// follow-up UPDATE. The second concurrent client should retrieve no
	// available test case
	waiter := make(chan error)
	result := make(chan *TestCase)
	defer close(waiter)
	defer close(result)
	selectedNextTestCase = func() {
		// We need to clear this callback for the 2nd call, but we need to do
		// that mutation *after* the first call has already entered the
		// callback. This first send signals that we've hit the callback.
		waiter <- nil

		// Now this second send will block until we receive from the channel
		// again in the outer test case
		waiter <- nil
	}
	go func() {
		innerActual, err := store.GetNextTestForExecution()

		// signal to the outside test case that the call has finished
		waiter <- err
		result <- innerActual
	}()

	// Clear the callback for the second call which should not pause
	// Block until the first call enters the callback
	<-waiter
	selectedNextTestCase = func() {}
	actual, err = store2.GetNextTestForExecution()
	assert.NoError(t, err)
	assert.Nil(t, actual)

	// Now receive from the channels to unblock the first call
	<-waiter
	err = <-waiter
	actual = <-result
	assert.NoError(t, err)
	assert.Equal(t, expectedCases[1], actual)
	actualCases, err = store.GetTestCases(nil)
	assert.NoError(t, err)
	expectedCases = map[int64]*TestCase{
		1: {
			Pk:                   1,
			TestCaseType:         "foo2",
			TestConfig:           []byte("fooconfig2"),
			IsCurrentlyExecuting: true,
			LastExecutionTime:    timestampProto(t, int64(frozenClock/time.Second)),
			State:                "anotherstate",
			Error:                "oh dang",
			NextScheduledTime:    timestampProto(t, int64((frozenClock-30*time.Minute)/time.Second)),
		},
	}
	assert.Equal(t, expectedCases, actualCases)

	// Advance the clock to timeout the test case
	frozenClock += 2*time.Hour + 30*time.Minute
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenClock))
	// Seed a second test case
	err = store.CreateOrUpdateTestCase(&MutableTestCase{Pk: 2, TestCaseType: "bar", TestConfig: []byte("barconfig")})
	assert.NoError(t, err)

	// Concurrency test: pause the client again, run a query on the second one
	// which should grab the second test case
	// We depend on postgres scanning the index on the Pk column in order to
	// get a deterministic result between the two clients here.
	selectedNextTestCase = func() {
		waiter <- nil
		waiter <- nil
	}
	go func() {
		innerActual, err := store.GetNextTestForExecution()
		waiter <- err
		result <- innerActual
	}()

	<-waiter
	selectedNextTestCase = func() {}
	actual, err = store2.GetNextTestForExecution()
	assert.NoError(t, err)
	expected = &TestCase{
		Pk:                   2,
		TestCaseType:         "bar",
		TestConfig:           []byte("barconfig"),
		IsCurrentlyExecuting: false,
		LastExecutionTime:    timestampProto(t, 0),
		State:                CommonStartState,
		Error:                "",
		NextScheduledTime:    timestampProto(t, 0),
	}
	assert.Equal(t, expected, actual)

	<-waiter
	err = <-waiter
	actual = <-result
	assert.NoError(t, err)
	assert.Equal(t, expectedCases[1], actual)
	actualCases, err = store.GetTestCases(nil)
	assert.NoError(t, err)
	expectedCases = map[int64]*TestCase{
		1: {
			Pk:                   1,
			TestCaseType:         "foo2",
			TestConfig:           []byte("fooconfig2"),
			IsCurrentlyExecuting: true,
			LastExecutionTime:    timestampProto(t, int64(frozenClock/time.Second)),
			State:                "anotherstate",
			Error:                "oh dang",
			// Since we didn't release this test case, this field never got
			// updated from its previous value of frozenClock - 30 minutes
			// Then we advanced the clock by 2hr30min to timeout the test,
			// so this should read 3 hours ago
			NextScheduledTime: timestampProto(t, int64((frozenClock-3*time.Hour)/time.Second)),
		},
		2: {
			Pk:                   2,
			TestCaseType:         "bar",
			TestConfig:           []byte("barconfig"),
			IsCurrentlyExecuting: true,
			LastExecutionTime:    timestampProto(t, int64(frozenClock/time.Second)),
			State:                CommonStartState,
			Error:                "",
			NextScheduledTime:    timestampProto(t, 0),
		},
	}
	assert.Equal(t, expectedCases, actualCases)
}

func timestampProto(t *testing.T, unix int64) *timestamp.Timestamp {
	ret, err := ptypes.TimestampProto(time.Unix(unix, 0))
	assert.NoError(t, err)
	return ret
}

func strPtr(s string) *string {
	return &s
}
