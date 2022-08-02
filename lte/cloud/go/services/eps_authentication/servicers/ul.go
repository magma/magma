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
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/magma/milenage"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/services/eps_authentication/metrics"
	"magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/orc8r/cloud/go/identity"
)

const (
	serviceSelection                  = "magma.ipv4"
	qosProfileClassID                 = 9
	qosProfilePriorityLevel           = 15
	qosProfilePreemptionCapability    = true
	qosProfilePreemptionVulnerability = false

	defaultMaxUlBitRate uint64 = 2000000000
	defaultMaxDlBitRate uint64 = 4000000000
)

var defaultSubscriberProfile = models.NetworkEpcConfigsSubProfilesAnon{
	MaxDlBitRate: defaultMaxUlBitRate,
	MaxUlBitRate: defaultMaxDlBitRate,
}

func (srv *EPSAuthServer) UpdateLocation(
	ctx context.Context, ulr *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error) {

	glog.V(2).Infof("received ULR from: %s", ulr.GetUserName())
	metrics.ULRequests.Inc()
	if err := validateULR(ulr); err != nil {
		glog.V(2).Infof("ULR is invalid: %v", err.Error())
		metrics.InvalidRequests.Inc()
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	networkID, gatewayID, err := identity.GetClientNetworkAndGatewayID(ctx)
	if err != nil {
		glog.V(2).Infof("could not lookup networkID: %v", err.Error())
		metrics.NetworkIDErrors.Inc()
		return nil, err
	}
	config, err := GetConfig(networkID)
	if err != nil {
		glog.V(2).Infof("could not lookup config for networkID '%s': %v", networkID, err.Error())
		metrics.ConfigErrors.Inc()
		return nil, err
	}
	subscriber, staticIps, apns, errorCode, err := srv.lookupSubscriberProfile(ulr.UserName, networkID)
	if err != nil {
		glog.V(2).Infof("failed to lookup subscriber '%s': %v", ulr.UserName, err.Error())
		metrics.UnknownSubscribers.Inc()
		return &protos.UpdateLocationAnswer{ErrorCode: errorCode}, err
	}
	profile := getSubProfile(subscriber.SubProfile, config)
	if profile == nil {
		glog.V(2).Infof("failed to find subscriber profile '%s'", subscriber.SubProfile)
		return &protos.UpdateLocationAnswer{ErrorCode: protos.ErrorCode_UNKNOWN_EPS_SUBSCRIPTION},
			status.Errorf(
				codes.FailedPrecondition,
				"unknown subscriber profile: %s and default profile was not initialized",
				subscriber.SubProfile,
			)
	}
	return &protos.UpdateLocationAnswer{
		ErrorCode: protos.ErrorCode_SUCCESS,
		TotalAmbr: &protos.UpdateLocationAnswer_AggregatedMaximumBitrate{
			MaxBandwidthUl: uint32(profile.MaxUlBitRate),
			MaxBandwidthDl: uint32(profile.MaxDlBitRate),
		},
		Apn: createApns(networkID, gatewayID, staticIps, apns, profile, config),
	}, nil
}

// getSubProfile looks up the subscription profile to be used for a subscriber.
func getSubProfile(profileName string, config *EpsAuthConfig) *models.NetworkEpcConfigsSubProfilesAnon {
	profile, ok := config.SubProfiles[profileName]
	if ok {
		return &profile
	}
	metrics.UnknownSubProfiles.Inc()

	profile, ok = config.SubProfiles["default"]
	if ok {
		glog.V(2).Infof("Subscriber profile '%s' not found, using 'default' network profile instead", profileName)
		return &profile
	}
	glog.V(2).Info("Network subscriber profile 'default' is not configured, using defaults")
	profile = defaultSubscriberProfile
	return &profile
}

// validateULR returns an error if the ULR is invalid.
func validateULR(ulr *protos.UpdateLocationRequest) error {
	if ulr == nil {
		return errors.New("received a nil UpdateLocationRequest")
	}
	if len(ulr.UserName) == 0 {
		return errors.New("user name was empty")
	}
	if len(ulr.VisitedPlmn) != milenage.ExpectedPlmnBytes {
		return fmt.Errorf("expected Visited PLMN to be %v bytes, but got %v bytes", milenage.ExpectedPlmnBytes, len(ulr.VisitedPlmn))
	}
	return nil
}

// createApns returns a list of APN Configs for ULA
func createApns(
	networkID, gatewayID string,
	staticIps map[string]string,
	apns []string,
	profile *models.NetworkEpcConfigsSubProfilesAnon,
	netCfg *EpsAuthConfig) []*protos.UpdateLocationAnswer_APNConfiguration {

	if netCfg == nil || len(netCfg.ApnConfigs) == 0 || len(apns) == 0 {
		return []*protos.UpdateLocationAnswer_APNConfiguration{
			{
				Ambr: &protos.UpdateLocationAnswer_AggregatedMaximumBitrate{
					MaxBandwidthUl: uint32(profile.MaxUlBitRate),
					MaxBandwidthDl: uint32(profile.MaxDlBitRate),
				},
				ServiceSelection: serviceSelection,
				Pdn:              protos.UpdateLocationAnswer_APNConfiguration_IPV4,
				QosProfile: &protos.UpdateLocationAnswer_APNConfiguration_QoSProfile{
					ClassId:                 qosProfileClassID,
					PriorityLevel:           qosProfilePriorityLevel,
					PreemptionCapability:    qosProfilePreemptionCapability,
					PreemptionVulnerability: qosProfilePreemptionVulnerability,
				},
			},
		}
	}
	if staticIps == nil {
		staticIps = map[string]string{}
	}
	res := []*protos.UpdateLocationAnswer_APNConfiguration{}
	gwApnResources := GetGwApnResources(networkID, gatewayID, netCfg.ApnResources, netCfg.ApnResourcesByName)

	for _, apnName := range apns {
		apn, found := netCfg.ApnConfigs[apnName]
		if !found {
			glog.Warningf("failed to find APN '%s' in network '%s'", apnName, networkID)
			continue
		}
		apnCfg := &protos.UpdateLocationAnswer_APNConfiguration{ServiceSelection: apnName}
		if sip, found := staticIps[apnName]; found {
			apnCfg.ServedPartyIpAddress = []string{sip}
		}
		if apn != nil {
			apnCfg.Pdn = protos.UpdateLocationAnswer_APNConfiguration_PDNType(apn.PdnType)
			if p := apn.QosProfile; p != nil {
				apnCfg.QosProfile = &protos.UpdateLocationAnswer_APNConfiguration_QoSProfile{}
				if p.ClassID != nil {
					apnCfg.QosProfile.ClassId = *p.ClassID
				}
				if p.PriorityLevel != nil {
					apnCfg.QosProfile.PriorityLevel = *p.PriorityLevel
				}
				if p.PreemptionCapability != nil {
					apnCfg.QosProfile.PreemptionCapability = *p.PreemptionCapability
				}
				if p.PreemptionVulnerability != nil {
					apnCfg.QosProfile.PreemptionVulnerability = *p.PreemptionVulnerability
				}
			}
			apnCfg.Ambr = &protos.UpdateLocationAnswer_AggregatedMaximumBitrate{
				MaxBandwidthUl: uint32(profile.MaxUlBitRate),
				MaxBandwidthDl: uint32(profile.MaxDlBitRate),
			}
			if a := apn.Ambr; a != nil {
				if a.MaxBandwidthDl != nil {
					apnCfg.Ambr.MaxBandwidthDl = *a.MaxBandwidthDl
				}
				if a.MaxBandwidthUl != nil {
					apnCfg.Ambr.MaxBandwidthUl = *a.MaxBandwidthUl
				}
			}
			if ar, found := gwApnResources[apnName]; found && ar != nil {
				apnCfg.Resource = &protos.UpdateLocationAnswer_APNConfiguration_APNResource{
					ApnName:    apnName,
					GatewayIp:  ar.GatewayIP.String(),
					GatewayMac: ar.GatewayMac.String(),
					VlanId:     ar.VlanID,
				}
			}
		}
		res = append(res, apnCfg)
	}
	return res
}
