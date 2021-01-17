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

package store

import (
	"fmt"
	"testing"
	"time"

	"magma/fbinternal/cloud/go/services/testcontroller/obsidian/models"
	"magma/fbinternal/cloud/go/services/testcontroller/utils"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const (
	DBType = "testcontroller"
)

type TestControllerStore struct {
	factory blobstore.BlobStorageFactory
	store   blobstore.TransactionalBlobStorage
	clock   clock
}

type clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time {
	return time.Now()
}

func NewTestControllerStore(factory blobstore.BlobStorageFactory) (*TestControllerStore, error) {
	s := &TestControllerStore{
		factory: factory,
		store:   nil,
		clock:   &realClock{},
	}
	return s, s.validate()
}

func (s *TestControllerStore) ShouldExecuteScript(networkID string, currentVersion string, minimumWaitTime int64) (executeScript bool, err error) {
	err = s.startTransaction()
	if err != nil {
		err = errors.Wrap(err, "Error starting transaction for testcontroller store")
		return
	}

	defer func() {
		switch err {
		case nil:
			err = s.store.Commit()
		default:
			if rollbackErr := s.store.Rollback(); rollbackErr != nil {
				glog.Errorf("Error rolling back transaction for testcontroller store: %s", rollbackErr)
			}
		}
	}()

	err = s.incrementVersion(networkID)
	if err != nil {
		err = errors.Wrap(err, "Error incrementing version in database query for testcontroller store")
		return
	}

	executeScript, err = s.shouldExecuteScript(networkID, currentVersion, minimumWaitTime)
	if err != nil {
		err = errors.Wrap(err, "Error checking if should execute script for testcontroller")
		return
	}
	return
}

func (s *TestControllerStore) shouldExecuteScript(networkID string, currentVersion string, minimumWaitTime int64) (bool, error) {
	version, timestamp, err := s.get(networkID)
	if err != nil {
		return false, err
	}

	noLatestScriptExecution := version == "" || timestamp == 0
	if noLatestScriptExecution {
		err := s.put(networkID, currentVersion, s.clock.Now().Unix())
		if err != nil {
			return false, err
		}
		return true, nil
	}

	if version == "0.0.0" {
		// put in form version-timestamp-hash
		version = version + "-0-0"
	}

	upgraded, err := utils.IsNewerVersion(version, currentVersion)
	if err != nil {
		return false, err
	}

	minTimeElapsed, err := s.isMinimumWaitTimeElapsed(timestamp, minimumWaitTime)
	if err != nil {
		return false, err
	}

	executeScript := minTimeElapsed && upgraded

	if executeScript {
		err := s.put(networkID, currentVersion, s.clock.Now().Unix())
		if err != nil {
			return false, err
		}
	}
	return executeScript, nil
}

func (s *TestControllerStore) validate() error {
	if s == nil {
		return fmt.Errorf("Nil TestControllerStore")
	}
	if s.factory == nil {
		return fmt.Errorf("Nil TestController blobstore factory")
	}
	return nil
}

func (s *TestControllerStore) startTransaction() error {
	// note: if multiple instances of the testcontroller service attempt to
	// start a transaction concurrently, the database might return an expected
	// serialization failure:
	// https://github.com/lib/pq/blob/master/error.go 40001
	store, err := s.factory.StartTransaction(&storage.TxOptions{Isolation: storage.LevelSerializable})
	if err != nil {
		return err
	}
	s.store = store
	return nil
}

func (s *TestControllerStore) isMinimumWaitTimeElapsed(timestamp int64, minimumWaitTime int64) (bool, error) {
	t := time.Unix(timestamp, 0)
	minWaitTime := time.Minute * time.Duration(minimumWaitTime)
	return s.clock.Now().After(t.Add(minWaitTime)), nil
}

func (s *TestControllerStore) get(networkID string) (string, int64, error) {
	typeAndKey := storage.TypeAndKey{
		Type: DBType,
		Key:  "gateway_version",
	}

	blob, err := s.store.Get(networkID, typeAndKey)
	if err != nil {
		return "", 0, err
	}

	execution, err := blobToExecStruct(blob)
	if err != nil {
		return "", 0, err
	}
	return execution.Version, execution.Timestamp, nil
}

func (s *TestControllerStore) put(networkID string, version string, timestamp int64) error {
	upgrade := &models.LatestScriptExecution{
		Version:   version,
		Timestamp: timestamp,
	}
	blob, err := execStructToBlob(upgrade)
	if err != nil {
		return err
	}
	err = s.store.CreateOrUpdate(networkID, blobstore.Blobs{blob})
	if err != nil {
		return err
	}
	return nil
}

func (s *TestControllerStore) incrementVersion(networkID string) error {
	typeAndKey := storage.TypeAndKey{
		Type: DBType,
		Key:  "gateway_version",
	}
	err := s.store.IncrementVersion(networkID, typeAndKey)
	if err != nil {
		return err
	}
	return nil
}

func blobToExecStruct(blob blobstore.Blob) (*models.LatestScriptExecution, error) {
	if len(blob.Value) == 0 {
		return &models.LatestScriptExecution{}, nil
	}
	execution := &models.LatestScriptExecution{}
	err := execution.UnmarshalBinary(blob.Value)
	if err != nil {
		return nil, err
	}
	return execution, nil
}

func execStructToBlob(execution *models.LatestScriptExecution) (blobstore.Blob, error) {
	marshaled, err := execution.MarshalBinary()
	if err != nil {
		return blobstore.Blob{}, err
	}
	return blobstore.Blob{
		Type:  DBType,
		Key:   "gateway_version",
		Value: marshaled,
	}, nil
}

// This method exists ONLY for testing - thus the required but unused *testing.T param
// DO NOT USE IN ANYTHING BUT TESTS
func (s *TestControllerStore) SetClock(_ *testing.T, mockClock clock) {
	s.clock = mockClock
}
