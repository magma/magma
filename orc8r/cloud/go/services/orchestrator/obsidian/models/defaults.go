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
	"magma/orc8r/cloud/go/models"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func NewDefaultDNSConfig() *NetworkDNSConfig {
	return &NetworkDNSConfig{
		EnableCaching: swag.Bool(true),
		LocalTTL:      swag.Uint32(60),
		Records: []*DNSConfigRecord{
			{
				ARecord:     []strfmt.IPv4{"192.88.99.142"},
				AaaaRecord:  []strfmt.IPv6{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
				CnameRecord: []string{"cname.example.com"},
				Domain:      "example.com",
			},
		},
	}
}

func NewDefaultFeaturesConfig() *NetworkFeatures {
	return &NetworkFeatures{Features: map[string]string{"foo": "bar"}}
}

func NewDefaultNetwork(networkID string, name string, description string) *Network {
	return &Network{
		ID:          models.NetworkID(networkID),
		Type:        "",
		Name:        models.NetworkName(name),
		Description: models.NetworkDescription(description),
		DNS:         NewDefaultDNSConfig(),
		Features:    NewDefaultFeaturesConfig(),
	}
}

func NewDefaultGatewayStatus(hardwareID string) *GatewayStatus {
	return &GatewayStatus{
		CheckinTime:        0,
		CertExpirationTime: 0,
		Meta:               map[string]string{"hello": "world"},
		SystemStatus: &SystemStatus{
			Time:       1495484735606,
			CPUUser:    31498,
			CPUSystem:  8361,
			CPUIdle:    1869111,
			MemTotal:   1016084,
			MemUsed:    54416,
			MemFree:    412772,
			UptimeSecs: 1234,
			SwapTotal:  1016081,
			SwapUsed:   54415,
			SwapFree:   412771,
			DiskPartitions: []*DiskPartition{
				{
					Device:     "/dev/sda1",
					MountPoint: "/",
					Total:      1,
					Used:       2,
					Free:       3,
				},
			},
		},
		PlatformInfo: &PlatformInfo{
			VpnIP: "facebook.com",
			Packages: []*Package{
				{
					Name:    "magma",
					Version: "0.0.0.0",
				},
			},
			KernelVersion:           "42",
			KernelVersionsInstalled: []string{"42", "43"},
			ConfigInfo: &ConfigInfo{
				MconfigCreatedAt: 1552968732,
			},
		},
		MachineInfo: &MachineInfo{
			CPUInfo: &MachineInfoCPUInfo{
				CoreCount:      4,
				ThreadsPerCore: 1,
				Architecture:   "x86_64",
				ModelName:      "Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz",
			},
			NetworkInfo: &MachineInfoNetworkInfo{
				NetworkInterfaces: []*NetworkInterface{
					{
						NetworkInterfaceID: "gtp_br0",
						Status:             NetworkInterfaceStatusUP,
						MacAddress:         "08:00:27:1e:8a:32",
						IPAddresses:        []string{"10.10.10.1"},
						IPV6Addresses:      []string{"fe80::a00:27ff:fe1e:8332"},
					},
				},
				RoutingTable: []*Route{
					{
						DestinationIP:      "0.0.0.0",
						GatewayIP:          "10.10.10.1",
						Genmask:            "255.255.255.0",
						NetworkInterfaceID: "eth0",
					},
				},
			},
		},
		HardwareID: hardwareID,
	}
}
