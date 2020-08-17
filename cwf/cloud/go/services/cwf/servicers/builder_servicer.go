/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"context"
	"log"
	"strings"

	"magma/cwf/cloud/go/cwf"
	cwf_mconfig "magma/cwf/cloud/go/protos/mconfig"
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

	network, err := (configurator.Network{}).FromStorageProto(request.Network)
	if err != nil {
		return nil, err
	}
	graph, err := (configurator.EntityGraph{}).FromStorageProto(request.Graph)
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

	vals, err := buildFromConfigs(nwConfig, gwConfig.Config.(*models.GatewayCwfConfigs))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	ret.ConfigsByKey, err = mconfig.MarshalConfigs(vals)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func buildFromConfigs(nwConfig *models.NetworkCarrierWifiConfigs, gwConfig *models.GatewayCwfConfigs) (map[string]proto.Message, error) {
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
	ipdrExportDst, err := getPipelineDIpdrExportDst(gwConfig.IPDRExportDst)
	if err != nil {
		return nil, err
	}

	eapAka := nwConfig.EapAka
	aaa := nwConfig.AaaServer
	if eapAka != nil {
		mc := &feg_mconfig.EapAkaConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(eapAka, mc)
		ret["eap_aka"] = mc
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
		LogLevel:     protos.LogLevel_INFO,
		RelayEnabled: true,
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
	if healthCfg != nil {
		mc := &cwf_mconfig.CwfGatewayHealthConfig{
			CpuUtilThresholdPct: healthCfg.CPUUtilThresholdPct,
			MemUtilThresholdPct: healthCfg.MemUtilThresholdPct,
			GreProbeInterval:    healthCfg.GreProbeIntervalSecs,
			IcmpProbePktCount:   healthCfg.IcmpProbePktCount,
		}
		protos.FillIn(healthCfg, mc)
		mc.GrePeers = getHealthServiceGrePeers(allowedGrePeers)
		ret["health"] = mc
	} else {
		mc := &cwf_mconfig.CwfGatewayHealthConfig{
			GrePeers: getHealthServiceGrePeers(allowedGrePeers),
		}
		ret["health"] = mc
	}

	return ret, nil
}

func getPipelineDAllowedGrePeers(allowedGrePeers models.AllowedGrePeers) ([]*lte_mconfig.PipelineD_AllowedGrePeer, error) {
	ues := make([]*lte_mconfig.PipelineD_AllowedGrePeer, 0, len(allowedGrePeers))
	for _, entry := range allowedGrePeers {
		ues = append(ues, &lte_mconfig.PipelineD_AllowedGrePeer{Ip: entry.IP, Key: swag.Uint32Value(entry.Key)})
	}
	return ues, nil
}

func getPipelineDIpdrExportDst(ipdrExportDst *models.IPDRExportDst) (*lte_mconfig.PipelineD_IPDRExportDst, error) {
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
