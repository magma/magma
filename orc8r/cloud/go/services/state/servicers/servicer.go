/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"fmt"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/protos"

	"golang.org/x/net/context"
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
func (srv *stateServicer) ReportStates(context context.Context, req *protos.ReportStatesRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if err := ValidateReportStatesRequest(req); err != nil {
		return nil, err
	}

	networkID, err := identity.GetClientNetworkID(context)
	if err != nil {
		return nil, err
	}
	states := protos.StatesToBlobs(req.GetStates())

	store, err := srv.factory.StartTransaction()
	if err != nil {
		return nil, err
	}
	err = store.CreateOrUpdate(networkID, states)
	if err != nil {
		store.Rollback()
		return ret, err
	}
	return ret, store.Commit()
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
