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

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/directoryd/storage"

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
	record, err := srv.storage.GetRecord(request.Table, request.Id)
	return record, err
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

	srv.storage.UpdateOrCreateRecord(request.Table, request.Id, request.Record)
	return ret, nil
}

func (srv *DirectoryServicer) DeleteLocation(ctx context.Context, request *protos.DeleteLocationRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if request == nil {
		return nil, errors.New("Empty DeleteLocationRequest")
	}

	err := srv.storage.DeleteRecord(request.Table, request.Id)
	return ret, err
}
