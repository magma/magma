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

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/fsclient"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/prometheus/alert"

	"github.com/prometheus/alertmanager/config"

	"gopkg.in/yaml.v2"
)

const (
	TenantBaseRoutePostfix = "tenant_base_route"
)

type AlertmanagerClient interface {
	CreateReceiver(tenantID string, rec Receiver) error
	GetReceivers(tenantID string) ([]Receiver, error)
	UpdateReceiver(tenantID string, newRec *Receiver) error
	DeleteReceiver(tenantID, receiverName string) error

	// ModifyNetworkRoute updates an existing routing tree for the given
	// tenant, or creates one if it already exists. Ensures that the base
	// route matches all alerts with label "tenantID" = <tenantID>.
	ModifyTenantRoute(tenantID string, route *config.Route) error

	// GetRoute returns the routing tree for the given tenantID
	GetRoute(tenantID string) (*config.Route, error)

	// ReloadAlertmanager triggers the alertmanager process to reload the
	// configuration file(s)
	ReloadAlertmanager() error

	Tenancy() *alert.TenancyConfig
}

// Client provides methods to create and read receiver configurations
type client struct {
	configPath      string
	alertmanagerURL string
	fsClient        fsclient.FSClient
	tenancy         *alert.TenancyConfig
	sync.RWMutex
}

func NewClient(configPath, alertmanagerURL string, tenancy *alert.TenancyConfig, fsClient fsclient.FSClient) AlertmanagerClient {
	return &client{
		configPath:      configPath,
		alertmanagerURL: alertmanagerURL,
		fsClient:        fsClient,
		tenancy:         tenancy,
	}
}

// CreateReceiver writes a new receiver to the config file with the tenantID
// prepended to the name so multiple tenants can be supported
func (c *client) CreateReceiver(tenantID string, rec Receiver) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}

	rec.Secure(tenantID)

	conf.Receivers = append(conf.Receivers, &rec)
	err = conf.Validate()
	if err != nil {
		return err
	}
	return c.writeConfigFile(conf)
}

// GetReceivers returns the receiver configs for the given tenantID
func (c *client) GetReceivers(tenantID string) ([]Receiver, error) {
	c.RLock()
	defer c.RUnlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return []Receiver{}, nil
	}

	recs := make([]Receiver, 0)
	for _, rec := range conf.Receivers {
		if strings.HasPrefix(rec.Name, receiverTenantPrefix(tenantID)) {
			rec.Unsecure(tenantID)
			recs = append(recs, *rec)
		}
	}
	return recs, nil
}

// UpdateReceiver modifies an existing receiver
func (c *client) UpdateReceiver(tenantID string, newRec *Receiver) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}

	newRec.Secure(tenantID)

	receiverIdx := -1
	for idx, rec := range conf.Receivers {
		if rec.Name == newRec.Name {
			receiverIdx = idx
			break
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
func (c *client) DeleteReceiver(tenantID, receiverName string) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}

	receiverToDelete := secureReceiverName(receiverName, tenantID)

	for idx, rec := range conf.Receivers {
		if rec.Name == receiverToDelete {
			conf.Receivers = append(conf.Receivers[:idx], conf.Receivers[idx+1:]...)
			return c.writeConfigFile(conf)
		}
	}

	return fmt.Errorf("Receiver '%s' does not exist", receiverName)
}

// ModifyTenantRoute takes a new route for a tenant and replaces the old one,
// ensuring that receivers are properly named and the resulting config is valid.
// Creates a new one if it doesn't already exist. If single-tenant client this
// just modifies the entire routing tree
func (c *client) ModifyTenantRoute(tenantID string, route *config.Route) error {
	c.Lock()
	defer c.Unlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return err
	}

	// ensure base route is valid base route for this tenant
	baseRoute := c.getBaseRouteForTenant(tenantID, conf)
	route.Receiver = baseRoute.Receiver
	if route.Match == nil {
		route.Match = map[string]string{}
	}

	route.Match[c.tenancy.RestrictorLabel] = tenantID

	for _, childRoute := range route.Routes {
		if childRoute == nil {
			continue
		}
		secureRoute(tenantID, childRoute)
	}

	tenantRouteIdx := conf.GetRouteIdx(makeBaseRouteName(tenantID))
	if tenantRouteIdx < 0 {
		err := conf.initializeNetworkBaseRoute(route, c.tenancy.RestrictorLabel, tenantID)
		if err != nil {
			return err
		}
	} else {
		conf.Route.Routes[tenantRouteIdx] = route
	}

	err = conf.Validate()
	if err != nil {
		return err
	}
	return c.writeConfigFile(conf)
}

// GetRoute returns the base route for the given tenantID
func (c *client) GetRoute(tenantID string) (*config.Route, error) {
	c.RLock()
	defer c.RUnlock()
	conf, err := c.readConfigFile()
	if err != nil {
		return &config.Route{}, err
	}

	routeIdx := conf.GetRouteIdx(makeBaseRouteName(tenantID))
	if routeIdx >= 0 {
		route := conf.Route.Routes[routeIdx]
		unsecureRoute(tenantID, route)
		return route, nil
	}
	return nil, fmt.Errorf("Route for tenant %s does not exist", tenantID)
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

func (c *client) Tenancy() *alert.TenancyConfig {
	return c.tenancy
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
// proper tenantID-prefixed receiver name
func secureRoute(tenantID string, route *config.Route) {
	route.Receiver = secureReceiverName(route.Receiver, tenantID)
	for _, childRoute := range route.Routes {
		secureRoute(tenantID, childRoute)
	}
}

// unsecureRoute traverses a routing tree and reverts receiver
// names to their non-prefixed original names
func unsecureRoute(tenantID string, route *config.Route) {
	route.Receiver = unsecureReceiverName(route.Receiver, tenantID)
	for _, childRoute := range route.Routes {
		unsecureRoute(tenantID, childRoute)
	}
}

func receiverTenantPrefix(tenantID string) string {
	return strings.Replace(tenantID, "_", "", -1) + "_"
}

func (c *client) getBaseRouteForTenant(tenantID string, conf *Config) *config.Route {
	baseRouteName := makeBaseRouteName(tenantID)
	for _, route := range conf.Route.Routes {
		if route.Receiver == baseRouteName {
			return route
		}
	}
	newBaseRoute := &config.Route{Receiver: makeBaseRouteName(tenantID)}
	return newBaseRoute
}

func makeBaseRouteName(tenantID string) string {
	return fmt.Sprintf("%s_%s", tenantID, TenantBaseRoutePostfix)
}
