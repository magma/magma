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

package protected

import (
	"context"
	"errors"

	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/state"
	servicers "magma/orc8r/cloud/go/services/state/servicers/southbound"
	"magma/orc8r/lib/go/protos"
)

type cloudStateServicer struct {
	factory blobstore.StoreFactory
}

// NewCloudStateServicer returns a state server backed by storage passed in.
func NewCloudStateServicer(factory blobstore.StoreFactory) (protos.CloudStateServiceServer, error) {
	if factory == nil {
		return nil, errors.New("storage factory is nil")
	}
	return &cloudStateServicer{factory}, nil
}

func (srv *cloudStateServicer) GetStates(ctx context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
	if err := servicers.ValidateGetStatesRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !funk.IsEmpty(req.Ids) {
		return srv.getStates(ctx, req)
	}
	return srv.searchStates(ctx, req)
}

func (srv *cloudStateServicer) getStates(_ context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, internalErr(err, "GetStates (get) blobstore start transaction")
	}

	ids := state.IdsToTKs(req.GetIds())
	blobs, err := store.GetMany(req.GetNetworkID(), ids)
	if err != nil {
		_ = store.Rollback()
		return nil, internalErr(err, "GetStates (get) blobstore get many")
	}

	err = store.Commit()
	if err != nil {
		return nil, internalErr(err, "GetStates (get) blobstore commit transaction")
	}

	return &protos.GetStatesResponse{States: state.BlobsToStates(blobs)}, nil
}

func (srv *cloudStateServicer) searchStates(_ context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, internalErr(err, "GetStates (search) blobstore start transaction")
	}

	var idPrefix *string
	if req.IdPrefix != "" {
		idPrefix = &req.IdPrefix
	}
	searchResults, err := store.Search(
		blobstore.CreateSearchFilter(&req.NetworkID, req.TypeFilter, req.IdFilter, idPrefix),
		blobstore.LoadCriteria{LoadValue: req.LoadValues},
	)
	if err != nil {
		_ = store.Rollback()
		return nil, internalErr(err, "GetStates (search) blobstore search")
	}

	err = store.Commit()
	if err != nil {
		return nil, internalErr(err, "GetStates (search) blobstore commit transaction")
	}

	return &protos.GetStatesResponse{States: state.BlobsToStates(searchResults[req.NetworkID])}, nil
}
