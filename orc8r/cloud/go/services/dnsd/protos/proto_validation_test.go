/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package protos_test

import (
	"testing"

	"magma/orc8r/cloud/go/services/dnsd/protos"

	"github.com/stretchr/testify/assert"
)

func TestValidateNetworkConfig(t *testing.T) {
	config := &protos.NetworkDNSConfig{EnableCaching: false, LocalTTL: 0}
	err := protos.ValidateNetworkConfig(config)
	assert.NoError(t, err)

	config.Records = []*protos.NetworkDNSConfigRecordsItems{
		{
			ARecord: []string{"192.168.88.99"},
			Domain:  "example.com",
		},
	}
	err = protos.ValidateNetworkConfig(config)
	assert.NoError(t, err)

	config.Records = []*protos.NetworkDNSConfigRecordsItems{
		{
			ARecord: []string{"123456789"},
			Domain:  "example.com",
		},
	}
	err = protos.ValidateNetworkConfig(config)
	assert.Error(t, err)

	aaaaRecord := []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"}
	config.Records = []*protos.NetworkDNSConfigRecordsItems{
		{
			AaaaRecord: aaaaRecord,
			Domain:     "example.com",
		},
	}
	err = protos.ValidateNetworkConfig(config)
	assert.NoError(t, err)

	failedAaaaRecord := []string{"123456789"}
	config.Records = []*protos.NetworkDNSConfigRecordsItems{
		{
			AaaaRecord: failedAaaaRecord,
			Domain:     "example.com",
		},
	}
	err = protos.ValidateNetworkConfig(config)
	assert.Error(t, err)

	config.Records = []*protos.NetworkDNSConfigRecordsItems{
		{
			Domain: "example.com",
		},
	}
	err = protos.ValidateNetworkConfig(config)
	assert.NoError(t, err)

	config.Records = []*protos.NetworkDNSConfigRecordsItems{
		{
			Domain: "",
		},
	}
	err = protos.ValidateNetworkConfig(config)
	assert.Error(t, err)

	// TODO: Test cname records
}
