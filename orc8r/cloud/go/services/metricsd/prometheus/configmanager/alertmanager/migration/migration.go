/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/alertmanager/config"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/configmanager/fsclient"

	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

const (
	defaultAlertmanagerConfigPath = "./alertmanager.yml"
)

func main() {
	alertmanagerConfPath := flag.String("alertmanager-conf", defaultAlertmanagerConfigPath, fmt.Sprintf("Path to alertmanager configuration file. Default is %s", defaultAlertmanagerConfigPath))
	flag.Parse()

	fsClient := fsclient.NewFSClient()

	// Read config file
	configFile := config.Config{}
	file, err := fsClient.ReadFile(*alertmanagerConfPath)
	if err != nil {
		glog.Fatalf("error reading config files: %v", err)
	}
	err = yaml.Unmarshal(file, &configFile)
	if err != nil {
		glog.Fatalf("error marshaling config file: %v", err)
	}

	// Do tenancy migration
	migrateToTenantBasedConfig(&configFile)

	// Write config file
	yamlFile, err := yaml.Marshal(configFile)
	if err != nil {
		glog.Fatalf("error marshaling config file: %v", err)
	}
	err = fsClient.WriteFile(*alertmanagerConfPath, yamlFile, 0660)
	if err != nil {
		glog.Fatalf("error writing config file: %v", err)
	}

	glog.Infof("Migrations completed successfully")
}

// This is necessary due to the change from 'network' based tenancy to 'tenant'
// based tenancy. Replaces 'network_base_route' with 'tenant_base_route'
const deprecatedTenancyPostfix = "network_base_route"

func migrateToTenantBasedConfig(conf *config.Config) {
	for _, route := range conf.Route.Routes {
		matched, _ := regexp.MatchString(fmt.Sprintf(".*_%s", deprecatedTenancyPostfix), route.Receiver)
		if matched {
			migratedName := strings.Replace(route.Receiver, deprecatedTenancyPostfix, config.TenantBaseRoutePostfix, 1)
			route.Receiver = migratedName
		}
	}
	for _, receiver := range conf.Receivers {
		matched, _ := regexp.MatchString(fmt.Sprintf(".*_%s", deprecatedTenancyPostfix), receiver.Name)
		if matched {
			migratedName := strings.Replace(receiver.Name, deprecatedTenancyPostfix, config.TenantBaseRoutePostfix, 1)
			receiver.Name = migratedName
		}
	}
}
