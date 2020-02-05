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
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
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

	if magmadGateway.Config != nil {
		magmadGatewayConfig := magmadGateway.Config.(*models.MagmadGatewayConfigs)
		mconfigOut["magmad"] = &mconfig.MagmaD{
			LogLevel:                protos.LogLevel_INFO,
			CheckinInterval:         int32(magmadGatewayConfig.CheckinInterval),
			CheckinTimeout:          int32(magmadGatewayConfig.CheckinTimeout),
			AutoupgradeEnabled:      swag.BoolValue(magmadGatewayConfig.AutoupgradeEnabled),
			AutoupgradePollInterval: magmadGatewayConfig.AutoupgradePollInterval,
			PackageVersion:          version,
			Images:                  images,
			DynamicServices:         magmadGatewayConfig.DynamicServices,
			FeatureFlags:            magmadGatewayConfig.FeatureFlags,
		}

		mconfigOut["td-agent-bit"] = getFluentBitMconfig(networkID, gatewayID, magmadGatewayConfig)
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
		retImages = append(retImages, &mconfig.ImageSpec{Name: swag.StringValue(image.Name), Order: swag.Int64Value(image.Order)})
	}
	return tierConfig.Version.ToString(), retImages, nil
}

func getFluentBitMconfig(networkID string, gatewayID string, mdGw *models.MagmadGatewayConfigs) *mconfig.FluentBit {
	ret := &mconfig.FluentBit{
		ExtraTags: map[string]string{
			"network_id": networkID,
			"gateway_id": gatewayID,
		},
		ThrottleRate:     1000,
		ThrottleWindow:   5,
		ThrottleInterval: "1m",
	}

	if mdGw.Logging != nil && mdGw.Logging.Aggregation != nil {
		ret.FilesByTag = mdGw.Logging.Aggregation.TargetFilesByTag
		if mdGw.Logging.Aggregation.ThrottleRate != nil {
			ret.ThrottleRate = *mdGw.Logging.Aggregation.ThrottleRate
		}
		if mdGw.Logging.Aggregation.ThrottleWindow != nil {
			ret.ThrottleWindow = *mdGw.Logging.Aggregation.ThrottleWindow
		}
		if mdGw.Logging.Aggregation.ThrottleInterval != nil {
			ret.ThrottleInterval = *mdGw.Logging.Aggregation.ThrottleInterval
		}
	}

	return ret
}

func (*DnsdMconfigBuilder) Build(networkID string, gatewayID string, graph configurator.EntityGraph, network configurator.Network, mconfigOut map[string]proto.Message) error {
	iConfig, found := network.Configs[orc8r.DnsdNetworkType]
	if !found {
		// fill out the dnsd mconfig with an empty struct if no network config
		mconfigOut["dnsd"] = &mconfig.DnsD{}
		return nil
	}

	dnsConfig := iConfig.(*models.NetworkDNSConfig)
	mconfigDnsd := &mconfig.DnsD{}
	protos.FillIn(dnsConfig, mconfigDnsd)
	mconfigDnsd.LocalTTL = int32(swag.Uint32Value(dnsConfig.LocalTTL))
	mconfigDnsd.EnableCaching = swag.BoolValue(dnsConfig.EnableCaching)
	mconfigDnsd.LogLevel = protos.LogLevel_INFO
	for _, record := range dnsConfig.Records {
		mconfigRecord := &mconfig.NetworkDNSConfigRecordsItems{}
		protos.FillIn(record, mconfigRecord)
		mconfigRecord.ARecord = funk.Map(record.ARecord, func(a strfmt.IPv4) string { return string(a) }).([]string)
		mconfigRecord.AaaaRecord = funk.Map(record.AaaaRecord, func(a strfmt.IPv6) string { return string(a) }).([]string)
		mconfigDnsd.Records = append(mconfigDnsd.Records, mconfigRecord)
	}

	mconfigOut["dnsd"] = mconfigDnsd
	return nil
}
