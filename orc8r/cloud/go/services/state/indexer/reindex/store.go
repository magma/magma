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
	"magma/orc8r/cloud/go/blobstore"
	state_types "magma/orc8r/cloud/go/services/state/types"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Store provides a cross-network DAO for local usage by the state service.
type Store interface {
	// GetAllIDs returns all IDs known to the state service, keyed by network ID.
	GetAllIDs() (state_types.IDsByNetwork, error)
}

type storeImpl struct {
	factory blobstore.BlobStorageFactory
}

func NewStore(factory blobstore.BlobStorageFactory) Store {
	return &storeImpl{factory: factory}
}

func (s *storeImpl) GetAllIDs() (state_types.IDsByNetwork, error) {
	store, err := s.factory.StartTransaction(nil)
	if err != nil {
		return nil, internalErr(err, "GetAllIDs blobstore start transaction")
	}

	blobsByNetwork, err := store.Search(
		blobstore.CreateSearchFilter(nil, nil, nil, nil),
		blobstore.LoadCriteria{LoadValue: false},
	)
	if err != nil {
		_ = store.Rollback()
		return nil, internalErr(err, "GetAllIDs blobstore search")
	}
	err = store.Commit()
	if err != nil {
		return nil, internalErr(err, "GetAllIDs blobstore commit transaction")
	}

	ids := blobsToIDs(blobsByNetwork)
	return ids, nil
}

func blobsToIDs(byNetwork map[string]blobstore.Blobs) state_types.IDsByNetwork {
	ids := state_types.IDsByNetwork{}
	for network, blobs := range byNetwork {
		for _, b := range blobs {
			ids[network] = append(ids[network], state_types.ID{Type: b.Type, DeviceID: b.Key})
		}
	}
	return ids
}

func internalErr(err error, wrap string) error {
	e := errors.Wrap(err, wrap)
	return status.Error(codes.Internal, e.Error())
}
