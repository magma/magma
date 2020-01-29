/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package receivers

import (
	"fmt"
	"strings"
	"time"

	"github.com/prometheus/alertmanager/config"
	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

// Config uses a custom receiver struct to avoid scrubbing of 'secrets' during
// marshaling
type Config struct {
	Global       *config.GlobalConfig  `yaml:"global,omitempty" json:"global,omitempty"`
	Route        *config.Route         `yaml:"route,omitempty" json:"route,omitempty"`
	InhibitRules []*config.InhibitRule `yaml:"inhibit_rules,omitempty" json:"inhibit_rules,omitempty"`
	Receivers    []*Receiver           `yaml:"receivers,omitempty" json:"receivers,omitempty"`
	Templates    []string              `yaml:"templates" json:"templates"`
}

// GetReceiver returns the receiver config with the given name
func (c *Config) GetReceiver(name string) *Receiver {
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

func (c *Config) initializeNetworkBaseRoute(route *config.Route, matcherLabel, tenantID string) error {
	baseRouteName := makeBaseRouteName(tenantID)
	if c.GetReceiver(baseRouteName) != nil {
		return fmt.Errorf("Base route for tenant %s already exists", tenantID)
	}

	c.Receivers = append(c.Receivers, &Receiver{Name: baseRouteName})
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

	err = yaml.Unmarshal(yamlData, &config.Config{})
	if err != nil {
		return err
	}
	return nil
}

// Receiver uses custom notifier configs to allow for marshaling of secrets.
type Receiver struct {
	Name string `yaml:"name" json:"name"`

	SlackConfigs   []*SlackConfig          `yaml:"slack_configs,omitempty" json:"slack_configs,omitempty"`
	WebhookConfigs []*config.WebhookConfig `yaml:"webhook_configs,omitempty" json:"webhook_configs,omitempty"`
	EmailConfigs   []*EmailConfig          `yaml:"email_configs,omitempty" json:"email_configs,omitempty"`
}

// Secure replaces the receiver's name with a tenantID prefix
func (r *Receiver) Secure(tenantID string) {
	r.Name = secureReceiverName(r.Name, tenantID)
}

// Unsecure removes the tenantID prefix from the receiver name
func (r *Receiver) Unsecure(tenantID string) {
	r.Name = unsecureReceiverName(r.Name, tenantID)
}

func secureReceiverName(name, tenantID string) string {
	return receiverTenantPrefix(tenantID) + name
}

func unsecureReceiverName(name, tenantID string) string {
	if strings.HasPrefix(name, receiverTenantPrefix(tenantID)) {
		return name[len(receiverTenantPrefix(tenantID)):]
	}
	return name
}

// SlackConfig uses string instead of SecretURL for the APIURL field so that it
// is marshaled as is instead of being obscured which is how alertmanager handles
// secrets
type SlackConfig struct {
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

// MarshalYAML implements the yaml.Marshaler interface for EmailConfig and
// forces RequireTLS to be false. RequireTLS must be false since we don't support
// storing certificate files.
func (e EmailConfig) MarshalYAML() (interface{}, error) {
	valFalse := false
	e.RequireTLS = &valFalse
	return e, nil
}

// RouteJSONWrapper Provides a struct to marshal/unmarshal into a rulefmt.Rule
// since rulefmt does not support json encoding
type RouteJSONWrapper struct {
	Receiver string `yaml:"receiver,omitempty" json:"receiver,omitempty"`

	GroupByStr []string          `yaml:"group_by,omitempty" json:"group_by,omitempty"`
	GroupBy    []model.LabelName `yaml:"-" json:"-"`
	GroupByAll bool              `yaml:"-" json:"-"`

	Match    map[string]string        `yaml:"match,omitempty" json:"match,omitempty"`
	MatchRE  map[string]config.Regexp `yaml:"match_re,omitempty" json:"match_re,omitempty"`
	Continue bool                     `yaml:"continue,omitempty" json:"continue,omitempty"`
	Routes   []*RouteJSONWrapper      `yaml:"routes,omitempty" json:"routes,omitempty"`

	GroupWait      string `yaml:"group_wait,omitempty" json:"group_wait,omitempty"`
	GroupInterval  string `yaml:"group_interval,omitempty" json:"group_interval,omitempty"`
	RepeatInterval string `yaml:"repeat_interval,omitempty" json:"repeat_interval,omitempty"`
}

// NewRouteJSONWrapper converts a config.Route to a json-compatible route
func NewRouteJSONWrapper(r config.Route) *RouteJSONWrapper {
	var childRoutes []*RouteJSONWrapper
	for _, child := range r.Routes {
		if child != nil {
			childRoutes = append(childRoutes, NewRouteJSONWrapper(*child))
		}
	}
	var groupWaitStr, groupIntervalStr, repeatIntervalStr string
	if r.GroupWait != nil {
		groupWaitStr = r.GroupWait.String()
	}
	if r.GroupInterval != nil {
		groupIntervalStr = r.GroupInterval.String()
	}
	if r.RepeatInterval != nil {
		repeatIntervalStr = r.RepeatInterval.String()
	}

	return &RouteJSONWrapper{
		Receiver:       r.Receiver,
		GroupByStr:     r.GroupByStr,
		GroupBy:        r.GroupBy,
		GroupByAll:     r.GroupByAll,
		Match:          r.Match,
		MatchRE:        r.MatchRE,
		Continue:       r.Continue,
		Routes:         childRoutes,
		GroupWait:      groupWaitStr,
		GroupInterval:  groupIntervalStr,
		RepeatInterval: repeatIntervalStr,
	}
}

// ToPrometheusConfig converts a json-compatible route specification to a
// prometheus route config
func (r *RouteJSONWrapper) ToPrometheusConfig() (config.Route, error) {
	var groupWait, groupInterval, repeatInterval model.Duration
	var groupWaitP, groupIntervalP, repeatIntervalP *model.Duration
	var err error

	if r.GroupWait != "" {
		groupWait, err = model.ParseDuration(r.GroupWait)
		if err != nil {
			return config.Route{}, fmt.Errorf("Invalid GroupWait '%s': %v", r.GroupWait, err)
		}
	}
	if r.GroupInterval != "" {
		groupInterval, err = model.ParseDuration(r.GroupInterval)
		if err != nil {
			return config.Route{}, fmt.Errorf("Invalid GroupInterval '%s': %v", r.GroupInterval, err)
		}
		if time.Duration(groupInterval) == time.Duration(0) {
			return config.Route{}, fmt.Errorf("GroupInterval cannot be 0")
		}
	}
	if r.RepeatInterval != "" {
		repeatInterval, err = model.ParseDuration(r.RepeatInterval)
		if err != nil {
			return config.Route{}, fmt.Errorf("Invalid RepeatInterval '%s': %v", r.RepeatInterval, err)
		}
		if time.Duration(repeatInterval) == time.Duration(0) {
			return config.Route{}, fmt.Errorf("RepeatInterval cannot be 0")
		}
	}
	groupWaitP = &groupWait
	if time.Duration(groupInterval) != time.Duration(0) {
		groupIntervalP = &groupInterval
	}
	if time.Duration(repeatInterval) != time.Duration(0) {
		repeatIntervalP = &repeatInterval
	}

	var configRoutes []*config.Route
	for _, childRoute := range r.Routes {
		route, err := childRoute.ToPrometheusConfig()
		if err != nil {
			return config.Route{}, fmt.Errorf("error converting child route: %v", err)
		}
		configRoutes = append(configRoutes, &route)
	}

	route := config.Route{
		Receiver:       r.Receiver,
		GroupByStr:     r.GroupByStr,
		GroupBy:        r.GroupBy,
		GroupByAll:     r.GroupByAll,
		Match:          r.Match,
		MatchRE:        r.MatchRE,
		Continue:       r.Continue,
		Routes:         configRoutes,
		GroupWait:      groupWaitP,
		GroupInterval:  groupIntervalP,
		RepeatInterval: repeatIntervalP,
	}
	return route, nil
}
