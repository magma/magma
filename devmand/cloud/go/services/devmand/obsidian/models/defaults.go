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
	"github.com/go-openapi/swag"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
)

func NewDefaultSymphonyNetwork() *SymphonyNetwork {
	return &SymphonyNetwork{
		ID:          "n1",
		Name:        "network_1",
		Description: "Network 1",
		Features:    models.NewDefaultFeaturesConfig(),
	}
}

func NewDefaultSymphonyAgent() *SymphonyAgent {
	return &SymphonyAgent{
		ID: "a1",
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		Name:        "agent_1",
		Description: "Agent 1",
		Tier:        "t1",
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		ManagedDevices: []string{"d1", "d2"},
	}
}

func NewDefaultSymphonyDevice() *SymphonyDevice {
	return &SymphonyDevice{
		Config:        NewDefaultSymphonyDeviceConfig(),
		ID:            "d1",
		Name:          "Device 1",
		ManagingAgent: "",
	}
}

func NewDefaultSymphonyDeviceConfig() *SymphonyDeviceConfig {
	return &SymphonyDeviceConfig{
		Channels: &SymphonyDeviceConfigChannels{
			SnmpChannel: &SnmpChannel{
				Community: "snmp community",
				Version:   "1",
			},
		},
		DeviceConfig: "{}",
		DeviceType:   []string{"device_type 1", "device_type 2"},
		Host:         "device_host",
		Platform:     "device_platform",
	}
}

func NewDefaultSymphonyDeviceState() *SymphonyDeviceState {
	return &SymphonyDeviceState{
		RawState: "{ SAMPLE_DEVICE: STATE }",
	}
}
