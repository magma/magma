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

package JsonStore

import (
	"magma/orc8r/cloud/go/storage"
	"sort"
	"strings"

	"github.com/thoas/go-funk"
)

//Json will carry a string for Storage.
type Json struct {
	Type  string 
	Key   string 
	Value string
	Version uint64
}

func (j Json) TK() storage.TK {
	return storage.TK{Type: j.Type, Key: j.Key}
}

type Jsons []Json

func (js Jsons) TKs() storage.TKs {
	tks := make(storage.TKs, 0, len(js))
	for _, json := range js {
		tks = append(tks, storage.TK{Type: json.Type, Key: json.Key})
	}
	return tks
}
// ByTK returns a computed view of a list of blobs as a map of
// blobs keyed by blob TK.
func (js Jsons) ByTK() map[storage.TK]Json {
	ret := make(map[storage.TK]Json, len(js))
	for _, blob := range js {
		ret[storage.TK{Type: blob.Type, Key: blob.Key}] = blob
	}
	return ret
}

func (js Jsons) keys() []string {
	var keys []string
	for _, b := range js {
		keys = append(keys, b.Key)

	}
	return keys
}

func CreateSearchFilter(networkID *string, types []string, keys []string, keyPrefix *string) SearchFilter {
	return SearchFilter{
		NetworkID: networkID,
		Types:     toMap(types),
		Keys:      toMap(keys),
		KeyPrefix: keyPrefix,
	}
}

type SearchFilter struct {
	NetworkID *string
	
	Types map[string]bool

	Keys map[string]bool

	KeyPrefix *string
}

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

func (sf SearchFilter) GetTypes() []string {
	ret := funk.Keys(sf.Types).([]string)
	sort.Strings(ret)
	return ret
}

func (sf SearchFilter) GetKeys() []string {
	ret := funk.Keys(sf.Keys).([]string)
	sort.Strings(ret)
	return ret
}

func GetDefaultLoadCriteria() LoadCriteria {
	return LoadCriteria{LoadValue: true}
}

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