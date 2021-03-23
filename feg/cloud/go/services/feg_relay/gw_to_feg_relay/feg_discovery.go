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

// gw_to_feg_relay is h2c & GRPC server serving requests from AGWs to FeG
package gw_to_feg_relay

import (
	"context"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/feg/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/lib/go/protos"
)

// RetrieveGatewayIdentity fetches Gateway Identity from ctx, validates & returns it
func RetrieveGatewayIdentity(inCtx context.Context) (*protos.Identity_Gateway, error) {
	identity := protos.GetClientIdentity(inCtx)
	if identity == nil {
		return nil, status.Errorf(codes.PermissionDenied, "Gateway Identity Metadata is missing")
	}
	gwId := identity.GetGateway()
	if gwId == nil {
		return nil, status.Errorf(codes.PermissionDenied, "'%s' is not a Gateway", identity.String())
	}
	if !gwId.Registered() {
		return nil, status.Errorf(codes.PermissionDenied, "Gateway '%s' is not registered", gwId.String())
	}
	return gwId, nil
}

// FindServingFeGHwId is the core of Neutral Host routing implementation, it
// 1) fetches the request's Gateway Network configs,
// 2) finds the serving Federation Gateway (FeG) Network,
// 3) retrieves the FeG Network's configuration,
// 4) if the FeG Network is a Neutral Host - FindServingFeGHwId finds the serving FeG network and the serving FeG ID
// 5) if not Neutral Host or user PLMN ID doesn't match any configured NH route - FindServingFeGHwId finds the serving
//    FeG ID of the serving FeG Network
// Returns serving Federation Gateway ID or error
func FindServingFeGHwId(agNwID, imsi string) (string, error) {
	cfg, err := configurator.LoadNetworkConfig(agNwID, feg.FederatedNetworkType, serdes.Network)
	if err != nil {
		return "", status.Errorf(
			codes.NotFound, "could not load federated network configs for access network %s: %s", agNwID, err)
	}
	federatedConfig, ok := cfg.(*models.FederatedNetworkConfigs)
	if !ok || federatedConfig == nil {
		return "", status.Errorf(codes.Internal, "invalid federated network config found for network: %s", agNwID)
	}
	if federatedConfig.FegNetworkID == nil || *federatedConfig.FegNetworkID == "" {
		return "", status.Errorf(codes.Internal, "FegNetworkID is empty in network config of network: %s", agNwID)
	}
	fegCfg, err := configurator.LoadNetworkConfig(*federatedConfig.FegNetworkID, feg.FegNetworkType, serdes.Network)
	if err != nil || fegCfg == nil {
		return "", status.Errorf(
			codes.Internal, "unable to retrieve config for federation network: %s", *federatedConfig.FegNetworkID)
	}
	networkFegConfigs, ok := fegCfg.(*models.NetworkFederationConfigs)
	if !ok || networkFegConfigs == nil {
		return "", status.Errorf(
			codes.Internal, "invalid federation network config found for network: %s", *federatedConfig.FegNetworkID)
	}
	servedNetworkIDs := networkFegConfigs.ServedNetworkIds
	for _, network := range servedNetworkIDs {
		if agNwID == network {
			servingFegNetwork := *federatedConfig.FegNetworkID
			// First check if the gateway's network is served by a Neutral Host FeG network (NhRoutes are configured &
			// IMSI is provided). If any step in finding serving NH FeG fails, try to relay the request to the
			// servingFegNetwork's FeG (non-NH logic)
			if len(networkFegConfigs.NhRoutes) > 0 && len(imsi) >= MinPlmnIdLen {
				// findServingNHFeg returns serving NH FeG Hardware ID if found or an empty string
				// given NH Routing map, calling NH network ID and user IMSI
				nhFegHwId := findServingNHFeg(networkFegConfigs.NhRoutes, servingFegNetwork, imsi)
				if len(nhFegHwId) > 0 {
					// Return here only if NH FeG was successfully found
					// in all other cases - fail back to the legacy logic and try to relay the request to the
					// servingFegNetwork's FeG
					glog.V(1).Infof(
						"routing IMSI %s request to NH FeG network: %s to FeG: %s",
						imsi, servingFegNetwork, nhFegHwId)
					return nhFegHwId, nil
				}
				glog.V(1).Infof("no NH route found for IMSI: %s", imsi)
			} else if glog.V(1) {
				if len(networkFegConfigs.NhRoutes) == 0 {
					glog.Infof("no NH route configured for Gateway Network: %s, IMSI: %s", agNwID, imsi)
				} else {
					glog.Infof("no valid IMSI (%s) for Gateway Network: %s", imsi, agNwID)
				}
			}
			return getActiveFeGForNetwork(servingFegNetwork)
		}
	}
	return "", status.Errorf(
		codes.FailedPrecondition,
		"federation network %s is not configured to serve network: %s", *federatedConfig.FegNetworkID, agNwID)
}

// findServingNHFeg returns serving NH FeG Hardware ID if found or an empty string given NH Routing map,
// calling NH network ID and user IMSI
func findServingNHFeg(routes models.NhRoutes, nhNetworkId, imsi string) string {
	var (
		servingFegNetworkId string
		found               bool
	)
	sanitizedImsi := strings.TrimPrefix(strings.TrimSpace(imsi), "IMSI")
	sanitizedLen := len(sanitizedImsi)
	if sanitizedLen < MinPlmnIdLen {
		glog.Errorf("invalid NH IMSI: '%s'", imsi)
		return ""
	}
	if sanitizedLen >= MaxPlmnIdLen {
		servingFegNetworkId, found = routes[sanitizedImsi[:MaxPlmnIdLen]]
	}
	if !found {
		if servingFegNetworkId, found = routes[sanitizedImsi[:MinPlmnIdLen]]; !found {
			return ""
		}
	}
	// verify that serving FeG network has the NH network in it's configuration
	fegCfg, err := configurator.LoadNetworkConfig(servingFegNetworkId, feg.FegNetworkType, serdes.Network)
	if err != nil || fegCfg == nil {
		glog.Errorf("unable to retrieve config for NH federation network: %s", servingFegNetworkId)
		return ""
	}
	networkFegConfigs, ok := fegCfg.(*models.NetworkFederationConfigs)
	if !ok || networkFegConfigs == nil {
		glog.Errorf("invalid federation network config found for NH network: %s", servingFegNetworkId)
		return ""
	}
	for _, network := range networkFegConfigs.ServedNhIds {
		if nhNetworkId == network {
			fegHwId, err := getActiveFeGForNetwork(servingFegNetworkId)
			if err != nil {
				glog.Errorf(
					"failed to find active FeG in '%s' NH network for IMSI: %s: %v", servingFegNetworkId, imsi, err)
				return ""
			}
			return fegHwId
		}
	}
	return ""
}
