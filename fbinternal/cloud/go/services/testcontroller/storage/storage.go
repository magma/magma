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
	"errors"
	"time"
)

// CommonStartState will be the value of a test case's State field if it has
// not been scheduled for execution since its creation. Every new test case
// starts off in this state.
const CommonStartState = "_test_controller_start_state"

// TestControllerStorage is the storage interface for managing end to end test
// cases and their execution.
type TestControllerStorage interface {
	// Init performs on-start initialization work such as table creation.
	Init() error

	// GetTestCases returns the requested test cases keyed by primary key.
	// If the `pks` parameter is empty, this will load and return all test
	// cases.
	GetTestCases(pks []int64) (map[int64]*TestCase, error)

	// CreateOrUpdateTestCase updates the configuration of a test case or
	// creates the test case if it doesn't yet exist.
	CreateOrUpdateTestCase(testCase *MutableTestCase) error

	// DeleteTestCase deletes the specified test case.
	DeleteTestCase(pk int64) error

	// GetNextTestForExecution loads a test case that is available for
	// continued execution, marks it as locked, and returns it.
	// A nil response with no error will indicate that there is no available
	// work to be done.
	GetNextTestForExecution() (*TestCase, error)

	// ReleaseTest will mark a test as unlocked and schedule it for further
	// execution. If a test case is not released by ReleaseTest, it will
	// automatically time out and be available for scheduling after a timeout
	// period.
	// `newState` is the name of the state that the test should transition into
	// `error` if non-nil is the error message to save
	// `nextSchedule` is how long to wait before making the test available for
	// further execution.
	ReleaseTest(pk int64, newState string, errorString *string, nextSchedule time.Duration) error
}

// ErrBadRelease indicates a bad argument to ReleaseNode
var ErrBadRelease = errors.New("no node matching the provided ID and lease ID was found")

// NodeLeasorStorage is the storage interface for managing baremetal CI worker
// nodes.
type NodeLeasorStorage interface {
	// Init performs on-start initialization work such as table creation.
	Init() error

	// GetNodes returns the requested CI nodes keyed by ID>
	// If the ids parameter is empty, this will return all nodes.
	// The tag parameter, if non-nil, will filter the returned nodes to only
	// those which match the tag. Tags can be 0-length so a pointer to an
	// empty string will result in a different query than a nil tag.
	GetNodes(ids []string, tag *string) (map[string]*CINode, error)

	// CreateOrUpdateNode updates the configuration of a node or creates it
	// if it doesn't exist.
	CreateOrUpdateNode(node *MutableCINode) error

	// DeleteNode deletes the identified node.
	DeleteNode(id string) error

	// LeaseNode obtains a lease on the next available worker node if one
	// exists, otherwise returns a nil NodeLease.
	// The tag parameter, if non-nil, will lease a node with the specified
	// tag. Otherwise, any node in the available pool will be returned.
	LeaseNode(tag *string) (*NodeLease, error)

	// ReserveNode obtains a lease on a specific node if it hasn't already been
	// leased (e.g. taking a node out of the pool for manual troubleshooting).
	ReserveNode(id string) (*NodeLease, error)

	// ReleaseNode releases (abandons) the lease on a worker node.
	// The provided leaseID argument is expected to match the leaseID of the
	// most recent lease of the node.
	ReleaseNode(id string, leaseID string) error
}
