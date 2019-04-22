/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 */
package protos

import (
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
)

func BlobsToStates(blobs []blobstore.Blob) []*State {
	states := make([]*State, 0, len(blobs))
	for _, blob := range blobs {
		states = append(states, &State{Type: blob.Type, DeviceID: blob.Key, Value: blob.Value})
	}
	return states
}

func StatesToBlobs(states []*State) []blobstore.Blob {
	blobs := []blobstore.Blob{}
	for _, state := range states {
		blobs = append(blobs, blobstore.Blob{Type: state.GetType(), Key: state.GetDeviceID(), Value: state.GetValue()})
	}
	return blobs
}

func StateIDsToTKs(IDs []*StateID) []storage.TypeAndKey {
	ids := []storage.TypeAndKey{}
	for _, id := range IDs {
		ids = append(ids, toStorageTK(id))
	}
	return ids
}

func toStorageTK(id *StateID) storage.TypeAndKey {
	return storage.TypeAndKey{Type: id.GetType(), Key: id.GetDeviceID()}
}
