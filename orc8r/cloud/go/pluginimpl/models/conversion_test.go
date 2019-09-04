/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package models_test

import (
	"testing"

	models1 "magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/pluginimpl/models"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/stretchr/testify/assert"
)

func Test_Conversions(t *testing.T) {
	cNetwork := configurator.Network{
		ID:          "test",
		Name:        "name",
		Type:        "type",
		Description: "desc",
		Configs: map[string]interface{}{
			orc8r.NetworkFeaturesConfig: models.NewDefaultFeaturesConfig(),
			orc8r.DnsdNetworkType:       models.NewDefaultDNSConfig(),
		},
	}
	generatedSNetwork := (&models.Network{}).FromConfiguratorNetwork(cNetwork)
	sNetwork := models.Network{
		ID:          models1.NetworkID("test"),
		Name:        models1.NetworkName("name"),
		Type:        "type",
		Description: models1.NetworkDescription("desc"),
		Features:    models.NewDefaultFeaturesConfig(),
		DNS:         models.NewDefaultDNSConfig(),
	}
	generatedCNetwork := sNetwork.ToConfiguratorNetwork()

	assert.Equal(t, sNetwork, *generatedSNetwork)
	assert.Equal(t, cNetwork, generatedCNetwork)
}
