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

	"github.com/thoas/go-funk"

	"magma/orc8r/cloud/go/storage"
)

// Blob encapsulates a blob for storage.
type Blob struct {
	Type    string
	Key     string
	Value   []byte
	Version uint64
}

// TK converts a blob to its associated type and key.
func (b Blob) TK() storage.TK {
	return storage.TK{Type: b.Type, Key: b.Key}
}

type Blobs []Blob

// TKs converts blobs to their associated type and key.
func (bs Blobs) TKs() storage.TKs {
	tks := make(storage.TKs, 0, len(bs))
	for _, blob := range bs {
		tks = append(tks, storage.TK{Type: blob.Type, Key: blob.Key})
	}
	return tks
}

// ByTK returns a computed view of a list of blobs as a map of
// blobs keyed by blob TK.
func (bs Blobs) ByTK() map[storage.TK]Blob {
	ret := make(map[storage.TK]Blob, len(bs))
	for _, blob := range bs {
		ret[storage.TK{Type: blob.Type, Key: blob.Key}] = blob
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
		Types:     toMap(types),
		Keys:      toMap(keys),
		KeyPrefix: keyPrefix,
	}
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
func (sf SearchFilter) DoesTKMatch(tk storage.TK) bool {
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

func GetDefaultLoadCriteria() LoadCriteria {
	return LoadCriteria{LoadValue: true}
}

// LoadCriteria specifies which fields of each blob should be loaded from the
// underlying store.
// Returned blobs will contain type-default values for non-loaded fields.
type LoadCriteria struct {
	// LoadValue specifies whether to load the value of a blob.
	// Set to false to only load blob metadata.
	LoadValue bool
}

func toMap(v []string) map[string]bool {
	present := map[string]bool{}
	for _, s := range v {
		present[s] = true
	}
	return present
}
