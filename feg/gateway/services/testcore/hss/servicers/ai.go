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

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	s6a "magma/feg/gateway/services/s6a_proxy/servicers"
	"magma/feg/gateway/services/testcore/hss/storage"
	lteprotos "magma/lte/cloud/go/protos"
)

// NewAIA outputs a authentication information answer (AIA) to reply to an
// authentication information request (AIR) message.
func NewAIA(srv *HomeSubscriberServer, msg *diam.Message) (*diam.Message, error) {
	if err := ValidateAIR(msg); err != nil {
		return msg.Answer(diam.MissingAVP), err
	}

	var air s6a.AIR
	if err := msg.Unmarshal(&air); err != nil {
		return msg.Answer(diam.UnableToComply), fmt.Errorf("AIR Unmarshal failed for message: %v failed: %v", msg, err)
	}

	subscriber, err := srv.store.GetSubscriberData(air.UserName)
	if err != nil {
		if _, ok := err.(storage.UnknownSubscriberError); ok {
			return ConstructFailureAnswer(msg, air.SessionID, srv.Config.Server, uint32(fegprotos.ErrorCode_USER_UNKNOWN)), err
		}
		return ConstructFailureAnswer(msg, air.SessionID, srv.Config.Server, uint32(fegprotos.ErrorCode_AUTHENTICATION_DATA_UNAVAILABLE)), err
	}

	subscriber.Lock()
	defer subscriber.Unlock()

	lteAuthNextSeq, err := ResyncLteAuthSeq(
		subscriber, air.RequestedEUTRANAuthInfo.ResyncInfo.Serialize(), srv.Config.LteAuthOp)
	if err != nil {
		return ConvertAuthErrorToFailureMessage(err, msg, air.SessionID, srv.Config.Server), err
	}
	if len(air.RequestedUtranGeranAuthInfo.ResyncInfo) > 0 {
		lteAuthNextUtranSeq, err := ResyncLteAuthSeq(
			subscriber, air.RequestedUtranGeranAuthInfo.ResyncInfo.Serialize(), srv.Config.LteAuthOp)
		if err != nil {
			return ConvertAuthErrorToFailureMessage(err, msg, air.SessionID, srv.Config.Server), err
		}
		if len(air.RequestedEUTRANAuthInfo.ResyncInfo) == 0 || lteAuthNextUtranSeq > lteAuthNextSeq {
			lteAuthNextSeq = lteAuthNextUtranSeq
		}
	}
	err = srv.setLteAuthNextSeq(subscriber, lteAuthNextSeq)
	if err != nil {
		return ConvertAuthErrorToFailureMessage(err, msg, air.SessionID, srv.Config.Server), err
	}

	const plmnOffsetBytes = 1
	plmn := air.VisitedPLMNID.Serialize()[plmnOffsetBytes:]

	vectors, utranVectors, lteAuthNextSeq, err := GenerateLteAuthVectors(
		uint32(air.RequestedEUTRANAuthInfo.NumVectors),
		uint32(air.RequestedUtranGeranAuthInfo.NumVectors),
		srv.Milenage, subscriber, plmn, srv.Config.LteAuthOp, srv.AuthSqnInd)
	if err == nil {
		err = srv.setLteAuthNextSeq(subscriber, lteAuthNextSeq)
	}
	if err != nil {
		return ConvertAuthErrorToFailureMessage(err, msg, air.SessionID, srv.Config.Server), err
	}

	return srv.NewSuccessfulAIA(msg, air.SessionID, vectors, utranVectors), nil
}

func (srv *HomeSubscriberServer) setLteAuthNextSeq(subscriber *lteprotos.SubscriberData, lteAuthNextSeq uint64) error {
	if subscriber.GetState() == nil {
		return NewAuthDataUnavailableError("subscriber state was nil")
	}
	subscriber.State.LteAuthNextSeq = lteAuthNextSeq
	return srv.store.UpdateSubscriber(subscriber)
}

// NewSuccessfulAIA outputs a successful authentication information answer (AIA) to reply to an
// authentication information request (AIR) message. It populates AIA with all of the mandatory fields
// and adds the authentication vectors.
func (srv *HomeSubscriberServer) NewSuccessfulAIA(
	msg *diam.Message,
	sessionID datatype.UTF8String,
	vectors []*milenage.EutranVector,
	utranVectors []*milenage.UtranVector) *diam.Message {

	answer := ConstructSuccessAnswer(msg, sessionID, srv.Config.Server, diam.TGPP_S6A_APP_ID)
	evs := []*diam.AVP{}
	for itemNumber, vector := range vectors {
		evs = append(evs, diam.NewAVP(avp.EUTRANVector, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.ItemNumber, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(itemNumber)),
				diam.NewAVP(avp.RAND, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Rand[:])),
				diam.NewAVP(avp.XRES, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Xres[:])),
				diam.NewAVP(avp.AUTN, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Autn[:])),
				diam.NewAVP(avp.KASME, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Kasme[:])),
			},
		}))
	}
	for itemNumber, vector := range utranVectors {
		evs = append(evs, diam.NewAVP(avp.UTRANVector, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.ItemNumber, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(itemNumber)),
				diam.NewAVP(avp.RAND, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Rand[:])),
				diam.NewAVP(avp.XRES, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Xres[:])),
				diam.NewAVP(avp.AUTN, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(vector.Autn[:])),
				diam.NewAVP(
					avp.ConfidentialityKey,
					avp.Mbit|avp.Vbit,
					diameter.Vendor3GPP,
					datatype.OctetString(vector.ConfidentialityKey[:])),
				diam.NewAVP(
					avp.IntegrityKey,
					avp.Mbit|avp.Vbit,
					diameter.Vendor3GPP,
					datatype.OctetString(vector.IntegrityKey[:])),
			},
		}))
	}
	answer.NewAVP(avp.AuthenticationInfo, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{AVP: evs})
	return answer
}

// ValidateAIR returns an error if the message is missing any mandatory AVPs.
// Mandatory AVPs are specified in 3GPP TS 29.272 Table 5.2.3.1.1/1
func ValidateAIR(msg *diam.Message) error {
	_, err := msg.FindAVP(avp.UserName, 0)
	if err != nil {
		return errors.New("Missing IMSI in message")
	}
	_, err = msg.FindAVP(avp.VisitedPLMNID, diameter.Vendor3GPP)
	if err != nil {
		return errors.New("Missing Visited PLMN ID in message")
	}
	_, err = msg.FindAVP(avp.RequestedEUTRANAuthenticationInfo, diameter.Vendor3GPP)
	if err != nil {
		_, err = msg.FindAVP(avp.RequestedUTRANGERANAuthenticationInfo, diameter.Vendor3GPP)
		if err != nil {
			return errors.New("Missing requested E-UTRAN and UTRAN/GERAN authentication info in message")
		}
	}
	_, err = msg.FindAVP(avp.SessionID, 0)
	if err != nil {
		return errors.New("Missing SessionID in message")
	}
	return nil
}
