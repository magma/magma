/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package receivers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"magma/orc8r/cloud/go/metrics"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/fsclient"

	"github.com/prometheus/alertmanager/config"

	"gopkg.in/yaml.v2"
)

const (
	networkBaseRoutePostfix = "network_base_route"
)

type AlertmanagerClient interface {
	CreateReceiver(networkID string, rec Receiver) error
	GetReceivers(networkID string) ([]Receiver, error)
	UpdateReceiver(networkID string, newRec *Receiver) error
	DeleteReceiver(networkID, receiverName string) error

	// ModifyNetworkRoute updates an existing routing tree for the given
	// network, or creates one if it already exists. Ensures that the base
	// route matches all alerts with label "networkID" = <networkID>.
	ModifyNetworkRoute(networkID string, route *config.Route) error

	// GetRoute returns the routing tree for the given networkID
	GetRoute(networkID string) (*config.Route, error)

	// ReloadAlertmanager triggers the alertmanager process to reload the
	// configuration file(s)
	ReloadAlertmanager() error
}

// Client provides methods to create and read receiver configurations
type client struct {
	configPath      string
	alertmanagerURL string
	fsClient        fsclient.FSClient
	sync.RWMutex
}

func NewClient(configPath, alertmanagerURL string, fsClient fsclient.FSClient) AlertmanagerClient {
	return &client{
		configPath:      configPath,
		alertmanagerURL: alertmanagerURL,
		fsClient:        fsClient,
	}
}

// CreateReceiver writes a new receiver to the config file with the networkID
// prepended to the name so multiple networks can be supported
func (c *client) CreateReceiver(networkID string, rec Receiver) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}

	rec.Secure(networkID)
	conf.Receivers = append(conf.Receivers, &rec)
	err = conf.Validate()
	if err != nil {
		return err
	}
	return c.writeConfigFile(conf)
}

// GetReceivers returns the receiver configs for the given networkID
func (c *client) GetReceivers(networkID string) ([]Receiver, error) {
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

// UpdateReceiver modifies an existing receiver
func (c *client) UpdateReceiver(networkID string, newRec *Receiver) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}

	newRec.Secure(networkID)
	receiverIdx := -1
	for idx, rec := range conf.Receivers {
		if rec.Name == newRec.Name {
			receiverIdx = idx
		}
	}
	if receiverIdx < 0 {
		return fmt.Errorf("Receiver '%s' not found", newRec.Name)
	}

	conf.Receivers[receiverIdx] = newRec
	err = conf.Validate()
	if err != nil {
		return fmt.Errorf("Error updating receiver: %v", err)
	}
	return c.writeConfigFile(conf)
}

// DeleteReceiver removes a receiver from the configuration
func (c *client) DeleteReceiver(networkID, receiverName string) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}

	receiverToDelete := secureReceiverName(receiverName, networkID)
	for idx, rec := range conf.Receivers {
		if rec.Name == receiverToDelete {
			conf.Receivers = append(conf.Receivers[:idx], conf.Receivers[idx+1:]...)
			return c.writeConfigFile(conf)
		}
	}

	return fmt.Errorf("Receiver '%s' does not exist", receiverName)
}

// ModifyNetworkRoute takes a new route for a network and replaces the old one,
// ensuring that receivers are properly named and the resulting config is valid.
// Creates a new one if it doesn't already exist
func (c *client) ModifyNetworkRoute(networkID string, route *config.Route) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}
	// ensure base route is valid base route for this network
	baseRoute := c.getBaseRouteForNetwork(networkID, conf)
	route.Receiver = baseRoute.Receiver
	if route.Match == nil {
		route.Match = map[string]string{}
	}

	route.Match[metrics.NetworkLabelName] = networkID

	for _, childRoute := range route.Routes {
		if childRoute == nil {
			continue
		}
		secureRoute(networkID, childRoute)
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
func (c *client) GetRoute(networkID string) (*config.Route, error) {
	c.RLock()
	defer c.RUnlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return &config.Route{}, err
	}

	routeIdx := conf.GetRouteIdx(makeBaseRouteName(networkID))
	if routeIdx >= 0 {
		route := conf.Route.Routes[routeIdx]
		unsecureRoute(networkID, route)
		return route, nil
	}
	return nil, fmt.Errorf("Route for network %s does not exist", networkID)
}

func (c *client) ReloadAlertmanager() error {
	resp, err := http.Post(fmt.Sprintf("http://%s%s", c.alertmanagerURL, "/-/reload"), "text/plain", &bytes.Buffer{})
	if err != nil {
		return fmt.Errorf("error reloading alertmanager: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("code: %d error reloading alertmanager: %s", resp.StatusCode, msg)
	}
	return nil
}

func (c *client) readConfigFile() (*Config, error) {
	configFile := Config{}
	file, err := c.fsClient.ReadFile(c.configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config files: %v", err)
	}
	err = yaml.Unmarshal(file, &configFile)
	return &configFile, err
}

func (c *client) writeConfigFile(conf *Config) error {
	yamlFile, err := yaml.Marshal(conf)
	if err != nil {
		return fmt.Errorf("error marshaling config file: %v", err)
	}
	err = c.fsClient.WriteFile(c.configPath, yamlFile, 0660)
	if err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}
	return nil
}

// secureRoute ensure that all receivers in the route have the
// proper networkID-prefixed receiver name
func secureRoute(networkID string, route *config.Route) {
	route.Receiver = secureReceiverName(route.Receiver, networkID)
	for _, childRoute := range route.Routes {
		secureRoute(networkID, childRoute)
	}
}

// unsecureRoute traverses a routing tree and reverts receiver
// names to their non-prefixed original names
func unsecureRoute(networkID string, route *config.Route) {
	route.Receiver = unsecureReceiverName(route.Receiver, networkID)
	for _, childRoute := range route.Routes {
		unsecureRoute(networkID, childRoute)
	}
}

func receiverNetworkPrefix(networkID string) string {
	return strings.Replace(networkID, "_", "", -1) + "_"
}

func (c *client) getBaseRouteForNetwork(networkID string, conf *Config) *config.Route {
	baseRouteName := makeBaseRouteName(networkID)
	for _, route := range conf.Route.Routes {
		if route.Receiver == baseRouteName {
			return route
		}
	}
	newBaseRoute := &config.Route{Receiver: makeBaseRouteName(networkID)}
	return newBaseRoute
}

func makeBaseRouteName(networkID string) string {
	return fmt.Sprintf("%s_%s", networkID, networkBaseRoutePostfix)
}
