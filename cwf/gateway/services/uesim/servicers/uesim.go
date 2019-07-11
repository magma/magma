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
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
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
	op    []byte
	amf   []byte
}

// NewUESimServer initializes a UESimServer with an empty store map.
// Output: a new UESimServer
func NewUESimServer(factory blobstore.BlobStorageFactory) (*UESimServer, error) {
	// TODO use config to assign these values
	Op := []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11")
	Amf := []byte("\x67\x41")

	return &UESimServer{
		store: factory,
		op:    Op,
		amf:   Amf,
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

// Converts a blob back into a UE config
func blobToUE(blob blobstore.Blob) (*cwfprotos.UEConfig, error) {
	ue := &cwfprotos.UEConfig{}
	err := protos.Unmarshal(blob.Value, ue)
	if err != nil {
		return nil, err
	}
	return ue, nil
}

// getUE gets the UE with the specified IMSI from the blobstore.
func getUE(blobStoreFactory blobstore.BlobStorageFactory, imsi string) (ue *cwfprotos.UEConfig, err error) {
	store, err := blobStoreFactory.StartTransaction()
	if err != nil {
		err = errors.Wrap(err, "Error while starting transaction")
		return
	}
	defer func() {
		switch err {
		case nil:
			if commitErr := store.Commit(); commitErr != nil {
				err = errors.Wrap(err, "Error while committing transaction")
			}
		default:
			if rollbackErr := store.Rollback(); rollbackErr != nil {
				glog.Errorf("Error while rolling back transaction: %s", err)
			}
		}
	}()

	blob, err := store.Get(networkIDPlaceholder, storage.TypeAndKey{Type: blobTypePlaceholder, Key: imsi})
	if err != nil {
		err = errors.Wrap(err, "Error getting UE with specified IMSI")
		return
	}
	ue, err = blobToUE(blob)
	return
}

// ConvertStorageErrorToGrpcStatus converts a UE error into a gRPC status error.
func ConvertStorageErrorToGrpcStatus(err error) error {
	if err == nil {
		return nil
	}
	return status.Errorf(codes.Unknown, err.Error())
}
