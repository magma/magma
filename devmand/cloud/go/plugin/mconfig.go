/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package plugin

import (
	"magma/orc8r/cloud/go/services/configurator"
	configuratorprotos "magma/orc8r/cloud/go/services/configurator/protos"
	merrors "magma/orc8r/lib/go/errors"
	"orc8r/devmand/cloud/go/devmand"
	models2 "orc8r/devmand/cloud/go/plugin/models"
	"orc8r/devmand/cloud/go/protos/mconfig"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/pkg/errors"
)

type Builder struct{}
type DevmandMconfigBuilderServicer struct{}

func (s *DevmandMconfigBuilderServicer) Build(
	request *configuratorprotos.BuildMconfigRequest,
) (*configuratorprotos.BuildMconfigResponse, error) {
	ret := &configuratorprotos.BuildMconfigResponse{
		ConfigsByKey: map[string]*any.Any{},
	}
	graph, err := (configurator.EntityGraph{}).FromStorageProto(request.GetEntityGraph())
	if err != nil {
		return ret, err
	}
	devmandAgent, err := graph.GetEntity(devmand.SymphonyAgentType, request.GetGatewayId())
	if err == merrors.ErrNotFound {
		return ret, nil
	}
	if err != nil {
		return ret, err
	}

	devices, err := graph.GetAllChildrenOfType(devmandAgent, devmand.SymphonyDeviceType)
	if err != nil {
		return ret, err
	}

	managedDevices := map[string]*mconfig.ManagedDevice{}
	for _, device := range devices {
		d_config := device.Config.(*models2.SymphonyDeviceConfig)
		var channels *mconfig.Channels
		if d_config.Channels != nil {
			var snmpChannel *mconfig.SNMPChannel
			if d_config.Channels.SnmpChannel != nil {
				s_c := d_config.Channels.SnmpChannel
				snmpChannel = &mconfig.SNMPChannel{
					Community: s_c.Community,
					Version:   s_c.Version,
				}
			}
			var frinxChannel *mconfig.FrinxChannel
			if d_config.Channels.FrinxChannel != nil {
				f_c := d_config.Channels.FrinxChannel
				frinxChannel = &mconfig.FrinxChannel{
					Authorization: f_c.Authorization,
					DeviceType:    f_c.DeviceType,
					DeviceVersion: f_c.DeviceVersion,
					FrinxPort:     f_c.FrinxPort,
					Host:          f_c.Host,
					Password:      f_c.Password,
					Port:          f_c.Port,
					TransportType: f_c.TransportType,
					Username:      f_c.Username,
				}
			}
			var cambiumChannel *mconfig.CambiumChannel
			if d_config.Channels.CambiumChannel != nil {
				c_c := d_config.Channels.CambiumChannel
				cambiumChannel = &mconfig.CambiumChannel{
					ClientId:     c_c.ClientID,
					ClientIp:     c_c.ClientIP,
					ClientMac:    c_c.ClientMac,
					ClientSecret: c_c.ClientSecret,
				}
			}
			var otherChannel *mconfig.OtherChannel
			if d_config.Channels.OtherChannel != nil {
				otherChannel = &mconfig.OtherChannel{
					ChannelProps: d_config.Channels.OtherChannel.ChannelProps,
				}
			}
			channels = &mconfig.Channels{
				SnmpChannel:    snmpChannel,
				FrinxChannel:   frinxChannel,
				CambiumChannel: cambiumChannel,
				OtherChannel:   otherChannel,
			}
		}

		deviceMconfig := &mconfig.ManagedDevice{
			DeviceConfig: d_config.DeviceConfig,
			Host:         d_config.Host,
			DeviceType:   d_config.DeviceType,
			Platform:     d_config.Platform,
			Channels:     channels,
		}
		managedDevices[device.Key] = deviceMconfig
	}
	devmandMconfig := &mconfig.DevmandGatewayConfig{
		ManagedDevices: managedDevices,
	}
	ret.ConfigsByKey["devmand"], err = ptypes.MarshalAny(devmandMconfig)
	return ret, err
}

func (*Builder) Build(networkID string, gatewayID string, graph configurator.EntityGraph, network configurator.Network, mconfigOut map[string]proto.Message) error {
	servicer := &DevmandMconfigBuilderServicer{}
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
