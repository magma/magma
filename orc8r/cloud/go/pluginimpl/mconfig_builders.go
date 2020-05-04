/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package pluginimpl

import (
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"
	configuratorprotos "magma/orc8r/cloud/go/services/configurator/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/protos/mconfig"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

type BaseOrchestratorMconfigBuilder struct{}
type DnsdMconfigBuilder struct{}

type BaseOrchestratorMconfigBuilderServicer struct{}
type DnsdMconfigBuilderServicer struct{}

func (s *BaseOrchestratorMconfigBuilderServicer) Build(
	request *configuratorprotos.BuildMconfigRequest,
) (*configuratorprotos.BuildMconfigResponse, error) {
	ret := &configuratorprotos.BuildMconfigResponse{
		ConfigsByKey: map[string]*any.Any{},
	}
	graph, err := (configurator.EntityGraph{}).FromStorageProto(request.GetEntityGraph())
	if err != nil {
		return ret, err
	}
	// get magmad gateway - this must be present in the graph
	magmadGateway, err := graph.GetEntity(orc8r.MagmadGatewayType, request.GatewayId)
	if err == merrors.ErrNotFound {
		return ret, errors.Errorf("could not find magmad gateway %s in graph", request.GetGatewayId())
	}
	if err != nil {
		return ret, err
	}

	version, images, err := getPackageVersionAndImages(magmadGateway, &graph)
	if err != nil {
		return ret, err
	}

	if magmadGateway.Config != nil {
		magmadGatewayConfig := magmadGateway.Config.(*models.MagmadGatewayConfigs)
		magmadMconfig := &mconfig.MagmaD{
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
		ret.ConfigsByKey["magmad"], err = ptypes.MarshalAny(magmadMconfig)
		if err != nil {
			return ret, err
		}
		fluentBitMconfig := getFluentBitMconfig(request.GetNetworkId(), request.GetGatewayId(), magmadGatewayConfig)
		ret.ConfigsByKey["td-agent-bit"], err = ptypes.MarshalAny(fluentBitMconfig)
		if err != nil {
			return ret, err
		}
	}
	controlProxyMconfig := &mconfig.ControlProxy{LogLevel: protos.LogLevel_INFO}
	ret.ConfigsByKey["control_proxy"], err = ptypes.MarshalAny(controlProxyMconfig)
	if err != nil {
		return ret, err
	}
	metricsdMconfig := &mconfig.MetricsD{LogLevel: protos.LogLevel_INFO}
	ret.ConfigsByKey["metricsd"], err = ptypes.MarshalAny(metricsdMconfig)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func (*BaseOrchestratorMconfigBuilder) Build(networkID string, gatewayID string, graph configurator.EntityGraph, network configurator.Network, mconfigOut map[string]proto.Message) error {
	servicer := &BaseOrchestratorMconfigBuilderServicer{}
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

	magmadAnyMconfig, ok := res.GetConfigsByKey()["magmad"]
	if ok {
		magmadMconfig := &mconfig.MagmaD{}
		err = ptypes.UnmarshalAny(magmadAnyMconfig, magmadMconfig)
		if err != nil {
			return errors.WithStack(err)
		}
		mconfigOut["magmad"] = magmadMconfig
	}
	fluentBitAnyMconfig, ok := res.GetConfigsByKey()["td-agent-bit"]
	if ok {
		fluentBitMconfig := &mconfig.FluentBit{}
		err = ptypes.UnmarshalAny(fluentBitAnyMconfig, fluentBitMconfig)
		if err != nil {
			return errors.WithStack(err)
		}
		mconfigOut["td-agent-bit"] = fluentBitMconfig
	}
	controlProxyMconfig := &mconfig.ControlProxy{}
	err = ptypes.UnmarshalAny(res.GetConfigsByKey()["control_proxy"], controlProxyMconfig)
	if err != nil {
		return errors.WithStack(err)
	}
	mconfigOut["control_proxy"] = controlProxyMconfig
	metricsdMconfig := &mconfig.MetricsD{}
	err = ptypes.UnmarshalAny(res.GetConfigsByKey()["metricsd"], metricsdMconfig)
	if err != nil {
		return errors.WithStack(err)
	}
	mconfigOut["metricsd"] = metricsdMconfig

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

func (s *DnsdMconfigBuilderServicer) Build(
	request *configuratorprotos.BuildMconfigRequest,
) (*configuratorprotos.BuildMconfigResponse, error) {
	ret := &configuratorprotos.BuildMconfigResponse{
		ConfigsByKey: map[string]*any.Any{},
	}
	network, err := (configurator.Network{}).FromStorageProto(request.GetNetwork())
	if err != nil {
		return ret, err
	}
	iConfig, found := network.Configs[orc8r.DnsdNetworkType]
	if !found {
		// fill out the dnsd mconfig with an empty struct if no network config
		ret.ConfigsByKey["dnsd"], err = ptypes.MarshalAny(&mconfig.DnsD{})
		return ret, err
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

	ret.ConfigsByKey["dnsd"], err = ptypes.MarshalAny(mconfigDnsd)
	return ret, err
}

func (*DnsdMconfigBuilder) Build(networkID string, gatewayID string, graph configurator.EntityGraph, network configurator.Network, mconfigOut map[string]proto.Message) error {
	servicer := &DnsdMconfigBuilderServicer{}
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

	dnsdMconfig := &mconfig.DnsD{}
	err = ptypes.UnmarshalAny(res.GetConfigsByKey()["dnsd"], dnsdMconfig)
	if err != nil {
		return errors.WithStack(err)
	}
	mconfigOut["dnsd"] = dnsdMconfig
	return nil
}
