/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"encoding/json"
	"fmt"
	"time"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/protos"
	stateService "magma/orc8r/cloud/go/services/state"

	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type stateServicer struct {
	factory blobstore.BlobStorageFactory
}

// NewStateServicer returns a state server backed by storage passed in
func NewStateServicer(factory blobstore.BlobStorageFactory) (protos.StateServiceServer, error) {
	if factory == nil {
		return nil, fmt.Errorf("Storage factory is nil")
	}
	return &stateServicer{factory}, nil
}

// GetStates retrieves states from blobstorage
func (srv *stateServicer) GetStates(context context.Context, req *protos.GetStatesRequest) (*protos.GetStatesResponse, error) {
	if err := ValidateGetStatesRequest(req); err != nil {
		return nil, err
	}

	ids := protos.StateIDsToTKs(req.GetIds())

	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	states, err := store.GetMany(req.GetNetworkID(), ids)
	if err != nil {
		store.Rollback()
		return nil, err
	}
	return &protos.GetStatesResponse{States: protos.BlobsToStates(states)}, store.Commit()
}

// ReportStates saves states into blobstorage
func (srv *stateServicer) ReportStates(context context.Context, req *protos.ReportStatesRequest) (*protos.ReportStatesResponse, error) {
	response := &protos.ReportStatesResponse{}
	validatedStates, invalidStates, err := PartitionStatesBySerializability(req)
	if err != nil {
		return response, err
	}
	response.UnreportedStates = invalidStates

	// Get gateway information from context
	gw := protos.GetClientGateway(context)
	if gw == nil {
		return response, status.Errorf(codes.PermissionDenied, "Missing Gateway Identity")
	}
	if !gw.Registered() {
		return response, status.Errorf(codes.PermissionDenied, "Gateway is not registered")
	}
	hwID := gw.HardwareId
	networkID := gw.NetworkId
	certExpiry := protos.GetClientCertExpiration(context)
	timeMs := uint64(clock.Now().UnixNano()) / uint64(time.Millisecond)

	states, err := addWrapperAndMakeBlobs(validatedStates, hwID, timeMs, certExpiry)
	if err != nil {
		return response, err
	}

	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return response, err
	}
	err = store.CreateOrUpdate(networkID, states)
	if err != nil {
		store.Rollback()
		return response, err
	}
	return response, store.Commit()
}

// DeleteStates deletes states from blobstorage
func (srv *stateServicer) DeleteStates(context context.Context, req *protos.DeleteStatesRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if len(req.GetNetworkID()) == 0 {
		// Get gateway information from context
		gw := protos.GetClientGateway(context)
		if gw == nil {
			return ret, status.Errorf(codes.PermissionDenied, "Missing networkID and Gateway Identity")
		}
		if !gw.Registered() {
			return ret, status.Errorf(codes.PermissionDenied, "Missing networkID and Gateway is not registered")
		}
		req.NetworkID = gw.NetworkId
	}
	if err := ValidateDeleteStatesRequest(req); err != nil {
		return ret, err
	}
	networkID := req.GetNetworkID()
	ids := protos.StateIDsToTKs(req.GetIds())

	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	err = store.Delete(networkID, ids)
	if err != nil {
		store.Rollback()
		return ret, err
	}
	return ret, store.Commit()
}

// SyncStates retrieves states from blobstorage, compares their versions to
// the states included in the request, and returns the IDAndVersions that differ
func (srv *stateServicer) SyncStates(
	context context.Context,
	req *protos.SyncStatesRequest,
) (*protos.SyncStatesResponse, error) {
	response := &protos.SyncStatesResponse{}
	if err := ValidateSyncStatesRequest(req); err != nil {
		return response, err
	}
	// Get gateway information from context
	gw := protos.GetClientGateway(context)
	if gw == nil {
		return response, status.Errorf(codes.PermissionDenied, "Missing Gateway Identity")
	}
	if !gw.Registered() {
		return response, status.Errorf(codes.PermissionDenied, "Gateway is not registered")
	}
	networkID := gw.NetworkId

	tkIds := protos.StateIDAndVersionsToTKs(req.GetStates())
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return response, err
	}
	blobs, err := store.GetMany(networkID, tkIds)
	if err != nil {
		store.Rollback()
		return response, err
	}
	// pre-sort the blobstore results for faster syncing
	statesByDeviceID := map[string][]*protos.State{}
	for _, blob := range blobs {
		state := &protos.State{Type: blob.Type, DeviceID: blob.Key, Version: blob.Version}
		statesByDeviceID[state.DeviceID] = append(statesByDeviceID[state.DeviceID], state)
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
	return &protos.SyncStatesResponse{UnsyncedStates: unsyncedStates}, store.Commit()
}

func isStateSynced(deviceIdToStates map[string][]*protos.State, reqIdAndVersion *protos.IDAndVersion) (bool, uint64) {
	statesForDevice, ok := deviceIdToStates[reqIdAndVersion.Id.DeviceID]
	if !ok {
		return false, 0
	}
	for _, state := range statesForDevice {
		if state.Type == reqIdAndVersion.Id.Type && state.Version == reqIdAndVersion.Version {
			return true, 0
		} else if state.Type == reqIdAndVersion.Id.Type {
			return false, state.Version
		}
	}
	return false, 0
}

func wrapStateWithAdditionalInfo(state *protos.State, hwID string, time uint64, certExpiry int64) ([]byte, error) {
	wrap := stateService.SerializedStateWithMeta{
		ReporterID:              hwID,
		TimeMs:                  time,
		CertExpirationTime:      certExpiry,
		SerializedReportedState: state.Value,
	}
	return json.Marshal(wrap)
}

func addWrapperAndMakeBlobs(states []*protos.State, hwID string, timeMs uint64, certExpiry int64) ([]blobstore.Blob, error) {
	blobs := []blobstore.Blob{}
	for _, state := range states {
		wrappedValue, err := wrapStateWithAdditionalInfo(state, hwID, timeMs, certExpiry)
		if err != nil {
			return nil, err
		}
		state.Value = wrappedValue
		blobs = append(blobs, state.ToBlob())
	}
	return blobs, nil
}
