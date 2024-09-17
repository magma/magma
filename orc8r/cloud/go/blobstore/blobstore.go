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

package blobstore

import (
	"magma/orc8r/cloud/go/storage"
)

// StoreFactory is an API to create a storage API bound to a transaction.
type StoreFactory interface {
	InitializeFactory() error
	// StartTransaction opens a transaction for all following blob storage
	// operations, and returns a Store instance tied to the
	// opened transaction.
	StartTransaction(opts *storage.TxOptions) (Store, error)
}

// Store is the client API for blob storage operations
// within the context of a transaction.
type Store interface {
	// Commit commits the existing transaction.
	// If an error is returned from the backing storage while committing,
	// the transaction will be rolled back.
	Commit() error

	// Rollback rolls back the existing transaction.
	// If the targeted transaction has already been committed,
	// rolling back has no effect and returns an error.
	Rollback() error

	// Search returns a filtered collection of blobs keyed by the network ID
	// to which they belong.
	// Blobs are filtered according to the search filter. Empty filter returns
	// all blobs. Blobs contents are loaded according to the load criteria.
	// Empty criteria loads all fields.
	Search(filter SearchFilter, criteria LoadCriteria) (map[string]Blobs, error)

	// Write blobs to the storage.
	// Blobs are either updated in-place or created. The Version field of
	// blobs passed here will be used if it is not set to 0, otherwise version
	// incrementation will be handled internally inside the storage
	// implementation.
	Write(networkID string, blobs Blobs) error

	// IncrementVersion is an atomic upsert (INSERT DO ON CONFLICT) that
	// increments the version column or inserts 1 if it does not exist.
	IncrementVersion(networkID string, id storage.TK) error

	// Delete deletes specified blobs from storage.
	Delete(networkID string, ids storage.TKs) error

	// TODO(4/9/2020): refactor Get-like methods into package-level defaults wrapping Search -- see e.g. ListKeysByNetwork

	// Get loads a specific blob from storage.
	// If there is no blob matching the given ID, ErrNotFound from
	// magma/orc8r/lib/go/merrors will be returned.
	Get(networkID string, id storage.TK) (Blob, error)

	// GetMany loads and returns a collection of blobs matching the
	// specifiedIDs.
	// If there is no blob corresponding to a TK, the returned list
	// will not have a corresponding Blob.
	GetMany(networkID string, ids storage.TKs) (Blobs, error)

	// GetExistingKeys takes in a list of keys and returns a list of keys that
	// exist from the input.
	// The filter specifies whether to look at the entire storage or just in
	// a network.
	GetExistingKeys(keys []string, filter SearchFilter) ([]string, error)
}

// GetAllOfType returns all blobs in the network of the passed type.
func GetAllOfType(store Store, networkID, typ string) (Blobs, error) {
	filter := CreateSearchFilter(nil, []string{typ}, nil, nil)
	criteria := LoadCriteria{LoadValue: true}

	blobsByNetwork, err := store.Search(filter, criteria)
	if err != nil {
		return nil, err
	}

	return blobsByNetwork[networkID], nil
}

func ListKeys(store Store, networkID string, typ string) ([]string, error) {
	filter := CreateSearchFilter(&networkID, []string{typ}, nil, nil)
	criteria := LoadCriteria{LoadValue: false}

	networkBlobs, err := store.Search(filter, criteria)
	if err != nil {
		return nil, err
	}
	return networkBlobs[networkID].Keys(), nil
}

// ListKeysByNetwork returns all blob keys, keyed by network ID.
func ListKeysByNetwork(store Store) (map[string]storage.TKs, error) {
	filter := CreateSearchFilter(nil, nil, nil, nil)
	criteria := LoadCriteria{LoadValue: false}

	blobsByNetwork, err := store.Search(filter, criteria)
	if err != nil {
		return nil, err
	}

	tks := map[string]storage.TKs{}
	for network, blobs := range blobsByNetwork {
		tks[network] = blobs.TKs()
	}

	return tks, nil
}
