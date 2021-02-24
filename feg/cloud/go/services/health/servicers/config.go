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

package servicers

import (
	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/feg/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/golang/glog"
)

const (
	defaultCpuUtilThreshold      = 0.75
	defaultMemAvailableThreshold = 0.90
	defaultStaleUpdateThreshold  = 30
)

var defaultServices = []string{"SWX_PROXY", "SESSION_PROXY"}

func GetHealthConfigForNetwork(networkID string) *healthConfig {
	defaultConfig := &healthConfig{
		services:              defaultServices,
		cpuUtilThreshold:      defaultCpuUtilThreshold,
		memAvailableThreshold: defaultMemAvailableThreshold,
		staleUpdateThreshold:  defaultStaleUpdateThreshold,
	}
	config, err := configurator.LoadNetworkConfig(networkID, feg.FegNetworkType, serdes.Network)
	if err != nil {
		glog.V(2).Infof("Using default health configuration for network %s; %s", networkID, err)
		return defaultConfig
	}
	cloudFegConfig, ok := config.(*models.NetworkFederationConfigs)
	if !ok {
		glog.V(2).Infof("Using default health configuration for network %s; Invalid config format", networkID)
		return defaultConfig
	}
	healthParams := cloudFegConfig.Health
	if healthParams == nil {
		glog.V(2).Infof("Using default health configuration for network %s; Health config not found", networkID)
		return defaultConfig
	}
	if healthParams.CPUUtilizationThreshold == 0 {
		glog.V(2).Infof("Using default health configuration for network %s; Cpu utilization threshold cannot be 0", networkID)
		return defaultConfig
	}
	if healthParams.MemoryAvailableThreshold == 0 {
		glog.V(2).Infof("Using default health configuration for network %s; Memory available threshold cannot be 0", networkID)
		return defaultConfig
	}
	staleUpdateThreshold := healthParams.UpdateFailureThreshold * healthParams.UpdateIntervalSecs
	if staleUpdateThreshold == 0 {
		glog.V(2).Infof("Using default health configuration for network %s; Stale update threshold cannot be 0", networkID)
		return defaultConfig
	}
	return &healthConfig{
		services:              healthParams.HealthServices,
		cpuUtilThreshold:      healthParams.CPUUtilizationThreshold,
		memAvailableThreshold: healthParams.MemoryAvailableThreshold,
		staleUpdateThreshold:  staleUpdateThreshold,
	}
}
