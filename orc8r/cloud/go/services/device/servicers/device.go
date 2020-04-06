/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"context"
	"fmt"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/device/protos"
	commonProtos "magma/orc8r/lib/go/protos"

	"github.com/thoas/go-funk"
)

type deviceServicer struct {
	factory blobstore.BlobStorageFactory
}

func NewDeviceServicer(factory blobstore.BlobStorageFactory) (protos.DeviceServer, error) {
	if factory == nil {
		return nil, fmt.Errorf("Storage is nil")
	}
	return &deviceServicer{factory: factory}, nil
}

func (srv *deviceServicer) RegisterDevices(ctx context.Context, req *protos.RegisterOrUpdateDevicesRequest) (*commonProtos.Void, error) {
	void := &commonProtos.Void{}
	if err := ValidateRegisterDevicesRequest(req); err != nil {
		return void, err
	}

	blobs := protos.EntitiesToBlobs(req.GetEntities())
	keys := funk.Map(blobs, func(b blobstore.Blob) string { return b.Key }).([]string)
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	existingKeys, err := store.GetExistingKeys(keys, blobstore.SearchFilter{})
	if err != nil {
		store.Rollback()
		return void, err
	}
	if len(existingKeys) > 0 {
		store.Rollback()
		return nil, fmt.Errorf("The following keys: %v are already registered", existingKeys)
	}
	err = store.CreateOrUpdate(req.NetworkID, blobs)
	if err != nil {
		store.Rollback()
		return void, err
	}
	return void, store.Commit()
}

func (srv *deviceServicer) UpdateDevices(ctx context.Context, req *protos.RegisterOrUpdateDevicesRequest) (*commonProtos.Void, error) {
	void := &commonProtos.Void{}
	if err := ValidateRegisterDevicesRequest(req); err != nil {
		return void, err
	}

	blobs := protos.EntitiesToBlobs(req.GetEntities())
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	err = store.CreateOrUpdate(req.NetworkID, blobs)
	if err != nil {
		store.Rollback()
		return void, err
	}
	return void, store.Commit()
}

func (srv *deviceServicer) GetDeviceInfo(ctx context.Context, req *protos.GetDeviceInfoRequest) (*protos.GetDeviceInfoResponse, error) {
	response := &protos.GetDeviceInfoResponse{}
	if err := ValidateGetDeviceInfoRequest(req); err != nil {
		return response, err
	}

	ids := protos.DeviceIDsToTypeAndKey(req.DeviceIDs)
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		store.Rollback()
		return nil, err
	}
	blobs, err := store.GetMany(req.NetworkID, ids)
	if err != nil {
		store.Rollback()
		return response, err
	}
	response.DeviceMap = protos.BlobsToEntityByDeviceID(blobs)
	return response, store.Commit()
}

func (srv *deviceServicer) DeleteDevices(ctx context.Context, req *protos.DeleteDevicesRequest) (*commonProtos.Void, error) {
	void := &commonProtos.Void{}
	if err := ValidateDeleteDevicesRequest(req); err != nil {
		return void, err
	}

	ids := protos.DeviceIDsToTypeAndKey(req.DeviceIDs)
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	err = store.Delete(req.NetworkID, ids)
	if err != nil {
		store.Rollback()
		return void, err
	}
	return void, store.Commit()
}
