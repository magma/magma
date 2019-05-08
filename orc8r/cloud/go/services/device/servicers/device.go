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
	commonProtos "magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/device/protos"
)

type deviceServicer struct {
	storage blobstore.BlobStorageFactory
}

func NewDeviceServicer(storage blobstore.BlobStorageFactory) (*deviceServicer, error) {
	if storage == nil {
		return nil, fmt.Errorf("Storage is nil")
	}
	return &deviceServicer{storage: storage}, nil
}
func (srv *deviceServicer) RegisterDevices(ctx context.Context, req *protos.RegisterDevicesRequest) (*commonProtos.Void, error) {
	return &commonProtos.Void{}, nil
}

func (srv *deviceServicer) GetDeviceInfo(ctx context.Context, req *protos.GetDeviceInfoRequest) (*protos.GetDeviceInfoResponse, error) {
	return &protos.GetDeviceInfoResponse{}, fmt.Errorf("GetDeviceInfo not yet implemented")
}
