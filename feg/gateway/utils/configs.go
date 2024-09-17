/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package utils

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/golang/glog"
)

// GetValueUint64 returns value of the flagValue if it exists, or defaultValue if not
func GetValueUint64(flagName string, defaultValue uint64) uint64 {
	if len(flagName) > 0 {
		value, err := getUint64FlagValue(flagName)
		if err != nil {
			return value
		}
	}
	glog.V(1).Infof("Using value: %v for flag: %s", defaultValue, flagName)
	return defaultValue
}

// GetValue returns value of the flagValue if it exists, or defaultValue if not
func GetValue(flagName, defaultValue string) string {
	flagValue := getFlagValue(flagName)
	if len(flagValue) != 0 {
		return flagValue
	}
	glog.V(1).Infof("Using value: %s for flag: %s", defaultValue, flagName)
	return defaultValue
}

// GetValueOrEnv returns value of the flagValue if it exists, then the environment
// variable if it exists, or defaultValue if not.
func GetValueOrEnv(flagName, envVariable, defaultValue string) string {
	flagValue := getFlagValue(flagName)
	if len(flagValue) != 0 {
		return flagValue
	}
	if len(envVariable) > 0 {
		envValue := os.Getenv(envVariable)
		if len(envValue) > 0 {
			glog.V(1).Infof(
				"Using Environment Parameter: %s => %s (default: '%s')", envVariable, envValue, defaultValue)
			return envValue
		}
	}
	glog.V(1).Infof("Using value: %s for flag: %s", defaultValue, flagName)
	return defaultValue
}

// GetBoolValueOrEnv returns value of the flagValue if it exists, then the environment
// variable if it exists, or defaultValue if not.
func GetBoolValueOrEnv(flagName string, envVariable string, defaultValue bool) bool {
	flagValue := getFlagValue(flagName)
	flagValueBool, err := strconv.ParseBool(flagValue)
	if len(flagValue) != 0 && err == nil {
		return flagValueBool
	}
	if len(envVariable) > 0 {
		envValue := os.Getenv(envVariable)
		envValueBool, err := strconv.ParseBool(envValue)
		if len(envValue) > 0 && err == nil {
			glog.V(1).Infof(
				"Using Environment Parameter: %s => %t (default: '%t')", envVariable, envValueBool, defaultValue)
			return envValueBool
		}
	}
	glog.V(1).Infof("Using value: %t for flag: %s", defaultValue, flagName)
	return defaultValue
}

// getUint64FlagValue looks up the flag and either returns its uint64 value
// or an error.
func getUint64FlagValue(flagName string) (uint64, error) {
	f := flag.Lookup(flagName)
	if f == nil {
		return 0, fmt.Errorf("Flag not found: %s", flagName)
	}
	if f.Value == nil {
		return 0, fmt.Errorf("Flag value is nil: %s", flagName)
	}
	getter, ok := f.Value.(flag.Getter)
	if !ok {
		return 0, fmt.Errorf("Flag value has no Getter: %s", flagName)
	}
	value, ok := getter.Get().(uint64)
	if !ok {
		return 0, fmt.Errorf("Flag value is not of type uint64: %s", flagName)
	}
	return value, nil
}

// getFlagValue returns the value of the flagValue if it exists, or an empty string if not
func getFlagValue(flagName string) string {
	var res string
	if len(flagName) > 0 {
		flag.Visit(func(f *flag.Flag) {
			if f.Name == flagName {
				res = f.Value.String()
				glog.V(1).Infof("Using runtime flag: %s => %s", flagName, res)
			}
		})
	}
	return res
}
