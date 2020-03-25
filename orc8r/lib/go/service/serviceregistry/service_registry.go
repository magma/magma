/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package serviceregistry

import (
	"fmt"
	"log"
	"strings"

	"magma/orc8r/lib/go/registry"
	"magma/orc8r/lib/go/service/config"
)

const (
	serviceRegistryFilename = "service_registry"
)

type rawMapType = map[interface{}]interface{}

// LoadServiceRegistryConfig reads service registry config file from /etc/magma/configs/{moduleName} or override config
func LoadServiceRegistryConfig(moduleName string) ([]registry.ServiceLocation, error) {
	config, err := config.GetServiceConfig(moduleName, serviceRegistryFilename)
	if err != nil {
		// file does not exist
		return nil, err
	}
	rawMap, err := getRawMap(config)
	if err != nil {
		// file is empty
		return nil, err
	}
	locations, err := convertToServiceLocations(rawMap)
	if err != nil {
		log.Printf("Failed to load in service registry for %s:%s.yml: %v", moduleName, serviceRegistryFilename, err)
	}
	return locations, err
}

func getProxyAliases(rawMap map[interface{}]interface{}) map[string]int {
	proxyAliases := map[string]int{}
	if val, ok := rawMap["proxy_aliases"]; ok {
		rawMap, _ := val.(rawMapType)
		for k, v := range rawMap {
			proxyName, _ := k.(string)
			portMap, _ := v.(rawMapType)
			port, _ := portMap["port"].(int)
			proxyAliases[proxyName] = port
		}
	}
	return proxyAliases
}

func getRawMap(serviceRegistry *config.ConfigMap) (map[interface{}]interface{}, error) {
	services, ok := serviceRegistry.RawMap["services"]
	if !ok {
		return nil, fmt.Errorf("The field:services does not exist")
	}
	rawMap, ok := services.(rawMapType)
	if !ok {
		return nil, fmt.Errorf("Unable to convert serviceRegistry to map")
	}
	return rawMap, nil
}

func convertToServiceLocations(rawMap rawMapType) ([]registry.ServiceLocation, error) {
	serviceLocations := make([]registry.ServiceLocation, 0, len(rawMap))
	for k, v := range rawMap {
		name, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("The name of the service is not a string: %v", k)
		}
		rawMap, ok := v.(rawMapType)
		if !ok {
			return nil, fmt.Errorf("The value associated with key:%v is not a map: %v", k, v)
		}
		configMap := &config.ConfigMap{RawMap: rawMap}
		host, err := configMap.GetStringParam("host")
		if err != nil {
			// Check old/py format: 'ip_address'
			var ipErr error
			if host, ipErr = configMap.GetStringParam("ip_address"); ipErr != nil {
				return nil, err
			}
		}
		port, err := configMap.GetIntParam("port")
		if err != nil {
			return nil, err
		}
		proxyAliases := getProxyAliases(rawMap)
		serviceLocations = append(serviceLocations, registry.ServiceLocation{Name: strings.ToUpper(name), Host: host, Port: port, ProxyAliases: proxyAliases})
	}
	return serviceLocations, nil
}
