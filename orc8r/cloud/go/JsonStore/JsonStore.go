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
)

type StoreFactory interface {
	InitializeFactory() error

	StartTransaction(opts *storage.TxOptions) (Store, error)
}

type Store interface {

	Commit() error

	Rollback() error

	Search(filter SearchFilter, criteria LoadCriteria) (map[string]Jsons, error)

	Write(networkId string, json Jsons) error

	IncrementVersion(networkID string, id storage.TK) error

	Delete(networkID string, ids storage.TKs) error

	Get(networkID string, id storage.TK) (Json, error)

	GetMany(networkID string, ids storage.TKs) (Jsons, error)

	GetExistingKeys(keys []string, filter SearchFilter) ([]string, error)

}

func GetAllOfType(store Store, networkID string, typ string) (Jsons, error) {
	filter := CreateSearchFilter(nil, []string{typ}, nil, nil)
	criteria := LoadCriteria{LoadValue : true}

	JsonsByNetwork, err := store.Search(filter, criteria)
	if err != nil {
		return nil, err
	}

	return JsonsByNetwork[networkID], nil
}

func ListKeys(store Store, networkID string, typ string) ([]string, error) {
	filter := CreateSearchFilter(&networkID, []string{typ}, nil, nil)
	criteria := LoadCriteria{LoadValue: false}

	networkJsons, err := store.Search(filter, criteria)
	if err != nil {
		return nil, err
	}
	return networkJsons[networkID].keys(), nil
}

func ListKeysByNetwork(store Store) (map[string]storage.TKs, error) {
	filter := CreateSearchFilter(nil, nil, nil ,nil)
	criteria := LoadCriteria{LoadValue: false}

	JsonsByNetwork, err := store.Search(filter, criteria)
	if err != nil {
		return nil, err
	}

	tks := map[string]storage.TKs{}
	for network, Jsons := range JsonsByNetwork {
		tks[network] = Jsons.TKs()
	}
	return tks, nil
}