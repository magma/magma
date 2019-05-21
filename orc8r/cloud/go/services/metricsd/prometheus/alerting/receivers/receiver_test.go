/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package receivers

import (
	"testing"

	"github.com/prometheus/alertmanager/config"
	"github.com/stretchr/testify/assert"
)

var (
	sampleRoute = config.Route{
		Receiver: "testReceiver",
	}
	sampleReceiver = Receiver{
		Name: "testReceiver",
	}
)

func TestConfig_Validate(t *testing.T) {
	defaultGlobalConf := config.DefaultGlobalConfig()
	validConfig := Config{
		Route:     &sampleRoute,
		Receivers: []*Receiver{&sampleReceiver},
		Global:    &defaultGlobalConf,
	}

	err := validConfig.Validate()
	assert.NoError(t, err)

	invalidConfig := Config{
		Route:     &sampleRoute,
		Receivers: []*Receiver{},
		Global:    &defaultGlobalConf,
	}

	err = invalidConfig.Validate()
	assert.Error(t, err)

	invalidSlackReceiver := Receiver{
		Name: "invalidSlack",
		SlackConfigs: []*SlackConfig{
			{
				APIURL: "invalidURL",
			},
		},
	}

	invalidSlackConfig := Config{
		Route: &config.Route{
			Receiver: "invalidSlack",
		},
		Receivers: []*Receiver{&invalidSlackReceiver},
		Global:    &defaultGlobalConf,
	}
	err = invalidSlackConfig.Validate()
	assert.Error(t, err)
}
