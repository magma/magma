/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	"orc8r/wifi/cloud/go/wifi"
)

func GetMeshUpdates(gID string, oldMeshID, newMeshID MeshID) []configurator.EntityUpdateCriteria {
	ret := []configurator.EntityUpdateCriteria{}

	if oldMeshID != "" && newMeshID == "" {
		return getMeshRemoveUpdates(gID, oldMeshID)
	}

	if oldMeshID == "" && newMeshID != "" {
		return getMeshAddUpdates(gID, newMeshID)
	}

	if newMeshID != oldMeshID {
		ret = append(ret, getMeshRemoveUpdates(gID, oldMeshID)...)
		ret = append(ret, getMeshAddUpdates(gID, newMeshID)...)
	}

	return ret
}

func getMeshAddUpdates(gID string, meshID MeshID) []configurator.EntityUpdateCriteria {
	return []configurator.EntityUpdateCriteria{
		configurator.EntityUpdateCriteria{
			Key:               string(meshID),
			Type:              wifi.MeshEntityType,
			AssociationsToAdd: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
		},
	}
}

func getMeshRemoveUpdates(gID string, meshID MeshID) []configurator.EntityUpdateCriteria {
	return []configurator.EntityUpdateCriteria{
		configurator.EntityUpdateCriteria{
			Key:                  string(meshID),
			Type:                 wifi.MeshEntityType,
			AssociationsToDelete: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: gID}},
		},
	}
}
