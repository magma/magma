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
	"sort"
	"strings"

	"magma/orc8r/cloud/go/storage"

	"github.com/thoas/go-funk"
)

// Blob encapsulates a blob for storage.
type Blob struct {
	Type    string
	Key     string
	Value   []byte
	Version uint64
}

type Blobs []Blob

// TKs converts blobs to their associated type and key.
func (bs Blobs) TKs() []storage.TypeAndKey {
	tks := make([]storage.TypeAndKey, 0, len(bs))
	for _, blob := range bs {
		tks = append(tks, storage.TypeAndKey{Type: blob.Type, Key: blob.Key})
	}
	return tks
}

// ByTK returns a computed view of a list of blobs as a map of
// blobs keyed by blob TypeAndKey.
func (bs Blobs) ByTK() map[storage.TypeAndKey]Blob {
	ret := make(map[storage.TypeAndKey]Blob, len(bs))
	for _, blob := range bs {
		ret[storage.TypeAndKey{Type: blob.Type, Key: blob.Key}] = blob
	}
	return ret
}

func (bs Blobs) Keys() []string {
	var keys []string
	for _, b := range bs {
		keys = append(keys, b.Key)
	}
	return keys
}

// CreateSearchFilter creates a search filter for the given criteria.
// Nil elements result in no filtering. If you prefer to instantiate string
// sets manually, you can also create a SearchFilter directly.
func CreateSearchFilter(networkID *string, types []string, keys []string, keyPrefix *string) SearchFilter {
	return SearchFilter{
		NetworkID: networkID,
		Types:     stringListToSet(types),
		Keys:      stringListToSet(keys),
		KeyPrefix: keyPrefix,
	}
}

func GetDefaultLoadCriteria() LoadCriteria {
	return LoadCriteria{LoadValue: true}
}

// BlobStorageFactory is an API to create a storage API bound to a transaction.
type BlobStorageFactory interface {
	InitializeFactory() error
	// StartTransaction opens a transaction for all following blob storage
	// operations, and returns a TransactionalBlobStorage instance tied to the
	// opened transaction.
	StartTransaction(opts *storage.TxOptions) (TransactionalBlobStorage, error)
}

// TransactionalBlobStorage is the client API for blob storage operations
// within the context of a transaction.
// TODO(4/9/2020): refactor Get-like methods into package-level defaults wrapping Search -- see e.g. ListKeysByNetwork
type TransactionalBlobStorage interface {
	// Commit commits the existing transaction.
	// If an error is returned from the backing storage while committing,
	// the transaction will be rolled back.
	Commit() error

	// Rollback rolls back the existing transaction.
	// If the targeted transaction has already been committed,
	// rolling back has no effect and returns an error.
	Rollback() error

	// Get loads a specific blob from storage.
	// If there is no blob matching the given ID, ErrNotFound from
	// magma/orc8r/lib/go/errors will be returned.
	Get(networkID string, id storage.TypeAndKey) (Blob, error)

	// GetMany loads and returns a collection of blobs matching the
	// specifiedIDs.
	// If there is no blob corresponding to a TypeAndKey, the returned list
	// will not have a corresponding Blob.
	GetMany(networkID string, ids []storage.TypeAndKey) (Blobs, error)

	// Search returns a filtered collection of blobs keyed by the network ID
	// to which they belong.
	// Blobs are filtered according to the search filter. Empty filter returns
	// all blobs. Blobs contents are loaded according to the load criteria.
	// Empty criteria loads all fields.
	Search(filter SearchFilter, criteria LoadCriteria) (map[string]Blobs, error)

	// CreateOrUpdate writes blobs to the storage.
	// Blobs are either updated in-place or created. The Version field of
	// blobs passed here will be used if it is not set to 0, otherwise version
	// incrementation will be handled internally inside the storage
	// implementation.
	CreateOrUpdate(networkID string, blobs Blobs) error

	// GetExistingKeys takes in a list of keys and returns a list of keys that
	// exist from the input.
	// The filter specifies whether to look at the entire storage or just in
	// a network.
	GetExistingKeys(keys []string, filter SearchFilter) ([]string, error)

	// Delete deletes specified blobs from storage.
	Delete(networkID string, ids []storage.TypeAndKey) error

	// IncrementVersion is an atomic upsert (INSERT DO ON CONFLICT) that
	// increments the version column or inserts 1 if it does not exist.
	IncrementVersion(networkID string, id storage.TypeAndKey) error
}

// GetAllOfType returns all blobs in the network of the passed type.
func GetAllOfType(store TransactionalBlobStorage, networkID, typ string) (Blobs, error) {
	filter := CreateSearchFilter(nil, []string{typ}, nil, nil)
	criteria := LoadCriteria{LoadValue: true}

	blobsByNetwork, err := store.Search(filter, criteria)
	if err != nil {
		return nil, err
	}

	return blobsByNetwork[networkID], nil
}

func ListKeys(store TransactionalBlobStorage, networkID string, typ string) ([]string, error) {
	filter := CreateSearchFilter(&networkID, []string{typ}, nil, nil)
	criteria := LoadCriteria{LoadValue: false}

	networkBlobs, err := store.Search(filter, criteria)
	if err != nil {
		return nil, err
	}
	return networkBlobs[networkID].Keys(), nil
}

// ListKeysByNetwork returns all blob keys, keyed by network ID.
func ListKeysByNetwork(store TransactionalBlobStorage) (map[string][]storage.TypeAndKey, error) {
	filter := CreateSearchFilter(nil, nil, nil, nil)
	criteria := LoadCriteria{LoadValue: false}

	blobsByNetwork, err := store.Search(filter, criteria)
	if err != nil {
		return nil, err
	}

	tks := map[string][]storage.TypeAndKey{}
	for network, blobs := range blobsByNetwork {
		tks[network] = blobs.TKs()
	}

	return tks, nil
}

// SearchFilter specifies search parameters.
// All fields are ANDed together in the final search that is performed.
type SearchFilter struct {
	// Optional network ID to search within
	NetworkID *string

	// Limit search to an OR matching any of the specified types
	Types map[string]bool
	// Limit search to an OR matching any of the specified keys
	// If the KeyPrefix of the search filter is specified, this argument will
	// be ignored by the blobstore.
	Keys map[string]bool
	// Prefix to match keys against. If this is specified (non-nil and non-
	// empty), the values of Keys will be ignored.
	KeyPrefix *string
}

// DoesTKMatch returns true if the given TK matches the search filter,
// false otherwise.
func (sf SearchFilter) DoesTKMatch(tk storage.TypeAndKey) bool {
	isTypesEmpty, isKeysEmpty, isPrefixEmpty := funk.IsEmpty(sf.Types), funk.IsEmpty(sf.Keys), funk.IsEmpty(sf.KeyPrefix)

	// Empty search filter matches everything
	if isTypesEmpty && isKeysEmpty && isPrefixEmpty {
		return true
	}

	if typeMatch := sf.Types[tk.Type]; !isTypesEmpty && !typeMatch {
		return false
	}

	// Key match: short-circuit if prefix is specified
	if !isPrefixEmpty {
		return strings.HasPrefix(tk.Key, *sf.KeyPrefix)
	}
	if keyMatch := sf.Keys[tk.Key]; !isKeysEmpty && !keyMatch {
		return false
	}
	return true
}

// GetTypes returns the types for this search filter sorted
func (sf SearchFilter) GetTypes() []string {
	ret := funk.Keys(sf.Types).([]string)
	sort.Strings(ret)
	return ret
}

// GetKeys returns the keys for this search filter sorted
func (sf SearchFilter) GetKeys() []string {
	ret := funk.Keys(sf.Keys).([]string)
	sort.Strings(ret)
	return ret
}

// LoadCriteria specifies which fields of each blob should be loaded from the
// underlying store.
// Returned blobs will contain type-default values for non-loaded fields.
type LoadCriteria struct {
	// LoadValue specifies whether to load the value of a blob.
	// Set to false to only load blob metadata.
	LoadValue bool
}

func stringListToSet(v []string) map[string]bool {
	ret := map[string]bool{}
	for _, s := range v {
		ret[s] = true
	}
	return ret
}
