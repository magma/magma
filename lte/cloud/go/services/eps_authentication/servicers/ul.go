/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"errors"
	"fmt"

	"magma/feg/cloud/go/protos"
	cellular "magma/lte/cloud/go/services/cellular/protos"
	"magma/lte/cloud/go/services/eps_authentication/crypto"
	"magma/lte/cloud/go/services/eps_authentication/metrics"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	serviceSelection                  = "oai.ipv4"
	qosProfileClassID                 = 9
	qosProfilePriorityLevel           = 15
	qosProfilePreemptionCapability    = true
	qosProfilePreemptionVulnerability = false
)

func (srv *EPSAuthServer) UpdateLocation(ctx context.Context, ulr *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error) {
	metrics.ULRequests.Inc()
	if err := validateULR(ulr); err != nil {
		metrics.InvalidRequests.Inc()
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	networkID, err := getNetworkID(ctx)
	if err != nil {
		metrics.NetworkIDErrors.Inc()
		return nil, err
	}
	config, err := getConfig(networkID)
	if err != nil {
		metrics.ConfigErrors.Inc()
		return nil, err
	}
	subscriber, errorCode, err := srv.lookupSubscriber(ulr.UserName, networkID)
	if err != nil {
		metrics.UnknownSubscribers.Inc()
		return &protos.UpdateLocationAnswer{ErrorCode: errorCode}, err
	}
	profile := getSubProfile(subscriber.SubProfile, config)
	if profile == nil {
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
		Apn: []*protos.UpdateLocationAnswer_APNConfiguration{
			&protos.UpdateLocationAnswer_APNConfiguration{
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
		},
	}, nil
}

// getSubProfile looks up the subscription profile to be used for a subscriber.
func getSubProfile(profileName string, config *EpsAuthConfig) *cellular.NetworkEPCConfig_SubscriptionProfile {
	profile, ok := config.SubProfiles[profileName]
	if ok && profile != nil {
		return profile
	}
	metrics.UnknownSubProfiles.Inc()

	profile, ok = config.SubProfiles["default"]
	if ok && profile != nil {
		glog.V(2).Infof("Subscriber profile '%s' not found, using default profile instead", profileName)
		return profile
	}

	return nil
}

// validateULR returns an error iff the ULR is invalid.
func validateULR(ulr *protos.UpdateLocationRequest) error {
	if ulr == nil {
		return errors.New("received a nil UpdateLocationRequest")
	}
	if len(ulr.UserName) == 0 {
		return errors.New("user name was empty")
	}
	if len(ulr.VisitedPlmn) != crypto.ExpectedPlmnBytes {
		return fmt.Errorf("expected Visited PLMN to be %v bytes, but got %v bytes", crypto.ExpectedPlmnBytes, len(ulr.VisitedPlmn))
	}
	return nil
}
