/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	"orc8r/devmand/cloud/go/devmand"
)

// DEVICE HELPERS
// Get the necessary updates for agents when changing a device's managing agent
func GetAgentUpdates(dID, oldAgentID, newAgentID string) []configurator.EntityUpdateCriteria {
	ret := []configurator.EntityUpdateCriteria{}
	if oldAgentID != "" && newAgentID == "" {
		return getAgentRemoveUpdates(dID, oldAgentID)
	}
	if oldAgentID == "" && newAgentID != "" {
		return getAgentAddUpdates(dID, newAgentID)
	}
	if newAgentID != oldAgentID {
		ret = append(ret, getAgentRemoveUpdates(dID, oldAgentID)...)
		ret = append(ret, getAgentAddUpdates(dID, newAgentID)...)
	}
	return ret
}

func getAgentAddUpdates(dID, aID string) []configurator.EntityUpdateCriteria {
	return []configurator.EntityUpdateCriteria{
		configurator.EntityUpdateCriteria{
			Key:               aID,
			Type:              devmand.SymphonyAgentType,
			AssociationsToAdd: []storage.TypeAndKey{{Type: devmand.SymphonyDeviceType, Key: dID}},
		},
	}
}

func getAgentRemoveUpdates(dID, aID string) []configurator.EntityUpdateCriteria {
	return []configurator.EntityUpdateCriteria{
		configurator.EntityUpdateCriteria{
			Key:                  aID,
			Type:                 devmand.SymphonyAgentType,
			AssociationsToDelete: []storage.TypeAndKey{{Type: devmand.SymphonyDeviceType, Key: dID}},
		},
	}
}
