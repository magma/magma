/*
 Copyright 2020 The Magma Authors.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package servicers

import (
	"context"

	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	merrors "magma/orc8r/lib/go/errors"
	wifi_mconfig "magma/wifi/cloud/go/protos/mconfig"
	"magma/wifi/cloud/go/serdes"
	"magma/wifi/cloud/go/services/wifi/obsidian/models"
	"magma/wifi/cloud/go/wifi"

	"github.com/golang/protobuf/proto"
)

type builderServicer struct{}

func NewBuilderServicer() builder_protos.MconfigBuilderServer {
	return &builderServicer{}
}

func (s *builderServicer) Build(ctx context.Context, request *builder_protos.BuildRequest) (*builder_protos.BuildResponse, error) {
	ret := &builder_protos.BuildResponse{ConfigsByKey: map[string][]byte{}}

	network, err := (configurator.Network{}).FromProto(request.Network, serdes.Network)
	if err != nil {
		return nil, err
	}
	graph, err := (configurator.EntityGraph{}).FromProto(request.Graph, serdes.Entity)
	if err != nil {
		return nil, err
	}
	gatewayID := request.GatewayId

	// Need a valid wifi nw config, wifi gw config, and mesh config to
	// construct the mconfig
	iwifiNwConfig, found := network.Configs[wifi.WifiNetworkType]
	if !found {
		return ret, nil
	}
	wifiNwConfig := iwifiNwConfig.(*models.NetworkWifiConfigs)

	wifiGw, err := graph.GetEntity(wifi.WifiGatewayType, gatewayID)
	if err == merrors.ErrNotFound {
		return ret, nil
	}
	if err != nil {
		return nil, err
	}
	if wifiGw.Config == nil {
		return ret, nil
	}
	wifiGwConfig := wifiGw.Config.(*models.GatewayWifiConfigs)

	mesh, err := graph.GetFirstAncestorOfType(wifiGw, wifi.MeshEntityType)
	if err == merrors.ErrNotFound {
		return ret, nil
	}
	if err != nil {
		return nil, err
	}
	if mesh.Config == nil {
		return ret, nil
	}
	meshConfig := mesh.Config.(*models.MeshWifiConfigs)

	apConfig := mergeOverrideSsid(meshConfig, wifiGwConfig)
	xwfConfig := mergeOverrideXwf(wifiNwConfig, meshConfig, wifiGwConfig)
	meshPeerGatewayIds := getMeshPeerIDs(mesh)

	vals := map[string]proto.Message{
		"hostapd": &wifi_mconfig.Hostapd{
			Ssid:                     apConfig.ssid,
			Password:                 apConfig.password,
			VlSsid:                   meshConfig.VlSsid,
			VlAuthServerAddr:         wifiNwConfig.VlAuthServerAddr,
			VlAuthServerPort:         wifiNwConfig.VlAuthServerPort,
			VlAuthServerSharedSecret: wifiNwConfig.VlAuthServerSharedSecret,
			WifiDisabled:             wifiGwConfig.WifiDisabled,
			ClientChannel:            wifiGwConfig.ClientChannel,
			XwfEnabled:               xwfConfig.xwfEnabled,
		},
		"linkstatsd": &wifi_mconfig.Linkstatsd{
			PingHostList:    wifiNwConfig.PingHostList,
			PingNumPackets:  wifiNwConfig.PingNumPackets,
			PingTimeoutSecs: wifiNwConfig.PingTimeoutSecs,
		},
		"meshd": &wifi_mconfig.Meshd{
			MeshRssiThreshold: wifiGwConfig.MeshRssiThreshold,
			MeshSsid:          meshConfig.MeshSsid,
			MeshFrequency:     meshConfig.MeshFrequency,
			MeshChannelType:   meshConfig.MeshChannelType,
		},
		"openr": &wifi_mconfig.Openr{
			OpenrEnabled: wifiNwConfig.OpenrEnabled,
		},
		"openvpn": &wifi_mconfig.Openvpn{
			MgmtVpnEnabled: wifiNwConfig.MgmtVpnEnabled,
			MgmtVpnProto:   wifiNwConfig.MgmtVpnProto,
			MgmtVpnRemote:  wifiNwConfig.MgmtVpnRemote,
		},
		"wifimetadata": &wifi_mconfig.WifiMetadata{
			Info:               wifiGwConfig.Info,
			Latitude:           wifiGwConfig.Latitude,
			Longitude:          wifiGwConfig.Longitude,
			NetworkId:          network.ID,
			MeshId:             mesh.Key,
			GatewayId:          gatewayID,
			MeshPeerGatewayIds: meshPeerGatewayIds,
			IsProduction:       wifiGwConfig.IsProduction,
		},
		"xwfchilli": &wifi_mconfig.Xwfchilli{
			XwfEnabled:            xwfConfig.xwfEnabled,
			XwfRadiusServer:       xwfConfig.xwfRadiusServer,
			XwfConfig:             xwfConfig.xwfConfig,
			XwfDhcpDns1:           xwfConfig.xwfDhcpDns1,
			XwfDhcpDns2:           xwfConfig.xwfDhcpDns2,
			XwfRadiusSharedSecret: xwfConfig.xwfRadiusSharedSecret,
			XwfRadiusAuthPort:     xwfConfig.xwfRadiusAuthPort,
			XwfRadiusAcctPort:     xwfConfig.xwfRadiusAcctPort,
			XwfUamSecret:          xwfConfig.xwfUamSecret,
			XwfPartnerName:        xwfConfig.xwfPartnerName,
			NetworkId:             network.ID,
			MeshId:                mesh.Key,
			GatewayId:             gatewayID,
		},
		"wifiproperties": &wifi_mconfig.WifiProperties{
			Info:               wifiGwConfig.Info,
			NetworkId:          network.ID,
			MeshId:             mesh.Key,
			GatewayId:          gatewayID,
			MeshPeerGatewayIds: meshPeerGatewayIds,
			NetworkProps:       wifiNwConfig.AdditionalProps,
			MeshProps:          meshConfig.AdditionalProps,
			GatewayProps:       wifiGwConfig.AdditionalProps,
		},
	}

	ret.ConfigsByKey, err = mconfig.MarshalConfigs(vals)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type apConfig struct {
	ssid     string
	password string
}

type xwfConfig struct {
	xwfEnabled            bool
	xwfRadiusServer       string
	xwfConfig             string
	xwfDhcpDns1           string
	xwfDhcpDns2           string
	xwfRadiusSharedSecret string
	xwfRadiusAuthPort     int32
	xwfRadiusAcctPort     int32
	xwfUamSecret          string
	xwfPartnerName        string
}

func mergeOverrideSsid(meshConfig *models.MeshWifiConfigs, wifiGwConfig *models.GatewayWifiConfigs) apConfig {
	if wifiGwConfig.UseOverrideSsid {
		return apConfig{
			ssid:     wifiGwConfig.OverrideSsid,
			password: wifiGwConfig.OverridePassword,
		}
	} else {
		return apConfig{
			ssid:     meshConfig.Ssid,
			password: meshConfig.Password,
		}
	}
}

func mergeOverrideXwf(wifiNwConfig *models.NetworkWifiConfigs, meshConfig *models.MeshWifiConfigs, wifiGwConfig *models.GatewayWifiConfigs) xwfConfig {
	if wifiGwConfig.UseOverrideXwf {
		return xwfConfig{
			xwfEnabled:            wifiGwConfig.OverrideXwfEnabled,
			xwfRadiusServer:       wifiGwConfig.OverrideXwfRadiusServer,
			xwfConfig:             wifiGwConfig.OverrideXwfConfig,
			xwfDhcpDns1:           wifiGwConfig.OverrideXwfDhcpDns1,
			xwfDhcpDns2:           wifiGwConfig.OverrideXwfDhcpDns2,
			xwfRadiusSharedSecret: wifiGwConfig.OverrideXwfRadiusSharedSecret,
			xwfRadiusAuthPort:     wifiGwConfig.OverrideXwfRadiusAuthPort,
			xwfRadiusAcctPort:     wifiGwConfig.OverrideXwfRadiusAcctPort,
			xwfUamSecret:          wifiGwConfig.OverrideXwfUamSecret,
			xwfPartnerName:        wifiGwConfig.OverrideXwfPartnerName,
		}
	} else {
		return xwfConfig{
			xwfEnabled:            meshConfig.XwfEnabled,
			xwfRadiusServer:       wifiNwConfig.XwfRadiusServer,
			xwfConfig:             wifiNwConfig.XwfConfig,
			xwfDhcpDns1:           wifiNwConfig.XwfDhcpDns1,
			xwfDhcpDns2:           wifiNwConfig.XwfDhcpDns2,
			xwfRadiusSharedSecret: wifiNwConfig.XwfRadiusSharedSecret,
			xwfRadiusAuthPort:     wifiNwConfig.XwfRadiusAuthPort,
			xwfRadiusAcctPort:     wifiNwConfig.XwfRadiusAcctPort,
			xwfUamSecret:          wifiNwConfig.XwfUamSecret,
			xwfPartnerName:        wifiNwConfig.XwfPartnerName,
		}
	}
}

func getMeshPeerIDs(mesh configurator.NetworkEntity) []string {
	ret := make([]string, 0, len(mesh.Associations))
	for _, tk := range mesh.Associations {
		ret = append(ret, tk.Key)
	}
	return ret
}
