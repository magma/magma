/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package plugin

import (
	"fmt"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/plugin/models"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorprotos "magma/orc8r/cloud/go/services/configurator/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

type Builder struct{}
type FegMconfigBuilderServicer struct{}

func (s *FegMconfigBuilderServicer) Build(
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
		return ret, nil
	}
	gwConfig, err := getFegConfig(request.GetGatewayId(), network, graph)
	if err == merrors.ErrNotFound {
		return ret, nil
	}
	if err != nil {
		return ret, err
	}
	// Network health config takes priority. Only use gw health config
	// if network health config is nil
	healthConfig, err := getHealthConfig(network)
	if err != nil {
		healthConfig = gwConfig.Health
	}

	s6ac := gwConfig.S6a
	gxc := gwConfig.Gx
	gyc := gwConfig.Gy
	hss := gwConfig.Hss
	swxc := gwConfig.Swx
	eapAka := gwConfig.EapAka
	aaa := gwConfig.AaaServer
	csfb := gwConfig.Csfb
	healthc := protos.SafeInit(healthConfig).(*models.Health)

	vals := map[string]proto.Message{}
	if s6ac != nil {
		vals["s6a_proxy"] = &mconfig.S6AConfig{
			LogLevel:                protos.LogLevel_INFO,
			Server:                  s6ac.Server.ToMconfig(),
			RequestFailureThreshold: healthc.RequestFailureThreshold,
			MinimumRequestThreshold: healthc.MinimumRequestThreshold,
		}
	}

	if gxc != nil || gyc != nil {
		mc := &mconfig.SessionProxyConfig{
			LogLevel:                protos.LogLevel_INFO,
			RequestFailureThreshold: healthc.RequestFailureThreshold,
			MinimumRequestThreshold: healthc.MinimumRequestThreshold,
		}
		// Servers include the content of server
		if gxc != nil {
			mc.Gx = &mconfig.GxConfig{
				OverwriteApn: gxc.OverwriteApn,
				Servers:      models.ToMultipleServersMconfig(gxc.Server, gxc.Servers),
			}
			// TODO: once backwards compatibility is not needed, remove server from swagger, remove server from mconfg
			if len(mc.Gx.Servers) > 0 {
				mc.Gx.Server = mc.Gx.Servers[0]
			}
		}
		if gyc != nil {
			mc.Gy = &mconfig.GyConfig{
				InitMethod:   getGyInitMethod(gyc.InitMethod),
				OverwriteApn: gyc.OverwriteApn,
				Servers:      models.ToMultipleServersMconfig(gyc.Server, gyc.Servers),
			}
			// TODO: once backwards compatibility is not needed, remove server from swagger, remove server from mconfg
			if len(mc.Gy.Servers) > 0 {
				mc.Gy.Server = mc.Gy.Servers[0]
			}
		}
		vals["session_proxy"] = mc
	}

	if healthConfig != nil {
		mc := &mconfig.GatewayHealthConfig{}
		protos.FillIn(healthc, mc)
		vals["health"] = mc
	}

	if hss != nil {
		mc := &mconfig.HSSConfig{
			SubProfiles: map[string]*mconfig.HSSConfig_SubscriptionProfile{}} // legacy: avoid nil map
		protos.FillIn(hss, mc)
		vals["hss"] = mc
	}

	if swxc != nil {
		mc := &mconfig.SwxConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(swxc, mc)
		vals["swx_proxy"] = mc
	}

	if eapAka != nil {
		mc := &mconfig.EapAkaConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(eapAka, mc)
		vals["eap_aka"] = mc
	}

	if aaa != nil {
		mc := &mconfig.AAAConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(aaa, mc)
		vals["aaa_server"] = mc
	}

	if csfb != nil {
		mc := &mconfig.CsfbConfig{LogLevel: protos.LogLevel_INFO}
		protos.FillIn(csfb, mc)
		vals["csfb"] = mc
	}
	for k, v := range vals {
		ret.ConfigsByKey[k], err = ptypes.MarshalAny(v)
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}
func (*Builder) Build(
	networkID string,
	gatewayID string,
	graph configurator.EntityGraph,
	network configurator.Network,
	mconfigOut map[string]proto.Message) error {
	servicer := &FegMconfigBuilderServicer{}
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

func getFegConfig(
	gatewayID string, network configurator.Network, graph configurator.EntityGraph) (*models.GatewayFederationConfigs, error) {

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

func getGyInitMethod(initMethod *uint32) mconfig.GyInitMethod {
	if initMethod == nil {
		return mconfig.GyInitMethod_RESERVED
	}
	return mconfig.GyInitMethod(*initMethod)
}
