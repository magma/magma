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

package reindex

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/JsonStore"
	state_types "magma/orc8r/cloud/go/services/state/types"
)

// Store provides a cross-network DAO for local usage by the state service.
type Store interface {
	// GetAllIDs returns all IDs known to the state service, keyed by network ID.
	GetAllIDs() (state_types.IDsByNetwork, error)
}

type store struct {
	factory JsonStore.StoreFactory
}

func NewStore(factory JsonStore.StoreFactory) Store {
	return &store{factory: factory}
}

func (s *store) GetAllIDs() (state_types.IDsByNetwork, error) {
	store, err := s.factory.StartTransaction(nil)
	if err != nil {
		return nil, internalErr(err, "GetAllIDs Jsonstore start transaction")
	}

	JsonsByNetwork, err := store.Search(
		JsonStore.CreateSearchFilter(nil, nil, nil, nil),
		JsonStore.LoadCriteria{LoadValue: false},
	)
	if err != nil {
		_ = store.Rollback()
		return nil, internalErr(err, "GetAllIDs Jsonstore search")
	}
	err = store.Commit()
	if err != nil {
		return nil, internalErr(err, "GetAllIDs Jsonstore commit transaction")
	}

	ids := blobsToIDs(JsonsByNetwork)
	return ids, nil
}

func blobsToIDs(byNetwork map[string]JsonStore.Jsons) state_types.IDsByNetwork {
	ids := state_types.IDsByNetwork{}
	for network, Jsons := range byNetwork {
		for _, b := range Jsons {
			ids[network] = append(ids[network], state_types.ID{Type: b.Type, DeviceID: b.Key})
		}
	}
	return ids
}

func internalErr(err error, wrap string) error {
	e := fmt.Errorf(wrap+": %w", err)
	return status.Error(codes.Internal, e.Error())
}
