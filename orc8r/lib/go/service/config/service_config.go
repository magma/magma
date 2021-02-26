/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	_ "magma/orc8r/lib/go/initflag"
)

var (
	// configDir is where the per-service configuration files are stored
	configDir         = "/etc/magma/configs"
	oldConfigDir      = "/etc/magma"
	configOverrideDir = "/var/opt/magma/configs"
	cfgDirMu          sync.RWMutex
)

// GetServiceConfig loads a config by name to a map of parameters
// Input: configName - name of config to load, e.g. control_proxy
// Output: map of parameters if it exists, error if not
func GetServiceConfig(moduleName string, serviceName string) (*ConfigMap, error) {
	cfgDirMu.RLock()
	main, legacy, overwrite := configDir, oldConfigDir, configOverrideDir
	cfgDirMu.RUnlock()
	return getServiceConfigImpl(moduleName, serviceName, main, legacy, overwrite)
}

// MustGetServiceConfig is same as GetServiceConfig but fails on errors.
func MustGetServiceConfig(moduleName string, serviceName string) *ConfigMap {
	cfg, err := GetServiceConfig(moduleName, serviceName)
	if err != nil {
		glog.Fatal(err)
	}
	return cfg
}

// GetServiceConfigs returns module-keyed configs for the named service
// from all known modules.
// The list of known modules is determined by listing all non-directory files
// under /etc/magma/configs.
func GetServiceConfigs(serviceName string) (map[string]*ConfigMap, error) {
	modules, err := getModules()
	if err != nil {
		return nil, err
	}

	ret := map[string]*ConfigMap{}
	for _, moduleName := range modules {
		cfg, err := GetServiceConfig(moduleName, serviceName)
		if err != nil {
			return nil, errors.Wrapf(err, "get service config for %v.%v", moduleName, serviceName)
		}
		ret[moduleName] = cfg
	}

	return ret, nil
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
		err = yaml.Unmarshal(yamlFileData, out)
		if err != nil {
			glog.Errorf("Structured CFG: Error Unmarshaling '%s' into type %T: %v", ymlFilePath, out, err)
		} else {
			glog.Infof("Successfully loaded structured '%s::%s' service configs from '%s'",
				moduleName, serviceName, ymlFilePath)
		}
	} else {
		glog.Errorf("Structured CFG: Error Reading '%s': %v", ymlFilePath, err)
	}
	// Overwrite params from override configs
	var oerr error
	ymlQWFilePath = filepath.Join(configOverrideDir, moduleName, fmt.Sprintf("%s.yml", serviceName))
	if fi, serr := os.Stat(ymlQWFilePath); serr == nil && !fi.IsDir() {
		yamlFileData, oerr = ioutil.ReadFile(ymlQWFilePath)
		if oerr == nil {
			oerr = yaml.Unmarshal(yamlFileData, out)
			if oerr != nil {
				glog.Errorf("Structured CFG: Error Unmarshaling Override file '%s' into type %T: %v",
					ymlQWFilePath, out, err)
			} else {
				glog.Infof("Successfully loaded Override configs for service %s:%s from '%s'",
					moduleName, serviceName, ymlQWFilePath)
			}
		} else {
			glog.Errorf("Structured CFG:Error Loading Override configs from '%s': %v", ymlQWFilePath, err)
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
		glog.Errorf("Error Loading %s::%s configs from '%s': %v", moduleName, serviceName, configFileName, err)
	} else {
		glog.Infof("Successfully loaded '%s::%s' service configs from '%s'", moduleName, serviceName, configFileName)
	}

	overrideFileName := filepath.Join(configOverrideDir, moduleName, fmt.Sprintf("%s.yml", serviceName))
	if fi, serr := os.Stat(overrideFileName); serr == nil && !fi.IsDir() {
		overrides, oerr := loadYamlFile(overrideFileName)
		if oerr != nil {
			glog.Errorf("Error Loading %s Override configs from '%s': %v", serviceName, overrideFileName, oerr)
			return config, err
		}
		config = updateMap(config, overrides)
		glog.Infof("Successfully loaded Override configs for service %s:%s from '%s'",
			moduleName, serviceName, overrideFileName)
		err = nil
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
			glog.Warningf("Cannot find '%s': %v, or Legacy Service Registry Configuration: '%s': %v",
				old, nerr, configFileName, err)
		}
	}
	return configFileName
}

func updateMap(baseMap, overrides *ConfigMap) *ConfigMap {
	for k, v := range overrides.RawMap {
		baseMap.RawMap[k] = v
	}
	return baseMap
}

// getModules returns the list of known modules.
// Each directory name in /etc/magma/configs is considered a known module.
func getModules() ([]string, error) {
	moduleFiles, err := ioutil.ReadDir(configDir)
	if err != nil {
		return nil, errors.Wrap(err, "read modules from config directory")
	}
	var modules []string
	for _, m := range moduleFiles {
		if m.IsDir() {
			modules = append(modules, m.Name())
		}
	}
	return modules, nil
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
	err = yaml.Unmarshal(yamlFile, &configMap)
	if err != nil {
		return nil, err
	}
	return &ConfigMap{configMap}, nil
}
