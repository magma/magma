/*
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
*/

package models

func NewDefaultManagedDeviceModel() *ManagedDevice {
	return &ManagedDevice{
		DeviceConfig: "config_json",
		DeviceType:   []string{"type_descriptor_1"},
		Channels:     NewDefaultChannels(),
		Host:         "hostname",
		Platform:     "platform_name",
	}
}

func NewDefaultCambiumChannel() *CambiumChannel {
	return &CambiumChannel{}
}

func NewDefaultChannels() *ManagedDeviceChannels {
	return &ManagedDeviceChannels{
		SnmpChannel:    NewDefaultSnmpChannel(),
		OtherChannel:   NewDefaultOtherChannel(),
		FrinxChannel:   NewDefaultFrinxChannel(),
		CambiumChannel: NewDefaultCambiumChannel(),
	}
}

func NewDefaultFrinxChannel() *FrinxChannel {
	return &FrinxChannel{}
}

func NewDefaultOtherChannel() *OtherChannel {
	return &OtherChannel{
		ChannelProps: map[string]string{},
	}
}

func NewDefaultSnmpChannel() *SnmpChannel {
	return &SnmpChannel{}
}
