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

package plugin_test

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	orc8rplugin "magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	"magma/wifi/cloud/go/plugin"
	"magma/wifi/cloud/go/protos/mconfig"
	"magma/wifi/cloud/go/services/wifi/obsidian/models"
	"magma/wifi/cloud/go/wifi"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestBuilder_Build_BaseCases(t *testing.T) {
	orc8rplugin.RegisterPluginForTests(t, &plugin.WifiOrchestratorPlugin{})
	builder := &plugin.Builder{}

	// empty case: no nw config
	nw := configurator.Network{ID: "n1"}
	gw := configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "gw1"}
	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gw},
	}

	actual := map[string]proto.Message{}
	expected := map[string]proto.Message{}
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// add nw config but no gw config
	nw.Configs = map[string]interface{}{
		wifi.WifiNetworkType: &models.NetworkWifiConfigs{
			AdditionalProps: map[string]string{"foo": "bar"},
		},
	}
	err = builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

	// nw config, gw config, no mesh
	wifigw := configurator.NetworkEntity{
		Type: wifi.WifiGatewayType,
		Key:  "gw1",
		Config: &models.GatewayWifiConfigs{
			AdditionalProps: map[string]string{"foo": "bar"},
		},
		ParentAssociations: []storage.TypeAndKey{{Type: orc8r.MagmadGatewayType, Key: "gw1"}},
	}
	gw.Associations = []storage.TypeAndKey{wifigw.GetTypeAndKey()}
	graph.Entities = []configurator.NetworkEntity{gw, wifigw}
	graph.Edges = []configurator.GraphEdge{
		{From: gw.GetTypeAndKey(), To: wifigw.GetTypeAndKey()},
	}
	err = builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestBuilder_Build(t *testing.T) {
	builder := &plugin.Builder{}

	// empty case: no nw config
	wifiNetworkConfigs := newDefaultNWConfig()
	wifiNetworkConfigs.VlAuthServerAddr = "20.1.1.8"
	wifiNetworkConfigs.VlAuthServerPort = 555
	wifiNetworkConfigs.VlAuthServerSharedSecret = "NewVlSecret"
	wifiNetworkConfigs.XwfRadiusServer = "vl.expresswifi.com"
	wifiNetworkConfigs.XwfConfig = "config line 1d\nconfig line 2d"
	wifiNetworkConfigs.XwfDhcpDns1 = "4.3.8.7"
	wifiNetworkConfigs.XwfDhcpDns2 = "8.2.3.3"
	wifiNetworkConfigs.XwfRadiusSharedSecret = "2231"
	wifiNetworkConfigs.XwfRadiusAuthPort = 2412
	wifiNetworkConfigs.XwfRadiusAcctPort = 2413
	wifiNetworkConfigs.XwfUamSecret = "233"
	wifiNetworkConfigs.XwfPartnerName = "xwffcfulld"
	wifiNetworkConfigs.AdditionalProps = map[string]string{"net1": "netval1", "net2": "netval2"}
	nw := configurator.Network{
		ID:      "n1",
		Configs: map[string]interface{}{wifi.WifiNetworkType: wifiNetworkConfigs},
	}

	gw := configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "gw1"}

	wifiGwConfigs := newDefaultGWConfig()
	wifiGwConfigs.MeshID = "m1"
	wifiGwConfigs.ClientChannel = "4"
	wifiGwConfigs.Latitude = 1.2
	wifiGwConfigs.Longitude = -4.4
	wifiGwConfigs.IsProduction = true
	wifiGwConfigs.AdditionalProps = map[string]string{"gw1": "gwval1", "gw2": "gwval2"}
	wifigw := configurator.NetworkEntity{
		Type:   wifi.WifiGatewayType,
		Key:    "gw1",
		Config: wifiGwConfigs,
	}

	mesh := configurator.NetworkEntity{
		Type: wifi.MeshEntityType,
		Key:  "m1",
		Config: &models.MeshWifiConfigs{
			Ssid:            "NewSsid",
			Password:        "NewPassword",
			VlSsid:          "NewVlSsid",
			XwfEnabled:      true,
			MeshSsid:        "OCP",
			MeshFrequency:   1987,
			MeshChannelType: "ED209",
			AdditionalProps: map[string]string{"mesh1": "meshval1", "mesh2": "meshval2"},
		},
		Associations: []storage.TypeAndKey{
			wifigw.GetTypeAndKey(),
			{Type: orc8r.MagmadGatewayType, Key: "gw2nd"},
		},
	}

	wifigw.ParentAssociations = []storage.TypeAndKey{mesh.GetTypeAndKey(), gw.GetTypeAndKey()}

	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gw, wifigw, mesh},
		Edges: []configurator.GraphEdge{
			{From: mesh.GetTypeAndKey(), To: gw.GetTypeAndKey()},
			{From: gw.GetTypeAndKey(), To: wifigw.GetTypeAndKey()},
		},
	}

	actual := map[string]proto.Message{}
	expected := map[string]proto.Message{
		"hostapd": &mconfig.Hostapd{
			Ssid:                     "NewSsid",
			Password:                 "NewPassword",
			VlSsid:                   "NewVlSsid",
			VlAuthServerAddr:         "20.1.1.8",
			VlAuthServerPort:         555,
			VlAuthServerSharedSecret: "NewVlSecret",
			WifiDisabled:             false,
			ClientChannel:            "4",
			XwfEnabled:               true,
		},
		"linkstatsd": &mconfig.Linkstatsd{
			PingHostList:    []string{"172.16.0.1", "www.facebook.com"},
			PingNumPackets:  5,
			PingTimeoutSecs: 6,
		},
		"meshd": &mconfig.Meshd{
			MeshRssiThreshold: -80,
			MeshSsid:          "OCP",
			MeshFrequency:     1987,
			MeshChannelType:   "ED209",
		},
		"openr": &mconfig.Openr{
			OpenrEnabled: false,
		},
		"openvpn": &mconfig.Openvpn{
			MgmtVpnEnabled: false,
		},
		"wifimetadata": &mconfig.WifiMetadata{
			Info:               "",
			Latitude:           1.2,
			Longitude:          -4.4,
			NetworkId:          "n1",
			MeshId:             "m1",
			GatewayId:          "gw1",
			MeshPeerGatewayIds: []string{"gw1", "gw2nd"},
			IsProduction:       true,
		},
		"xwfchilli": &mconfig.Xwfchilli{
			XwfEnabled:            true,
			XwfRadiusServer:       "vl.expresswifi.com",
			XwfConfig:             "config line 1d\nconfig line 2d",
			XwfDhcpDns1:           "4.3.8.7",
			XwfDhcpDns2:           "8.2.3.3",
			XwfRadiusSharedSecret: "2231",
			XwfRadiusAuthPort:     2412,
			XwfRadiusAcctPort:     2413,
			XwfUamSecret:          "233",
			XwfPartnerName:        "xwffcfulld",
			NetworkId:             "n1",
			MeshId:                "m1",
			GatewayId:             "gw1",
		},
		"wifiproperties": &mconfig.WifiProperties{
			Info:               "",
			NetworkId:          "n1",
			MeshId:             "m1",
			GatewayId:          "gw1",
			MeshPeerGatewayIds: []string{"gw1", "gw2nd"},
			NetworkProps:       map[string]string{"net1": "netval1", "net2": "netval2"},
			MeshProps:          map[string]string{"mesh1": "meshval1", "mesh2": "meshval2"},
			GatewayProps:       map[string]string{"gw1": "gwval1", "gw2": "gwval2"},
		},
	}
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestBuilder_Build_OverrideSsid(t *testing.T) {
	builder := &plugin.Builder{}

	// empty case: no nw config
	wifiNetworkConfigs := newDefaultNWConfig()
	wifiNetworkConfigs.VlAuthServerAddr = "20.1.1.8"
	wifiNetworkConfigs.VlAuthServerPort = 555
	wifiNetworkConfigs.VlAuthServerSharedSecret = "NewVlSecret"
	wifiNetworkConfigs.XwfRadiusServer = "vl.expresswifi.com"
	wifiNetworkConfigs.XwfConfig = "config line 1d\nconfig line 2d"
	wifiNetworkConfigs.XwfDhcpDns1 = "4.3.8.7"
	wifiNetworkConfigs.XwfDhcpDns2 = "8.2.3.3"
	wifiNetworkConfigs.XwfRadiusSharedSecret = "2231"
	wifiNetworkConfigs.XwfRadiusAuthPort = 2412
	wifiNetworkConfigs.XwfRadiusAcctPort = 2413
	wifiNetworkConfigs.XwfUamSecret = "233"
	wifiNetworkConfigs.XwfPartnerName = "xwffcfulld"
	nw := configurator.Network{
		ID:      "n1",
		Configs: map[string]interface{}{wifi.WifiNetworkType: wifiNetworkConfigs},
	}

	gw := configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "gw1"}

	wifiGwConfigs := newDefaultGWConfig()
	wifiGwConfigs.MeshID = "m1"
	wifiGwConfigs.ClientChannel = "4"
	wifiGwConfigs.Latitude = 1.2
	wifiGwConfigs.Longitude = -4.4
	// test override ssid
	wifiGwConfigs.UseOverrideSsid = true
	wifiGwConfigs.OverrideSsid = "overridden ssid"
	wifiGwConfigs.OverridePassword = "overridden password"
	wifigw := configurator.NetworkEntity{
		Type:   wifi.WifiGatewayType,
		Key:    "gw1",
		Config: wifiGwConfigs,
	}

	mesh := configurator.NetworkEntity{
		Type: wifi.MeshEntityType,
		Key:  "m1",
		Config: &models.MeshWifiConfigs{
			Ssid:            "NewSsid",
			Password:        "NewPassword",
			VlSsid:          "NewVlSsid",
			XwfEnabled:      true,
			MeshSsid:        "OCP",
			MeshFrequency:   1987,
			MeshChannelType: "ED209",
			AdditionalProps: nil,
		},
		Associations: []storage.TypeAndKey{
			wifigw.GetTypeAndKey(),
		},
	}

	wifigw.ParentAssociations = []storage.TypeAndKey{mesh.GetTypeAndKey(), gw.GetTypeAndKey()}

	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gw, wifigw, mesh},
		Edges: []configurator.GraphEdge{
			{From: mesh.GetTypeAndKey(), To: gw.GetTypeAndKey()},
			{From: gw.GetTypeAndKey(), To: wifigw.GetTypeAndKey()},
		},
	}

	actual := map[string]proto.Message{}
	expected := map[string]proto.Message{
		"hostapd": &mconfig.Hostapd{
			Ssid:                     "overridden ssid",
			Password:                 "overridden password",
			VlSsid:                   "NewVlSsid",
			VlAuthServerAddr:         "20.1.1.8",
			VlAuthServerPort:         555,
			VlAuthServerSharedSecret: "NewVlSecret",
			WifiDisabled:             false,
			ClientChannel:            "4",
			XwfEnabled:               true,
		},
		"linkstatsd": &mconfig.Linkstatsd{
			PingHostList:    []string{"172.16.0.1", "www.facebook.com"},
			PingNumPackets:  5,
			PingTimeoutSecs: 6,
		},
		"meshd": &mconfig.Meshd{
			MeshRssiThreshold: -80,
			MeshSsid:          "OCP",
			MeshFrequency:     1987,
			MeshChannelType:   "ED209",
		},
		"openr": &mconfig.Openr{
			OpenrEnabled: false,
		},
		"openvpn": &mconfig.Openvpn{
			MgmtVpnEnabled: false,
		},
		"wifimetadata": &mconfig.WifiMetadata{
			Info:               "",
			Latitude:           1.2,
			Longitude:          -4.4,
			NetworkId:          "n1",
			MeshId:             "m1",
			GatewayId:          "gw1",
			MeshPeerGatewayIds: []string{"gw1"},
			IsProduction:       false,
		},
		"xwfchilli": &mconfig.Xwfchilli{
			XwfEnabled:            true,
			XwfRadiusServer:       "vl.expresswifi.com",
			XwfConfig:             "config line 1d\nconfig line 2d",
			XwfDhcpDns1:           "4.3.8.7",
			XwfDhcpDns2:           "8.2.3.3",
			XwfRadiusSharedSecret: "2231",
			XwfRadiusAuthPort:     2412,
			XwfRadiusAcctPort:     2413,
			XwfUamSecret:          "233",
			XwfPartnerName:        "xwffcfulld",
			NetworkId:             "n1",
			MeshId:                "m1",
			GatewayId:             "gw1",
		},
		"wifiproperties": &mconfig.WifiProperties{
			Info:               "",
			NetworkId:          "n1",
			MeshId:             "m1",
			GatewayId:          "gw1",
			MeshPeerGatewayIds: []string{"gw1"},
			NetworkProps:       nil,
			MeshProps:          nil,
			GatewayProps:       nil,
		},
	}
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestBuilder_Build_OverrideXwf(t *testing.T) {
	builder := &plugin.Builder{}

	// empty case: no nw config
	wifiNetworkConfigs := newDefaultNWConfig()
	wifiNetworkConfigs.VlAuthServerAddr = "20.1.1.8"
	wifiNetworkConfigs.VlAuthServerPort = 555
	wifiNetworkConfigs.VlAuthServerSharedSecret = "NewVlSecret"
	wifiNetworkConfigs.XwfRadiusServer = "vl.expresswifi.com"
	wifiNetworkConfigs.XwfConfig = "config line 1d\nconfig line 2d"
	wifiNetworkConfigs.XwfDhcpDns1 = "4.3.8.7"
	wifiNetworkConfigs.XwfDhcpDns2 = "8.2.3.3"
	wifiNetworkConfigs.XwfRadiusSharedSecret = "2231"
	wifiNetworkConfigs.XwfRadiusAuthPort = 2412
	wifiNetworkConfigs.XwfRadiusAcctPort = 2413
	wifiNetworkConfigs.XwfUamSecret = "233"
	wifiNetworkConfigs.XwfPartnerName = "xwffcfulld"
	nw := configurator.Network{
		ID:      "n1",
		Configs: map[string]interface{}{wifi.WifiNetworkType: wifiNetworkConfigs},
	}

	gw := configurator.NetworkEntity{Type: orc8r.MagmadGatewayType, Key: "gw1"}

	wifiGwConfigs := newDefaultGWConfig()
	wifiGwConfigs.MeshID = "m1"
	wifiGwConfigs.ClientChannel = "4"
	wifiGwConfigs.Latitude = 1.2
	wifiGwConfigs.Longitude = -4.4
	// test override xwf configs
	wifiGwConfigs.UseOverrideXwf = true
	wifiGwConfigs.OverrideXwfEnabled = true
	wifiGwConfigs.OverrideXwfRadiusServer = "override radius"
	wifiGwConfigs.OverrideXwfConfig = "override config\n2\n3"
	wifiGwConfigs.OverrideXwfDhcpDns1 = "1.2.3.4"
	wifiGwConfigs.OverrideXwfDhcpDns2 = "5.6.7.8"
	wifiGwConfigs.OverrideXwfRadiusSharedSecret = "crash override"
	wifiGwConfigs.OverrideXwfRadiusAuthPort = 4
	wifiGwConfigs.OverrideXwfRadiusAcctPort = 6
	wifiGwConfigs.OverrideXwfUamSecret = "secret override"
	wifiGwConfigs.OverrideXwfPartnerName = "godzilla"
	wifigw := configurator.NetworkEntity{
		Type:   wifi.WifiGatewayType,
		Key:    "gw1",
		Config: wifiGwConfigs,
	}

	mesh := configurator.NetworkEntity{
		Type: wifi.MeshEntityType,
		Key:  "m1",
		Config: &models.MeshWifiConfigs{
			Ssid:            "NewSsid",
			Password:        "NewPassword",
			VlSsid:          "NewVlSsid",
			XwfEnabled:      true,
			MeshSsid:        "OCP",
			MeshFrequency:   1987,
			MeshChannelType: "ED209",
			AdditionalProps: map[string]string{},
		},
		Associations: []storage.TypeAndKey{
			wifigw.GetTypeAndKey(),
		},
	}

	wifigw.ParentAssociations = []storage.TypeAndKey{mesh.GetTypeAndKey(), gw.GetTypeAndKey()}

	graph := configurator.EntityGraph{
		Entities: []configurator.NetworkEntity{gw, wifigw, mesh},
		Edges: []configurator.GraphEdge{
			{From: mesh.GetTypeAndKey(), To: gw.GetTypeAndKey()},
			{From: gw.GetTypeAndKey(), To: wifigw.GetTypeAndKey()},
		},
	}

	actual := map[string]proto.Message{}
	expected := map[string]proto.Message{
		"hostapd": &mconfig.Hostapd{
			Ssid:                     "NewSsid",
			Password:                 "NewPassword",
			VlSsid:                   "NewVlSsid",
			VlAuthServerAddr:         "20.1.1.8",
			VlAuthServerPort:         555,
			VlAuthServerSharedSecret: "NewVlSecret",
			WifiDisabled:             false,
			ClientChannel:            "4",
			XwfEnabled:               true,
		},
		"linkstatsd": &mconfig.Linkstatsd{
			PingHostList:    []string{"172.16.0.1", "www.facebook.com"},
			PingNumPackets:  5,
			PingTimeoutSecs: 6,
		},
		"meshd": &mconfig.Meshd{
			MeshRssiThreshold: -80,
			MeshSsid:          "OCP",
			MeshFrequency:     1987,
			MeshChannelType:   "ED209",
		},
		"openr": &mconfig.Openr{
			OpenrEnabled: false,
		},
		"openvpn": &mconfig.Openvpn{
			MgmtVpnEnabled: false,
		},
		"wifimetadata": &mconfig.WifiMetadata{
			Info:               "",
			Latitude:           1.2,
			Longitude:          -4.4,
			NetworkId:          "n1",
			MeshId:             "m1",
			GatewayId:          "gw1",
			MeshPeerGatewayIds: []string{"gw1"},
			IsProduction:       false,
		},
		"xwfchilli": &mconfig.Xwfchilli{
			XwfEnabled:            true,
			XwfRadiusServer:       "override radius",
			XwfConfig:             "override config\n2\n3",
			XwfDhcpDns1:           "1.2.3.4",
			XwfDhcpDns2:           "5.6.7.8",
			XwfRadiusSharedSecret: "crash override",
			XwfRadiusAuthPort:     4,
			XwfRadiusAcctPort:     6,
			XwfUamSecret:          "secret override",
			XwfPartnerName:        "godzilla",
			NetworkId:             "n1",
			MeshId:                "m1",
			GatewayId:             "gw1",
		},
		"wifiproperties": &mconfig.WifiProperties{
			Info:               "",
			NetworkId:          "n1",
			MeshId:             "m1",
			GatewayId:          "gw1",
			MeshPeerGatewayIds: []string{"gw1"},
			NetworkProps:       nil,
			MeshProps:          nil,
			GatewayProps:       nil,
		},
	}
	err := builder.Build("n1", "gw1", graph, nw, actual)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func newDefaultNWConfig() *models.NetworkWifiConfigs {
	return &models.NetworkWifiConfigs{
		PingHostList:    []string{"172.16.0.1", "www.facebook.com"},
		PingNumPackets:  5,
		PingTimeoutSecs: 6,

		XwfConfig:             "config line 1\nconfig line 2",
		XwfDhcpDns1:           "8.8.8.8",
		XwfDhcpDns2:           "8.8.4.4",
		XwfRadiusSharedSecret: "1234",
		XwfRadiusAuthPort:     1812,
		XwfRadiusAcctPort:     1813,
		XwfUamSecret:          "1234",
		XwfPartnerName:        "xwfcfull",

		MgmtVpnEnabled: false,
		OpenrEnabled:   false,

		AdditionalProps: map[string]string{},
	}
}

func newDefaultGWConfig() *models.GatewayWifiConfigs {
	return &models.GatewayWifiConfigs{
		Info:              "",
		MeshID:            "",
		WifiDisabled:      false,
		MeshRssiThreshold: -80,
		ClientChannel:     "11",
		UseOverrideXwf:    false,
		UseOverrideSsid:   false,
		IsProduction:      false,
		AdditionalProps:   nil,
	}
}
