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

	"magma/cwf/cloud/go/cwf"
	"magma/cwf/cloud/go/services/carrier_wifi/obsidian/models"
	fegmconfig "magma/feg/cloud/go/protos/mconfig"
	ltemconfig "magma/lte/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
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
	emptyRet := map[string]proto.Message{}
	pipelineDServices, err := getPipelineDServicesConfig(nwConfig.NetworkServices)
	if err != nil {
		return emptyRet, err
	}

	eapAka := nwConfig.EapAka
	aaa := nwConfig.AaaServer
	vals := map[string]proto.Message{
		"eap_aka": &fegmconfig.EapAkaConfig{
			LogLevel: protos.LogLevel_INFO,
			Timeout: &fegmconfig.EapAkaConfig_Timeouts{
				ChallengeMs:            eapAka.Timeout.ChallengeMs,
				ErrorNotificationMs:    eapAka.Timeout.ErrorNotificationMs,
				SessionMs:              eapAka.Timeout.SessionMs,
				SessionAuthenticatedMs: eapAka.Timeout.SessionAuthenticatedMs,
			},
			PlmnIds: eapAka.PlmnIds,
		},
		"aaa_server": &fegmconfig.AAAConfig{
			LogLevel:             protos.LogLevel_INFO,
			IdleSessionTimeoutMs: aaa.IDLESessionTimeoutMs,
			AccountingEnabled:    aaa.AccountingEnabled,
			CreateSessionOnAuth:  aaa.CreateSessionOnAuth,
		},
		"pipelined": &ltemconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     DefaultUeIpBlock, // Unused by CWF
			NatEnabled:    nwConfig.NatEnabled,
			DefaultRuleId: nwConfig.DefaultRuleID,
			RelayEnabled:  nwConfig.RelayEnabled,
			Services:      pipelineDServices,
		},
		"sessiond": &ltemconfig.SessionD{
			LogLevel:     protos.LogLevel_INFO,
			RelayEnabled: nwConfig.RelayEnabled,
		},
	}
	return vals, err
}

func getPipelineDServicesConfig(networkServices []string) ([]ltemconfig.PipelineD_NetworkServices, error) {
	if networkServices == nil || len(networkServices) == 0 {
		return []ltemconfig.PipelineD_NetworkServices{
			ltemconfig.PipelineD_METERING,
			ltemconfig.PipelineD_DPI,
			ltemconfig.PipelineD_ENFORCEMENT,
		}, nil
	}
	apps := make([]ltemconfig.PipelineD_NetworkServices, 0, len(networkServices))
	for _, service := range networkServices {
		mc, found := networkServicesByName[service]
		if !found {
			return nil, errors.Errorf("unknown network service name %s", service)
		}
		apps = append(apps, mc)
	}
	return apps, nil
}
