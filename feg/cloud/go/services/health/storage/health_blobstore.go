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

	"github.com/golang/glog"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/services/health"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
)

type healthBlobstore struct {
	factory blobstore.StoreFactory
}

// NewHealthBlobstore creates a new HealthBlobstore using the provided
// blobstore factory for the underlying storage functionality.
func NewHealthBlobstore(factory blobstore.StoreFactory) (HealthBlobstore, error) {
	if factory == nil {
		return nil, fmt.Errorf("Storage factory is nil")
	}
	return &healthBlobstore{
		factory,
	}, nil
}

// GetHealth fetches health status for the given networkID and gatewayID from
// the blobstore.
func (h *healthBlobstore) GetHealth(networkID string, gatewayID string) (*fegprotos.HealthStats, error) {
	store, err := h.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	healthTK := storage.TK{
		Type: health.HealthStatusType,
		Key:  gatewayID,
	}
	healthBlob, err := store.Get(networkID, healthTK)
	if err != nil {
		store.Rollback()
		return nil, err
	}
	retHealth := &fegprotos.HealthStats{}
	err = protos.Unmarshal(healthBlob.Value, retHealth)
	if err != nil {
		store.Rollback()
		return retHealth, err
	}
	return retHealth, store.Commit()
}

// UpdateHealth updates the given gateway's health status in the
// blobstore.
func (h *healthBlobstore) UpdateHealth(networkID string, gatewayID string, healthStats *fegprotos.HealthStats) error {
	healthBlob, err := HealthToBlob(gatewayID, healthStats)
	if err != nil {
		return err
	}
	store, err := h.factory.StartTransaction(nil)
	if err != nil {
		return err
	}
	err = store.Write(networkID, blobstore.Blobs{healthBlob})
	if err != nil {
		store.Rollback()
		return err
	}
	return store.Commit()
}

// UpdateClusterState updates the given cluster's state in the
// blobstore.
func (h *healthBlobstore) UpdateClusterState(networkID string, clusterID string, logicalID string) error {
	clusterBlob, err := ClusterToBlob(clusterID, logicalID)
	if err != nil {
		return err
	}
	store, err := h.factory.StartTransaction(nil)
	if err != nil {
		return err
	}
	err = store.Write(networkID, blobstore.Blobs{clusterBlob})
	if err != nil {
		store.Rollback()
		return err
	}
	return store.Commit()
}

// GetClusterState retrieves the stored clusterState for the provided networkID
// and logicalID from the blobstore. The clusterState is
// initialized if it doesn't already exist.
func (h *healthBlobstore) GetClusterState(networkID string, logicalID string) (*fegprotos.ClusterState, error) {
	keys := []string{networkID}
	filter := blobstore.SearchFilter{
		NetworkID: &networkID,
	}
	store, err := h.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	foundKeys, err := store.GetExistingKeys(keys, filter)
	if err != nil {
		store.Rollback()
		return nil, err
	}
	if len(foundKeys) == 0 {
		err = h.initializeCluster(store, networkID, networkID, logicalID)
		if err != nil {
			store.Rollback()
			return nil, err
		}
	}
	clusterID := networkID
	clusterTK := storage.TK{
		Type: health.ClusterStatusType,
		Key:  clusterID,
	}
	clusterBlob, err := store.Get(networkID, clusterTK)
	if err != nil {
		store.Rollback()
		return nil, err
	}
	retClusterState := &fegprotos.ClusterState{}
	err = protos.Unmarshal(clusterBlob.Value, retClusterState)
	if err != nil {
		store.Rollback()
		return retClusterState, err
	}
	return retClusterState, store.Commit()
}

func (h *healthBlobstore) initializeCluster(store blobstore.Store, networkID string, clusterID string, logicalID string) error {
	glog.V(2).Infof("Initializing clusterState for networkID: %s with active: %s", networkID, logicalID)
	clusterBlob, err := ClusterToBlob(networkID, logicalID)
	if err != nil {
		return err
	}
	return store.Write(networkID, blobstore.Blobs{clusterBlob})
}
