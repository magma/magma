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

	"github.com/emakeev/milenage"
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	swx "magma/feg/gateway/services/swx_proxy/servicers"
	"magma/feg/gateway/services/testcore/hss/storage"
	lteprotos "magma/lte/cloud/go/protos"
)

// NewMAA outputs a multimedia authentication answer (MAA) to reply to a multimedia
// authentication request (MAR) message.
func NewMAA(srv *HomeSubscriberServer, msg *diam.Message) (*diam.Message, error) {
	err := ValidateMAR(msg)
	if err != nil {
		return msg.Answer(diam.MissingAVP), err
	}

	var mar swx.MAR
	if err := msg.Unmarshal(&mar); err != nil {
		return msg.Answer(diam.UnableToComply), fmt.Errorf("MAR Unmarshal failed for message: %v failed: %v", msg, err)
	}

	subscriber, err := srv.store.GetSubscriberData(mar.UserName)
	if err != nil {
		if _, ok := err.(storage.UnknownSubscriberError); ok {
			return ConstructFailureAnswer(msg, mar.SessionID, srv.Config.Server, uint32(fegprotos.ErrorCode_USER_UNKNOWN)), err
		}
		return ConstructFailureAnswer(msg, mar.SessionID, srv.Config.Server, uint32(diam.UnableToComply)), err
	}

	subscriber.Lock()
	defer subscriber.Unlock()

	if !isRATTypeAllowed(uint32(mar.RATType)) {
		answer := ConstructFailureAnswer(msg, mar.SessionID, srv.Config.Server, uint32(fegprotos.ErrorCode_RAT_NOT_ALLOWED))
		return answer, fmt.Errorf("RAT-Type not allowed: %v", uint32(mar.RATType))
	}

	aaaServer := datatype.DiameterIdentity(subscriber.GetState().GetTgppAaaServerName())
	if len(aaaServer) == 0 {
		err = srv.set3GPPAAAServerName(subscriber, mar.OriginHost)
		if err != nil {
			return ConstructFailureAnswer(msg, mar.SessionID, srv.Config.Server, uint32(diam.UnableToComply)), err
		}
	} else if aaaServer != mar.OriginHost {
		err = errors.New("diameter identity for AAA server already registered")
		return getRedirectMessage(msg, mar.SessionID, srv.Config.Server, aaaServer), err
	}

	lteAuthNextSeq, err := ResyncLteAuthSeq(subscriber, mar.AuthData.Authorization.Serialize(), srv.Config.LteAuthOp)
	if err == nil {
		err = srv.setLteAuthNextSeq(subscriber, lteAuthNextSeq)
	}
	if err != nil {
		return ConvertAuthErrorToFailureMessage(err, msg, mar.SessionID, srv.Config.Server), err
	}

	if mar.AuthData.AuthScheme != swx.SipAuthScheme_EAP_AKA {
		err = fmt.Errorf("Unsupported SIP authentication scheme: %s", mar.AuthData.AuthScheme)
		return ConstructFailureAnswer(msg, mar.SessionID, srv.Config.Server, uint32(diam.UnableToComply)), err
	}

	vectors, lteAuthNextSeq, err := srv.GenerateSIPAuthVectors(subscriber, mar.NumberAuthItems)
	if err == nil {
		err = srv.setLteAuthNextSeq(subscriber, lteAuthNextSeq)
	}
	if err != nil {
		// If we generated any auth vectors successfully, then we can return them.
		// Otherwise, we must signal an error.
		// See 3GPP TS 29.273 section 8.1.2.1.2.
		if len(vectors) == 0 {
			return ConvertAuthErrorToFailureMessage(err, msg, mar.SessionID, srv.Config.Server), err
		}
	}

	return srv.NewSuccessfulMAA(msg, mar.SessionID, datatype.UTF8String(mar.UserName), vectors), nil
}

// NewSuccessfulMAA outputs a successful multimedia authentication answer (MAA) to reply to an
// multimedia authentication request (MAR) message. It populates the MAA with all of the mandatory fields
// and adds the authentication vectors. See 3GPP TS 29.273 table 8.1.2.1.1/5.
func (srv *HomeSubscriberServer) NewSuccessfulMAA(msg *diam.Message, sessionID datatype.UTF8String, userName datatype.UTF8String, vectors []*milenage.SIPAuthVector) *diam.Message {
	maa := ConstructSuccessAnswer(msg, sessionID, srv.Config.Server, diam.TGPP_SWX_APP_ID)
	for itemNumber, vector := range vectors {
		authenticate := append(vector.Rand[:], vector.Autn[:]...)
		maa.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.SIPItemNumber, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(itemNumber)),
				diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String(swx.SipAuthScheme_EAP_AKA)),
				diam.NewAVP(avp.SIPAuthenticate, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(authenticate)),
				diam.NewAVP(avp.SIPAuthorization, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Xres[:])),
				diam.NewAVP(avp.ConfidentialityKey, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.ConfidentialityKey[:])),
				diam.NewAVP(avp.IntegrityKey, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.IntegrityKey[:])),
			},
		})
	}
	maa.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(len(vectors)))
	maa.NewAVP(avp.UserName, avp.Mbit, 0, userName)
	return maa
}

// GenerateSIPAuthVectors generates `numVectors` SIP auth vectors for the subscriber.
// The vectors and the next value of lteAuthNextSeq are returned (or an error).
func (srv *HomeSubscriberServer) GenerateSIPAuthVectors(subscriber *lteprotos.SubscriberData, numVectors uint32) ([]*milenage.SIPAuthVector, uint64, error) {
	var vectors = make([]*milenage.SIPAuthVector, 0, numVectors)
	lteAuthNextSeq := subscriber.GetState().GetLteAuthNextSeq()
	for i := uint32(0); i < numVectors; i++ {
		vector, nextSeq, err := srv.GenerateSIPAuthVector(subscriber)
		if err != nil {
			if i == 0 {
				return nil, 0, err
			}
			glog.Errorf("failed to generate SWX auth vector: %v", err)
			break
		}
		lteAuthNextSeq = nextSeq
		subscriber.State.LteAuthNextSeq = lteAuthNextSeq
		vectors = append(vectors, vector)
	}
	return vectors, lteAuthNextSeq, nil
}

// GenerateSIPAuthVector returns the SIP auth vector and the next value of lteAuthNextSeq for the subscriber (or an error).
func (srv *HomeSubscriberServer) GenerateSIPAuthVector(subscriber *lteprotos.SubscriberData) (*milenage.SIPAuthVector, uint64, error) {
	lte := subscriber.Lte
	if err := ValidateLteSubscription(lte); err != nil {
		return nil, 0, NewAuthRejectedError(err.Error())
	}
	if subscriber.State == nil {
		return nil, 0, NewAuthRejectedError("Subscriber data missing subscriber state")
	}

	opc, err := GetOrGenerateOpc(lte, srv.Config.LteAuthOp)
	if err != nil {
		return nil, 0, err
	}

	sqn := SeqToSqn(subscriber.State.LteAuthNextSeq, srv.AuthSqnInd)
	vector, err := srv.Milenage.GenerateSIPAuthVector(lte.AuthKey, opc, sqn)
	if err != nil {
		return nil, 0, NewAuthRejectedError(err.Error())
	}
	return vector, subscriber.State.LteAuthNextSeq + 1, err
}

// ValidateMAR returns an error if the message is missing any mandatory AVPs.
// Mandatory AVPs are specified in 3GPP TS 29.273 Table 8.1.2.1.1/1.
func ValidateMAR(msg *diam.Message) error {
	if msg == nil {
		return errors.New("Message is nil")
	}
	_, err := msg.FindAVP(avp.UserName, 0)
	if err != nil {
		return errors.New("Missing IMSI in message")
	}
	_, err = msg.FindAVP(avp.SIPNumberAuthItems, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing SIP-Number-Auth-Items in message")
	}
	_, err = msg.FindAVP(avp.SIPAuthDataItem, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing SIP-Auth-Data-Item in message")
	}
	_, err = msg.FindAVP(avp.SIPAuthenticationScheme, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing SIP-Authentication-Scheme in message")
	}
	_, err = msg.FindAVP(avp.RATType, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing RAT type in message")
	}
	return nil
}

// set3GPPAAAServerName sets the 3GPP AAA Server stored inside of a SubscriberData proto.
func (srv *HomeSubscriberServer) set3GPPAAAServerName(subscriber *lteprotos.SubscriberData, serverName datatype.DiameterIdentity) error {
	if subscriber.State == nil {
		subscriber.State = &lteprotos.SubscriberState{}
	}
	subscriber.State.TgppAaaServerName = string(serverName)
	subscriber.State.TgppAaaServerRegistered = false
	return srv.store.UpdateSubscriber(subscriber)
}
