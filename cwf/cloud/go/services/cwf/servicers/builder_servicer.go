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
	"log"
	"strings"

	"magma/cwf/cloud/go/cwf"
	cwf_mconfig "magma/cwf/cloud/go/protos/mconfig"
	"magma/cwf/cloud/go/serdes"
	"magma/cwf/cloud/go/services/cwf/obsidian/models"
	feg_mconfig "magma/feg/cloud/go/protos/mconfig"
	lte_mconfig "magma/lte/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	orc8r_mconfig "magma/orc8r/lib/go/protos/mconfig"

	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	DefaultUeIpBlock = "192.168.128.0/24"
)

var networkServicesByName = map[string]lte_mconfig.PipelineD_NetworkServices{
	"metering":           lte_mconfig.PipelineD_METERING,
	"dpi":                lte_mconfig.PipelineD_DPI,
	"policy_enforcement": lte_mconfig.PipelineD_ENFORCEMENT,
}

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

	// Only build an mconfig if carrier_wifi network configs exist
	inwConfig, found := network.Configs[cwf.CwfNetworkType]
	if !found || inwConfig == nil {
		return ret, nil
	}
	nwConfig := inwConfig.(*models.NetworkCarrierWifiConfigs)
	gwConfig, err := graph.GetEntity(cwf.CwfGatewayType, request.GetGatewayId())
	if err == merrors.ErrNotFound {
		return ret, nil
	}
	if err != nil {
		return nil, err
	}
	if gwConfig.Config == nil {
		return ret, nil
	}

	var haPairConfigs *models.CwfHaPairConfigs
	haPairEnt, err := graph.GetFirstAncestorOfType(gwConfig, cwf.CwfHAPairType)
	if err != nil {
		haPairConfigs = nil
	} else {
		haPairConfigs = haPairEnt.Config.(*models.CwfHaPairConfigs)
	}
	vals, err := buildFromConfigs(nwConfig, gwConfig.Config.(*models.GatewayCwfConfigs), haPairConfigs)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ret.ConfigsByKey, err = mconfig.MarshalConfigs(vals)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func buildFromConfigs(nwConfig *models.NetworkCarrierWifiConfigs, gwConfig *models.GatewayCwfConfigs, haPairConfigs *models.CwfHaPairConfigs) (map[string]proto.Message, error) {
	ret := map[string]proto.Message{}
	if nwConfig == nil {
		return ret, nil
	}
	pipelineDServices, err := getPipelineDServicesConfig(nwConfig.NetworkServices)
	if err != nil {
		return nil, err
	}
	allowedGrePeers, err := getPipelineDAllowedGrePeers(gwConfig.AllowedGrePeers)
	if err != nil {
		return nil, err
	}
	liUes, err := getPipelineDLiUes(nwConfig.LiUes)
	if err != nil {
		return nil, err
	}
	ipdrExportDst, err := getPipelineDIpdrExportDst(gwConfig.IpdrExportDst)
	if err != nil {
		return nil, err
	}
	eapAka := nwConfig.EapAka
	eapSim := nwConfig.EapSim
	aaa := nwConfig.AaaServer
	if eapAka != nil {
		mc := &feg_mconfig.EapAkaConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(eapAka, mc)
		ret["eap_aka"] = mc
	}
	if eapSim != nil {
		mc := &feg_mconfig.EapSimConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(eapSim, mc)
		ret["eap_sim"] = mc
	}
	if aaa != nil {
		mc := &feg_mconfig.AAAConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(aaa, mc)
		ret["aaa_server"] = mc
	}
	ret["pipelined"] = &lte_mconfig.PipelineD{
		LogLevel:        protos.LogLevel_INFO,
		UeIpBlock:       DefaultUeIpBlock, // Unused by CWF
		NatEnabled:      false,
		DefaultRuleId:   swag.StringValue(nwConfig.DefaultRuleID),
		Services:        pipelineDServices,
		AllowedGrePeers: allowedGrePeers,
		LiUes:           liUes,
		IpdrExportDst:   ipdrExportDst,
	}
	ret["sessiond"] = &lte_mconfig.SessionD{
		LogLevel:         protos.LogLevel_INFO,
		GxGyRelayEnabled: true,
		WalletExhaustDetection: &lte_mconfig.WalletExhaustDetection{
			TerminateOnExhaust: true,
			Method:             lte_mconfig.WalletExhaustDetection_GxTrackedRules,
		},
	}
	ret["redirectd"] = &lte_mconfig.RedirectD{
		LogLevel: protos.LogLevel_INFO,
	}
	ret["directoryd"] = &orc8r_mconfig.DirectoryD{
		LogLevel: protos.LogLevel_INFO,
	}
	healthCfg := gwConfig.GatewayHealthConfigs
	mc := &cwf_mconfig.CwfGatewayHealthConfig{
		GrePeers: getHealthServiceGrePeers(allowedGrePeers),
	}
	if haPairConfigs != nil {
		mc.ClusterVirtualIp = haPairConfigs.TransportVirtualIP
	}
	if healthCfg != nil {
		protos.FillIn(healthCfg, mc)
	}
	ret["health"] = mc

	return ret, nil
}

func getPipelineDAllowedGrePeers(allowedGrePeers models.AllowedGrePeers) ([]*lte_mconfig.PipelineD_AllowedGrePeer, error) {
	ues := make([]*lte_mconfig.PipelineD_AllowedGrePeer, 0, len(allowedGrePeers))
	for _, entry := range allowedGrePeers {
		ues = append(ues, &lte_mconfig.PipelineD_AllowedGrePeer{Ip: entry.IP, Key: swag.Uint32Value(entry.Key)})
	}
	return ues, nil
}

func getPipelineDIpdrExportDst(ipdrExportDst *models.IpdrExportDst) (*lte_mconfig.PipelineD_IPDRExportDst, error) {
	if ipdrExportDst == nil {
		return nil, nil
	}
	dst := &lte_mconfig.PipelineD_IPDRExportDst{
		Ip:   ipdrExportDst.IP.String(),
		Port: ipdrExportDst.Port,
	}
	return dst, nil
}

func getPipelineDServicesConfig(networkServices []string) ([]lte_mconfig.PipelineD_NetworkServices, error) {
	apps := make([]lte_mconfig.PipelineD_NetworkServices, 0, len(networkServices))
	for _, service := range networkServices {
		mc, found := networkServicesByName[strings.ToLower(service)]
		if !found {
			log.Printf("CWAG: unknown network service name %s", service)
		} else {
			apps = append(apps, mc)
		}
	}
	return apps, nil
}

func getHealthServiceGrePeers(pipelinedPeers []*lte_mconfig.PipelineD_AllowedGrePeer) []*cwf_mconfig.CwfGatewayHealthConfigGrePeer {
	var healthPeers []*cwf_mconfig.CwfGatewayHealthConfigGrePeer
	if pipelinedPeers == nil {
		return healthPeers
	}
	for _, peer := range pipelinedPeers {
		healthPeer := &cwf_mconfig.CwfGatewayHealthConfigGrePeer{
			Ip: peer.Ip,
		}
		healthPeers = append(healthPeers, healthPeer)
	}
	return healthPeers
}

func getPipelineDLiUes(liUes *models.LiUes) (*lte_mconfig.PipelineD_LiUes, error) {
	if liUes == nil {
		return nil, nil
	}
	dst := &lte_mconfig.PipelineD_LiUes{
		Imsis:   liUes.Imsis,
		Msisdns: liUes.Msisdns,
		Macs:    liUes.Macs,
		Ips:     liUes.Ips,
	}
	return dst, nil
}
