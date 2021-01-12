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

package test_utils

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"

	"github.com/stretchr/testify/assert"
)

func RegisterNetwork(t *testing.T, networkID string, networkName string) {
	err := configurator.CreateNetwork(configurator.Network{ID: networkID, Name: networkName}, nil)
	assert.NoError(t, err)
}

func RegisterGateway(t *testing.T, networkID string, gatewayID string, record *models.GatewayDevice) {
	RegisterGatewayWithName(t, networkID, gatewayID, "", record)
}

func RegisterGatewayWithName(t *testing.T, networkID string, gatewayID string, name string, record *models.GatewayDevice) {
	var gwEntity configurator.NetworkEntity
	if record != nil {
		if exists, _ := device.DoesDeviceExist(networkID, orc8r.AccessGatewayRecordType, record.HardwareID); exists {
			t.Fatalf("Hwid is already registered %s", record.HardwareID)
		}
		// write into device
		err := device.RegisterDevice(networkID, orc8r.AccessGatewayRecordType, record.HardwareID, record, serdes.Device)
		assert.NoError(t, err)

		gwEntity = configurator.NetworkEntity{
			Type:       orc8r.MagmadGatewayType,
			Key:        gatewayID,
			Name:       name,
			PhysicalID: record.HardwareID,
		}
	} else {
		gwEntity = configurator.NetworkEntity{
			Type: orc8r.MagmadGatewayType,
			Key:  gatewayID,
			Name: name,
		}
	}
	_, err := configurator.CreateEntity(networkID, gwEntity, serdes.Entity)
	assert.NoError(t, err)
}

// RemoveGateway assumes there is a device entity corresponding to the
// configurator entity
func RemoveGateway(t *testing.T, networkID, gatewayID string) {
	physicalID, err := configurator.GetPhysicalIDOfEntity(networkID, orc8r.MagmadGatewayType, gatewayID)
	assert.NoError(t, err)
	assert.NoError(t, device.DeleteDevice(networkID, orc8r.AccessGatewayRecordType, physicalID))
	assert.NoError(t, configurator.DeleteEntity(networkID, orc8r.MagmadGatewayType, gatewayID))
}
