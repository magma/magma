/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config_test

import (
	"testing"

	"magma/orc8r/cloud/go/plugin"
	"magma/orc8r/cloud/go/pluginimpl"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config"
	config_test_init "magma/orc8r/cloud/go/services/config/test_init"
	dnsd_config "magma/orc8r/cloud/go/services/dnsd/config"
	dnsd_protos "magma/orc8r/cloud/go/services/dnsd/protos"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func TestDNSDBuilder_Build(t *testing.T) {
	plugin.RegisterPluginForTests(t, &pluginimpl.BaseOrchestratorPlugin{})
	config_test_init.StartTestService(t)

	builder := &dnsd_config.DnsdMconfigBuilder{}
	actual, err := builder.Build("network", "gw")
	assert.NoError(t, err)
	assert.Equal(t, map[string]proto.Message{}, actual)

	err = config.CreateConfig("network", dnsd_config.DnsdNetworkType, "network", &dnsd_protos.NetworkDNSConfig{EnableCaching: false, LocalTTL: 0})
	assert.NoError(t, err)

	actual, err = builder.Build("network", "gw")
	assert.NoError(t, err)
	expected := map[string]proto.Message{
		"dnsd": &mconfig.DnsD{
			EnableCaching: false,
			LocalTTL:      0,
			Records:       []*mconfig.NetworkDNSConfigRecordsItems{},
		},
	}
	assert.Equal(t, expected, actual)

	networkUpdate := &dnsd_protos.NetworkDNSConfig{
		EnableCaching: true,
	}
	config.UpdateConfig("network", dnsd_config.DnsdNetworkType, "network", networkUpdate)

	actual, err = builder.Build("network", "gw")
	assert.NoError(t, err)
	expected = map[string]proto.Message{
		"dnsd": &mconfig.DnsD{
			EnableCaching: true,
			LocalTTL:      0,
			Records:       []*mconfig.NetworkDNSConfigRecordsItems{},
		},
	}
	assert.Equal(t, expected, actual)
}
