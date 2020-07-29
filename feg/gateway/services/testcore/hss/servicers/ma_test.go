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

package servicers_test

import (
	"context"
	"testing"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	definitions "magma/feg/gateway/services/swx_proxy/servicers"
	hss "magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/servicers/test_utils"
	"magma/feg/gateway/services/testcore/hss/storage"
	"magma/lte/cloud/go/crypto"
	lteprotos "magma/lte/cloud/go/protos"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/stretchr/testify/assert"
)

func TestNewMAA_SuccessfulResponse(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)
	testNewMAASuccessfulResponse(t, server)
}

func TestNewMAA_UnknownIMSI(t *testing.T) {
	mar := createMARWithSingleAuthItem("sub_unknown")
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewMAA(server, mar)
	assert.Exactly(t, storage.NewUnknownSubscriberError("sub_unknown"), err)

	// Check that the MAA is a failure message.
	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	assert.Equal(t, uint32(fegprotos.ErrorCode_USER_UNKNOWN), maa.ExperimentalResult.ExperimentalResultCode)
}

func TestNewMAA_MissingAuthKey(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)

	mar := createMARWithSingleAuthItem("missing_auth_key")
	response, err := hss.NewMAA(server, mar)
	assert.Exactly(t, hss.NewAuthRejectedError("incorrect key size. Expected 16 bytes, but got 0 bytes"), err)

	// Check that the MAA has the expected error.
	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	assert.Equal(t, uint32(fegprotos.ErrorCode_AUTHORIZATION_REJECTED), maa.ExperimentalResult.ExperimentalResultCode)
	checkSIPAuthVectors(t, maa, 0)
	assert.Equal(t, "magma;123_1234", maa.SessionID)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginRealm)
}

func TestNewMAA_Redirect(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)
	set3GPPAAAServerName(t, server, "sub1", "different_server")

	mar := createMARWithSingleAuthItem("sub1")
	response, err := hss.NewMAA(server, mar)
	assert.EqualError(t, err, "diameter identity for AAA server already registered")

	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", maa.SessionID)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginRealm)
	assert.Equal(t, uint32(diam.RedirectIndication), maa.ResultCode)
	assert.Equal(t, uint32(fegprotos.SwxErrorCode_IDENTITY_ALREADY_REGISTERED), maa.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, datatype.DiameterIdentity("different_server"), maa.AAAServerName)
}

func TestNewMAA_StoreAAAServerName(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)
	set3GPPAAAServerName(t, server, "sub1", "")
	testNewMAASuccessfulResponse(t, server)

	subscriber, err := server.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)
	assert.Equal(t, "magma.com", subscriber.State.TgppAaaServerName)
}

func TestNewMAA_MultipleVectors(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)
	mar := createMARExtended("sub1", 3, definitions.RadioAccessTechnologyType_WLAN)
	response, err := hss.NewMAA(server, mar)
	assert.NoError(t, err)

	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	checkSIPAuthVectors(t, maa, 3)
}

func TestNewMAA_MissingAVP(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(definitions.RadioAccessTechnologyType_WLAN))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})

	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewMAA(server, mar)
	assert.EqualError(t, err, "Missing IMSI in message")

	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	assert.Equal(t, uint32(diam.MissingAVP), maa.ResultCode)
}

func TestNewMAA_RATTypeNotAllowed(t *testing.T) {
	mar := createMARExtended("sub1", 1, 20)
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewMAA(server, mar)
	assert.EqualError(t, err, "RAT-Type not allowed: 20")

	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", maa.SessionID)
	assert.Equal(t, uint32(fegprotos.ErrorCode_RAT_NOT_ALLOWED), maa.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginRealm)
	checkSIPAuthVectors(t, maa, 0)
}

func TestValidateMAR_Success(t *testing.T) {
	mar := createMARWithSingleAuthItem("sub1")
	err := hss.ValidateMAR(mar)
	assert.NoError(t, err)
}

func TestValidateMAR_MissingUserName(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(definitions.RadioAccessTechnologyType_WLAN))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})

	err := hss.ValidateMAR(mar)
	assert.EqualError(t, err, "Missing IMSI in message")
}

func TestValidateMAR_MissingSIPNumberAuthItems(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("sub1"))
	mar.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(definitions.RadioAccessTechnologyType_WLAN))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})

	err := hss.ValidateMAR(mar)
	assert.EqualError(t, err, "Missing SIP-Number-Auth-Items in message")
}

func TestValidateMAR_MissingSIPAuthDataItem(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("sub1"))
	mar.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(definitions.RadioAccessTechnologyType_WLAN))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA"))

	err := hss.ValidateMAR(mar)
	assert.EqualError(t, err, "Missing SIP-Auth-Data-Item in message")
}

func TestValidateMAR_MissingSIPAuthenticationScheme(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("sub1"))
	mar.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(definitions.RadioAccessTechnologyType_WLAN))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{})

	err := hss.ValidateMAR(mar)
	assert.EqualError(t, err, "Missing SIP-Authentication-Scheme in message")
}

func TestValidateMAR_MissingRATType(t *testing.T) {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("sub1"))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})

	err := hss.ValidateMAR(mar)
	assert.EqualError(t, err, "Missing RAT type in message")
}

func TestValidateMAR_NilMessage(t *testing.T) {
	err := hss.ValidateMAR(nil)
	assert.EqualError(t, err, "Message is nil")
}

func createBaseMAR() *diam.Message {
	mar := diameter.NewProxiableRequest(diam.MultimediaAuthentication, diam.TGPP_SWX_APP_ID, dict.Default)
	mar.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	mar.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	mar.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	return mar
}

func createMARWithSingleAuthItem(userName string) *diam.Message {
	return createMARExtended(userName, 1, definitions.RadioAccessTechnologyType_WLAN)
}

func createMARExtended(userName string, numberAuthItems uint32, ratType uint32) *diam.Message {
	mar := createBaseMAR()
	mar.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(userName))
	mar.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(ratType))
	mar.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(numberAuthItems))
	mar.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String("EAP-AKA")),
		},
	})
	return mar
}

func checkSIPAuthVectors(t *testing.T, maa definitions.MAA, expectedNumVectors uint32) {
	assert.Equal(t, int(expectedNumVectors), len(maa.SIPAuthDataItems))
	assert.Equal(t, expectedNumVectors, maa.SIPNumberAuthItems)

	for _, vector := range maa.SIPAuthDataItems {
		assert.Equal(t, definitions.SipAuthScheme_EAP_AKA, vector.AuthScheme)
		assert.Equal(t, crypto.RandChallengeBytes+crypto.AutnBytes, len(vector.Authenticate))
		assert.Equal(t, crypto.XresBytes, len(vector.Authorization))
		assert.Equal(t, crypto.ConfidentialityKeyBytes, len(vector.ConfidentialityKey))
		assert.Equal(t, crypto.IntegrityKeyBytes, len(vector.IntegrityKey))
	}
}

func set3GPPAAAServerName(t *testing.T, server *hss.HomeSubscriberServer, imsi string, serverName string) {
	id := &lteprotos.SubscriberID{Id: imsi}
	subscriber, err := server.GetSubscriberData(context.Background(), id)
	assert.NoError(t, err)
	subscriber.State.TgppAaaServerName = serverName
	_, err = server.UpdateSubscriber(context.Background(), subscriber)
	assert.NoError(t, err)
}

func testNewMAASuccessfulResponse(t *testing.T, server *hss.HomeSubscriberServer) {
	mar := createMARWithSingleAuthItem("sub1")
	response, err := hss.NewMAA(server, mar)
	assert.NoError(t, err)

	var maa definitions.MAA
	err = response.Unmarshal(&maa)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", maa.SessionID)
	assert.Equal(t, diam.Success, int(maa.ResultCode))
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), maa.OriginRealm)
	assert.Equal(t, int32(definitions.AuthSessionState_NO_STATE_MAINTAINED), maa.AuthSessionState)
	checkSIPAuthVectors(t, maa, 1)

	subscriber, err := server.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)
	assert.Equal(t, uint64(7351), subscriber.State.LteAuthNextSeq)
}
