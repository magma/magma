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
	"fmt"

	"magma/feg/cloud/go/feg"
	feg_mconfig "magma/feg/cloud/go/protos/mconfig"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/feg/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
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
	gwConfig, err := getFegConfig(request.GatewayId, network, graph)
	if err == merrors.ErrNotFound {
		return ret, nil
	}
	if err != nil {
		return nil, err
	}
	// Network health config takes priority. Only use gw health config
	// if network health config is nil
	healthConfig, err := getHealthConfig(network)
	if err != nil {
		healthConfig = gwConfig.Health
	}

	s6ac := gwConfig.S6a
	s8c := gwConfig.S8
	gxc := gwConfig.Gx
	gyc := gwConfig.Gy
	hss := gwConfig.Hss
	swxc := gwConfig.Swx
	eapAka := gwConfig.EapAka
	eapSim := gwConfig.EapSim
	aaa := gwConfig.AaaServer
	csfb := gwConfig.Csfb
	healthc := protos.SafeInit(healthConfig).(*models.Health)

	vals := map[string]proto.Message{}

	if s6ac != nil {
		mc := &feg_mconfig.S6AConfig{
			LogLevel:                protos.LogLevel_INFO,
			RequestFailureThreshold: healthc.RequestFailureThreshold,
			MinimumRequestThreshold: healthc.MinimumRequestThreshold,
		}
		protos.FillIn(s6ac, mc)
		vals["s6a_proxy"] = mc
	}

	if s8c != nil {
		mc := &feg_mconfig.S8Config{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(s8c, mc)
		vals["s8_proxy"] = mc
	}

	if gxc != nil || gyc != nil {
		mc := &feg_mconfig.SessionProxyConfig{
			LogLevel:                protos.LogLevel_INFO,
			RequestFailureThreshold: healthc.RequestFailureThreshold,
			MinimumRequestThreshold: healthc.MinimumRequestThreshold,
		}
		// Servers include the content of server
		if gxc != nil {
			mc.Gx = &feg_mconfig.GxConfig{
				DisableGx:       swag.BoolValue(gxc.DisableGx),
				OverwriteApn:    gxc.OverwriteApn,
				Servers:         models.ToMultipleServersMconfig(gxc.Server, gxc.Servers),
				VirtualApnRules: models.ToVirtualApnRuleMconfig(gxc.VirtualApnRules),
			}
			// TODO(uri200): 5/7/20 remove this once backwards compatibility is not needed for the field server, remove server from swagger and mconfig
			if len(mc.Gx.Servers) > 0 {
				mc.Gx.Server = mc.Gx.Servers[0]
			}
		}
		if gyc != nil {
			mc.Gy = &feg_mconfig.GyConfig{
				DisableGy:       swag.BoolValue(gyc.DisableGy),
				InitMethod:      getGyInitMethod(gyc.InitMethod),
				OverwriteApn:    gyc.OverwriteApn,
				Servers:         models.ToMultipleServersMconfig(gyc.Server, gyc.Servers),
				VirtualApnRules: models.ToVirtualApnRuleMconfig(gyc.VirtualApnRules),
			}
			// TODO(uri200): 5/7/20 remove this once backwards compatibility is not needed for the field server, remove server from swagger and mconfig
			if len(mc.Gy.Servers) > 0 {
				mc.Gy.Server = mc.Gy.Servers[0]
			}
		}
		vals["session_proxy"] = mc
	}

	if healthConfig != nil {
		mc := &feg_mconfig.GatewayHealthConfig{}
		protos.FillIn(healthc, mc)
		vals["health"] = mc
	}

	if hss != nil {
		mc := &feg_mconfig.HSSConfig{
			SubProfiles: map[string]*feg_mconfig.HSSConfig_SubscriptionProfile{}} // legacy: avoid nil map
		protos.FillIn(hss, mc)
		vals["hss"] = mc
	}

	if swxc != nil {
		mc := &feg_mconfig.SwxConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(swxc, mc)

		// TODO(uri200): 5/7/20 remove this once backwards compatibility is not needed for the field server, remove server from swagger and mconfig
		mc.Servers = models.ToMultipleServersMconfig(swxc.Server, swxc.Servers)
		vals["swx_proxy"] = mc
	}

	if eapAka != nil {
		mc := &feg_mconfig.EapAkaConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(eapAka, mc)
		vals["eap_aka"] = mc
	}

	if eapSim != nil {
		mc := &feg_mconfig.EapSimConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(eapSim, mc)
		vals["eap_sim"] = mc
	}

	if aaa != nil {
		mc := &feg_mconfig.AAAConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(aaa, mc)
		vals["aaa_server"] = mc
	}

	if csfb != nil {
		mc := &feg_mconfig.CsfbConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(csfb, mc)
		vals["csfb"] = mc
	}

	ret.ConfigsByKey, err = mconfig.MarshalConfigs(vals)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
func getFegConfig(gatewayID string, network configurator.Network, graph configurator.EntityGraph) (*models.GatewayFederationConfigs, error) {
	fegGW, err := graph.GetEntity(feg.FegGatewayType, gatewayID)
	if err != nil && err != merrors.ErrNotFound {
		return nil, errors.WithStack(err)
	}
	// err can only be merrors.ErrNotFound at this point - if it's nil, we'll
	// just return the feg gateway config if it exists
	if err == nil && fegGW.Config != nil {
		return fegGW.Config.(*models.GatewayFederationConfigs), nil
	}

	inwConfig, found := network.Configs[feg.FegNetworkType]
	if !found || inwConfig == nil {
		return nil, merrors.ErrNotFound
	}
	nwConfig := inwConfig.(*models.NetworkFederationConfigs)
	return (*models.GatewayFederationConfigs)(nwConfig), nil
}

func getHealthConfig(network configurator.Network) (*models.Health, error) {
	inwConfig, found := network.Configs[feg.FegNetworkType]
	if !found || inwConfig == nil {
		return nil, merrors.ErrNotFound
	}
	nwConfig := inwConfig.(*models.NetworkFederationConfigs)
	config := (*models.GatewayFederationConfigs)(nwConfig)
	if config != nil && config.Health != nil {
		return config.Health, nil
	}
	return nil, fmt.Errorf("network health config is nil")
}

func getGyInitMethod(initMethod *uint32) feg_mconfig.GyInitMethod {
	if initMethod == nil {
		return feg_mconfig.GyInitMethod_RESERVED
	}
	return feg_mconfig.GyInitMethod(*initMethod)
}
