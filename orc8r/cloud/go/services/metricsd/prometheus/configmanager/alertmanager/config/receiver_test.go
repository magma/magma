/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package config_test

import (
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/config"
	tc "magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/test_common"

	"strings"
	"testing"

	amconfig "github.com/prometheus/alertmanager/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const testNID = "test"

func TestConfig_Validate(t *testing.T) {
	defaultGlobalConf := config.DefaultGlobalConfig()
	validConfig := config.Config{
		Route:     &tc.SampleRoute,
		Receivers: []*config.Receiver{&tc.SampleReceiver, &tc.SampleSlackReceiver},
		Global:    &defaultGlobalConf,
	}
	err := validConfig.Validate()
	assert.NoError(t, err)

	invalidConfig := config.Config{
		Route:     &tc.SampleRoute,
		Receivers: []*config.Receiver{},
		Global:    &defaultGlobalConf,
	}
	err = invalidConfig.Validate()
	assert.EqualError(t, err, `undefined receiver "testReceiver" used in route`)

	invalidSlackReceiver := config.Receiver{
		Name: "invalidSlack",
		SlackConfigs: []*config.SlackConfig{
			{
				APIURL: "invalidURL",
			},
		},
	}

	invalidSlackConfig := config.Config{
		Route: &config.Route{
			Receiver: "invalidSlack",
		},
		Receivers: []*config.Receiver{&invalidSlackReceiver},
		Global:    &defaultGlobalConf,
	}
	err = invalidSlackConfig.Validate()
	assert.EqualError(t, err, `unsupported scheme "" for URL`)

	// Fail if action is missing a type
	invalidSlackAction := config.Config{
		Route: &config.Route{
			Receiver: "invalidSlackAction",
		},
		Receivers: []*config.Receiver{{
			Name: "invalidSlackAction",
			SlackConfigs: []*config.SlackConfig{{
				APIURL: "http://slack.com",
				Actions: []*amconfig.SlackAction{{
					URL:  "test.com",
					Text: "test",
				}},
			}},
		}},
	}
	err = invalidSlackAction.Validate()
	assert.EqualError(t, err, `missing type in Slack action configuration`)
}

func TestConfig_GetReceiver(t *testing.T) {
	rec := tc.SampleConfig.GetReceiver("testReceiver")
	assert.NotNil(t, rec)

	rec = tc.SampleConfig.GetReceiver("slack_receiver")
	assert.NotNil(t, rec)

	rec = tc.SampleConfig.GetReceiver("webhook_receiver")
	assert.NotNil(t, rec)

	rec = tc.SampleConfig.GetReceiver("email_receiver")
	assert.NotNil(t, rec)

	rec = tc.SampleConfig.GetReceiver("nonRoute")
	assert.Nil(t, rec)
}

func TestConfig_GetRouteIdx(t *testing.T) {
	idx := tc.SampleConfig.GetRouteIdx("testReceiver")
	assert.Equal(t, 0, idx)

	idx = tc.SampleConfig.GetRouteIdx("slack_receiver")
	assert.Equal(t, 1, idx)

	idx = tc.SampleConfig.GetRouteIdx("nonRoute")
	assert.Equal(t, -1, idx)
}

func TestReceiver_Secure(t *testing.T) {
	rec := config.Receiver{Name: "receiverName"}
	rec.Secure(testNID)
	assert.Equal(t, "test_receiverName", rec.Name)
}

func TestReceiver_Unsecure(t *testing.T) {
	rec := config.Receiver{Name: "receiverName"}
	rec.Secure(testNID)
	assert.Equal(t, "test_receiverName", rec.Name)

	rec.Unsecure(testNID)
	assert.Equal(t, "receiverName", rec.Name)
}

// TestMarshalYamlEmailConfig checks that all EmailConfigs are marshaled with
// requireTLS set to false
func TestMarshalYamlEmailConfig(t *testing.T) {
	valTrue := true
	emailConf := config.EmailConfig{
		To:         "test@mail.com",
		RequireTLS: &valTrue,
		Headers:    map[string]string{"test": "true", "new": "old"},
	}
	ymlData, err := yaml.Marshal(emailConf)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(ymlData), "require_tls: false"))
	assert.False(t, strings.Contains(string(ymlData), "require_tls: true"))
}
