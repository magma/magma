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

package plugin

import (
	"magma/orc8r/cloud/go/services/configurator"
	configuratorprotos "magma/orc8r/cloud/go/services/configurator/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/wifi/cloud/go/protos/mconfig"
	"magma/wifi/cloud/go/services/wifi/obsidian/models"
	"magma/wifi/cloud/go/wifi"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

type Builder struct{}
type WifiMconfigBuilderServicer struct{}

// Build builds the wifi mconfig for a given networkID and gatewayID. It returns
// the mconfig as a map of config keys to mconfig messages.
func (s *WifiMconfigBuilderServicer) Build(
	request *configuratorprotos.BuildMconfigRequest,
) (*configuratorprotos.BuildMconfigResponse, error) {
	ret := &configuratorprotos.BuildMconfigResponse{
		ConfigsByKey: map[string]*any.Any{},
	}
	network, err := (configurator.Network{}).FromStorageProto(request.GetNetwork())
	if err != nil {
		return ret, err
	}
	graph, err := (configurator.EntityGraph{}).FromStorageProto(request.GetEntityGraph())
	if err != nil {
		return ret, err
	}
	gatewayID := request.GetGatewayId()
	networkID := request.GetNetworkId()
	// We need a valid wifi nw config, wifi gw config, and mesh config to
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
		return ret, err
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
		return ret, err
	}
	if mesh.Config == nil {
		return ret, nil
	}
	meshConfig := mesh.Config.(*models.MeshWifiConfigs)

	// process gwconfig ssid/password override
	apConfig := mergeOverrideSsid(meshConfig, wifiGwConfig)
	// process xwf config override
	xwfConfig := mergeOverrideXwf(wifiNwConfig, meshConfig, wifiGwConfig)
	meshPeerGatewayIds := getMeshPeerIDs(mesh)

	newConfigs := map[string]proto.Message{
		"hostapd": &mconfig.Hostapd{
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
		"linkstatsd": &mconfig.Linkstatsd{
			PingHostList:    wifiNwConfig.PingHostList,
			PingNumPackets:  wifiNwConfig.PingNumPackets,
			PingTimeoutSecs: wifiNwConfig.PingTimeoutSecs,
		},
		"meshd": &mconfig.Meshd{
			MeshRssiThreshold: wifiGwConfig.MeshRssiThreshold,
			MeshSsid:          meshConfig.MeshSsid,
			MeshFrequency:     meshConfig.MeshFrequency,
			MeshChannelType:   meshConfig.MeshChannelType,
		},
		"openr": &mconfig.Openr{
			OpenrEnabled: wifiNwConfig.OpenrEnabled,
		},
		"openvpn": &mconfig.Openvpn{
			MgmtVpnEnabled: wifiNwConfig.MgmtVpnEnabled,
			MgmtVpnProto:   wifiNwConfig.MgmtVpnProto,
			MgmtVpnRemote:  wifiNwConfig.MgmtVpnRemote,
		},
		"wifimetadata": &mconfig.WifiMetadata{
			Info:               wifiGwConfig.Info,
			Latitude:           wifiGwConfig.Latitude,
			Longitude:          wifiGwConfig.Longitude,
			NetworkId:          networkID,
			MeshId:             mesh.Key,
			GatewayId:          gatewayID,
			MeshPeerGatewayIds: meshPeerGatewayIds,
			IsProduction:       wifiGwConfig.IsProduction,
		},
		"xwfchilli": &mconfig.Xwfchilli{
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
			NetworkId:             networkID,
			MeshId:                mesh.Key,
			GatewayId:             gatewayID,
		},
		"wifiproperties": &mconfig.WifiProperties{
			Info:               wifiGwConfig.Info,
			NetworkId:          networkID,
			MeshId:             mesh.Key,
			GatewayId:          gatewayID,
			MeshPeerGatewayIds: meshPeerGatewayIds,
			NetworkProps:       wifiNwConfig.AdditionalProps,
			MeshProps:          meshConfig.AdditionalProps,
			GatewayProps:       wifiGwConfig.AdditionalProps,
		},
	}
	for k, v := range newConfigs {
		ret.ConfigsByKey[k], err = ptypes.MarshalAny(v)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func (*Builder) Build(networkID string, gatewayID string, graph configurator.EntityGraph, network configurator.Network, mconfigOut map[string]proto.Message) error {
	servicer := &WifiMconfigBuilderServicer{}
	networkProto, err := network.ToStorageProto()
	if err != nil {
		return errors.WithStack(err)
	}
	graphProto, err := graph.ToStorageProto()
	if err != nil {
		return errors.WithStack(err)
	}
	request := &configuratorprotos.BuildMconfigRequest{
		NetworkId:   networkID,
		GatewayId:   gatewayID,
		EntityGraph: graphProto,
		Network:     networkProto,
	}
	res, err := servicer.Build(request)
	if err != nil {
		return errors.WithStack(err)
	}
	for k, v := range res.GetConfigsByKey() {
		mconfigMessage, err := ptypes.Empty(v)
		if err != nil {
			return err
		}
		err = ptypes.UnmarshalAny(v, mconfigMessage)
		if err != nil {
			return err
		}
		mconfigOut[k] = mconfigMessage
	}
	return nil
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
