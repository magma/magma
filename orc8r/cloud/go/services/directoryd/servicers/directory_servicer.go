/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Servicer: implementation of the DirectoryService service

package servicers

import (
	"errors"
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DirectoryServicer struct {
	storage storage.DirectorydPersistenceService
}

func NewDirectoryServicer(storage storage.DirectorydPersistenceService) (*DirectoryServicer, error) {
	srv := &DirectoryServicer{storage: storage}
	return srv, nil
}

func (srv *DirectoryServicer) GetLocation(ctx context.Context, request *protos.GetLocationRequest) (*protos.LocationRecord, error) {
	if request == nil {
		return nil, errors.New("Empty GetLocationRequest")
	}
	glog.V(2).Infof("get location request: %v\n", request)
	if request.GetTable() == protos.TableID_HWID_TO_HOSTNAME {
		return srv.storage.GetRecord(request.Table, request.Id)
	}
	// Fetch IMSI directory records from state service
	directoryState, err := state.GetState(request.GetNetworkID(), orc8r.DirectoryRecordType, request.GetId())
	// In case of a legacy gateway, try retrieving using legacy db
	if err != nil {
		return srv.storage.GetRecord(request.Table, request.Id)
	}
	directoryRecord, ok := directoryState.ReportedState.(*directoryd.DirectoryRecord)
	if !ok || len(directoryRecord.LocationHistory) == 0 {
		return &protos.LocationRecord{}, fmt.Errorf("stored directory record is not properly formatted")
	}
	return &protos.LocationRecord{
		Location: directoryRecord.LocationHistory[0],
	}, nil
}

func (srv *DirectoryServicer) UpdateLocation(ctx context.Context, request *protos.UpdateDirectoryLocationRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if request == nil {
		return ret, errors.New("Empty UpdateLocationRequest")
	}
	if request.Table == protos.TableID_IMSI_TO_HWID {
		glog.V(2).Infof("update location request for IMSI->HWID table: %v\n", request)
		gw := protos.GetClientGateway(ctx)
		if gw == nil {
			return nil, status.Errorf(
				codes.PermissionDenied, "Missing Gateway Identity")
		}
		if !gw.Registered() {
			return nil, status.Errorf(
				codes.PermissionDenied, "Gateway is not registered")
		}
		request.Record = &protos.LocationRecord{Location: gw.HardwareId}
	}

	err := srv.storage.UpdateOrCreateRecord(request.Table, request.Id, request.Record)
	return ret, err
}

func (srv *DirectoryServicer) DeleteLocation(ctx context.Context, request *protos.DeleteLocationRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if request == nil {
		return nil, errors.New("Empty DeleteLocationRequest")
	}
	if request.GetTable() == protos.TableID_HWID_TO_HOSTNAME {
		return ret, srv.storage.DeleteRecord(request.Table, request.Id)
	}
	// Delete IMSI directory records from state service
	reqId := state.StateID{
		DeviceID: request.GetId(),
		Type:     orc8r.DirectoryRecordType,
	}
	err := state.DeleteStates(request.GetNetworkID(), []state.StateID{reqId})
	// In case of a legacy gateway, try retrieving using legacy db
	if err != nil {
		return ret, srv.storage.DeleteRecord(request.Table, request.Id)
	}
	return ret, err
}
