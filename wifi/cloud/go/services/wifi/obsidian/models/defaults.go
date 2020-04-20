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
	models2 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"

	"github.com/go-openapi/swag"
)

func NewDefaultWifiNetwork() *WifiNetwork {
	return &WifiNetwork{
		ID:          "n1",
		Name:        "network_1",
		Description: "Network 1",
		Features:    models.NewDefaultFeaturesConfig(),
		Wifi:        NewDefaultWifiNetworkConfig(),
	}
}

func NewDefaultWifiNetworkConfig() *NetworkWifiConfigs {
	return &NetworkWifiConfigs{
		VlAuthServerAddr:         "192.168.1.1",
		VlAuthServerPort:         1234,
		VlAuthServerSharedSecret: "ssssh",
		PingHostList:             []string{"172.16.0.1", "www.facebook.com"},
		PingNumPackets:           10,
		PingTimeoutSecs:          15,
		XwfRadiusServer:          "radiusnow",
		XwfConfig:                "line 1a\nline 2b",
		XwfDhcpDns1:              "4.8.8.7",
		XwfDhcpDns2:              "8.8.3.3",
		XwfRadiusSharedSecret:    "1231",
		XwfRadiusAuthPort:        2812,
		XwfRadiusAcctPort:        2813,
		XwfUamSecret:             "1233",
		XwfPartnerName:           "xwffcfull",
		MgmtVpnEnabled:           true,
		MgmtVpnProto:             "cows",
		MgmtVpnRemote:            "are yummy",
		OpenrEnabled:             true,
		AdditionalProps:          map[string]string{"prop1": "val1", "prop2": "val2", "prop3": "val3"},
	}
}

func NewDefaultWifiGateway() *WifiGateway {
	return &WifiGateway{
		Name:        "gateway_1",
		Description: "gateway 1",
		ID:          "g1",
		Device: &models.GatewayDevice{
			HardwareID: "hw1",
			Key:        &models.ChallengeKey{KeyType: "ECHO"},
		},
		Magmad: &models.MagmadGatewayConfigs{
			AutoupgradeEnabled:      swag.Bool(true),
			AutoupgradePollInterval: 300,
			CheckinInterval:         15,
			CheckinTimeout:          5,
		},
		Tier: "t1",
		Wifi: NewDefaultWifiGatewayConfig(),
	}
}

func NewDefaultWifiGatewayConfig() *GatewayWifiConfigs {
	return &GatewayWifiConfigs{
		AdditionalProps: map[string]string{
			"gwprop1": "gwvalue1",
			"gwprop2": "gwvalue2",
		},
		ClientChannel:                 "11",
		Info:                          "GatewayInfo",
		IsProduction:                  false,
		Latitude:                      37.48497,
		Longitude:                     -122.148284,
		MeshID:                        "m1",
		MeshRssiThreshold:             -80,
		OverridePassword:              "password",
		OverrideSsid:                  "SuperFastWifiNetwork",
		OverrideXwfConfig:             "xwf config",
		OverrideXwfDhcpDns1:           "8.8.8.8",
		OverrideXwfDhcpDns2:           "8.8.4.4",
		OverrideXwfEnabled:            false,
		OverrideXwfPartnerName:        "xwfcfull",
		OverrideXwfRadiusAcctPort:     1813,
		OverrideXwfRadiusAuthPort:     1812,
		OverrideXwfRadiusServer:       "gradius.example.com",
		OverrideXwfRadiusSharedSecret: "xwfisgood",
		OverrideXwfUamSecret:          "theuamsecret",
		UseOverrideSsid:               false,
		UseOverrideXwf:                false,
		WifiDisabled:                  false,
	}
}

func NewDefaultWifiMesh() *WifiMesh {
	return &WifiMesh{
		ID:         "m1",
		Name:       MeshName("mesh_1"),
		Config:     NewDefaultMeshWifiConfigs(),
		GatewayIds: []models2.GatewayID{"g1"},
	}
}

func NewDefaultMeshWifiConfigs() *MeshWifiConfigs {
	return &MeshWifiConfigs{
		AdditionalProps: map[string]string{
			"mesh_prop1": "mesh_value1",
			"mesh_prop2": "mesh_value2",
		},
		MeshChannelType: "HT20",
		MeshFrequency:   5825,
		MeshSsid:        "mesh_ssid",
		Password:        "password",
		Ssid:            "ssid",
		VlSsid:          "vl_ssid",
		XwfEnabled:      false,
	}
}
