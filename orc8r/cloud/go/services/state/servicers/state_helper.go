/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */
package servicers

import (
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
)

func BlobsToStates(blobs []blobstore.Blob) []*protos.State {
	states := make([]*protos.State, 0, len(blobs))
	for _, blob := range blobs {
		states = append(states, &protos.State{Type: blob.Type, DeviceID: blob.Key,
			Value: blob.Value, Version: blob.Version})
	}
	return states
}

func ToBlob(state *protos.State) blobstore.Blob {
	return blobstore.Blob{
		Type:    state.GetType(),
		Key:     state.GetDeviceID(),
		Value:   state.GetValue(),
		Version: state.GetVersion(),
	}
}

func StateIDsToTKs(IDs []*protos.StateID) []storage.TypeAndKey {
	ids := []storage.TypeAndKey{}
	for _, id := range IDs {
		ids = append(ids, toStorageTK(id))
	}
	return ids
}

func StateIDAndVersionsToTKs(IDs []*protos.IDAndVersion) []storage.TypeAndKey {
	ids := []storage.TypeAndKey{}
	for _, idAndVersion := range IDs {
		ids = append(ids, toStorageTK(idAndVersion.Id))
	}
	return ids
}

func toStorageTK(id *protos.StateID) storage.TypeAndKey {
	return storage.TypeAndKey{Type: id.GetType(), Key: id.GetDeviceID()}
}
