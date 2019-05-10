/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos

import (
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
)

// EntitiesToBlobs maps a list of PhysicalEntity to a list of blobstore.Blob
// by using each entity's networkID as blob value
func EntitiesToBlobs(entities []*PhysicalEntity) []blobstore.Blob {
	blobs := []blobstore.Blob{}
	for _, entity := range entities {
		blobs = append(blobs, entityToBlob(entity))
	}
	return blobs
}

// DeviceIDsToTypeAndKey maps a list of DeviceID to a list of
// storage.TypeAndKey.
func DeviceIDsToTypeAndKey(deviceIDs []*DeviceID) []storage.TypeAndKey {
	tks := []storage.TypeAndKey{}
	for _, id := range deviceIDs {
		tks = append(tks, deviceIDToTK(id))
	}
	return tks
}

// BlobsToEntityByDeviceID maps a list of blobstore.Blob to map[deviceID]PhysicalEntity
func BlobsToEntityByDeviceID(entities []blobstore.Blob) map[string]*PhysicalEntity {
	ret := map[string]*PhysicalEntity{}
	for _, blob := range entities {
		ret[blob.Key] = blobToEntity(blob)
	}
	return ret
}

func entityToBlob(entity *PhysicalEntity) blobstore.Blob {
	return blobstore.Blob{
		Key:   entity.GetDeviceID(),
		Type:  entity.GetType(),
		Value: entity.GetInfo(),
	}
}

func blobToEntity(blob blobstore.Blob) *PhysicalEntity {
	return &PhysicalEntity{
		Type:     blob.Type,
		DeviceID: blob.Key,
		Info:     blob.Value,
	}
}

func deviceIDToTK(id *DeviceID) storage.TypeAndKey {
	return storage.TypeAndKey{Type: id.Type, Key: id.DeviceID}
}
