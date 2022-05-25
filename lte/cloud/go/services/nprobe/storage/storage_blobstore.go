/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package storage

import (
	"fmt"

	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
)

// NProbeBlobType is the blobstore type field for nprobe service
const NProbeBlobType = "nprobe"

// NewNProbeBlobstore returns a nprobe storage implementation
// backed by the provided blobstore factory.
func NewNProbeBlobstore(factory blobstore.StoreFactory) NProbeStorage {
	return &nprobeBlobStore{factory: factory}
}

type nprobeBlobStore struct {
	factory blobstore.StoreFactory
}

// StoreNProbeData stores current state for a given networkID and taskID
func (c *nprobeBlobStore) StoreNProbeData(networkID, taskID string, data models.NetworkProbeData) error {
	store, err := c.factory.StartTransaction(nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer store.Rollback()

	dataBlob, err := nprobeDataToBlob(taskID, data)
	if err != nil {
		return err
	}

	err = store.Write(networkID, blobstore.Blobs{dataBlob})
	if err != nil {
		return fmt.Errorf("failed to store nprobe data  %s: %w", taskID, err)
	}
	return store.Commit()
}

// GetNProbeData returns the state keyed by networkID and taskID
func (c *nprobeBlobStore) GetNProbeData(networkID, taskID string) (*models.NetworkProbeData, error) {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer store.Rollback()

	blob, err := store.Get(
		networkID,
		storage.TK{Type: NProbeBlobType, Key: taskID},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get nprobe data %s: %w", taskID, err)
	}

	data, err := nprobeDataFromBlob(blob)
	if err != nil {
		return nil, err
	}
	return &data, store.Commit()
}

// DeleteNProbeData returns the state keyed by networkID and taskID
func (c *nprobeBlobStore) DeleteNProbeData(networkID, taskID string) error {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer store.Rollback()

	err = store.Delete(
		networkID,
		storage.TKs{
			{Type: NProbeBlobType, Key: taskID},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to delete nprobe data %s: %w", taskID, err)
	}
	return store.Commit()
}

func nprobeDataToBlob(taskID string, data models.NetworkProbeData) (blobstore.Blob, error) {
	marshaledData, err := data.MarshalBinary()
	if err != nil {
		return blobstore.Blob{}, fmt.Errorf("Error marshaling NetworkProbeData: %w", err)
	}
	return blobstore.Blob{
		Type:  NProbeBlobType,
		Key:   taskID,
		Value: marshaledData,
	}, nil
}

func nprobeDataFromBlob(blob blobstore.Blob) (models.NetworkProbeData, error) {
	data := models.NetworkProbeData{}
	err := data.UnmarshalBinary(blob.Value)
	if err != nil {
		return models.NetworkProbeData{}, fmt.Errorf("Error unmarshaling NetworkProbeData: %w", err)
	}
	return data, nil
}
