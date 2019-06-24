/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin

import (
	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/cloud/go/services/controller/obsidian/models"
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type Builder struct{}

func (*Builder) Build(networkID string, gatewayID string, graph configurator.EntityGraph, network configurator.Network, mconfigOut map[string]proto.Message) error {
	gwConfig, err := getFegConfig(gatewayID, network, graph)
	if err == merrors.ErrNotFound {
		return nil
	}
	if err != nil {
		return errors.WithStack(err)
	}

	s6ac := gwConfig.S6a
	gxc := gwConfig.Gx
	gyc := gwConfig.Gy
	hss := gwConfig.Hss
	swxc := gwConfig.Swx
	eapAka := gwConfig.EapAka
	aaa := gwConfig.AaaServer
	healthc := gwConfig.Health

	mconfigOut["s6a_proxy"] = &mconfig.S6AConfig{
		LogLevel:                protos.LogLevel_INFO,
		Server:                  diamClientConfigToMconfig(s6ac.Server),
		RequestFailureThreshold: healthc.RequestFailureThreshold,
		MinimumRequestThreshold: healthc.MinimumRequestThreshold,
	}

	mconfigOut["session_proxy"] = &mconfig.SessionProxyConfig{
		LogLevel: protos.LogLevel_INFO,
		Gx: &mconfig.GxConfig{
			Server: diamClientConfigToMconfig(gxc.Server),
		},
		Gy: &mconfig.GyConfig{
			Server:     diamClientConfigToMconfig(gyc.Server),
			InitMethod: mconfig.GyInitMethod(*gyc.InitMethod),
		},
		RequestFailureThreshold: healthc.RequestFailureThreshold,
		MinimumRequestThreshold: healthc.MinimumRequestThreshold,
	}

	mconfigOut["health"] = &mconfig.GatewayHealthConfig{
		RequiredServices:          healthc.HealthServices,
		UpdateIntervalSecs:        healthc.UpdateIntervalSecs,
		UpdateFailureThreshold:    healthc.UpdateFailureThreshold,
		CloudDisconnectPeriodSecs: healthc.CloudDisablePeriodSecs,
		LocalDisconnectPeriodSecs: healthc.LocalDisablePeriodSecs,
	}

	hssSubProfile := map[string]*mconfig.HSSConfig_SubscriptionProfile{}
	for imsi, profile := range hss.SubProfiles {
		hssSubProfile[imsi] = subProfileToMconfig(&profile)
	}
	mconfigOut["hss"] = &mconfig.HSSConfig{
		Server:            diamServerConfigToMconfig(hss.Server),
		LteAuthOp:         hss.LteAuthOp,
		LteAuthAmf:        hss.LteAuthAmf,
		DefaultSubProfile: subProfileToMconfig(hss.DefaultSubProfile),
		SubProfiles:       hssSubProfile,
		StreamSubscribers: hss.StreamSubscribers,
	}

	mconfigOut["swx_proxy"] = &mconfig.SwxConfig{
		LogLevel:            protos.LogLevel_INFO,
		Server:              diamClientConfigToMconfig(swxc.Server),
		VerifyAuthorization: swxc.VerifyAuthorization,
		CacheTTLSeconds:     swxc.CacheTTLSeconds,
	}

	mconfigOut["eap_aka"] = &mconfig.EapAkaConfig{
		LogLevel: protos.LogLevel_INFO,
		Timeout: &mconfig.EapAkaConfig_Timeouts{
			ChallengeMs:            eapAka.Timeout.ChallengeMs,
			ErrorNotificationMs:    eapAka.Timeout.ErrorNotificationMs,
			SessionMs:              eapAka.Timeout.SessionMs,
			SessionAuthenticatedMs: eapAka.Timeout.SessionAuthenticatedMs,
		},
		PlmnIds: eapAka.PlmnIds,
	}

	mconfigOut["aaa_server"] = &mconfig.AAAConfig{
		LogLevel:             protos.LogLevel_INFO,
		IdleSessionTimeoutMs: aaa.IDLESessionTimeoutMs,
		AccountingEnabled:    aaa.AccountingEnabled,
		CreateSessionOnAuth:  aaa.CreateSessionOnAuth,
	}

	return nil
}

func getFegConfig(gatewayID string, network configurator.Network, graph configurator.EntityGraph) (*models.GatewayFegConfigs, error) {
	fegGW, err := graph.GetEntity(feg.FegGatewayType, gatewayID)
	if err != nil && err != merrors.ErrNotFound {
		return nil, errors.WithStack(err)
	}
	// err can only be merrors.ErrNotFound at this point - if it's nil, we'll
	// just return the feg gateway config if it exists
	if err == nil && fegGW.Config != nil {
		return fegGW.Config.(*models.GatewayFegConfigs), nil
	}

	inwConfig, found := network.Configs[feg.FegNetworkType]
	if !found || inwConfig == nil {
		return nil, merrors.ErrNotFound
	}
	nwConfig := inwConfig.(*models.NetworkFederationConfigs)
	return &models.GatewayFegConfigs{NetworkFederationConfigs: *nwConfig}, nil
}

func subProfileToMconfig(profile *models.SubscriptionProfile) *mconfig.HSSConfig_SubscriptionProfile {
	return &mconfig.HSSConfig_SubscriptionProfile{
		MaxUlBitRate: profile.MaxUlBitRate,
		MaxDlBitRate: profile.MaxDlBitRate,
	}
}

func diamClientConfigToMconfig(config *models.DiameterClientConfigs) *mconfig.DiamClientConfig {
	return &mconfig.DiamClientConfig{
		Protocol:         config.Protocol,
		Address:          config.Address,
		Retransmits:      config.Retransmits,
		WatchdogInterval: config.WatchdogInterval,
		RetryCount:       config.RetryCount,
		LocalAddress:     config.LocalAddress,
		ProductName:      config.ProductName,
		Realm:            config.Realm,
		Host:             config.Host,
		DestRealm:        config.DestRealm,
		DestHost:         config.DestHost,
	}
}

func diamServerConfigToMconfig(config *models.DiameterServerConfigs) *mconfig.DiamServerConfig {
	return &mconfig.DiamServerConfig{
		Protocol:     config.Protocol,
		Address:      config.Address,
		LocalAddress: config.LocalAddress,
		DestRealm:    config.DestRealm,
		DestHost:     config.DestHost,
	}
}
