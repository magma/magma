/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"fmt"

	"magma/feg/cloud/go/protos/mconfig"
	fegprotos "magma/feg/cloud/go/services/controller/protos"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/streaming"

	"github.com/golang/protobuf/ptypes"
)

// Subset of mconfig fields that this streamer manages
var managedFields = []string{
	"s6a_proxy",
	"session_proxy",
	"swx_proxy",
	"hss",
}

type FegStreamer struct{}

func (*FegStreamer) GetSubscribedConfigTypes() []string {
	// IMPORTANT: for now, the feg gateway config type is blacklisted.
	// Eventually, when we figure out the final format of the user-facing
	// config format, we can refactor this streamer and backfill the views.
	return []string{FegNetworkType}
}

func (fs *FegStreamer) SeedNewGatewayMconfig(
	networkId string,
	gatewayId string,
	mconfigOut *protos.GatewayConfigs, // output parameter
) error {
	// For a new gateway, we'll fill in network configs first
	cfg, err := config.GetConfig(networkId, FegNetworkType, networkId)
	if err != nil {
		return err
	}
	if cfg == nil {
		return nil
	}
	return fs.applyNwConfigUpdate(streaming.CreateOperation, cfg.(*fegprotos.Config), mconfigOut)
}

func (fs *FegStreamer) ApplyMconfigUpdate(
	update *streaming.ConfigUpdate,
	oldMconfigsByGatewayId map[string]*protos.GatewayConfigs,
) (map[string]*protos.GatewayConfigs, error) {
	if update.ConfigType != FegNetworkType {
		return oldMconfigsByGatewayId, fmt.Errorf("feg mconfig streamer received unsubscribed type %s", update.ConfigType)
	}

	newConfig := castConfigValueToNwConfig(update.NewValue)
	for _, mconfigValue := range oldMconfigsByGatewayId {
		err := fs.applyNwConfigUpdate(update.Operation, newConfig, mconfigValue)
		if err != nil {
			return oldMconfigsByGatewayId, err
		}
	}
	return oldMconfigsByGatewayId, nil
}

func castConfigValueToNwConfig(v interface{}) *fegprotos.Config {
	if v == nil {
		return nil
	}
	return v.(*fegprotos.Config)
}

func (*FegStreamer) applyNwConfigUpdate(
	operation streaming.ChangeOperation,
	newConfig *fegprotos.Config,
	mconfigOut *protos.GatewayConfigs, // output param
) error {
	switch operation {
	case streaming.DeleteOperation:
		for _, field := range managedFields {
			delete(mconfigOut.ConfigsByKey, field)
		}
		return nil
	case streaming.ReadOperation, streaming.UpdateOperation, streaming.CreateOperation:
		s6ac := newConfig.GetS6A()
		gxc := newConfig.GetGx()
		gyc := newConfig.GetGy()
		hssc := newConfig.GetHss()
		swxc := newConfig.GetSwx()

		s6aMconfig := &mconfig.S6AConfig{
			LogLevel: protos.LogLevel_INFO,
			Server:   s6ac.GetServer().ToMconfig(),
		}
		sessionProxyMconfig := &mconfig.SessionProxyConfig{
			LogLevel: protos.LogLevel_INFO,
			Gx: &mconfig.GxConfig{
				Server: gxc.GetServer().ToMconfig(),
			},
			Gy: &mconfig.GyConfig{
				Server:     gyc.GetServer().ToMconfig(),
				InitMethod: mconfig.GyInitMethod(gyc.GetInitMethod()),
			},
		}
		hssSubProfileMconfig := make(map[string]*mconfig.HSSConfig_SubscriptionProfile)
		for imsi, profile := range hssc.GetSubProfiles() {
			hssSubProfileMconfig[imsi] = profile.ToMconfig()
		}
		hssMconfig := &mconfig.HSSConfig{
			Server:            hssc.GetServer().ToMconfig(),
			LteAuthOp:         hssc.GetLteAuthOp(),
			LteAuthAmf:        hssc.GetLteAuthAmf(),
			DefaultSubProfile: hssc.GetDefaultSubProfile().ToMconfig(),
			SubProfiles:       hssSubProfileMconfig,
		}
		swxProxyMconfig := &mconfig.SwxConfig{
			LogLevel: protos.LogLevel_INFO,
			Server:   swxc.GetServer().ToMconfig(),
		}

		s6aAny, err := ptypes.MarshalAny(s6aMconfig)
		if err != nil {
			return err
		}
		sessionProxyAny, err := ptypes.MarshalAny(sessionProxyMconfig)
		if err != nil {
			return err
		}
		hssAny, err := ptypes.MarshalAny(hssMconfig)
		if err != nil {
			return err
		}
		swxProxyAny, err := ptypes.MarshalAny(swxProxyMconfig)
		if err != nil {
			return err
		}

		mconfigOut.ConfigsByKey["s6a_proxy"] = s6aAny
		mconfigOut.ConfigsByKey["session_proxy"] = sessionProxyAny
		mconfigOut.ConfigsByKey["hss"] = hssAny
		mconfigOut.ConfigsByKey["swx_proxy"] = swxProxyAny
		return nil
	default:
		return fmt.Errorf("Unrecognized stream change operation %s", operation)
	}
}
