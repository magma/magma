/*
 *  Copyright 2020 The Magma Authors.
 *
 *  This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package storage_test

import (
	"testing"
	"time"

	"magma/lte/cloud/go/services/smsd/storage"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

// I'm not totally sold on sqlmock unit tests for a storage impl that's this
// complex, seems like it could be a waste of time. E.g. look at the
// configurator sql unit tests, those are a nightmare
// In any case, I'm leaving the scaffolding up here for now.

// TODO: maybe fill in some happy-path test cases for the important methods
//  using sqlMock

func TestSqlSMSStorage_GetSMSs(t *testing.T) {

}

func TestSQLSMSStorage_GetSMSsToDeliver(t *testing.T) {

}

func TestSQLSMSStorage_CreateSMS(t *testing.T) {

}

func TestSQLSMSStorage_DeleteSMSs(t *testing.T) {

}

func TestSQLSMSStorage_ReportDelivery(t *testing.T) {

}

type testCase struct {
	setup func(sqlmock.Sqlmock)
	run   func(storage storage.SMSStorage) (interface{}, error)

	expectedError  error
	expectedResult interface{}
}

func runCase(t *testing.T, test *testCase) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer func() {
		_ = db.Close()
	}()

	mock.ExpectBegin()
	test.setup(mock)

	store := storage.NewSQLSMSStorage(db, sqorc.GetSqlBuilder(), &mockRefCounter{numRefs: 1}, &mockIDGenerator{})
	actual, err := test.run(store)

	if test.expectedError != nil {
		assert.EqualError(t, err, test.expectedError.Error())
	} else {
		assert.NoError(t, err)
	}

	if test.expectedResult != nil {
		assert.Equal(t, test.expectedResult, actual)
	}
}

func tsProto(t *testing.T, ts int64) *timestamp.Timestamp {
	ret, err := ptypes.TimestampProto(time.Unix(ts, 0))
	assert.NoError(t, err)
	return ret
}
