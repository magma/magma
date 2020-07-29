/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package models

import (
	"magma/devmand/cloud/go/devmand"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
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
