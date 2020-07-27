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
	"fmt"
	"strings"
	"testing"
	"time"

	"magma/fbinternal/cloud/go/services/testcontroller/storage"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var expectedNodeCols = []string{
	"pk", "vpn_ip", "tag", "available", "last_leased_sec",
}
var nodeColsJoined = strings.Join(expectedNodeCols, ", ")
var expectedNodeTable = "testcontroller_nodes"

func TestSqlNodeLeasorStorage_GetNodes(t *testing.T) {
	happyPath := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", nodeColsJoined, expectedNodeTable)).
				WithArgs("foo", "bar").
				WillReturnRows(
					sqlmock.NewRows(expectedNodeCols).
						AddRow("foo", "192.168.100.1", "", true, 0).
						AddRow("bar", "10.0.2.1", "tag", false, 100),
				)
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return s.GetNodes([]string{"foo", "bar"}, nil)
		},
		expectedResult: map[string]*storage.CINode{
			"foo": {
				Id:            "foo",
				VpnIp:         "192.168.100.1",
				Available:     true,
				LastLeaseTime: timestampProto(t, 0),
			},
			"bar": {
				Id:            "bar",
				Tag:           "tag",
				VpnIp:         "10.0.2.1",
				Available:     false,
				LastLeaseTime: timestampProto(t, 100),
			},
		},
	}

	loadAll := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", nodeColsJoined, expectedNodeTable)).
				WillReturnRows(
					sqlmock.NewRows(expectedNodeCols).
						AddRow("foo", "192.168.100.1", "", true, 0).
						AddRow("bar", "10.0.2.1", "", false, 100),
				)
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return s.GetNodes(nil, nil)
		},
		expectedResult: map[string]*storage.CINode{
			"foo": {
				Id:            "foo",
				VpnIp:         "192.168.100.1",
				Available:     true,
				LastLeaseTime: timestampProto(t, 0),
			},
			"bar": {
				Id:            "bar",
				VpnIp:         "10.0.2.1",
				Available:     false,
				LastLeaseTime: timestampProto(t, 100),
			},
		},
	}

	loadTag := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", nodeColsJoined, expectedNodeTable)).
				WithArgs("tag").
				WillReturnRows(
					sqlmock.NewRows(expectedNodeCols).
						AddRow("bar", "10.0.2.1", "tag", false, 100),
				)
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return s.GetNodes(nil, strPtr("tag"))
		},
		expectedResult: map[string]*storage.CINode{
			"bar": {
				Id:            "bar",
				Tag:           "tag",
				VpnIp:         "10.0.2.1",
				Available:     false,
				LastLeaseTime: timestampProto(t, 100),
			},
		},
	}

	queryError := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", nodeColsJoined, expectedNodeTable)).
				WillReturnError(errors.New("mock query error"))
			m.ExpectRollback()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return s.GetNodes(nil, nil)
		},
		expectedError: errors.New("failed to retrieve nodes: mock query error"),
	}

	runNodeCase(t, happyPath)
	runNodeCase(t, loadAll)
	runNodeCase(t, loadTag)
	runNodeCase(t, queryError)
}

func TestSqlNodeLeasorStorage_CreateOrUpdateNode(t *testing.T) {
	happyPath := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("INSERT INTO %s", expectedNodeTable)).
				WithArgs("foo", "tag", "192.168.100.1", "tag", "192.168.100.1").
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return nil, s.CreateOrUpdateNode(&storage.MutableCINode{Id: "foo", Tag: "tag", VpnIP: "192.168.100.1"})
		},
	}

	errorCase := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("INSERT INTO %s", expectedNodeTable)).
				WithArgs("foo", "", "192.168.100.1", "", "192.168.100.1").
				WillReturnError(errors.New("mock exec error"))
			m.ExpectRollback()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return nil, s.CreateOrUpdateNode(&storage.MutableCINode{Id: "foo", VpnIP: "192.168.100.1"})
		},
		expectedError: errors.New("failed to write update to node: mock exec error"),
	}

	runNodeCase(t, happyPath)
	runNodeCase(t, errorCase)
}

func TestSqlNodeLeasorStorage_DeleteNode(t *testing.T) {
	happyPath := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("DELETE FROM %s", expectedNodeTable)).
				WithArgs("foo").
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return nil, s.DeleteNode("foo")
		},
	}

	errorCase := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("DELETE FROM %s", expectedNodeTable)).
				WithArgs("foo").
				WillReturnError(errors.New("mock exec error"))
			m.ExpectRollback()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return nil, s.DeleteNode("foo")
		},
		expectedError: errors.New("failed to delete node: mock exec error"),
	}

	runNodeCase(t, happyPath)
	runNodeCase(t, errorCase)
}

func TestSqlNodeLeasorStorage_LeaseNode(t *testing.T) {
	frozenTime := 4 * time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenTime))
	defer clock.UnfreezeClock(t)

	happyPath := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", nodeColsJoined, expectedNodeTable)).
				WithArgs(true, false, (frozenTime-2*time.Hour)/time.Second).
				WillReturnRows(sqlmock.NewRows(expectedNodeCols).AddRow("foo", "192.168.100.1", "", true, 0))
			m.ExpectExec(fmt.Sprintf("UPDATE %s", expectedNodeTable)).
				WithArgs(false, frozenTime/time.Second, "1", "foo").
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return s.LeaseNode(nil)
		},
		expectedResult: &storage.NodeLease{
			Id:      "foo",
			VpnIP:   "192.168.100.1",
			LeaseID: "1",
		},
	}

	withTag := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", nodeColsJoined, expectedNodeTable)).
				WithArgs(true, false, (frozenTime-2*time.Hour)/time.Second, "tag").
				WillReturnRows(sqlmock.NewRows(expectedNodeCols).AddRow("foo", "192.168.100.1", "tag", true, 0))
			m.ExpectExec(fmt.Sprintf("UPDATE %s", expectedNodeTable)).
				WithArgs(false, frozenTime/time.Second, "1", "foo").
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return s.LeaseNode(strPtr("tag"))
		},
		expectedResult: &storage.NodeLease{
			Id:      "foo",
			VpnIP:   "192.168.100.1",
			LeaseID: "1",
		},
	}

	emptySelect := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", nodeColsJoined, expectedNodeTable)).
				WithArgs(true, false, (frozenTime-2*time.Hour)/time.Second).
				WillReturnRows(sqlmock.NewRows(expectedNodeCols))
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return s.LeaseNode(nil)
		},
	}

	// add cases for select and update query errors later

	runNodeCase(t, happyPath)
	runNodeCase(t, withTag)
	runNodeCase(t, emptySelect)
}

func TestSqlNodeLeasor_ReserveNode(t *testing.T) {
	frozenTime := 4 * time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenTime))
	defer clock.UnfreezeClock(t)

	happyPath := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", nodeColsJoined, expectedNodeTable)).
				WithArgs("foo").
				WillReturnRows(sqlmock.NewRows(expectedNodeCols).AddRow("foo", "192.168.100.1", "", true, 0))
			m.ExpectExec(fmt.Sprintf("UPDATE %s", expectedNodeTable)).
				WithArgs(false, frozenTime/time.Second, "manual", "foo").
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return s.ReserveNode("foo")
		},
		expectedResult: &storage.NodeLease{
			Id:      "foo",
			VpnIP:   "192.168.100.1",
			LeaseID: "manual",
		},
	}

	emptySelect := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectQuery(fmt.Sprintf("SELECT %s FROM %s", nodeColsJoined, expectedNodeTable)).
				WithArgs("foo").
				WillReturnRows(sqlmock.NewRows(expectedNodeCols))
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return s.ReserveNode("foo")
		},
	}

	// add cases for select and update query errors later

	runNodeCase(t, happyPath)
	runNodeCase(t, emptySelect)
}

func TestSqlNodeLeasorStorage_ReleaseNode(t *testing.T) {
	frozenTime := 4 * time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenTime))
	defer clock.UnfreezeClock(t)

	happyPath := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("UPDATE %s", expectedNodeTable)).
				WithArgs(true, "foo", "fooLease").
				WillReturnResult(sqlmock.NewResult(1, 1))
			m.ExpectCommit()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return nil, s.ReleaseNode("foo", "fooLease")
		},
	}

	emptySelect := &nodeTestCase{
		setup: func(m sqlmock.Sqlmock) {
			m.ExpectExec(fmt.Sprintf("UPDATE %s", expectedNodeTable)).
				WithArgs(true, "foo", "fooLease").
				WillReturnResult(sqlmock.NewResult(1, 0))
			m.ExpectRollback()
		},
		run: func(s storage.NodeLeasorStorage) (interface{}, error) {
			return nil, s.ReleaseNode("foo", "fooLease")
		},
		expectedError: errors.New("no node matching the provided ID and lease ID was found"),
	}

	runNodeCase(t, happyPath)
	runNodeCase(t, emptySelect)

}

type nodeTestCase struct {
	setup func(m sqlmock.Sqlmock)

	run func(store storage.NodeLeasorStorage) (interface{}, error)

	expectedError  error
	expectedResult interface{}
}

func runNodeCase(t *testing.T, test *nodeTestCase) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}

	mock.ExpectBegin()
	test.setup(mock)

	store := storage.NewSQLNodeLeasorStorage(db, &mockIDGenerator{}, sqorc.GetSqlBuilder())
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

type mockIDGenerator struct {
	current uint64
}

func (m *mockIDGenerator) New() string {
	m.current++
	return fmt.Sprintf("%d", m.current)
}
