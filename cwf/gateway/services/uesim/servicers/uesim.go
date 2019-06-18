/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/protos"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	networkIDPlaceholder = "magma"
	blobTypePlaceholder  = "uesim"
)

// UESimServer tracks all the UEs being simulated.
type UESimServer struct {
	store blobstore.BlobStorageFactory
}

// NewUESimServer initializes a UESimServer with an empty store map.
// Output: a new UESimServer
func NewUESimServer(factory blobstore.BlobStorageFactory) (*UESimServer, error) {
	return &UESimServer{
		store: factory,
	}, nil
}

// AddUE tries to add this UE to the server.
// Input: The UE data which will be added.
func (srv *UESimServer) AddUE(ctx context.Context, ue *cwfprotos.UEConfig) (ret *protos.Void, err error) {
	ret = &protos.Void{}

	err = validateUEData(ue)
	if err != nil {
		err = ConvertStorageErrorToGrpcStatus(err)
		return
	}
	blob, err := ueToBlob(ue)
	store, err := srv.store.StartTransaction()
	if err != nil {
		err = errors.Wrap(err, "Error while starting transaction")
		err = ConvertStorageErrorToGrpcStatus(err)
		return
	}
	defer func() {
		switch err {
		case nil:
			if commitErr := store.Commit(); commitErr != nil {
				err = errors.Wrap(err, "Error while committing transaction")
				err = ConvertStorageErrorToGrpcStatus(err)
			}
		default:
			if rollbackErr := store.Rollback(); rollbackErr != nil {
				err = errors.Wrap(err, "Error while rolling back transaction")
				err = ConvertStorageErrorToGrpcStatus(err)
			}
		}
	}()

	err = store.CreateOrUpdate(networkIDPlaceholder, []blobstore.Blob{blob})
	return
}

// Converts UE data to a blob for storage.
func ueToBlob(ue *cwfprotos.UEConfig) (blobstore.Blob, error) {
	marshaledUE, err := protos.Marshal(ue)
	if err != nil {
		return blobstore.Blob{}, err
	}
	return blobstore.Blob{
		Type:  blobTypePlaceholder,
		Key:   ue.GetImsi(),
		Value: marshaledUE,
	}, nil
}

// ConvertStorageErrorToGrpcStatus converts a UE error into a gRPC status error.
func ConvertStorageErrorToGrpcStatus(err error) error {
	if err == nil {
		return nil
	}
	return status.Errorf(codes.Unknown, err.Error())
}
