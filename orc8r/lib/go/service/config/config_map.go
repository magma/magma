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
	"errors"
	"fmt"
	"reflect"

	"github.com/golang/glog"
)

// ConfigMap represents a map generated from a service YML file.
type ConfigMap struct {
	RawMap map[interface{}]interface{}
}

// NewConfigMap creates a new ConfigMap based on the input map.
func NewConfigMap(config map[interface{}]interface{}) *ConfigMap {
	return &ConfigMap{config}
}

// GetInt retrieves the int parameter keyed by the passed key.
func (c *ConfigMap) GetInt(key string) (int, error) {
	if err := c.Validate(); err != nil {
		return 0, err
	}

	paramIface, ok := c.RawMap[key]
	if !ok {
		return 0, fmt.Errorf("key '%s' not found in: %+v", key, c.RawMap)
	}
	param, ok := paramIface.(int)
	if !ok {
		return 0, fmt.Errorf("could not convert to integer for key %s", key)
	}
	return param, nil
}

// GetBool retrieves the bool parameter keyed by the passed key.
func (c *ConfigMap) GetBool(key string) (bool, error) {
	if err := c.Validate(); err != nil {
		return false, err
	}

	paramIface, ok := c.RawMap[key]
	if !ok {
		return false, fmt.Errorf("key '%s' not found in: %+v", key, c.RawMap)
	}
	param, ok := paramIface.(bool)
	if !ok {
		return false, fmt.Errorf("could not convert to bool for key %s", key)
	}
	return param, nil
}

// GetString retrieves the string parameter keyed by the passed key.
func (c *ConfigMap) GetString(key string) (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}

	paramIface, ok := c.RawMap[key]
	if !ok {
		return "", fmt.Errorf("key '%s' not found in: %+v", key, c.RawMap)
	}
	param, ok := paramIface.(string)
	if !ok {
		return "", fmt.Errorf("could not convert to string for key %s", key)
	}
	return param, nil
}

// GetStrings retrieves the []string parameter keyed by the passed key.
func (c *ConfigMap) GetStrings(key string) ([]string, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	paramIface, ok := c.RawMap[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found in: %+v", key, c.RawMap)
	}
	strs, err := makeStrs(paramIface)
	if err != nil {
		return nil, fmt.Errorf("could not convert to []string for key %s", key)
	}
	return strs, nil
}

// GetMap retrieves the map[interface{}]interface{} parameter keyed by the passed key.
func (c *ConfigMap) GetMap(key string) (map[interface{}]interface{}, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	paramIface, ok := c.RawMap[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found in: %+v", key, c.RawMap)
	}
	param, ok := paramIface.(map[interface{}]interface{})
	if !ok {
		return nil, fmt.Errorf("could not convert to map for key %s (value type was %T)", key, paramIface)
	}
	return param, nil
}

// MustGetInt is same as GetInt but fails on errors and when the value does not exist.
func (c *ConfigMap) MustGetInt(key string) int {
	param, err := c.GetInt(key)
	if err != nil {
		glog.Fatalf("Error retrieving %s: %v\n", key, err)
	}
	return param
}

// MustGetBool is same as GetBool but fails on errors and when the value does not exist.
func (c *ConfigMap) MustGetBool(key string) bool {
	param, err := c.GetBool(key)
	if err != nil {
		glog.Fatalf("Error retrieving %s: %v\n", key, err)
	}
	return param
}

// MustGetString is same as GetString but fails on errors and when the value does not exist.
func (c *ConfigMap) MustGetString(key string) string {
	str, err := c.GetString(key)
	if err != nil {
		glog.Fatalf("Error retrieving %s: %v\n", key, err)
	}
	return str
}

// MustGetStrings is same as GetStrings but fails on errors and when the value does not exist.
func (c *ConfigMap) MustGetStrings(key string) []string {
	param, err := c.GetStrings(key)
	if err != nil {
		glog.Fatalf("Error retrieving %s: %v\n", key, err)
	}
	return param
}

// MustGetMap is same as GetMap but fails on errors and when the value does not exist.
func (c *ConfigMap) MustGetMap(key string) map[interface{}]interface{} {
	param, err := c.GetMap(key)
	if err != nil {
		glog.Fatalf("Error retrieving %s: %v\n", key, err)
	}
	return param
}

// SetString tries to set set the given string pointer to config value
// keyed by the passed key.
func (c *ConfigMap) SetString(strPtr *string, key string) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if strPtr == nil {
		return fmt.Errorf("validate update string: strPtr cannot be nil (key: %s)", key)
	}

	str, err := c.GetString(key)
	if err == nil && len(str) > 0 {
		*strPtr = str
	}
	return err
}

// Validate the config map.
func (c *ConfigMap) Validate() error {
	if c == nil {
		return errors.New("validate config map: map cannot be nil")
	}
	return nil
}

// makeStrs tries to convert the config map param into a string slice.
func makeStrs(param interface{}) ([]string, error) {
	if reflect.TypeOf(param).Kind() != reflect.Slice {
		return nil, errors.New("param type is not slice")
	}
	v := reflect.ValueOf(param)
	strs := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		strs[i] = v.Index(i).Interface().(string) // panic on failed type assertion
	}
	return strs, nil
}
