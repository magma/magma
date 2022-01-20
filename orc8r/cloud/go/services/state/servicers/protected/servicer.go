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

package servicers

import (
	"context"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/state/protos"
	common "magma/orc8r/cloud/go/services/state/servicers"
	"magma/orc8r/cloud/go/storage"
	external_protos "magma/orc8r/lib/go/protos"
)

var (
	errMissingGateway       = status.Error(codes.PermissionDenied, "missing gateway identity")
	errGatewayNotRegistered = status.Error(codes.PermissionDenied, "gateway not registered")
)

type stateInternalServicer struct {
	factory blobstore.StoreFactory
}

// NewStateServicer returns a state server backed by storage passed in.
func NewStateServicer(factory blobstore.StoreFactory) (protos.StateInternalServiceServer, error) {
	if factory == nil {
		return nil, errors.New("storage factory is nil")
	}
	return &stateInternalServicer{factory}, nil
}

func (srv *stateInternalServicer) GetStates(ctx context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
	if err := validateGetStatesRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !funk.IsEmpty(req.Ids) {
		return srv.getStates(ctx, req)
	}
	return srv.searchStates(ctx, req)
}

func (srv *stateInternalServicer) getStates(_ context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, common.InternalErr(err, "GetStates (get) blobstore start transaction")
	}

	ids := idsToTKs(req.GetIds())
	blobs, err := store.GetMany(req.GetNetworkID(), ids)
	if err != nil {
		_ = store.Rollback()
		return nil, common.InternalErr(err, "GetStates (get) blobstore get many")
	}

	err = store.Commit()
	if err != nil {
		return nil, common.InternalErr(err, "GetStates (get) blobstore commit transaction")
	}

	return &protos.GetStatesResponse{States: blobsToStates(blobs)}, nil
}

func (srv *stateInternalServicer) searchStates(_ context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, common.InternalErr(err, "GetStates (search) blobstore start transaction")
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
		return nil, common.InternalErr(err, "GetStates (search) blobstore search")
	}

	err = store.Commit()
	if err != nil {
		return nil, common.InternalErr(err, "GetStates (search) blobstore commit transaction")
	}

	return &protos.GetStatesResponse{States: blobsToStates(searchResults[req.NetworkID])}, nil
}

func idToTK(id *external_protos.StateID) storage.TK {
	return storage.TK{Type: id.GetType(), Key: id.GetDeviceID()}
}

func idsToTKs(ids []*external_protos.StateID) storage.TKs {
	var tks storage.TKs
	for _, id := range ids {
		tks = append(tks, idToTK(id))
	}
	return tks
}

func blobsToStates(blobs blobstore.Blobs) []*external_protos.State {
	var states []*external_protos.State
	for _, b := range blobs {
		st := &external_protos.State{
			Type:     b.Type,
			DeviceID: b.Key,
			Value:    b.Value,
			Version:  b.Version,
		}
		states = append(states, st)
	}
	return states
}
