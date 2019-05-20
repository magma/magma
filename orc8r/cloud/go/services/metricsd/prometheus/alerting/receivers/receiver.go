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

	"magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"

	"github.com/prometheus/alertmanager/config"
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

func (c *Config) initializeNetworkBaseRoute(route *config.Route, networkID string) error {
	baseRouteName := makeBaseRouteName(networkID)
	if c.GetReceiver(baseRouteName) != nil {
		return fmt.Errorf("Base route for network %s already exists", networkID)
	}

	c.Receivers = append(c.Receivers, &Receiver{Name: baseRouteName})
	route.Receiver = baseRouteName
	route.Match = map[string]string{exporters.NetworkLabelNetwork: networkID}

	c.Route.Routes = append(c.Route.Routes, route)

	return c.Validate()
}

// Validate makes sure that the config is properly formed. Have to do this here
// since alertmanager only does validation during unmarshaling
func (c *Config) Validate() error {
	receiverNames := map[string]struct{}{}

	for _, rcv := range c.Receivers {
		if _, ok := receiverNames[rcv.Name]; ok {
			return fmt.Errorf("notification config name %s is not unique", rcv.Name)
		}
		for _, sc := range rcv.SlackConfigs {
			err := validateURL(sc.APIURL)
			if err != nil {
				return err
			}
		}
		receiverNames[rcv.Name] = struct{}{}
	}
	if c.Route == nil {
		return fmt.Errorf("no route provided")
	}
	if len(c.Route.Receiver) == 0 {
		return fmt.Errorf("root route must specify a default receiver")
	}
	if len(c.Route.Match) > 0 || len(c.Route.MatchRE) > 0 {
		return fmt.Errorf("root route must not have any matchers")
	}

	// check that all receivers used in routing tree are defined
	return checkReceiver(c.Route, receiverNames)
}

func validateURL(url string) error {
	if !strings.HasPrefix(url, "http") {
		return fmt.Errorf("invalid url: %s", url)
	}
	return nil
}

// checkReceiver returns an error if a node in the routing tree
// references a receiver not in the given map.
func checkReceiver(r *config.Route, receivers map[string]struct{}) error {
	for _, sr := range r.Routes {
		if err := checkReceiver(sr, receivers); err != nil {
			return err
		}
	}
	if r.Receiver == "" {
		return nil
	}
	if _, ok := receivers[r.Receiver]; !ok {
		return fmt.Errorf("undefined receiver %q used in route", r.Receiver)
	}
	return nil
}

// Receiver uses custom notifier configs to allow for marshaling of secrets.
type Receiver struct {
	Name string `yaml:"name" json:"name"`

	SlackConfigs []*SlackConfig `yaml:"slack_configs,omitempty" json:"slack_configs,omitempty"`
}

// Secure replaces the receiver's name with a networkID prefix
func (r *Receiver) Secure(networkID string) {
	r.Name = secureReceiverName(r.Name, networkID)
}

// Unsecure removes the networkID prefix from the receiver name
func (r *Receiver) Unsecure(networkID string) {
	r.Name = unsecureReceiverName(r.Name, networkID)
}

func secureReceiverName(name, networkID string) string {
	return receiverNetworkPrefix(networkID) + name
}

func unsecureReceiverName(name, networkID string) string {
	if strings.HasPrefix(name, receiverNetworkPrefix(networkID)) {
		return name[len(receiverNetworkPrefix(networkID)):]
	}
	return name
}

// SlackConfig uses string instead of SecretURL for the APIURL field so that it
// is marshaled as is instead of being obscured which is how alertmanager handles
// secrets
type SlackConfig struct {
	APIURL   string `yaml:"api_url" json:"api_url"`
	Channel  string `yaml:"channel" json:"channel"`
	Username string `yaml:"username" json:"username"`
}
