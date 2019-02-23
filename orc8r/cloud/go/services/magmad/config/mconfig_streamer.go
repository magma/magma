/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"fmt"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config/streaming"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"

	"github.com/golang/protobuf/ptypes"
)

// Subset of mconfig fields that this streamer manages
var managedFields = []string{
	"control_proxy",
	"magmad",
	"metricsd",
}

type MagmadStreamer struct{}

func (*MagmadStreamer) GetSubscribedConfigTypes() []string {
	return []string{MagmadGatewayType}
}

func (*MagmadStreamer) SeedNewGatewayMconfig(
	networkId string,
	gatewayId string,
	mconfigOut *protos.GatewayConfigs, // output parameter
) error {
	// No network configs to propagate here
	return nil
}

func (*MagmadStreamer) ApplyMconfigUpdate(
	update *streaming.ConfigUpdate,
	oldMconfigsByGatewayId map[string]*protos.GatewayConfigs,
) (map[string]*protos.GatewayConfigs, error) {
	if update.ConfigType != MagmadGatewayType {
		return oldMconfigsByGatewayId, fmt.Errorf("magmad mconfig streamer received unsubscribed type %s", update.ConfigType)
	}

	switch update.Operation {
	case streaming.DeleteOperation:
		for _, mconfigValue := range oldMconfigsByGatewayId {
			for _, fieldName := range managedFields {
				delete(mconfigValue.ConfigsByKey, fieldName)
			}
		}
		return oldMconfigsByGatewayId, nil

	case streaming.CreateOperation, streaming.ReadOperation, streaming.UpdateOperation:
		if update.NewValue == nil {
			return oldMconfigsByGatewayId, nil
		}
		newValueCasted := update.NewValue.(*magmadprotos.MagmadGatewayConfig)

		for _, mconfigValue := range oldMconfigsByGatewayId {
			err := updateMconfig(update.NetworkId, newValueCasted, mconfigValue)
			if err != nil {
				return oldMconfigsByGatewayId, err
			}
		}

		return oldMconfigsByGatewayId, nil

	default:
		return oldMconfigsByGatewayId, fmt.Errorf("Unrecognized streaming operation: %s", update.Operation)
	}
}

func updateMconfig(networkId string, newGwConfig *magmadprotos.MagmadGatewayConfig, mconfigOut *protos.GatewayConfigs) error {
	packageVersion, images, err := getPackageVersionAndImagesForGateway(networkId, newGwConfig.Tier)
	if err != nil {
		return err
	}

	magmadMconfig := &mconfig.MagmaD{
		LogLevel:                protos.LogLevel_INFO,
		CheckinInterval:         newGwConfig.CheckinInterval,
		CheckinTimeout:          newGwConfig.CheckinTimeout,
		AutoupgradeEnabled:      newGwConfig.AutoupgradeEnabled,
		AutoupgradePollInterval: newGwConfig.AutoupgradePollInterval,
		PackageVersion:          packageVersion,
		Images:                  images,
		TierId:                  newGwConfig.Tier,
		DynamicServices:         newGwConfig.DynamicServices,
		FeatureFlags:            newGwConfig.FeatureFlags,
	}

	magmadAny, err := ptypes.MarshalAny(magmadMconfig)
	if err != nil {
		return err
	}
	controlProxyAny, err := ptypes.MarshalAny(&mconfig.ControlProxy{LogLevel: protos.LogLevel_INFO})
	if err != nil {
		return err
	}
	metricsdAny, err := ptypes.MarshalAny(&mconfig.MetricsD{LogLevel: protos.LogLevel_INFO})
	if err != nil {
		return err
	}

	mconfigOut.ConfigsByKey["magmad"] = magmadAny
	mconfigOut.ConfigsByKey["control_proxy"] = controlProxyAny
	mconfigOut.ConfigsByKey["metricsd"] = metricsdAny
	return nil
}
