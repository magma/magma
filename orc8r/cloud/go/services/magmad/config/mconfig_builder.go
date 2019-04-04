/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"fmt"
	"reflect"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config"
	magmad_protos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/upgrade"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type MagmadMconfigBuilder struct{}

func (builder *MagmadMconfigBuilder) Build(networkId string, gatewayId string) (map[string]proto.Message, error) {
	magmadGatewayConfig, err := getMagmadGatewayConfig(networkId, gatewayId)
	if err != nil {
		return nil, err
	}
	if magmadGatewayConfig == nil {
		return map[string]proto.Message{}, nil
	}

	packageVersion, images, err := getPackageVersionAndImagesForGateway(networkId, magmadGatewayConfig.GetTier())
	if err != nil {
		return nil, err
	}

	return map[string]proto.Message{
		"control_proxy": &mconfig.ControlProxy{LogLevel: protos.LogLevel_INFO},
		"magmad": &mconfig.MagmaD{
			LogLevel:                protos.LogLevel_INFO,
			CheckinInterval:         magmadGatewayConfig.GetCheckinInterval(),
			CheckinTimeout:          magmadGatewayConfig.GetCheckinTimeout(),
			AutoupgradeEnabled:      magmadGatewayConfig.GetAutoupgradeEnabled(),
			AutoupgradePollInterval: magmadGatewayConfig.GetAutoupgradePollInterval(),
			PackageVersion:          packageVersion,
			Images:                  images,
			TierId:                  magmadGatewayConfig.Tier,
			DynamicServices:         magmadGatewayConfig.DynamicServices,
			FeatureFlags:            magmadGatewayConfig.FeatureFlags,
		},
		"metricsd": &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO},
	}, nil
}

func getMagmadGatewayConfig(networkId string, logicalId string) (*magmad_protos.MagmadGatewayConfig, error) {
	iGatewayConfigs, err := config.GetConfig(networkId, MagmadGatewayType, logicalId)
	if err != nil || iGatewayConfigs == nil {
		return nil, err
	}
	gatewayConfigs, ok := iGatewayConfigs.(*magmad_protos.MagmadGatewayConfig)
	if !ok {
		return nil, fmt.Errorf(
			"received unexpected type for gateway record. "+
				"Expected *MagmadGatewayConfig but got %s",
			reflect.TypeOf(iGatewayConfigs),
		)
	}
	return gatewayConfigs, nil
}

// Returns 0.0.0-0 if a nonexistent tier is queried because we don't validate
// tier IDs in magmad configs yet
func getPackageVersionAndImagesForGateway(networkId string, tierId string) (string, []*mconfig.ImageSpec, error) {
	// Load all tiers so the request doesn't error out if we're looking for
	// a nonexistent tier. Tier scale for a network will be small so this
	// should be fine from a performance standpoint.
	tiers, err := upgrade.GetTiers(networkId, []string{})
	if err != nil {
		return "0.0.0-0", []*mconfig.ImageSpec{}, err
	}
	tier, ok := tiers[tierId]
	if !ok {
		glog.V(2).Infof("Unable to load tier %s, returning 0.0.0-0 as target version.", tierId)
		return "0.0.0-0", []*mconfig.ImageSpec{}, nil
	}

	retImages := make([]*mconfig.ImageSpec, 0, len(tier.GetImages()))
	for _, image := range tier.GetImages() {
		retImages = append(retImages, &mconfig.ImageSpec{Name: image.GetName(), Order: image.GetOrder()})
	}
	return tier.GetVersion(), retImages, nil
}
