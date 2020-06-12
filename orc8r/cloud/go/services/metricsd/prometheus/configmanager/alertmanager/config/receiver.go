/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package config

import (
	"strings"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/common"

	"github.com/prometheus/alertmanager/config"

	"github.com/prometheus/common/model"
)

// Receiver uses custom notifier configs to allow for marshaling of secrets.
type Receiver struct {
	Name string `yaml:"name" json:"name"`

	SlackConfigs     []*SlackConfig     `yaml:"slack_configs,omitempty" json:"slack_configs,omitempty"`
	WebhookConfigs   []*WebhookConfig   `yaml:"webhook_configs,omitempty" json:"webhook_configs,omitempty"`
	EmailConfigs     []*EmailConfig     `yaml:"email_configs,omitempty" json:"email_configs,omitempty"`
	PagerDutyConfigs []*PagerDutyConfig `yaml:"pagerduty_configs,omitempty" json:"pagerduty_configs,omitempty"`
	PushoverConfigs  []*PushoverConfig  `yaml:"pushover_configs,omitempty" json:"pushover_configs,omitempty"`
}

// ReceiverJSONWrapper uses custom (JSON compatible) notifier configs to allow
// for marshaling of secrets.
type ReceiverJSONWrapper struct {
	Name string `yaml:"name" json:"name"`

	SlackConfigs     []*SlackConfig         `yaml:"slack_configs,omitempty" json:"slack_configs,omitempty"`
	WebhookConfigs   []*WebhookConfig       `yaml:"webhook_configs,omitempty" json:"webhook_configs,omitempty"`
	EmailConfigs     []*EmailConfig         `yaml:"email_configs,omitempty" json:"email_configs,omitempty"`
	PagerDutyConfigs []*PagerDutyConfig     `yaml:"pagerduty_configs,omitempty" json:"pagerduty_configs,omitempty"`
	PushoverConfigs  []*PushoverJSONWrapper `yaml:"pushover_configs,omitempty" json:"pushover_configs,omitempty"`
}

// Secure replaces the receiver's name with a tenantID prefix
func (r *Receiver) Secure(tenantID string) {
	r.Name = SecureReceiverName(r.Name, tenantID)
}

// Unsecure removes the tenantID prefix from the receiver name
func (r *Receiver) Unsecure(tenantID string) {
	r.Name = UnsecureReceiverName(r.Name, tenantID)
}

func SecureReceiverName(name, tenantID string) string {
	return ReceiverTenantPrefix(tenantID) + name
}

func UnsecureReceiverName(name, tenantID string) string {
	if strings.HasPrefix(name, ReceiverTenantPrefix(tenantID)) {
		return name[len(ReceiverTenantPrefix(tenantID)):]
	}
	return name
}

// SlackConfig uses string instead of SecretURL for the APIURL field so that it
// is marshaled as is instead of being obscured which is how alertmanager handles
// secrets
type SlackConfig struct {
	HTTPConfig *common.HTTPConfig `yaml:"http_config,omitempty" json:"http_config,omitempty"`

	APIURL      string                `yaml:"api_url" json:"api_url"`
	Channel     string                `yaml:"channel" json:"channel"`
	Username    string                `yaml:"username" json:"username"`
	Color       string                `yaml:"color,omitempty" json:"color,omitempty"`
	Title       string                `yaml:"title,omitempty" json:"title,omitempty"`
	TitleLink   string                `yaml:"title_link,omitempty" json:"title_link,omitempty"`
	Pretext     string                `yaml:"pretext,omitempty" json:"pretext,omitempty"`
	Text        string                `yaml:"text,omitempty" json:"text,omitempty"`
	Fields      []*config.SlackField  `yaml:"fields,omitempty" json:"fields,omitempty"`
	ShortFields bool                  `yaml:"short_fields,omitempty" json:"short_fields,omitempty"`
	Footer      string                `yaml:"footer,omitempty" json:"footer,omitempty"`
	Fallback    string                `yaml:"fallback,omitempty" json:"fallback,omitempty"`
	CallbackID  string                `yaml:"callback_id,omitempty" json:"callback_id,omitempty"`
	IconEmoji   string                `yaml:"icon_emoji,omitempty" json:"icon_emoji,omitempty"`
	IconURL     string                `yaml:"icon_url,omitempty" json:"icon_url,omitempty"`
	ImageURL    string                `yaml:"image_url,omitempty" json:"image_url,omitempty"`
	ThumbURL    string                `yaml:"thumb_url,omitempty" json:"thumb_url,omitempty"`
	LinkNames   bool                  `yaml:"link_names,omitempty" json:"link_names,omitempty"`
	Actions     []*config.SlackAction `yaml:"actions,omitempty" json:"actions,omitempty"`
}

// EmailConfig uses string instead of Secret for the AuthPassword and AuthSecret
// field so that it is marshaled as is instead of being obscured which is how
// alertmanager handles secrets. Otherwise the secrets would be obscured on write
// to the yml file, making it unusable.
type EmailConfig struct {
	config.NotifierConfig `yaml:",inline" json:",inline"`

	To           string            `yaml:"to,omitempty" json:"to,omitempty"`
	From         string            `yaml:"from,omitempty" json:"from,omitempty"`
	Hello        string            `yaml:"hello,omitempty" json:"hello,omitempty"`
	Smarthost    string            `yaml:"smarthost,omitempty" json:"smarthost,omitempty"`
	AuthUsername string            `yaml:"auth_username,omitempty" json:"auth_username,omitempty"`
	AuthPassword string            `yaml:"auth_password,omitempty" json:"auth_password,omitempty"`
	AuthSecret   string            `yaml:"auth_secret,omitempty" json:"auth_secret,omitempty"`
	AuthIdentity string            `yaml:"auth_identity,omitempty" json:"auth_identity,omitempty"`
	Headers      map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	HTML         string            `yaml:"html,omitempty" json:"html,omitempty"`
	Text         string            `yaml:"text,omitempty" json:"text,omitempty"`
	RequireTLS   *bool             `yaml:"require_tls,omitempty" json:"require_tls,omitempty"`
}

// PagerDutyConfig uses string instead of Secret for the RoutingKey and ServiceKey
// field so that it is mashaled as is instead of being obscured which is how
// alertmanager handles secrets. Otherwise the secrets would be obscured on
// write to the yml file, making it unusable.
type PagerDutyConfig struct {
	HTTPConfig *common.HTTPConfig `yaml:"http_config,omitempty" json:"http_config,omitempty"`

	RoutingKey  string                   `yaml:"routing_key,omitempty" json:"routing_key,omitempty"`
	ServiceKey  string                   `yaml:"service_key,omitempty" json:"service_key,omitempty"`
	URL         string                   `yaml:"url,omitempty" json:"url,omitempty"`
	Client      string                   `yaml:"client,omitempty" json:"client,omitempty"`
	ClientURL   string                   `yaml:"client_url,omitempty" json:"client_url,omitempty"`
	Description string                   `yaml:"description,omitempty" json:"description,omitempty"`
	Severity    string                   `yaml:"severity,omitempty" json:"severity,omitempty"`
	Details     map[string]string        `yaml:"details,omitempty" json:"details,omitempty"`
	Images      []*config.PagerdutyImage `yaml:"images,omitempty" json:"images,omitempty"`
	Links       []*config.PagerdutyLink  `yaml:"links,omitempty" json:"links,omitempty"`
}

// PushoverConfig uses string instead of Secret for the UserKey and Token
// field so that it is mashaled as is instead of being obscured which is how
// alertmanager handles secrets. Otherwise the secrets would be obscured on
// write to the yml file, making it unusable.
type PushoverConfig struct {
	HTTPConfig *common.HTTPConfig `yaml:"http_config,omitempty" json:"http_config,omitempty"`

	UserKey  string         `yaml:"user_key" json:"user_key"`
	Token    string         `yaml:"token" json:"token"`
	Title    string         `yaml:"title,omitempty" json:"title,omitempty"`
	Message  string         `yaml:"message,omitempty" json:"message,omitempty"`
	URL      string         `yaml:"url,omitempty" json:"url,omitempty"`
	Priority string         `yaml:"priority,omitempty" json:"priority,omitempty"`
	Retry    model.Duration `yaml:"retry,omitempty" json:"retry,omitempty"`
	Expire   model.Duration `yaml:"expire,omitempty" json:"expire,omitempty"`
}

// PushoverJSONWrapper uses strings instead of duration objects to allow easier
// input and interaction through the API.
type PushoverJSONWrapper struct {
	HTTPConfig *common.HTTPConfig `yaml:"http_config,omitempty" json:"http_config,omitempty"`

	UserKey  string `yaml:"user_key" json:"user_key"`
	Token    string `yaml:"token" json:"token"`
	Title    string `yaml:"title,omitempty" json:"title,omitempty"`
	Message  string `yaml:"message,omitempty" json:"message,omitempty"`
	URL      string `yaml:"url,omitempty" json:"url,omitempty"`
	Priority string `yaml:"priority,omitempty" json:"priority,omitempty"`
	Retry    string `yaml:"retry,omitempty" json:"retry,omitempty"`
	Expire   string `yaml:"expire,omitempty" json:"expire,omitempty"`
}

// ToReceiverFmt convers the JSONWrapper object to a true Receiver object. This will
// only be necessary when dealing with Pushover objects for the time being (due to
// complexities surrounding JSON unmarshalling)
func (r *ReceiverJSONWrapper) ToReceiverFmt() (Receiver, error) {
	receiver := Receiver{
		Name:             r.Name,
		SlackConfigs:     r.SlackConfigs,
		WebhookConfigs:   r.WebhookConfigs,
		EmailConfigs:     r.EmailConfigs,
		PagerDutyConfigs: r.PagerDutyConfigs,
	}

	for _, p := range r.PushoverConfigs {
		pushoverConf := PushoverConfig{
			HTTPConfig: p.HTTPConfig,
			UserKey:    p.UserKey,
			Token:      p.Token,
			Title:      p.Title,
			Message:    p.Message,
			URL:        p.URL,
			Priority:   p.Priority,
		}
		if p.Retry != "" {
			modelRetry, err := model.ParseDuration(p.Retry)
			if err != nil {
				return receiver, err
			}
			pushoverConf.Retry = modelRetry
		}
		if p.Expire != "" {
			modelExpire, err := model.ParseDuration(p.Expire)
			if err != nil {
				return receiver, err
			}
			pushoverConf.Expire = modelExpire
		}
		receiver.PushoverConfigs = append(receiver.PushoverConfigs, &pushoverConf)
	}

	return receiver, nil
}

// MarshalYAML implements the yaml.Marshaler interface for EmailConfig and
// forces RequireTLS to be false. RequireTLS must be false since we don't support
// storing certificate files.
func (e EmailConfig) MarshalYAML() (interface{}, error) {
	valFalse := false
	e.RequireTLS = &valFalse
	return e, nil
}

// WebhookConfig is a copy of prometheus/alertmanager/config.WebhookConfig with
// alertmanager-configurer's custom HTTPConfig
type WebhookConfig struct {
	config.NotifierConfig

	HTTPConfig *common.HTTPConfig `yaml:"http_config,omitempty" json:"http_config,omitempty"`

	URL *config.URL `yaml:"url" json:"url"`
}

func ReceiverTenantPrefix(tenantID string) string {
	return strings.Replace(tenantID, "_", "", -1) + "_"
}
