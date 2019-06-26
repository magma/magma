/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package pluginimpl

import (
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/configurator"
	models3 "magma/orc8r/cloud/go/services/dnsd/obsidian/models"
	models2 "magma/orc8r/cloud/go/services/magmad/obsidian/models"
	"magma/orc8r/cloud/go/services/upgrade/obsidian/models"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type BaseOrchestratorMconfigBuilder struct{}
type DnsdMconfigBuilder struct{}

func (*BaseOrchestratorMconfigBuilder) Build(networkID string, gatewayID string, graph configurator.EntityGraph, network configurator.Network, mconfigOut map[string]proto.Message) error {
	// get magmad gateway - this must be present in the graph
	magmadGateway, err := graph.GetEntity(orc8r.MagmadGatewayType, gatewayID)
	if err == merrors.ErrNotFound {
		return errors.Errorf("could not find magmad gateway %s in graph", gatewayID)
	}
	if err != nil {
		return errors.WithStack(err)
	}

	version, images, err := getPackageVersionAndImages(magmadGateway, &graph)
	if err != nil {
		return errors.WithStack(err)
	}

	magmadGatewayConfig := magmadGateway.Config.(*models2.MagmadGatewayConfig)
	mconfigOut["magmad"] = &mconfig.MagmaD{
		LogLevel:                protos.LogLevel_INFO,
		CheckinInterval:         magmadGatewayConfig.CheckinInterval,
		CheckinTimeout:          magmadGatewayConfig.CheckinTimeout,
		AutoupgradeEnabled:      magmadGatewayConfig.AutoupgradeEnabled,
		AutoupgradePollInterval: magmadGatewayConfig.AutoupgradePollInterval,
		PackageVersion:          version,
		Images:                  images,
		TierId:                  magmadGatewayConfig.Tier,
		DynamicServices:         magmadGatewayConfig.DynamicServices,
		FeatureFlags:            magmadGatewayConfig.FeatureFlags,
	}

	mconfigOut["control_proxy"] = &mconfig.ControlProxy{LogLevel: protos.LogLevel_INFO}
	mconfigOut["metricsd"] = &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO}

	return nil
}

func getPackageVersionAndImages(magmadGateway configurator.NetworkEntity, graph *configurator.EntityGraph) (string, []*mconfig.ImageSpec, error) {
	tier, err := graph.GetFirstAncestorOfType(magmadGateway, orc8r.UpgradeTierEntityType)
	if err == merrors.ErrNotFound {
		return "0.0.0-0", []*mconfig.ImageSpec{}, nil
	}
	if err != nil {
		return "0.0.0-0", []*mconfig.ImageSpec{}, errors.Wrap(err, "failed to load upgrade tier")
	}

	tierConfig := tier.Config.(*models.Tier)
	retImages := make([]*mconfig.ImageSpec, 0, len(tierConfig.Images))
	for _, image := range tierConfig.Images {
		retImages = append(retImages, &mconfig.ImageSpec{Name: image.Name, Order: image.Order})
	}
	return tierConfig.Version, retImages, nil
}

func (*DnsdMconfigBuilder) Build(networkID string, gatewayID string, graph configurator.EntityGraph, network configurator.Network, mconfigOut map[string]proto.Message) error {
	iConfig, found := network.Configs[orc8r.DnsdNetworkType]
	if !found {
		// fill out the dnsd mconfig with an empty struct if no network config
		mconfigOut["dnsd"] = &mconfig.DnsD{}
		return nil
	}

	dnsConfig := iConfig.(*models3.NetworkDNSConfig)
	mconfigDnsd := &mconfig.DnsD{}
	protos.FillIn(dnsConfig, mconfigDnsd)
	mconfigDnsd.LogLevel = protos.LogLevel_INFO
	for _, record := range dnsConfig.Records {
		mconfigRecord := &mconfig.NetworkDNSConfigRecordsItems{}
		protos.FillIn(record, mconfigRecord)
		mconfigDnsd.Records = append(mconfigDnsd.Records, mconfigRecord)
	}

	mconfigOut["dnsd"] = mconfigDnsd
	return nil
}
