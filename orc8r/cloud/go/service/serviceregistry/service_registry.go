/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package serviceregistry

import (
	"fmt"
	"strings"

	"magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/service/config"
)

const (
	serviceRegistryFilename = "service_registry"
)

// LoadServiceRegistryConfig reads service registry config file from /etc/magma/configs/{moduleName} or override config
func LoadServiceRegistryConfig(moduleName string) ([]registry.ServiceLocation, error) {
	config, err := config.GetServiceConfig(moduleName, serviceRegistryFilename)
	if err != nil {
		return nil, err
	}

	return convertToServiceLocations(config)
}

func convertToServiceLocations(serviceRegistry *config.ConfigMap) ([]registry.ServiceLocation, error) {
	services := serviceRegistry.RawMap["services"]
	rawMap, ok := services.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("Unable to convert map %v", rawMap)
	}

	serviceLocations := make([]registry.ServiceLocation, len(serviceRegistry.RawMap))
	for k, v := range rawMap {
		name, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("Unable convert key:%v to string", k)
		}

		rawMap, ok := v.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("Unable to convert map %v", rawMap)
		}
		configMap := &config.ConfigMap{RawMap: rawMap}

		host, err := configMap.GetStringParam("host")
		if err != nil {
			return nil, err
		}

		port, err := configMap.GetIntParam("port")
		if err != nil {
			return nil, err
		}

		serviceLocations = append(serviceLocations, registry.ServiceLocation{Name: strings.ToUpper(name), Host: host, Port: port})
	}
	return serviceLocations, nil
}
