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
	"errors"
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/s6a_proxy/servicers"
	"magma/feg/gateway/services/testcore/hss/storage"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"
)

const (
	allAPNConfigurationsIncludedIndicator = 0
	apnContextIdentifier                  = 0
	apnNetworkAccessMode                  = 2
	apnServiceSelection                   = "oai.ipv4"
	apnQoSClassIdentifier                 = 9
	epsPriorityLevel                      = 15
	epsPreemptionCapability               = 1
	epsPreemptionVulnerability            = 0
	ulaFlags                              = 0
	msisdn                                = "12345"
	accessRestrictionData                 = 47
	subscriberStatus                      = 0
)

// NewULA outputs a update location answer (ULA) to reply to an
// update location request (ULR) message.
func NewULA(srv *HomeSubscriberServer, msg *diam.Message) (*diam.Message, error) {
	err := ValidateULR(msg)
	if err != nil {
		return msg.Answer(diam.MissingAVP), err
	}

	var ulr servicers.ULR
	if err := msg.Unmarshal(&ulr); err != nil {
		return msg.Answer(diam.UnableToComply), fmt.Errorf("ULR Unmarshal failed for message: %v failed: %v", msg, err)
	}

	subscriber, err := srv.store.GetSubscriberData(string(ulr.UserName))
	if err != nil {
		if _, ok := err.(storage.UnknownSubscriberError); ok {
			return ConstructFailureAnswer(msg, ulr.SessionID, srv.Config.Server, uint32(protos.ErrorCode_USER_UNKNOWN)), err
		}
		return ConstructFailureAnswer(msg, ulr.SessionID, srv.Config.Server, uint32(protos.ErrorCode_AUTHENTICATION_DATA_UNAVAILABLE)), err
	}

	profile, ok := srv.Config.SubProfiles[subscriber.SubProfile]
	if !ok || profile == nil {
		profile = srv.Config.DefaultSubProfile
		if profile == nil {
			answer := ConstructFailureAnswer(msg, ulr.SessionID, srv.Config.Server, uint32(protos.ErrorCode_UNKNOWN_EPS_SUBSCRIPTION))
			return answer, fmt.Errorf("unknown subscriber profile: %s and default profile was not initialized", subscriber.SubProfile)
		} else {
			glog.V(2).Infof("Subscriber profile '%s' not found, using default profile instead", subscriber.SubProfile)
		}
	}

	if !isRATTypeAllowed(uint32(ulr.RATType)) {
		answer := ConstructFailureAnswer(msg, ulr.SessionID, srv.Config.Server, uint32(protos.ErrorCode_RAT_NOT_ALLOWED))
		return answer, fmt.Errorf("RAT-Type not allowed: %v", uint32(ulr.RATType))
	}

	return srv.NewSuccessfulULA(msg, ulr.SessionID, profile), nil
}

// NewSuccessfulULA outputs a successful update location answer (ULA) to reply to an
// update location request (ULR) message. It populates the ULA with all of the mandatory fields
// and adds the subscriber profile information.
func (srv *HomeSubscriberServer) NewSuccessfulULA(msg *diam.Message, sessionID datatype.UTF8String, profile *mconfig.HSSConfig_SubscriptionProfile) *diam.Message {
	ula := ConstructSuccessAnswer(msg, sessionID, srv.Config.Server, diam.TGPP_S6A_APP_ID)
	ula.NewAVP(avp.ULAFlags, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(ulaFlags))
	ula.NewAVP(avp.SubscriptionData, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.MSISDN, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(msisdn)),
			diam.NewAVP(avp.AccessRestrictionData, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(accessRestrictionData)),
			diam.NewAVP(avp.SubscriberStatus, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(subscriberStatus)),
			diam.NewAVP(avp.NetworkAccessMode, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(apnNetworkAccessMode)),
			diam.NewAVP(avp.APNConfigurationProfile, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.ContextIdentifier, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(apnContextIdentifier)),
					diam.NewAVP(avp.AllAPNConfigurationsIncludedIndicator, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(allAPNConfigurationsIncludedIndicator)),
					diam.NewAVP(avp.APNConfiguration, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
						AVP: []*diam.AVP{
							diam.NewAVP(avp.ContextIdentifier, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(apnContextIdentifier)),
							diam.NewAVP(avp.PDNType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(protos.UpdateLocationAnswer_APNConfiguration_IPV4)),
							diam.NewAVP(avp.ServiceSelection, avp.Mbit, 0, datatype.UTF8String(apnServiceSelection)),
							diam.NewAVP(avp.EPSSubscribedQoSProfile, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
								AVP: []*diam.AVP{
									diam.NewAVP(avp.QoSClassIdentifier, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(apnQoSClassIdentifier)),
									diam.NewAVP(avp.AllocationRetentionPriority, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
										AVP: []*diam.AVP{
											diam.NewAVP(avp.PriorityLevel, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(epsPriorityLevel)),
											diam.NewAVP(avp.PreemptionCapability, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(epsPreemptionCapability)),
											diam.NewAVP(avp.PreemptionVulnerability, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(epsPreemptionVulnerability)),
										},
									}),
								},
							}),
							diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
								AVP: []*diam.AVP{
									diam.NewAVP(avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(profile.MaxDlBitRate)),
									diam.NewAVP(avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(profile.MaxUlBitRate)),
								},
							}),
						},
					}),
				},
			}),
			diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(profile.MaxDlBitRate)),
					diam.NewAVP(avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(profile.MaxUlBitRate)),
				},
			}),
		},
	})
	return ula
}

// ValidateULR returns an error if the message is missing any mandatory AVPs.
// Mandatory AVPs are specified in 3GPP TS 29.272 Table 5.2.1.1.1/1
func ValidateULR(msg *diam.Message) error {
	_, err := msg.FindAVP(avp.UserName, 0)
	if err != nil {
		return errors.New("Missing IMSI in message")
	}
	_, err = msg.FindAVP(avp.VisitedPLMNID, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing Visited PLMN ID in message")
	}
	_, err = msg.FindAVP(avp.ULRFlags, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing ULR flags in message")
	}
	_, err = msg.FindAVP(avp.RATType, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing RAT type in message")
	}
	_, err = msg.FindAVP(avp.SessionID, 0)
	if err != nil {
		return errors.New("Missing SessionID in message")
	}
	return nil
}
