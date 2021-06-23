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
	"encoding/json"
	"time"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/state/indexer/index"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errMissingGateway       = status.Error(codes.PermissionDenied, "missing gateway identity")
	errGatewayNotRegistered = status.Error(codes.PermissionDenied, "gateway not registered")
)

type stateServicer struct {
	factory blobstore.BlobStorageFactory
}

// NewStateServicer returns a state server backed by storage passed in.
func NewStateServicer(factory blobstore.BlobStorageFactory) (protos.StateServiceServer, error) {
	if factory == nil {
		return nil, errors.New("storage factory is nil")
	}
	return &stateServicer{factory}, nil
}

func (srv *stateServicer) GetStates(ctx context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
	if err := validateGetStatesRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if !funk.IsEmpty(req.Ids) {
		return srv.getStates(ctx, req)
	}
	return srv.searchStates(ctx, req)
}

// ReportStates from a gateway.
// Always reports UnreportedStates as empty.
func (srv *stateServicer) ReportStates(ctx context.Context, req *protos.ReportStatesRequest) (*protos.ReportStatesResponse, error) {
	if err := validateReportStatesRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Get gateway information from context
	gw := protos.GetClientGateway(ctx)
	if gw == nil {
		return nil, errMissingGateway
	}
	if !gw.Registered() {
		return nil, errGatewayNotRegistered
	}
	hwID := gw.HardwareId
	networkID := gw.NetworkId
	certExpiry := protos.GetClientCertExpiration(ctx)
	timeMs := uint64(clock.Now().UnixNano()) / uint64(time.Millisecond)

	states, err := addWrapperAndMakeBlobs(req.States, hwID, timeMs, certExpiry)
	if err != nil {
		return nil, internalErr(err, "ReportStates convert to blobs")
	}

	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, internalErr(err, "ReportStates blobstore start transaction")
	}
	err = store.CreateOrUpdate(networkID, states)
	if err != nil {
		_ = store.Rollback()
		return nil, internalErr(err, "ReportStates blobstore create or update")
	}
	err = store.Commit()
	if err != nil {
		return nil, internalErr(err, "ReportStates blobstore commit transaction")
	}

	byID, err := state_types.MakeSerializedStatesByID(req.States)
	if err != nil {
		return nil, internalErr(err, "ReportStates make states by ID")
	}
	go index.MustIndex(networkID, byID)

	return &protos.ReportStatesResponse{}, nil
}

func (srv *stateServicer) DeleteStates(ctx context.Context, req *protos.DeleteStatesRequest) (*protos.Void, error) {
	if len(req.GetNetworkID()) == 0 {
		// Get gateway information from context
		gw := protos.GetClientGateway(ctx)
		if gw == nil {
			return nil, status.Error(codes.PermissionDenied, "missing network and missing gateway identity")
		}
		if !gw.Registered() {
			return nil, status.Error(codes.PermissionDenied, "missing network and gateway not registered")
		}
		req.NetworkID = gw.NetworkId
	}
	if err := validateDeleteStatesRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	networkID := req.GetNetworkID()
	ids := idsToTKs(req.GetIds())

	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, internalErr(err, "DeleteStates blobstore start transaction")
	}
	err = store.Delete(networkID, ids)
	if err != nil {
		_ = store.Rollback()
		return nil, internalErr(err, "DeleteStates blobstore delete")
	}
	err = store.Commit()
	if err != nil {
		return nil, internalErr(err, "DeleteStates blobstore commit transaction")
	}

	return &protos.Void{}, nil
}

func (srv *stateServicer) SyncStates(ctx context.Context, req *protos.SyncStatesRequest) (*protos.SyncStatesResponse, error) {
	if err := validateSyncStatesRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	// Get gateway information from context
	gw := protos.GetClientGateway(ctx)
	if gw == nil {
		return nil, errMissingGateway
	}
	if !gw.Registered() {
		return nil, errGatewayNotRegistered
	}
	networkID := gw.NetworkId

	tkIds := idAndVersionsToTKs(req.GetStates())
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, internalErr(err, "SyncStates blobstore start transaction")
	}
	blobs, err := store.GetMany(networkID, tkIds)
	if err != nil {
		_ = store.Rollback()
		return nil, internalErr(err, "SyncStates blobstore get many")
	}
	// Pre-sort the blobstore results for faster syncing
	statesByDeviceID := map[string][]*protos.State{}
	for _, blob := range blobs {
		st := &protos.State{Type: blob.Type, DeviceID: blob.Key, Version: blob.Version}
		statesByDeviceID[st.DeviceID] = append(statesByDeviceID[st.DeviceID], st)
	}
	var unsyncedStates []*protos.IDAndVersion
	for _, reqIdAndVersion := range req.GetStates() {
		isStateSynced, unsyncedVersion := isStateSynced(statesByDeviceID, reqIdAndVersion)
		if isStateSynced {
			continue
		}
		unsyncedState := &protos.IDAndVersion{
			Id:      reqIdAndVersion.Id,
			Version: unsyncedVersion,
		}
		unsyncedStates = append(unsyncedStates, unsyncedState)
	}
	err = store.Commit()
	if err != nil {
		return nil, internalErr(err, "SyncStates blobstore commit transaction")
	}

	return &protos.SyncStatesResponse{UnsyncedStates: unsyncedStates}, nil
}

func (srv *stateServicer) getStates(_ context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, internalErr(err, "GetStates (get) blobstore start transaction")
	}

	ids := idsToTKs(req.GetIds())
	blobs, err := store.GetMany(req.GetNetworkID(), ids)
	if err != nil {
		_ = store.Rollback()
		return nil, internalErr(err, "GetStates (get) blobstore get many")
	}

	err = store.Commit()
	if err != nil {
		return nil, internalErr(err, "GetStates (get) blobstore commit transaction")
	}

	return &protos.GetStatesResponse{States: blobsToStates(blobs)}, nil
}

func (srv *stateServicer) searchStates(_ context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
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

	return &protos.GetStatesResponse{States: blobsToStates(searchResults[req.NetworkID])}, nil
}

func isStateSynced(deviceIdToStates map[string][]*protos.State, reqIdAndVersion *protos.IDAndVersion) (bool, uint64) {
	statesForDevice, ok := deviceIdToStates[reqIdAndVersion.Id.DeviceID]
	if !ok {
		return false, 0
	}
	for _, st := range statesForDevice {
		if st.Type == reqIdAndVersion.Id.Type && st.Version == reqIdAndVersion.Version {
			return true, 0
		} else if st.Type == reqIdAndVersion.Id.Type {
			return false, st.Version
		}
	}
	return false, 0
}

func wrapStateWithAdditionalInfo(st *protos.State, hwID string, time uint64, certExpiry int64) ([]byte, error) {
	wrap := state_types.SerializedState{
		ReporterID:              hwID,
		TimeMs:                  time,
		SerializedReportedState: st.Value,
	}
	ret, err := json.Marshal(wrap)
	if err != nil {
		return nil, errors.Wrap(err, "json marshal state with meta")
	}
	return ret, nil
}

func addWrapperAndMakeBlobs(states []*protos.State, hwID string, timeMs uint64, certExpiry int64) (blobstore.Blobs, error) {
	var blobs blobstore.Blobs
	for _, st := range states {
		wrappedValue, err := wrapStateWithAdditionalInfo(st, hwID, timeMs, certExpiry)
		if err != nil {
			return nil, err
		}
		st.Value = wrappedValue
		blobs = append(blobs, stateToBlob(st))
	}
	return blobs, nil
}

func idToTK(id *protos.StateID) storage.TypeAndKey {
	return storage.TypeAndKey{Type: id.GetType(), Key: id.GetDeviceID()}
}

func idsToTKs(ids []*protos.StateID) []storage.TypeAndKey {
	var tks []storage.TypeAndKey
	for _, id := range ids {
		tks = append(tks, idToTK(id))
	}
	return tks
}

func idAndVersionsToTKs(IDs []*protos.IDAndVersion) []storage.TypeAndKey {
	var ids []storage.TypeAndKey
	for _, idAndVersion := range IDs {
		ids = append(ids, idToTK(idAndVersion.Id))
	}
	return ids
}

func blobsToStates(blobs blobstore.Blobs) []*protos.State {
	var states []*protos.State
	for _, b := range blobs {
		st := &protos.State{
			Type:     b.Type,
			DeviceID: b.Key,
			Value:    b.Value,
			Version:  b.Version,
		}
		states = append(states, st)
	}
	return states
}

func stateToBlob(state *protos.State) blobstore.Blob {
	return blobstore.Blob{
		Type:    state.GetType(),
		Key:     state.GetDeviceID(),
		Value:   state.GetValue(),
		Version: state.GetVersion(),
	}
}

func internalErr(err error, wrap string) error {
	e := errors.Wrap(err, wrap)
	return status.Error(codes.Internal, e.Error())
}
