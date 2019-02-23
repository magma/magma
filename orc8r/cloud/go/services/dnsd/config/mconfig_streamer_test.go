/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config_test

import (
	"testing"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config/streaming"
	"magma/orc8r/cloud/go/services/dnsd/config"
	dnsprotos "magma/orc8r/cloud/go/services/dnsd/protos"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestDnsdStreamer_ApplyMconfigUpdate(t *testing.T) {
	s := &config.DnsdStreamer{}

	// Create a network config
	inputMconfigs := map[string]*protos.GatewayConfigs{
		"gw1": {ConfigsByKey: map[string]*any.Any{}},
		"gw2": {ConfigsByKey: map[string]*any.Any{}},
	}
	update := &streaming.ConfigUpdate{
		ConfigType: config.DnsdNetworkType,
		ConfigKey:  "nw",
		NewValue: &dnsprotos.NetworkDNSConfig{
			EnableCaching: true,
			LocalTTL:      1,
			Records: []*dnsprotos.NetworkDNSConfigRecordsItems{
				{
					ARecord:     []string{"A"},
					AaaaRecord:  []string{"aaaa"},
					Domain:      "facebook.com",
					CnameRecord: []string{"cname"},
				},
			},
		},
		Operation: streaming.CreateOperation,
	}

	_, err := s.ApplyMconfigUpdate(update, inputMconfigs)
	assert.NoError(t, err)
	expected := map[string]proto.Message{
		"dnsd": &mconfig.DnsD{
			LogLevel:      protos.LogLevel_INFO,
			EnableCaching: true,
			LocalTTL:      1,
			Records: []*mconfig.NetworkDNSConfigRecordsItems{
				{
					ARecord:     []string{"A"},
					AaaaRecord:  []string{"aaaa"},
					Domain:      "facebook.com",
					CnameRecord: []string{"cname"},
				},
			},
		},
	}
	expectedMconfig := getExpectedMconfig(t, expected)
	assert.Equal(
		t,
		map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedMconfig},
		inputMconfigs,
	)

	// Delete the network config
	update = &streaming.ConfigUpdate{
		ConfigType: config.DnsdNetworkType,
		ConfigKey:  "nw",
		NewValue:   nil,
		Operation:  streaming.DeleteOperation,
	}

	_, err = s.ApplyMconfigUpdate(update, inputMconfigs)
	assert.NoError(t, err)
	expectedMconfig = &protos.GatewayConfigs{ConfigsByKey: map[string]*any.Any{}}
	assert.Equal(t, map[string]*protos.GatewayConfigs{"gw1": expectedMconfig, "gw2": expectedMconfig}, inputMconfigs)
}

func getExpectedMconfig(t *testing.T, expected map[string]proto.Message) *protos.GatewayConfigs {
	ret := &protos.GatewayConfigs{ConfigsByKey: map[string]*any.Any{}}
	for k, v := range expected {
		vAny, err := ptypes.MarshalAny(v)
		assert.NoError(t, err)
		ret.ConfigsByKey[k] = vAny
	}
	return ret
}
