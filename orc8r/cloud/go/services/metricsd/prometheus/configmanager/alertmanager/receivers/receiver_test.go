/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package receivers

import (
	"net/url"
	"strings"
	"testing"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var (
	sampleURL, _ = url.Parse("http://test.com")
	sampleRoute  = config.Route{
		Receiver: "testReceiver",
		Routes: []*config.Route{
			{
				Receiver: "testReceiver",
			},
			{
				Receiver: "slack_receiver",
			},
		},
	}
	sampleReceiver = Receiver{
		Name: "testReceiver",
	}
	sampleSlackReceiver = Receiver{
		Name: "slack_receiver",
		SlackConfigs: []*SlackConfig{{
			APIURL:   "http://slack.com/12345",
			Username: "slack_user",
			Channel:  "slack_alert_channel",
		}},
	}
	sampleWebhookReceiver = Receiver{
		Name: "webhook_receiver",
		WebhookConfigs: []*config.WebhookConfig{{
			URL: &config.URL{
				URL: sampleURL,
			},
			NotifierConfig: config.NotifierConfig{
				VSendResolved: true,
			},
		}},
	}
	sampleEmailReceiver = Receiver{
		Name: "email_receiver",
		EmailConfigs: []*EmailConfig{{
			To:        "test@mail.com",
			From:      "sampleUser",
			Headers:   map[string]string{"header": "value"},
			Smarthost: "http://mail-server.com",
		}},
	}
	sampleConfig = Config{
		Route: &sampleRoute,
		Receivers: []*Receiver{
			&sampleSlackReceiver, &sampleReceiver, &sampleWebhookReceiver, &sampleEmailReceiver,
		},
	}
)

func TestConfig_Validate(t *testing.T) {
	defaultGlobalConf := config.DefaultGlobalConfig()
	validConfig := Config{
		Route:     &sampleRoute,
		Receivers: []*Receiver{&sampleReceiver, &sampleSlackReceiver},
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
	assert.EqualError(t, err, `undefined receiver "testReceiver" used in route`)

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
	assert.EqualError(t, err, `unsupported scheme "" for URL`)

	// Fail if action is missing a type
	invalidSlackAction := Config{
		Route: &config.Route{
			Receiver: "invalidSlackAction",
		},
		Receivers: []*Receiver{{
			Name: "invalidSlackAction",
			SlackConfigs: []*SlackConfig{{
				APIURL: "http://slack.com",
				Actions: []*config.SlackAction{{
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
	rec := sampleConfig.GetReceiver("testReceiver")
	assert.NotNil(t, rec)

	rec = sampleConfig.GetReceiver("slack_receiver")
	assert.NotNil(t, rec)

	rec = sampleConfig.GetReceiver("webhook_receiver")
	assert.NotNil(t, rec)

	rec = sampleConfig.GetReceiver("email_receiver")
	assert.NotNil(t, rec)

	rec = sampleConfig.GetReceiver("nonRoute")
	assert.Nil(t, rec)
}

func TestConfig_GetRouteIdx(t *testing.T) {
	idx := sampleConfig.GetRouteIdx("testReceiver")
	assert.Equal(t, 0, idx)

	idx = sampleConfig.GetRouteIdx("slack_receiver")
	assert.Equal(t, 1, idx)

	idx = sampleConfig.GetRouteIdx("nonRoute")
	assert.Equal(t, -1, idx)
}

func TestReceiver_Secure(t *testing.T) {
	rec := Receiver{Name: "receiverName"}
	rec.Secure(testNID)
	assert.Equal(t, "test_receiverName", rec.Name)
}

func TestReceiver_Unsecure(t *testing.T) {
	rec := Receiver{Name: "receiverName"}
	rec.Secure(testNID)
	assert.Equal(t, "test_receiverName", rec.Name)

	rec.Unsecure(testNID)
	assert.Equal(t, "receiverName", rec.Name)
}

func TestRouteJSONWrapper_ToPrometheusConfig(t *testing.T) {
	jsonRoute := RouteJSONWrapper{
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

	expectedRoute := config.Route{
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

	badGroupWait := RouteJSONWrapper{
		Receiver:  "receiver",
		GroupWait: "abcd",
	}
	route, err = badGroupWait.ToPrometheusConfig()
	assert.EqualError(t, err, "Invalid GroupWait 'abcd': not a valid duration string: \"abcd\"")

	badGroupInterval := RouteJSONWrapper{
		Receiver:      "receiver",
		GroupInterval: "abcd",
	}
	route, err = badGroupInterval.ToPrometheusConfig()
	assert.EqualError(t, err, "Invalid GroupInterval 'abcd': not a valid duration string: \"abcd\"")

	zeroGroupInterval := RouteJSONWrapper{
		Receiver:      "receiver",
		GroupInterval: "0s",
	}
	route, err = zeroGroupInterval.ToPrometheusConfig()
	assert.EqualError(t, err, "GroupInterval cannot be 0")

	badRepeatInterval := RouteJSONWrapper{
		Receiver:       "receiver",
		RepeatInterval: "abcd",
	}
	route, err = badRepeatInterval.ToPrometheusConfig()
	assert.EqualError(t, err, "Invalid RepeatInterval 'abcd': not a valid duration string: \"abcd\"")

	zeroRepeatInterval := RouteJSONWrapper{
		Receiver:       "receiver",
		RepeatInterval: "0s",
	}
	route, err = zeroRepeatInterval.ToPrometheusConfig()
	assert.EqualError(t, err, "RepeatInterval cannot be 0")

	childRoutes := RouteJSONWrapper{
		Receiver: "parent",
		Routes:   []*RouteJSONWrapper{{Receiver: "child1"}, {Receiver: "child2"}},
	}
	route, err = childRoutes.ToPrometheusConfig()
	assert.Equal(t, 2, len(route.Routes))
	assert.NoError(t, err)

	childrenWithErrors := RouteJSONWrapper{
		Receiver: "parent",
		Routes:   []*RouteJSONWrapper{{Receiver: "child", RepeatInterval: "0s"}},
	}
	route, err = childrenWithErrors.ToPrometheusConfig()
	assert.EqualError(t, err, "error converting child route: RepeatInterval cannot be 0")
}

func TestNewRouteJSONWrapper(t *testing.T) {
	fiveSeconds, _ := model.ParseDuration("5s")
	sixSeconds, _ := model.ParseDuration("6s")
	sevenSeconds, _ := model.ParseDuration("7s")

	origRoute := config.Route{
		Receiver:       "receiver",
		GroupByStr:     []string{"groupBy"},
		Match:          map[string]string{"match": "value"},
		Continue:       true,
		GroupWait:      &fiveSeconds,
		GroupInterval:  &sixSeconds,
		RepeatInterval: &sevenSeconds,
		Routes:         []*config.Route{{Receiver: "child"}},
	}

	expectedJSONRoute := RouteJSONWrapper{
		Receiver:       "receiver",
		GroupByStr:     []string{"groupBy"},
		Match:          map[string]string{"match": "value"},
		Continue:       true,
		GroupWait:      "5s",
		GroupInterval:  "6s",
		RepeatInterval: "7s",
		Routes:         []*RouteJSONWrapper{{Receiver: "child"}},
	}
	wrappedRoute := NewRouteJSONWrapper(origRoute)
	assert.Equal(t, expectedJSONRoute, *wrappedRoute)
}

// TestMarshalYamlEmailConfig checks that all EmailConfigs are marshaled with
// requireTLS set to false
func TestMarshalYamlEmailConfig(t *testing.T) {
	valTrue := true
	emailConf := EmailConfig{
		To:         "test@mail.com",
		RequireTLS: &valTrue,
		Headers:    map[string]string{"test": "true", "new": "old"},
	}
	ymlData, err := yaml.Marshal(emailConf)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(ymlData), "require_tls: false"))
	assert.False(t, strings.Contains(string(ymlData), "require_tls: true"))
}
