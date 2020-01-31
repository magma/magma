/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package config

import (
	"fmt"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/common"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/receivers"

	amconfig "github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

const (
	TenantBaseRoutePostfix = "tenant_base_route"
)

// Config uses a custom receiver struct to avoid scrubbing of 'secrets' during
// marshaling
type Config struct {
	Global       *GlobalConfig           `yaml:"global,omitempty" json:"global,omitempty"`
	Route        *amconfig.Route         `yaml:"route,omitempty" json:"route,omitempty"`
	InhibitRules []*amconfig.InhibitRule `yaml:"inhibit_rules,omitempty" json:"inhibit_rules,omitempty"`
	Receivers    []*receivers.Receiver   `yaml:"receivers,omitempty" json:"receivers,omitempty"`
	Templates    []string                `yaml:"templates" json:"templates"`
}

// GetReceiver returns the receiver config with the given name
func (c *Config) GetReceiver(name string) *receivers.Receiver {
	for _, rec := range c.Receivers {
		if rec.Name == name {
			return rec
		}
	}
	return nil
}

func (c *Config) GetRouteIdx(name string) int {
	for idx, route := range c.Route.Routes {
		if route.Receiver == name {
			return idx
		}
	}
	return -1
}

func (c *Config) InitializeNetworkBaseRoute(route *amconfig.Route, matcherLabel, tenantID string) error {
	baseRouteName := MakeBaseRouteName(tenantID)
	if c.GetReceiver(baseRouteName) != nil {
		return fmt.Errorf("Base route for tenant %s already exists", tenantID)
	}

	c.Receivers = append(c.Receivers, &receivers.Receiver{Name: baseRouteName})
	route.Receiver = baseRouteName

	if matcherLabel != "" {
		route.Match = map[string]string{matcherLabel: tenantID}
	}

	c.Route.Routes = append(c.Route.Routes, route)

	return c.Validate()
}

// Validate makes sure that the config is properly formed. Unmarshal the yaml
// data into an alertmanager Config struct to ensure that it is properly formed
func (c *Config) Validate() error {
	yamlData, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlData, &amconfig.Config{})
	if err != nil {
		return err
	}
	return nil
}

// GlobalConfig is a copy of prometheus/alertmanager/config.GlobalConfig with
// `Secret` fields replaced with strings to enable marshaling without obfuscation
type GlobalConfig struct {
	// ResolveTimeout is the time after which an alert is declared resolved
	// if it has not been updated.
	ResolveTimeout string `yaml:"resolve_timeout" json:"resolve_timeout"`

	HTTPConfig *common.HTTPConfig `yaml:"http_config,omitempty" json:"http_config,omitempty"`

	SMTPFrom         string        `yaml:"smtp_from,omitempty" json:"smtp_from,omitempty"`
	SMTPHello        string        `yaml:"smtp_hello,omitempty" json:"smtp_hello,omitempty"`
	SMTPSmarthost    string        `yaml:"smtp_smarthost,omitempty" json:"smtp_smarthost,omitempty"`
	SMTPAuthUsername string        `yaml:"smtp_auth_username,omitempty" json:"smtp_auth_username,omitempty"`
	SMTPAuthPassword string        `yaml:"smtp_auth_password,omitempty" json:"smtp_auth_password,omitempty"`
	SMTPAuthSecret   string        `yaml:"smtp_auth_secret,omitempty" json:"smtp_auth_secret,omitempty"`
	SMTPAuthIdentity string        `yaml:"smtp_auth_identity,omitempty" json:"smtp_auth_identity,omitempty"`
	SMTPRequireTLS   bool          `yaml:"smtp_require_tls,omitempty" json:"smtp_require_tls,omitempty"`
	SlackAPIURL      *amconfig.URL `yaml:"slack_api_url,omitempty" json:"slack_api_url,omitempty"`
	PagerdutyURL     *amconfig.URL `yaml:"pagerduty_url,omitempty" json:"pagerduty_url,omitempty"`
	HipchatAPIURL    *amconfig.URL `yaml:"hipchat_api_url,omitempty" json:"hipchat_api_url,omitempty"`
	HipchatAuthToken string        `yaml:"hipchat_auth_token,omitempty" json:"hipchat_auth_token,omitempty"`
	OpsGenieAPIURL   *amconfig.URL `yaml:"opsgenie_api_url,omitempty" json:"opsgenie_api_url,omitempty"`
	OpsGenieAPIKey   string        `yaml:"opsgenie_api_key,omitempty" json:"opsgenie_api_key,omitempty"`
	WeChatAPIURL     *amconfig.URL `yaml:"wechat_api_url,omitempty" json:"wechat_api_url,omitempty"`
	WeChatAPISecret  string        `yaml:"wechat_api_secret,omitempty" json:"wechat_api_secret,omitempty"`
	WeChatAPICorpID  string        `yaml:"wechat_api_corp_id,omitempty" json:"wechat_api_corp_id,omitempty"`
	VictorOpsAPIURL  *amconfig.URL `yaml:"victorops_api_url,omitempty" json:"victorops_api_url,omitempty"`
	VictorOpsAPIKey  string        `yaml:"victorops_api_key,omitempty" json:"victorops_api_key,omitempty"`
}

func DefaultGlobalConfig() GlobalConfig {
	return GlobalConfig{
		ResolveTimeout: model.Duration(5 * time.Minute).String(),
		HTTPConfig:     &common.HTTPConfig{},

		SMTPHello:      "localhost",
		SMTPRequireTLS: false,
	}
}

func MakeBaseRouteName(tenantID string) string {
	return fmt.Sprintf("%s_%s", tenantID, TenantBaseRoutePostfix)
}
