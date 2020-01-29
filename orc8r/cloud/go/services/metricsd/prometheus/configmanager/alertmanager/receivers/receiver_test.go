/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package receivers_test

import (
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/config"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/receivers"
	tc "magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/test_common"

	"strings"
	"testing"

	amconfig "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const testNID = "test"

func TestConfig_Validate(t *testing.T) {
	defaultGlobalConf := config.DefaultGlobalConfig()
	validConfig := config.Config{
		Route:     &tc.SampleRoute,
		Receivers: []*receivers.Receiver{&tc.SampleReceiver, &tc.SampleSlackReceiver},
		Global:    &defaultGlobalConf,
	}
	err := validConfig.Validate()
	assert.NoError(t, err)

	invalidConfig := config.Config{
		Route:     &tc.SampleRoute,
		Receivers: []*receivers.Receiver{},
		Global:    &defaultGlobalConf,
	}
	err = invalidConfig.Validate()
	assert.EqualError(t, err, `undefined receiver "testReceiver" used in route`)

	invalidSlackReceiver := receivers.Receiver{
		Name: "invalidSlack",
		SlackConfigs: []*receivers.SlackConfig{
			{
				APIURL: "invalidURL",
			},
		},
	}

	invalidSlackConfig := config.Config{
		Route: &amconfig.Route{
			Receiver: "invalidSlack",
		},
		Receivers: []*receivers.Receiver{&invalidSlackReceiver},
		Global:    &defaultGlobalConf,
	}
	err = invalidSlackConfig.Validate()
	assert.EqualError(t, err, `unsupported scheme "" for URL`)

	// Fail if action is missing a type
	invalidSlackAction := config.Config{
		Route: &amconfig.Route{
			Receiver: "invalidSlackAction",
		},
		Receivers: []*receivers.Receiver{{
			Name: "invalidSlackAction",
			SlackConfigs: []*receivers.SlackConfig{{
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
	rec := receivers.Receiver{Name: "receiverName"}
	rec.Secure(testNID)
	assert.Equal(t, "test_receiverName", rec.Name)
}

func TestReceiver_Unsecure(t *testing.T) {
	rec := receivers.Receiver{Name: "receiverName"}
	rec.Secure(testNID)
	assert.Equal(t, "test_receiverName", rec.Name)

	rec.Unsecure(testNID)
	assert.Equal(t, "receiverName", rec.Name)
}

func TestRouteJSONWrapper_ToPrometheusConfig(t *testing.T) {
	jsonRoute := receivers.RouteJSONWrapper{
		Receiver:       "receiver",
		GroupByStr:     []string{"groupBy"},
		Match:          map[string]string{"match": "value"},
		Continue:       true,
		GroupWait:      "5s",
		GroupInterval:  "6s",
		RepeatInterval: "7s",
	}

	fiveSeconds, _ := model.ParseDuration("5s")
	sixSeconds, _ := model.ParseDuration("6s")
	sevenSeconds, _ := model.ParseDuration("7s")

	expectedRoute := amconfig.Route{
		Receiver:       "receiver",
		GroupByStr:     []string{"groupBy"},
		Match:          map[string]string{"match": "value"},
		Continue:       true,
		GroupWait:      &fiveSeconds,
		GroupInterval:  &sixSeconds,
		RepeatInterval: &sevenSeconds,
	}

	route, err := jsonRoute.ToPrometheusConfig()
	assert.NoError(t, err)
	assert.Equal(t, expectedRoute, route)

	badGroupWait := receivers.RouteJSONWrapper{
		Receiver:  "receiver",
		GroupWait: "abcd",
	}
	route, err = badGroupWait.ToPrometheusConfig()
	assert.EqualError(t, err, "Invalid GroupWait 'abcd': not a valid duration string: \"abcd\"")

	badGroupInterval := receivers.RouteJSONWrapper{
		Receiver:      "receiver",
		GroupInterval: "abcd",
	}
	route, err = badGroupInterval.ToPrometheusConfig()
	assert.EqualError(t, err, "Invalid GroupInterval 'abcd': not a valid duration string: \"abcd\"")

	zeroGroupInterval := receivers.RouteJSONWrapper{
		Receiver:      "receiver",
		GroupInterval: "0s",
	}
	route, err = zeroGroupInterval.ToPrometheusConfig()
	assert.EqualError(t, err, "GroupInterval cannot be 0")

	badRepeatInterval := receivers.RouteJSONWrapper{
		Receiver:       "receiver",
		RepeatInterval: "abcd",
	}
	route, err = badRepeatInterval.ToPrometheusConfig()
	assert.EqualError(t, err, "Invalid RepeatInterval 'abcd': not a valid duration string: \"abcd\"")

	zeroRepeatInterval := receivers.RouteJSONWrapper{
		Receiver:       "receiver",
		RepeatInterval: "0s",
	}
	route, err = zeroRepeatInterval.ToPrometheusConfig()
	assert.EqualError(t, err, "RepeatInterval cannot be 0")

	childRoutes := receivers.RouteJSONWrapper{
		Receiver: "parent",
		Routes:   []*receivers.RouteJSONWrapper{{Receiver: "child1"}, {Receiver: "child2"}},
	}
	route, err = childRoutes.ToPrometheusConfig()
	assert.Equal(t, 2, len(route.Routes))
	assert.NoError(t, err)

	childrenWithErrors := receivers.RouteJSONWrapper{
		Receiver: "parent",
		Routes:   []*receivers.RouteJSONWrapper{{Receiver: "child", RepeatInterval: "0s"}},
	}
	route, err = childrenWithErrors.ToPrometheusConfig()
	assert.EqualError(t, err, "error converting child route: RepeatInterval cannot be 0")
}

func TestNewRouteJSONWrapper(t *testing.T) {
	fiveSeconds, _ := model.ParseDuration("5s")
	sixSeconds, _ := model.ParseDuration("6s")
	sevenSeconds, _ := model.ParseDuration("7s")

	origRoute := amconfig.Route{
		Receiver:       "receiver",
		GroupByStr:     []string{"groupBy"},
		Match:          map[string]string{"match": "value"},
		Continue:       true,
		GroupWait:      &fiveSeconds,
		GroupInterval:  &sixSeconds,
		RepeatInterval: &sevenSeconds,
		Routes:         []*amconfig.Route{{Receiver: "child"}},
	}

	expectedJSONRoute := receivers.RouteJSONWrapper{
		Receiver:       "receiver",
		GroupByStr:     []string{"groupBy"},
		Match:          map[string]string{"match": "value"},
		Continue:       true,
		GroupWait:      "5s",
		GroupInterval:  "6s",
		RepeatInterval: "7s",
		Routes:         []*receivers.RouteJSONWrapper{{Receiver: "child"}},
	}
	wrappedRoute := receivers.NewRouteJSONWrapper(origRoute)
	assert.Equal(t, expectedJSONRoute, *wrappedRoute)
}

// TestMarshalYamlEmailConfig checks that all EmailConfigs are marshaled with
// requireTLS set to false
func TestMarshalYamlEmailConfig(t *testing.T) {
	valTrue := true
	emailConf := receivers.EmailConfig{
		To:         "test@mail.com",
		RequireTLS: &valTrue,
		Headers:    map[string]string{"test": "true", "new": "old"},
	}
	ymlData, err := yaml.Marshal(emailConf)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(ymlData), "require_tls: false"))
	assert.False(t, strings.Contains(string(ymlData), "require_tls: true"))
}
