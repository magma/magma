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

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/swx_proxy/servicers"
	"magma/feg/gateway/services/testcore/hss/storage"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"
)

// NewSAA outputs a server assignment answer (SAA) to reply to a server
// assignment request (SAR) message. See 3GPP TS 29.273 section 8.1.2.2.2.2.
func NewSAA(srv *HomeSubscriberServer, msg *diam.Message) (*diam.Message, error) {
	err := ValidateSAR(msg)
	if err != nil {
		return msg.Answer(diam.MissingAVP), err
	}

	var sar servicers.SAR
	if err := msg.Unmarshal(&sar); err != nil {
		return msg.Answer(diam.UnableToComply), fmt.Errorf("SAR Unmarshal failed for message: %v failed: %v", msg, err)
	}

	subscriber, err := srv.store.GetSubscriberData(string(sar.UserName))
	if err != nil {
		if _, ok := err.(storage.UnknownSubscriberError); ok {
			return ConstructFailureAnswer(msg, sar.SessionID, srv.Config.Server, uint32(fegprotos.ErrorCode_USER_UNKNOWN)), err
		}
		return ConstructFailureAnswer(msg, sar.SessionID, srv.Config.Server, uint32(diam.UnableToComply)), err
	}

	aaaServer := datatype.DiameterIdentity(subscriber.GetState().GetTgppAaaServerName())
	if len(aaaServer) == 0 {
		err = errors.New("no 3GPP AAA server is already serving the user")
		return ConstructFailureAnswer(msg, sar.SessionID, srv.Config.Server, uint32(diam.UnableToComply)), err
	}
	if aaaServer != sar.OriginHost {
		err = errors.New("diameter identity for AAA server already registered")
		return getRedirectMessage(msg, sar.SessionID, srv.Config.Server, aaaServer), err
	}

	if subscriber.GetNon_3Gpp().GetApnConfig() == nil {
		subscriber.State.TgppAaaServerName = ""
		err = srv.store.UpdateSubscriber(subscriber)
		if err != nil {
			glog.Errorf("Failed to remove the 3GPP AAA Server name: %v", err)
		}
		err = errors.New("User has no non 3GPP subscription")
		return ConstructFailureAnswer(msg, sar.SessionID, srv.Config.Server, uint32(fegprotos.SwxErrorCode_USER_NO_NON_3GPP_SUBSCRIPTION)), err
	}

	answer := ConstructSuccessAnswer(msg, sar.SessionID, srv.Config.Server, diam.TGPP_SWX_APP_ID)
	answer.NewAVP(avp.UserName, avp.Mbit, 0, sar.UserName)

	switch sar.ServerAssignmentType {
	case servicers.ServerAssignnmentType_USER_DEREGISTRATION:
		subscriber.State.TgppAaaServerName = ""
		subscriber.State.TgppAaaServerRegistered = false
		err = srv.store.UpdateSubscriber(subscriber)
		if err != nil {
			err = fmt.Errorf("Failed to deregister 3GPP AAA server: %v", err)
			return ConstructFailureAnswer(msg, sar.SessionID, srv.Config.Server, uint32(diam.UnableToComply)), err
		}

	case servicers.ServerAssignmentType_REGISTRATION:
		subscriber.State.TgppAaaServerRegistered = true
		err = srv.store.UpdateSubscriber(subscriber)
		if err != nil {
			err = fmt.Errorf("Failed to register 3GPP AAA server: %v", err)
			return ConstructFailureAnswer(msg, sar.SessionID, srv.Config.Server, uint32(diam.UnableToComply)), err
		}
		answer.NewAVP(avp.TGPPAAAServerName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, aaaServer)
		answer.AddAVP(getNon3GPPUserDataAVP(subscriber.GetNon_3Gpp()))

	case servicers.ServerAssignmentType_AAA_USER_DATA_REQUEST:
		answer.NewAVP(avp.TGPPAAAServerName, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, aaaServer)
		answer.AddAVP(getNon3GPPUserDataAVP(subscriber.GetNon_3Gpp()))

	default:
		err = fmt.Errorf("server assignment type not implemented: %v", sar.ServerAssignmentType)
		return ConstructFailureAnswer(msg, sar.SessionID, srv.Config.Server, uint32(diam.UnableToComply)), err
	}
	return answer, nil
}

// getNon3GPPUserDataAVP converts a Non3GPPUserProfile proto to a diameter AVP.
func getNon3GPPUserDataAVP(profile *lteprotos.Non3GPPUserProfile) *diam.AVP {
	apnConfig := profile.GetApnConfig()
	qosProfile := apnConfig[0].GetQosProfile()
	qosProfileAvp := diam.NewAVP(avp.EPSSubscribedQoSProfile, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.QoSClassIdentifier, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(qosProfile.GetClassId())),
			diam.NewAVP(avp.AllocationRetentionPriority, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.PriorityLevel, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(qosProfile.GetPriorityLevel())),
					diam.NewAVP(avp.PreemptionCapability, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(BoolToInt(qosProfile.GetPreemptionCapability()))),
					diam.NewAVP(avp.PreemptionVulnerability, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(BoolToInt(qosProfile.GetPreemptionVulnerability()))),
				},
			}),
		},
	})
	apnConfigAvp := diam.NewAVP(avp.APNConfiguration, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.ContextIdentifier, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(apnConfig[0].GetContextId())),
			diam.NewAVP(avp.PDNType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(apnConfig[0].GetPdn())),
			diam.NewAVP(avp.ServiceSelection, avp.Mbit, 0, datatype.UTF8String(apnConfig[0].GetServiceSelection())),
			diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(apnConfig[0].GetAmbr().GetMaxBandwidthUl())),
					diam.NewAVP(avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(apnConfig[0].GetAmbr().GetMaxBandwidthDl())),
				},
			}),
			qosProfileAvp,
		},
	})
	return diam.NewAVP(avp.Non3GPPUserData, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Enumerated(servicers.END_USER_E164)),
					diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(profile.GetMsisdn())),
				},
			}),
			diam.NewAVP(avp.Non3GPPIPAccess, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(profile.GetNon_3GppIpAccess())),
			diam.NewAVP(avp.Non3GPPIPAccessAPN, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(profile.GetNon_3GppIpAccessApn())),
			diam.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(servicers.RadioAccessTechnologyType_WLAN)),
			diam.NewAVP(avp.AMBR, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.MaxRequestedBandwidthUL, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(profile.GetAmbr().GetMaxBandwidthUl())),
					diam.NewAVP(avp.MaxRequestedBandwidthDL, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(profile.GetAmbr().GetMaxBandwidthDl())),
				},
			}),
			apnConfigAvp,
		},
	})
}

// ValidateSAR returns an error if the message is missing any mandatory AVPs.
// Mandatory AVPs are specified in 3GPP TS 29.273 Table 8.1.2.2.2.1/1.
func ValidateSAR(msg *diam.Message) error {
	if msg == nil {
		return errors.New("Message is nil")
	}
	_, err := msg.FindAVP(avp.UserName, 0)
	if err != nil {
		return errors.New("Missing IMSI in message")
	}
	_, err = msg.FindAVP(avp.ServerAssignmentType, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing server assignment type in message")
	}
	return nil
}
