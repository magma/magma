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
	"magma/orc8r/cloud/go/protos"
	stateservice "magma/orc8r/cloud/go/services/state"

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

	store, err := srv.factory.StartTransaction()
	if err != nil {
		return nil, err
	}
	states, err := store.GetMany(req.GetNetworkID(), ids)
	store.Commit()
	return &protos.GetStatesResponse{States: protos.BlobsToStates(states)}, nil
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
	time := uint64(time.Now().UnixNano()) / uint64(time.Millisecond)

	states, err := addWrapperAndMakeBlobs(validatedStates, hwID, time, certExpiry)
	if err != nil {
		return response, err
	}

	store, err := srv.factory.StartTransaction()
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
	if err := ValidateDeleteStatesRequest(req); err != nil {
		return ret, err
	}
	networkID := req.GetNetworkID()
	ids := protos.StateIDsToTKs(req.GetIds())

	store, err := srv.factory.StartTransaction()
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

func addAdditionalInfo(state *protos.State, hwID string, time uint64, certExpiry int64) ([]byte, error) {
	wrap := stateservice.StateValue{
		ReporterID:         hwID,
		Time:               time,
		CertExpirationTime: certExpiry,
		ReportedValue:      state.Value,
	}
	return json.Marshal(wrap)
}

func addWrapperAndMakeBlobs(states []*protos.State, hwID string, time uint64, certExpiry int64) ([]blobstore.Blob, error) {
	blobs := []blobstore.Blob{}
	for _, state := range states {
		wrappedValue, err := addAdditionalInfo(state, hwID, time, certExpiry)
		if err != nil {
			return nil, err
		}
		state.Value = wrappedValue
		blobs = append(blobs, state.ToBlob())
	}
	return blobs, nil
}
