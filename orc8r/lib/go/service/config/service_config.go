/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"gopkg.in/yaml.v2"
)

var (
	// ConfigDir is where the per-service configuration files are stored
	configDir         = "/etc/magma/configs"
	oldConfigDir      = "/etc/magma"
	configOverrideDir = "/var/opt/magma/configs/"
	cfgDirMu          sync.RWMutex
)

// ConfigMap is a struct for representing a map generated from a service YML file
type ConfigMap struct {
	RawMap map[interface{}]interface{}
}

// NewConfigMap creates a new ConfigMap based on the input map
func NewConfigMap(config map[interface{}]interface{}) *ConfigMap {
	return &ConfigMap{config}
}

// GetServiceConfig loads a config by name to a map of parameters
// Input: configName - name of config to load, e.g. control_proxy
// Output: map of parameters if it exists, error if not
func GetServiceConfig(moduleName string, serviceName string) (*ConfigMap, error) {
	cfgDirMu.RLock()
	main, legacy, overwrite := configDir, oldConfigDir, configOverrideDir
	cfgDirMu.RUnlock()
	return getServiceConfigImpl(moduleName, serviceName, main, legacy, overwrite)
}

// GetStringParam is used to retrieve a string param from a YML file and returns error if param does not exist/ill-formed
func (cfgMap *ConfigMap) GetStringParam(key string) (string, error) {
	return getStringParamImpl(cfgMap, key)
}

// GetRequiredStringParam is same as GetStringParam but fails when the string does not exist
func (cfgMap *ConfigMap) GetRequiredStringParam(key string) string {
	str, err := getStringParamImpl(cfgMap, key)
	if err != nil {
		log.Fatalf("Error retrieving %s: %v\n", key, err)
	}
	return str
}

// UpdateStringParam updates the given string pointed by strPtr if the config param exists and non-empty
func (cfgMap *ConfigMap) UpdateStringParam(strPtr *string, key string) error {
	if strPtr == nil {
		return fmt.Errorf("Nil pointer for parameter with key: %s", key)
	}
	str, err := getStringParamImpl(cfgMap, key)
	if err == nil && len(str) > 0 {
		*strPtr = str
	}
	return err
}

// GetIntParam is used to retrieve an int param from a YML file
func (cfgMap *ConfigMap) GetIntParam(key string) (int, error) {
	return getIntParamImpl(cfgMap, key)
}

// GetRequiredIntParam is same as GetIntParam but fails when the int does not exist
func (cfgMap *ConfigMap) GetRequiredIntParam(key string) int {
	param, err := getIntParamImpl(cfgMap, key)
	if err != nil {
		log.Fatalf("Error retrieving %s: %v\n", key, err)
	}
	return param
}

// GetBoolParam is used to retrieve a bool param from a YML file
func (cfgMap *ConfigMap) GetBoolParam(key string) (bool, error) {
	paramIface, ok := cfgMap.RawMap[key]
	if !ok {
		return false, fmt.Errorf("Could not find key %s", key)
	}
	param, ok := paramIface.(bool)
	if !ok {
		return false, fmt.Errorf("Could not convert param to bool for key %s", key)
	}
	return param, nil
}

// GetStringArrayParam finds and returns a string list/array parameter or error
func (cfgMap *ConfigMap) GetStringArrayParam(key string) ([]string, error) {
	return getStringArrayParamImpl(cfgMap, key)
}

// GetRequiredStringArrayParam finds and returns a string list/array parameter or "fatals"
func (cfgMap *ConfigMap) GetRequiredStringArrayParam(key string) []string {
	param, err := getStringArrayParamImpl(cfgMap, key)
	if err != nil {
		log.Fatalf("Error retrieving %s: %v\n", key, err)
	}
	return param
}

// GetMapParam finds and returns map parameter or error
func (cfgMap *ConfigMap) GetMapParam(key string) (map[interface{}]interface{}, error) {
	return cfgMap.getMapParamImpl(key)
}

// GetRequiredMapParam finds and returns a map parameter or "fatals"
func (cfgMap *ConfigMap) GetRequiredMapParam(key string) map[interface{}]interface{} {
	param, err := cfgMap.getMapParamImpl(key)
	if err != nil {
		log.Fatalf("Error retrieving %s: %v\n", key, err)
	}
	return param
}

// GetStructuredServiceConfig updates 'out' structure with configs from the configs YML
// If successful, GetStructuredServiceConfig returns used YML file path & used overwrite cfg YML file path
func GetStructuredServiceConfig(moduleName string, serviceName string, out interface{}) (string, string, error) {
	cfgDirMu.RLock()
	main, legacy, overwrite := configDir, oldConfigDir, configOverrideDir
	cfgDirMu.RUnlock()
	return GetStructuredServiceConfigExt(moduleName, serviceName, main, legacy, overwrite, out)
}

// GetStructuredServiceConfigExt is an extended version of GetStructuredServiceConfig, it allows to pass config
// directory names
func GetStructuredServiceConfigExt(
	moduleName,
	serviceName,
	configDir,
	oldConfigDir,
	configOverrideDir string,
	out interface{}) (ymlFilePath, ymlQWFilePath string, err error) {

	if out == nil {
		return ymlFilePath, ymlQWFilePath, fmt.Errorf("Structured CFG: Invalid (nil) output parameter")
	}

	moduleName, serviceName = strings.ToLower(moduleName), strings.ToLower(serviceName)
	ymlFilePath = getServiceConfigFilePath(moduleName, serviceName, configDir, oldConfigDir)
	yamlFileData, err := ioutil.ReadFile(ymlFilePath)
	if err == nil {
		err = yaml.Unmarshal([]byte(yamlFileData), out)
		if err != nil {
			log.Printf("Structured CFG: Error Unmarshaling '%s' into type %T: %v", ymlFilePath, out, err)
		} else {
			log.Printf("Successfully loaded structured '%s::%s' service configs from '%s'",
				moduleName, serviceName, ymlFilePath)
		}
	} else {
		log.Printf("Structured CFG: Error Reading '%s': %v", ymlFilePath, err)
	}
	// Overwrite params from override configs
	var oerr error
	ymlQWFilePath = filepath.Join(configOverrideDir, moduleName, fmt.Sprintf("%s.yml", serviceName))
	if fi, serr := os.Stat(ymlQWFilePath); serr == nil && !fi.IsDir() {
		yamlFileData, oerr = ioutil.ReadFile(ymlQWFilePath)
		if oerr == nil {
			oerr = yaml.Unmarshal([]byte(yamlFileData), out)
			if oerr != nil {
				log.Printf("Structured CFG: Error Unmarshaling Override file '%s' into type %T: %v",
					ymlQWFilePath, out, err)
			} else {
				log.Printf("Successfully loaded Override configs for service %s:%s from '%s'",
					moduleName, serviceName, ymlQWFilePath)
			}
		} else {
			log.Printf("Structured CFG:Error Loading Override configs from '%s': %v", ymlQWFilePath, err)
		}
	} else {
		ymlQWFilePath = ""
	}
	if err == nil || oerr == nil { // fully or partially succeeded
		return ymlFilePath, ymlQWFilePath, nil
	}
	if err != nil {
		return
	}
	return ymlFilePath, ymlQWFilePath, oerr
}

// GetCurrentConfigDirectories returns currently used service YML configuration locations
func GetCurrentConfigDirectories() (main, legacy, overwrite string) {
	cfgDirMu.RLock()
	defer cfgDirMu.RUnlock()
	return configDir, oldConfigDir, configOverrideDir
}

// SetConfigDirectories sets main, legacy, overwrite config directories to be used
func SetConfigDirectories(main, legacy, overwrite string) {
	cfgDirMu.Lock()
	configDir, oldConfigDir, configOverrideDir = main, legacy, overwrite
	cfgDirMu.Unlock()
}

func getServiceConfigImpl(moduleName, serviceName, configDir, oldConfigDir, configOverrideDir string) (*ConfigMap, error) {
	moduleName, serviceName = strings.ToLower(moduleName), strings.ToLower(serviceName)
	configFileName := getServiceConfigFilePath(moduleName, serviceName, configDir, oldConfigDir)
	config, err := loadYamlFile(configFileName)
	if err != nil {
		// If error - try Override cfg
		config = &ConfigMap{RawMap: map[interface{}]interface{}{}}
		log.Printf("Error Loading %s::%s configs from '%s': %v", moduleName, serviceName, configFileName, err)
	} else {
		log.Printf("Successfully loaded '%s::%s' service configs from '%s'", moduleName, serviceName, configFileName)
	}

	overrideFileName := filepath.Join(configOverrideDir, moduleName, fmt.Sprintf("%s.yml", serviceName))
	if fi, serr := os.Stat(overrideFileName); serr == nil && !fi.IsDir() {
		overrides, oerr := loadYamlFile(overrideFileName)
		if oerr != nil {
			log.Printf("Error Loading %s Override configs from '%s': %v", serviceName, overrideFileName, oerr)
			return config, err
		}
		config = updateMap(config, overrides)
		log.Printf("Successfully loaded Override configs for service %s:%s from '%s'",
			moduleName, serviceName, overrideFileName)
	}
	return config, err
}

func getServiceConfigFilePath(moduleName, serviceName, configDir, oldConfigDir string) string {
	// Filenames should be lower case
	moduleName = strings.ToLower(moduleName)
	serviceName = strings.ToLower(serviceName)

	configFileName := filepath.Join(configDir, moduleName, fmt.Sprintf("%s.yml", serviceName))
	if fi, nerr := os.Stat(configFileName); nerr != nil || fi.IsDir() {
		old := configFileName
		configFileName = filepath.Join(oldConfigDir, moduleName, fmt.Sprintf("%s.yml", serviceName))
		if fi, err := os.Stat(configFileName); err != nil || fi.IsDir() {
			log.Printf("Cannot find '%s': %v, or Legacy Service Registry Configuration: '%s': %v",
				old, nerr, configFileName, err)
		}
	}
	return configFileName
}

func updateMap(baseMap, overrides *ConfigMap) *ConfigMap {
	for k, v := range overrides.RawMap {
		if _, ok := baseMap.RawMap[k]; ok {
			baseMap.RawMap[k] = v
		}
	}
	return baseMap
}

// loadYamlFile loads a config by file name to a map of parameters
// Input: configFileName - name of config file to load, e.g. /etc/magma/control_proxy.yml
// Output: map of parameters if it exists, error if not
func loadYamlFile(configFileName string) (*ConfigMap, error) {
	yamlFile, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}
	configMap := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(yamlFile), &configMap)
	if err != nil {
		return nil, err
	}
	return &ConfigMap{configMap}, nil
}

// getStringParamImpl retrieves a string param from a ConfigMap
func getStringParamImpl(cfgMap *ConfigMap, key string) (string, error) {
	paramIface, ok := cfgMap.RawMap[key]
	if !ok {
		return "", fmt.Errorf("Key '%s' is Not Found in: %+v", key, cfgMap.RawMap)
	}
	param, ok := paramIface.(string)
	if !ok {
		return "", fmt.Errorf("Could not convert param to string for key %s", key)
	}
	return param, nil
}

// getIntParamImpl retrieves an int param from a ConfigMap
func getIntParamImpl(cfgMap *ConfigMap, key string) (int, error) {
	paramIface, ok := cfgMap.RawMap[key]
	if !ok {
		return 0, fmt.Errorf("Could not find key %s", key)
	}
	param, ok := paramIface.(int)
	if !ok {
		return 0, fmt.Errorf("Could not convert param to integer for key %s", key)
	}
	return param, nil
}

func getStringArrayParamImpl(cfgMap *ConfigMap, key string) ([]string, error) {
	paramIface, ok := cfgMap.RawMap[key]
	if !ok {
		return []string{}, fmt.Errorf("Could not find key %s", key)
	}
	var strings []string
	if reflect.TypeOf(paramIface).Kind() == reflect.Slice {
		v := reflect.ValueOf(paramIface)
		for i := 0; i < v.Len(); i++ {
			strings = append(strings, v.Index(i).Interface().(string))
		}
	} else {
		return []string{}, fmt.Errorf("could not convert param to string array for key %s", key)
	}
	return strings, nil
}

func (cfgMap *ConfigMap) getMapParamImpl(key string) (map[interface{}]interface{}, error) {
	if cfgMap == nil {
		return map[interface{}]interface{}{}, fmt.Errorf("Invalid (nil) ConfigMap")
	}
	paramIface, ok := cfgMap.RawMap[key]
	if !ok {
		return map[interface{}]interface{}{}, fmt.Errorf("Could not find key %s", key)
	}
	param, ok := paramIface.(map[interface{}]interface{})
	if !ok {
		return map[interface{}]interface{}{},
			fmt.Errorf("Could not convert %T param to map for key %s", paramIface, key)
	}
	return param, nil
}
