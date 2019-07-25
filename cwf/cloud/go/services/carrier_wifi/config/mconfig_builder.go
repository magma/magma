/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package config

import (
	"fmt"
	"log"
	"strings"

	"magma/cwf/cloud/go/cwf"
	"magma/cwf/cloud/go/services/carrier_wifi/obsidian/models"
	fegmconfig "magma/feg/cloud/go/protos/mconfig"
	ltemconfig "magma/lte/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/golang/protobuf/proto"
)

const (
	DefaultUeIpBlock = "192.168.128.0/24"
)

var networkServicesByName = map[string]ltemconfig.PipelineD_NetworkServices{
	"metering":           ltemconfig.PipelineD_METERING,
	"dpi":                ltemconfig.PipelineD_DPI,
	"policy_enforcement": ltemconfig.PipelineD_ENFORCEMENT,
}

type Builder struct{}

func (builder *Builder) Build(networkId string, gatewayId string) (map[string]proto.Message, error) {
	emptyRet := map[string]proto.Message{}

	// Get Network configs to fill in GW configs
	netCfg, err := configurator.GetNetworkConfigsByType(networkId, cwf.CwfNetworkType)
	if err != nil || netCfg == nil {
		return emptyRet, err
	}
	cwfCfg, ok := netCfg.(*models.NetworkCarrierWifiConfigs)
	if !ok {
		return emptyRet, fmt.Errorf(
			"received unexpected type for network record. Expected *NetworkCarrierWifiConfigs but got %T", netCfg)
	}
	return BuildFromNetworkConfig(cwfCfg)
}

func BuildFromNetworkConfig(nwConfig *models.NetworkCarrierWifiConfigs) (map[string]proto.Message, error) {
	ret := map[string]proto.Message{}
	if nwConfig == nil {
		return ret, nil
	}
	pipelineDServices, err := getPipelineDServicesConfig(nwConfig.NetworkServices)
	if err != nil {
		return ret, err
	}

	eapAka := nwConfig.EapAka
	aaa := nwConfig.AaaServer
	if eapAka != nil {
		mc := &fegmconfig.EapAkaConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(eapAka, mc)
		ret["eap_aka"] = mc
	}
	if aaa != nil {
		mc := &fegmconfig.AAAConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(aaa, mc)
		ret["aaa_server"] = mc
	}
	ret["pipelined"] = &ltemconfig.PipelineD{
		LogLevel:      protos.LogLevel_INFO,
		UeIpBlock:     DefaultUeIpBlock, // Unused by CWF
		NatEnabled:    nwConfig.NatEnabled,
		DefaultRuleId: nwConfig.DefaultRuleID,
		RelayEnabled:  nwConfig.RelayEnabled,
		Services:      pipelineDServices,
	}
	ret["sessiond"] = &ltemconfig.SessionD{
		LogLevel:     protos.LogLevel_INFO,
		RelayEnabled: nwConfig.RelayEnabled,
	}
	return ret, err
}

func getPipelineDServicesConfig(networkServices []string) ([]ltemconfig.PipelineD_NetworkServices, error) {
	apps := make([]ltemconfig.PipelineD_NetworkServices, 0, len(networkServices))
	for _, service := range networkServices {
		mc, found := networkServicesByName[strings.ToLower(service)]
		if !found {
			log.Printf("CWAG: unknown network service name %s", service)
		} else {
			apps = append(apps, mc)
		}
	}
	if len(apps) == 0 {
		apps = []ltemconfig.PipelineD_NetworkServices{
			ltemconfig.PipelineD_METERING,
			ltemconfig.PipelineD_DPI,
			ltemconfig.PipelineD_ENFORCEMENT,
		}
	}
	return apps, nil
}
