/*
Copyright 2022 The Magma Authors.

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
	"context"
	"sync"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
)

// ConfigCacheTTL is a TTL of cached configs instances
const ConfigCacheTTL = time.Minute * 10

// EpsAuthConfig stores network related configs needed
type EpsAuthConfig struct {
	LteAuthOp          []byte
	LteAuthAmf         []byte
	SubProfiles        map[string]models.NetworkEpcConfigsSubProfilesAnon
	ApnConfigs         map[string]*models.ApnConfiguration
	ApnResources       map[string]*models.ApnResource
	ApnResourcesByName map[string]*models.ApnResource
	lastUpdate         time.Time
}

// epsCfgCache is a cache of network related configs for the service
type epsCfgCache struct {
	sync.RWMutex
	configs map[string]*EpsAuthConfig
}

// GwApnResources stores gateway related configs needed
type GwApnResources struct {
	apnResources map[string]*models.ApnResource
	lastUpdate   time.Time
}

// GwApnResourcesKey is a key for gateway related configs cache
type GwApnResourcesKey struct {
	networkID, gatewayID string
}

// epsGwCfgCache is a cache of gateway related configs for the service
type epsGwCfgCache struct {
	sync.RWMutex
	resources map[GwApnResourcesKey]*GwApnResources
}

var (
	cfgCache          = epsCfgCache{configs: map[string]*EpsAuthConfig{}}
	apnResourcesCache = epsGwCfgCache{resources: map[GwApnResourcesKey]*GwApnResources{}}
)

// GetConfig returns the EpsAuthConfig config for a given networkId.
func GetConfig(networkID string) (*EpsAuthConfig, error) {
	now := time.Now()
	cfgCache.RLock()
	cfg, found := cfgCache.configs[networkID]
	cfgCache.RUnlock()

	if found && cfg != nil && now.Before(cfg.lastUpdate.Add(ConfigCacheTTL)) {
		return cfg, nil
	}

	iCellularConfigs, err := configurator.LoadNetworkConfig(
		context.Background(), networkID, lte.CellularNetworkConfigType, serdes.Network)
	if err != nil {
		return nil, err
	}
	if iCellularConfigs == nil {
		return nil, status.Error(codes.NotFound, "got nil when looking up config")
	}
	cellularConfig, ok := iCellularConfigs.(*models.NetworkCellularConfigs)
	if !ok {
		return nil, status.Error(codes.FailedPrecondition, "failed to convert config")
	}
	epc := cellularConfig.Epc

	apnCfgs, err := loadApnsByName(networkID)
	if err != nil {
		glog.Errorf("failed to get APNs by name for network '%s': %v", networkID, err)
	}
	networkApnResources, networkApnResourcesByName, err := loadApnsResources(networkID)
	if err != nil {
		glog.Errorf("failed to get APN Resources for network '%s': %v", networkID, err)
	}
	cfg = &EpsAuthConfig{
		LteAuthOp:          epc.LteAuthOp,
		LteAuthAmf:         epc.LteAuthAmf,
		SubProfiles:        epc.SubProfiles,
		ApnConfigs:         apnCfgs,
		ApnResources:       networkApnResources,
		ApnResourcesByName: networkApnResourcesByName,
		lastUpdate:         now,
	}
	cfgCache.Lock()
	cfgCache.configs[networkID] = cfg
	cfgCache.Unlock()

	return cfg, nil
}

// GetGwApnResources returns the APN Resources configured for the given AGW & network
func GetGwApnResources(
	networkID, gwID string,
	networkApnResources, networkApnResourcesByName map[string]*models.ApnResource) map[string]*models.ApnResource {

	key := GwApnResourcesKey{networkID: networkID, gatewayID: gwID}
	now := time.Now()
	apnResourcesCache.RLock()
	rsrs, found := apnResourcesCache.resources[key]
	apnResourcesCache.RUnlock()

	if found && rsrs != nil && rsrs.apnResources != nil && now.Before(rsrs.lastUpdate.Add(ConfigCacheTTL)) {
		return rsrs.apnResources
	}
	lteGateway, err := configurator.LoadEntity(
		context.Background(), networkID, lte.CellularGatewayEntityType, gwID,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true}, serdes.Entity)
	if err != nil {
		glog.Errorf("load cellular gateway for network:gateway '%s:%s': %v", networkID, gwID, err)
		// in case of an error return all network APN Resources
		return networkApnResourcesByName
	}
	gwApnResourceIds := lteGateway.Associations.Filter(lte.APNResourceEntityType).Keys()
	res := map[string]*models.ApnResource{}
	for _, resourceId := range gwApnResourceIds {
		if resource, found := networkApnResources[resourceId]; found && resource != nil {
			res[string(resource.ApnName)] = resource
		} else {
			glog.Errorf("unmatched APN resource ID '%s' for network:gateway %s:%s", resourceId, networkID, gwID)
		}
	}
	apnResourcesCache.Lock()
	apnResourcesCache.resources[key] = &GwApnResources{apnResources: res, lastUpdate: now}
	apnResourcesCache.Unlock()
	return res
}

func loadApnsByName(networkID string) (map[string]*models.ApnConfiguration, error) {
	apns, _, err := configurator.LoadAllEntitiesOfType(
		context.Background(),
		networkID, lte.APNEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	apnsByName := map[string]*models.ApnConfiguration{}
	if err != nil {
		return apnsByName, err
	}
	for _, ent := range apns {
		apn, ok := ent.Config.(*models.ApnConfiguration)
		if !ok {
			glog.Errorf("attempt to convert entity %+v of type %T into ApnConfiguration failed.", ent.Key, ent)
			continue
		}
		apnsByName[ent.Key] = apn
	}
	return apnsByName, err
}

// loadApnsResources returns two ApnResource maps, first - keyed by the resource name & second - keyed by apn name
func loadApnsResources(networkID string) (map[string]*models.ApnResource, map[string]*models.ApnResource, error) {
	apns, _, err := configurator.LoadAllEntitiesOfType(
		context.Background(),
		networkID, lte.APNResourceEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	apnsResourcesByResourceId, apnsResourcesByApnName :=
		map[string]*models.ApnResource{}, map[string]*models.ApnResource{}
	if err != nil {
		return apnsResourcesByResourceId, apnsResourcesByApnName, err
	}
	for _, ent := range apns {
		apnRes, ok := ent.Config.(*models.ApnResource)
		if !ok {
			glog.Errorf("attempt to convert entity %+v of type %T into ApnResource failed.", ent.Key, ent)
			continue
		}
		apnsResourcesByResourceId[ent.Key] = apnRes
		apnsResourcesByApnName[string(apnRes.ApnName)] = apnRes
	}
	return apnsResourcesByResourceId, apnsResourcesByApnName, err
}
