/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package plugin

import (
	merrors "magma/orc8r/cloud/go/errors"
	"magma/orc8r/cloud/go/services/configurator"
	"orc8r/devmand/cloud/go/devmand"
	"orc8r/devmand/cloud/go/protos/mconfig"
	"orc8r/devmand/cloud/go/services/devmand/obsidian/models"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type Builder struct{}

func (*Builder) Build(networkID string, gatewayID string, graph configurator.EntityGraph, network configurator.Network, mconfigOut map[string]proto.Message) error {
	devmandGW, err := graph.GetEntity(devmand.DevmandGatewayType, gatewayID)
	if err == merrors.ErrNotFound {
		return nil
	}
	if err != nil {
		return errors.WithStack(err)
	}

	devices, err := graph.GetAllChildrenOfType(devmandGW, devmand.DeviceType)
	if err != nil {
		return errors.WithStack(err)
	}

	managedDevices := map[string]*mconfig.ManagedDevice{}
	for _, device := range devices {
		d := device.Config.(*models.ManagedDevice)
		var channels *mconfig.Channels
		if d.Channels != nil {
			var snmpChannel *mconfig.SNMPChannel
			if d.Channels.SnmpChannel != nil {
				s_c := d.Channels.SnmpChannel
				snmpChannel = &mconfig.SNMPChannel{
					Version:   s_c.Version,
					Community: s_c.Community,
				}
			}
			var frinxChannel *mconfig.FrinxChannel
			if d.Channels.FrinxChannel != nil {
				f_c := d.Channels.FrinxChannel
				frinxChannel = &mconfig.FrinxChannel{
					FrinxPort:     f_c.FrinxPort,
					Authorization: f_c.Authorization,
					Host:          f_c.Host,
					Port:          f_c.Port,
					TransportType: f_c.TransportType,
					DeviceType:    f_c.DeviceType,
					DeviceVersion: f_c.DeviceVersion,
					Username:      f_c.Username,
					Password:      f_c.Password,
				}
			}
			var cambiumChannel *mconfig.CambiumChannel
			if d.Channels.CambiumChannel != nil {
				c_c := d.Channels.CambiumChannel
				cambiumChannel = &mconfig.CambiumChannel{
					ClientId:     c_c.ClientID,
					ClientSecret: c_c.ClientSecret,
					ClientMac:    c_c.ClientMac,
					ClientIp:     c_c.ClientIP,
				}
			}
			var otherChannel *mconfig.OtherChannel
			if d.Channels.OtherChannel != nil {
				otherChannel = &mconfig.OtherChannel{
					ChannelProps: d.Channels.OtherChannel.ChannelProps,
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
			DeviceConfig: d.DeviceConfig,
			Host:         d.Host,
			DeviceType:   d.DeviceType,
			Platform:     d.Platform,
			Channels:     channels,
		}
		managedDevices[device.Key] = deviceMconfig
	}

	mconfigOut["devmand"] = &mconfig.DevmandGatewayConfig{
		ManagedDevices: managedDevices,
	}
	return nil
}
