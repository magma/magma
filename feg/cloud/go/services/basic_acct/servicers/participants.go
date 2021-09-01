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

package servicers

import (
	"context"
	"strings"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/feg"
	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/serdes"
	"magma/feg/cloud/go/services/feg/obsidian/models"
	"magma/feg/cloud/go/services/feg_relay/gw_to_feg_relay"
	"magma/orc8r/cloud/go/services/configurator"
)

const (
	// Minimal length of PLMNID
	MinPlmnIdLen = 5
	MaxPlmnIdLen = 6
)

// RetrieveParticipants retrieves GW Network, FeG Network & GW logical ID from the caller's context
func RetrieveParticipants(ctx context.Context, session *protos.AcctSession) (provider, consumer, gw string, err error) {
	if session == nil {
		err = status.Errorf(codes.InvalidArgument, "nil Session")
		return
	}
	gwId, err := gw_to_feg_relay.RetrieveGatewayIdentity(ctx)
	if err != nil {
		return
	}
	gw = gwId.GetLogicalId()
	provider = gwId.GetNetworkId()
	imsi := session.GetIMSI()
	if len(imsi) == 0 {
		err = status.Errorf(
			codes.InvalidArgument, "no IMSI for session: %s, serving network: %s", session.GetSessionId(), provider)
		return
	}
	consumer, err = FindServingFeGNetworkId(provider, imsi)
	return
}

// FindServingFeGNetworkId
// 1) fetches the request's Gateway Network configs,
// 2) finds the serving Federation Gateway (FeG) Network,
// 3) retrieves the FeG Network's configuration,
// 4) if the FeG Network is a Neutral Host - FindServingFeGNetworkId finds the serving FeG network
// 5) if not Neutral Host or user PLMN ID doesn't match any configured NH route - FindServingFeGNetworkId finds the
//    serving FeG Network
// Returns serving Federation Gateway Network ID or error
func FindServingFeGNetworkId(agNwID, imsi string) (string, error) {
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
			// IMSI is provided). If any step in finding serving NH FeG fails, servingFegNetwork's FeG would auth the
			// user (non-NH logic)
			if len(networkFegConfigs.NhRoutes) > 0 && len(imsi) >= MinPlmnIdLen {
				nhFegNetwork := findServingNHFegNetwork(networkFegConfigs.NhRoutes, servingFegNetwork, imsi)
				if len(nhFegNetwork) > 0 {
					// Return here only if NH FeG Network was successfully found & verified
					// in all other cases - fail back to the legacy logic and try to relay the request to the
					// servingFegNetwork's FeG
					glog.V(1).Infof("serving IMSI %s for NH FeG network: %s", imsi, servingFegNetwork)
					return nhFegNetwork, nil
				}
				glog.V(1).Infof("no NH route found for IMSI: %s", imsi)
			} else if glog.V(1) {
				if len(networkFegConfigs.NhRoutes) == 0 {
					glog.Infof("no NH route configured for Gateway Network: %s, IMSI: %s", agNwID, imsi)
				} else {
					glog.Infof("no valid IMSI (%s) for Gateway Network: %s", imsi, agNwID)
				}
			}
			return servingFegNetwork, nil
		}
	}
	return "", status.Errorf(
		codes.FailedPrecondition,
		"federation network %s is not configured to serve network: %s", *federatedConfig.FegNetworkID, agNwID)
}

// verifyServingNHFegNetwork returns true if selected NH FeG network has corresponding configuration for given IMSI
func findServingNHFegNetwork(routes models.NhRoutes, nhNetworkId, imsi string) (servingFeg string) {
	var (
		servingFegNetworkId string
		found               bool
	)
	sanitizedImsi := strings.TrimPrefix(strings.TrimSpace(imsi), "IMSI")
	sanitizedLen := len(sanitizedImsi)
	if sanitizedLen < MinPlmnIdLen {
		glog.Errorf("invalid NH IMSI: '%s'", imsi)
		return
	}
	if sanitizedLen >= MaxPlmnIdLen {
		servingFegNetworkId, found = routes[sanitizedImsi[:MaxPlmnIdLen]]
	}
	if !found {
		if servingFegNetworkId, found = routes[sanitizedImsi[:MinPlmnIdLen]]; !found {
			return
		}
	}
	// verify that serving FeG network has the NH network in it's configuration
	fegCfg, err := configurator.LoadNetworkConfig(servingFegNetworkId, feg.FegNetworkType, serdes.Network)
	if err != nil || fegCfg == nil {
		glog.Errorf("unable to retrieve config for NH federation network: %s", servingFegNetworkId)
		return
	}
	networkFegConfigs, ok := fegCfg.(*models.NetworkFederationConfigs)
	if !ok || networkFegConfigs == nil {
		glog.Errorf("invalid federation network config found for NH network: %s", servingFegNetworkId)
		return
	}
	for _, network := range networkFegConfigs.ServedNhIds {
		if nhNetworkId == network {
			return servingFegNetworkId
		}
	}
	return
}
