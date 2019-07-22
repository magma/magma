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
	"io/ioutil"
	"strings"
	"sync"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/exporters"

	"github.com/prometheus/alertmanager/config"

	"gopkg.in/yaml.v2"
)

const (
	networkBaseRoutePostfix = "network_base_route"
)

// Client provides methods to create and read receiver configurations
type Client struct {
	configPath string
	sync.RWMutex
}

func NewClient(configPath string) *Client {
	return &Client{
		configPath: configPath,
	}
}

// CreateReceiver writes a new receiver to the config file with the networkID
// prepended to the name so multiple networks can be supported
func (c *Client) CreateReceiver(rec *Receiver, networkID string) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}

	rec.Secure(networkID)
	conf.Receivers = append(conf.Receivers, rec)
	err = conf.Validate()
	if err != nil {
		return err
	}
	return c.writeConfigFile(conf)
}

// GetReceivers returns the receiver configs for the given networkID
func (c *Client) GetReceivers(networkID string) ([]Receiver, error) {
	c.RLock()
	defer c.RUnlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return []Receiver{}, nil
	}

	recs := make([]Receiver, 0)
	for _, rec := range conf.Receivers {
		if strings.HasPrefix(rec.Name, receiverNetworkPrefix(networkID)) {
			rec.Unsecure(networkID)
			recs = append(recs, *rec)
		}
	}
	return recs, nil
}

// ModifyNetworkRoute takes a new route for a network and replaces the old one,
// ensuring that receivers are properly named and the resulting config is valid.
// Creates a new one if it doesn't already exist
func (c *Client) ModifyNetworkRoute(route *config.Route, networkID string) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}
	// ensure base route is valid base route for this network
	route.Receiver = makeBaseRouteName(networkID)
	if route.Match == nil {
		route.Match = map[string]string{}
	}
	route.Match[exporters.NetworkLabelNetwork] = networkID

	for _, childRoute := range route.Routes {
		if childRoute == nil {
			continue
		}
		secureRoute(childRoute, networkID)
	}

	networkRouteIdx := conf.GetRouteIdx(makeBaseRouteName(networkID))
	if networkRouteIdx < 0 {
		err := conf.initializeNetworkBaseRoute(route, networkID)
		if err != nil {
			return err
		}
	} else {
		conf.Route.Routes[networkRouteIdx] = route
	}

	err = conf.Validate()
	if err != nil {
		return err
	}
	return c.writeConfigFile(conf)
}

// GetRoute returns the base route for the given networkID
func (c *Client) GetRoute(networkID string) (*config.Route, error) {
	c.RLock()
	defer c.RUnlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return &config.Route{}, err
	}

	routeIdx := conf.GetRouteIdx(makeBaseRouteName(networkID))
	if routeIdx >= 0 {
		route := conf.Route.Routes[routeIdx]
		unsecureRoute(route, networkID)
		return route, nil
	}
	return nil, fmt.Errorf("Route for network %s does not exist", networkID)
}

func (c *Client) readConfigFile() (*Config, error) {
	configFile := Config{}
	file, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config files: %v", err)
	}
	err = yaml.Unmarshal(file, &configFile)
	return &configFile, err
}

func (c *Client) writeConfigFile(conf *Config) error {
	yamlFile, err := yaml.Marshal(conf)
	if err != nil {
		return fmt.Errorf("error marshaling config file: %v", err)
	}
	err = ioutil.WriteFile(c.configPath, yamlFile, 0660)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}
	return nil
}

// secureRoute ensure that all receivers in the route have the
// proper networkID-prefixed receiver name
func secureRoute(route *config.Route, networkID string) {
	route.Receiver = secureReceiverName(route.Receiver, networkID)
	for _, childRoute := range route.Routes {
		secureRoute(childRoute, networkID)
	}
}

// unsecureRoute traverses a routing tree and reverts receiver
// names to their non-prefixed original names
func unsecureRoute(route *config.Route, networkID string) {
	route.Receiver = unsecureReceiverName(route.Receiver, networkID)
	for _, childRoute := range route.Routes {
		unsecureRoute(childRoute, networkID)
	}
}

func receiverNetworkPrefix(networkID string) string {
	return strings.Replace(networkID, "_", "", -1) + "_"
}

func (c *Client) getBaseRouteForNetwork(networkID string, conf *Config) (*config.Route, error) {
	baseRouteName := makeBaseRouteName(networkID)
	for _, route := range conf.Route.Routes {
		if route.Receiver == baseRouteName {
			return route, nil
		}
	}
	return nil, fmt.Errorf("base route for %s not found", networkID)
}

func makeBaseRouteName(networkID string) string {
	return fmt.Sprintf("%s_%s", networkID, networkBaseRoutePostfix)
}
