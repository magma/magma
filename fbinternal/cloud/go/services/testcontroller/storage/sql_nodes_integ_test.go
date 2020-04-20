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
	"fmt"
	"testing"
	"time"

	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLNodeLeasorStorage_Integration(t *testing.T) {
	const dbName = "testcontroller__storage__sql_nodes_integ_test"
	db1 := sqorc.OpenCleanForTest(t, dbName, sqorc.PostgresDriver)
	db2 := sqorc.OpenForTest(t, dbName, sqorc.PostgresDriver)
	defer db1.Close()
	defer db2.Close()
	_, err := db1.Exec("DROP TABLE IF EXISTS testcontroller_nodes")
	assert.NoError(t, err)

	idgen := &mockIDGenerator{}
	store := NewSQLNodeLeasorStorage(db1, idgen, sqorc.GetSqlBuilder())
	err = store.Init()
	assert.NoError(t, err)
	// 2nd store instance using different DB conn for concurrency tests
	store2 := NewSQLNodeLeasorStorage(db2, idgen, sqorc.GetSqlBuilder())

	frozenClock := 1000 * time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenClock))
	defer clock.UnfreezeClock(t)

	// Empty cases
	actual, err := store.LeaseNode(nil)
	assert.NoError(t, err)
	assert.Nil(t, actual)

	actualNodes, err := store.GetNodes(nil, nil)
	assert.NoError(t, err)
	assert.Empty(t, actualNodes)

	// Basic CRUD: create 2, update 1, delete 1
	err = store.CreateOrUpdateNode(&MutableCINode{Id: "node1", VpnIP: "192.168.100.1"})
	assert.NoError(t, err)
	err = store.CreateOrUpdateNode(&MutableCINode{Id: "node2", VpnIP: "192.168.100.2"})
	assert.NoError(t, err)
	actualNodes, err = store.GetNodes(nil, nil)
	assert.NoError(t, err)
	expectedNodes := map[string]*CINode{
		"node1": {
			Id:            "node1",
			VpnIp:         "192.168.100.1",
			Available:     true,
			LastLeaseTime: timestampProto(t, 0),
		},
		"node2": {
			Id:            "node2",
			VpnIp:         "192.168.100.2",
			Available:     true,
			LastLeaseTime: timestampProto(t, 0),
		},
	}
	assert.Equal(t, expectedNodes, actualNodes)
	actualNodes, err = store.GetNodes([]string{"node1", "node2"}, nil)
	assert.NoError(t, err)
	assert.Equal(t, expectedNodes, actualNodes)

	err = store.CreateOrUpdateNode(&MutableCINode{Id: "node2", VpnIP: "10.0.2.1"})
	assert.NoError(t, err)
	actualNodes, err = store.GetNodes([]string{"node2", "node1"}, nil)
	assert.NoError(t, err)
	expectedNodes["node2"].VpnIp = "10.0.2.1"
	assert.Equal(t, expectedNodes, actualNodes)

	err = store.DeleteNode("node1")
	assert.NoError(t, err)
	actualNodes, err = store.GetNodes(nil, nil)
	assert.NoError(t, err)
	delete(expectedNodes, "node1")
	assert.Equal(t, expectedNodes, actualNodes)

	// Reserve a specific node, further lease requests should return empty
	// Release that node
	// Lease 1 node, further lease requests should return empty
	// Read should show that the node is unavailable
	// Release the node so we can do concurrency test
	actual, err = store.ReserveNode("node2")
	assert.NoError(t, err)
	expected := &NodeLease{
		Id:      "node2",
		LeaseID: "manual",
		VpnIP:   "10.0.2.1",
	}
	assert.Equal(t, expected, actual)

	actual, err = store.LeaseNode(nil)
	assert.NoError(t, err)
	assert.Nil(t, actual)

	err = store.ReleaseNode("node2", "manual")
	assert.NoError(t, err)

	actual, err = store.LeaseNode(nil)
	assert.NoError(t, err)
	expected = &NodeLease{
		Id:      "node2",
		LeaseID: "1",
		VpnIP:   "10.0.2.1",
	}
	assert.Equal(t, expected, actual)

	actual, err = store.LeaseNode(nil)
	assert.NoError(t, err)
	assert.Nil(t, actual)

	// Manual lease should override any existing lease
	actual, err = store.ReserveNode("node2")
	assert.NoError(t, err)
	expected = &NodeLease{
		Id:      "node2",
		LeaseID: "manual",
		VpnIP:   "10.0.2.1",
	}
	assert.Equal(t, expected, actual)

	actualNodes, err = store.GetNodes(nil, nil)
	assert.NoError(t, err)
	expectedNodes["node2"].Available, expectedNodes["node2"].LastLeaseTime = false, timestampProto(t, int64(frozenClock/time.Second))
	assert.Equal(t, expectedNodes, actualNodes)

	err = store.ReleaseNode("node2", "12345")
	assert.EqualError(t, err, "no node matching the provided ID and lease ID was found")
	err = store.ReleaseNode("node2", "1")
	assert.Error(t, err, "no node matching the provided ID and lease ID was found")
	err = store.ReleaseNode("node2", "manual")
	assert.NoError(t, err)

	frozenClock += 30 * time.Minute
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenClock))

	// Concurrency test: pause one client between SELECT FOR UPDATE and
	// follow-up UPDATE. The second concurrent client should retrieve no
	// available node.
	// see sql_integ_test.go from this package for comments on how the
	// callback works
	assert.NoError(t, err)
	waiter := make(chan error)
	result := make(chan *NodeLease)
	selectedNextNode = func() {
		waiter <- nil
		waiter <- nil
	}
	go func() {
		innerActual, err := store.LeaseNode(nil)
		waiter <- err
		result <- innerActual
	}()

	<-waiter
	selectedNextNode = func() {}
	actual, err = store2.LeaseNode(nil)
	assert.NoError(t, err)
	assert.Nil(t, actual)

	<-waiter
	err = <-waiter
	actual = <-result
	assert.NoError(t, err)
	expected = &NodeLease{
		Id:      "node2",
		VpnIP:   "10.0.2.1",
		LeaseID: "2",
	}
	assert.Equal(t, expected, actual)
	actualNodes, err = store.GetNodes(nil, nil)
	assert.NoError(t, err)
	expectedNodes["node2"].LastLeaseTime = timestampProto(t, int64(frozenClock/time.Second))
	assert.Equal(t, expectedNodes, actualNodes)

	// Timeout the lease by advancing clock by 3 hours, we should get a lease
	frozenClock += 3 * time.Hour
	clock.SetAndFreezeClock(t, time.Unix(0, 0).Add(frozenClock))
	actual, err = store.LeaseNode(nil)
	assert.NoError(t, err)
	expected.LeaseID = "3"
	assert.Equal(t, expected, actual)

	actualNodes, err = store.GetNodes(nil, nil)
	assert.NoError(t, err)
	expectedNodes["node2"].LastLeaseTime = timestampProto(t, int64(frozenClock/time.Second))
	assert.Equal(t, expectedNodes, actualNodes)

	err = store.ReleaseNode("node2", "3")
	assert.NoError(t, err)

	// Another concurrency test, this time with a second client which should
	// receive another lease
	// The correctness of this test case depends on the specific order in which
	// postgres selects the rows from the table. If this starts flaking, we
	// can make the test case order-independent. The behavior appears
	// deterministic for now.
	err = store.CreateOrUpdateNode(&MutableCINode{Id: "node1", VpnIP: "192.168.100.1"})
	assert.NoError(t, err)

	selectedNextNode = func() {
		waiter <- nil
		waiter <- nil
	}
	go func() {
		innerActual, err := store.LeaseNode(nil)
		waiter <- err
		result <- innerActual
	}()

	<-waiter
	selectedNextNode = func() {}
	actual, err = store2.LeaseNode(nil)
	assert.NoError(t, err)
	expected = &NodeLease{
		Id:      "node1",
		VpnIP:   "192.168.100.1",
		LeaseID: "4",
	}
	assert.Equal(t, expected, actual)

	<-waiter
	err = <-waiter
	actual = <-result
	assert.NoError(t, err)
	expected = &NodeLease{
		Id:      "node2",
		VpnIP:   "10.0.2.1",
		LeaseID: "5",
	}
	assert.Equal(t, expected, actual)

	// Release
	err = store.ReleaseNode("node2", "5")
	assert.NoError(t, err)
	err = store.ReleaseNode("node1", "4")
	assert.NoError(t, err)

	// Basic tagged tests
	err = store.CreateOrUpdateNode(&MutableCINode{Id: "node3", Tag: "tag", VpnIP: "10.0.2.1"})
	assert.NoError(t, err)
	actualNodes, err = store.GetNodes(nil, strPtr("tag"))
	assert.NoError(t, err)
	expectedNodes = map[string]*CINode{
		"node3": {
			Id:            "node3",
			VpnIp:         "10.0.2.1",
			Tag:           "tag",
			Available:     true,
			LastLeaseTime: timestampProto(t, 0),
		},
	}
	assert.Equal(t, expectedNodes, actualNodes)

	actualNodes, err = store.GetNodes(nil, strPtr(""))
	assert.NoError(t, err)
	expectedNodes = map[string]*CINode{
		"node1": {
			Id:            "node1",
			VpnIp:         "192.168.100.1",
			Available:     true,
			LastLeaseTime: timestampProto(t, int64(frozenClock/time.Second)),
		},
		"node2": {
			Id:            "node2",
			VpnIp:         "10.0.2.1",
			Available:     true,
			LastLeaseTime: timestampProto(t, int64(frozenClock/time.Second)),
		},
	}
	assert.Equal(t, expectedNodes, actualNodes)

	actual, err = store.LeaseNode(strPtr("tag"))
	assert.NoError(t, err)
	expected = &NodeLease{
		Id:      "node3",
		LeaseID: "6",
		VpnIP:   "10.0.2.1",
	}
	assert.Equal(t, expected, actual)

	actual, err = store.LeaseNode(strPtr("tag"))
	assert.NoError(t, err)
	assert.Nil(t, actual)

	actualNodes, err = store.GetNodes(nil, nil)
	assert.NoError(t, err)
	expectedNodes["node3"] = &CINode{
		Id:            "node3",
		VpnIp:         "10.0.2.1",
		Tag:           "tag",
		Available:     false,
		LastLeaseTime: timestampProto(t, int64(frozenClock/time.Second)),
	}
	assert.Equal(t, expectedNodes, actualNodes)
}

type mockIDGenerator struct {
	current uint64
}

func (m *mockIDGenerator) New() string {
	m.current++
	return fmt.Sprintf("%d", m.current)
}
