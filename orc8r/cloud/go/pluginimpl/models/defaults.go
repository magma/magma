/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models

import (
	"magma/orc8r/cloud/go/models"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func NewDefaultDNSConfig() *NetworkDNSConfig {
	return &NetworkDNSConfig{
		EnableCaching: swag.Bool(true),
		LocalTTL:      swag.Uint32(60),
		Records: []*DNSConfigRecord{
			{
				ARecord:     []strfmt.IPv4{"192.88.99.142"},
				AaaaRecord:  []strfmt.IPv6{"2001:0db8:85a3:0000:0000:8a2e:0370:7334"},
				CnameRecord: []string{"cname.example.com"},
				Domain:      "example.com",
			},
		},
	}
}

func NewDefaultFeaturesConfig() *NetworkFeatures {
	return &NetworkFeatures{Features: map[string]string{"foo": "bar"}}
}

func NewDefaultNetwork(networkID string, name string, description string) *Network {
	return &Network{
		ID:          models.NetworkID(networkID),
		Type:        "",
		Name:        models.NetworkName(name),
		Description: models.NetworkDescription(description),
		DNS:         NewDefaultDNSConfig(),
		Features:    NewDefaultFeaturesConfig(),
	}
}
