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
	hss "magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/servicers/test_utils"
	lteprotos "magma/lte/cloud/go/protos"

	definitions "magma/feg/gateway/services/swx_proxy/servicers"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/stretchr/testify/assert"
)

func TestNewSAA_SuccessfulRegistration(t *testing.T) {
	sar := createSAR("sub1", definitions.ServerAssignmentType_REGISTRATION)
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewSAA(server, sar)
	assert.NoError(t, err)
	checkSAASuccess(t, response)

	subscriber, err := server.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)
	assert.True(t, subscriber.GetState().GetTgppAaaServerRegistered())
}

func TestNewSAA_SuccessfulDeregistration(t *testing.T) {
	sar := createSAR("sub1", definitions.ServerAssignnmentType_USER_DEREGISTRATION)
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewSAA(server, sar)
	assert.NoError(t, err)
	checkSAASuccessDeregistration(t, response)

	subscriber, err := server.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)
	assert.False(t, subscriber.GetState().GetTgppAaaServerRegistered())
	assert.Empty(t, subscriber.GetState().TgppAaaServerName)
}

func TestNewSAA_SuccessfulUserDataRequest(t *testing.T) {
	sar := createSAR("sub1", definitions.ServerAssignmentType_AAA_USER_DATA_REQUEST)
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewSAA(server, sar)
	assert.NoError(t, err)
	checkSAASuccess(t, response)

	subscriber, err := server.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)
	assert.False(t, subscriber.GetState().GetTgppAaaServerRegistered())
}

func TestNewSAA_UnknownIMSI(t *testing.T) {
	sar := createSAR("sub_unknown", definitions.ServerAssignmentType_REGISTRATION)
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewSAA(server, sar)
	assert.EqualError(t, err, "Subscriber 'sub_unknown' not found")

	saa := testUnmarshalSAA(t, response)
	assert.Equal(t, uint32(fegprotos.ErrorCode_USER_UNKNOWN), saa.ExperimentalResult.ExperimentalResultCode)
}

func TestNewSAA_No3GPPAAAServer(t *testing.T) {
	sar := createSAR("empty_sub", definitions.ServerAssignmentType_REGISTRATION)
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewSAA(server, sar)
	assert.EqualError(t, err, "no 3GPP AAA server is already serving the user")

	saa := testUnmarshalSAA(t, response)
	assert.Equal(t, uint32(diam.UnableToComply), saa.ExperimentalResult.ExperimentalResultCode)
}

func TestNewSAA_Redirect(t *testing.T) {
	sar := createSARExtended("sub1", definitions.ServerAssignmentType_REGISTRATION, "different_host")
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewSAA(server, sar)
	assert.EqualError(t, err, "diameter identity for AAA server already registered")

	saa := testUnmarshalSAA(t, response)
	assert.Equal(t, uint32(diam.RedirectIndication), saa.ResultCode)
	assert.Equal(t, uint32(fegprotos.SwxErrorCode_IDENTITY_ALREADY_REGISTERED), saa.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), saa.AAAServerName)
}

func TestNewSAA_MissingAPNConfig(t *testing.T) {
	sar := createSAR("missing_auth_key", definitions.ServerAssignmentType_REGISTRATION)
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewSAA(server, sar)
	assert.EqualError(t, err, "User has no non 3GPP subscription")

	saa := testUnmarshalSAA(t, response)
	assert.Equal(t, uint32(fegprotos.SwxErrorCode_USER_NO_NON_3GPP_SUBSCRIPTION), saa.ExperimentalResult.ExperimentalResultCode)
}

func TestNewSAA_MissingAVP(t *testing.T) {
	sar := createBaseSAR()
	sar.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("sub1"))
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewSAA(server, sar)
	assert.EqualError(t, err, "Missing server assignment type in message")

	var saa definitions.SAA
	err = response.Unmarshal(&saa)
	assert.NoError(t, err)
	assert.Equal(t, diam.MissingAVP, int(saa.ResultCode))
}

func TestValidateSAR_MissingUserName(t *testing.T) {
	sar := createBaseSAR()
	sar.NewAVP(avp.ServerAssignmentType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(definitions.ServerAssignmentType_REGISTRATION))
	err := hss.ValidateSAR(sar)
	assert.EqualError(t, err, "Missing IMSI in message")
}

func TestValidateSAR_MissingServerAssignmentType(t *testing.T) {
	sar := createBaseSAR()
	sar.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("sub1"))
	err := hss.ValidateSAR(sar)
	assert.EqualError(t, err, "Missing server assignment type in message")
}

func TestValidateSAR_NilMessage(t *testing.T) {
	err := hss.ValidateSAR(nil)
	assert.EqualError(t, err, "Message is nil")
}

func TestValidateSAR_Success(t *testing.T) {
	sar := createSAR("sub1", definitions.ServerAssignmentType_REGISTRATION)
	err := hss.ValidateSAR(sar)
	assert.NoError(t, err)
}

func createBaseSAR() *diam.Message {
	sar := diameter.NewProxiableRequest(diam.ServerAssignment, diam.TGPP_SWX_APP_ID, dict.Default)
	sar.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	sar.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	sar.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	return sar
}

func createSAR(userName string, serverAssignmentType int) *diam.Message {
	return createSARExtended(userName, serverAssignmentType, "magma.com")
}

func createSARExtended(userName string, serverAssignmentType int, originHost string) *diam.Message {
	sar := diameter.NewProxiableRequest(diam.ServerAssignment, diam.TGPP_SWX_APP_ID, dict.Default)
	sar.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	sar.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(originHost))
	sar.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	sar.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(userName))
	sar.NewAVP(avp.ServerAssignmentType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(serverAssignmentType))
	return sar
}

// checkSAASuccess ensures that a successful SAA contains all the expected data
func checkSAASuccess(t *testing.T, response *diam.Message) {
	saa := testUnmarshalSAA(t, response)
	assert.Equal(t, diam.Success, int(saa.ResultCode))
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), saa.AAAServerName)
	assert.Equal(t, datatype.UTF8String("sub1"), saa.UserName)
	assert.Equal(t, int32(definitions.AuthSessionState_NO_STATE_MAINTAINED), saa.AuthSessionState)

	profile := saa.UserData
	assert.Equal(t, datatype.Enumerated(lteprotos.Non3GPPUserProfile_NON_3GPP_SUBSCRIPTION_ALLOWED), profile.Non3GPPIPAccess)
	assert.Equal(t, datatype.Enumerated(definitions.END_USER_E164), profile.SubscriptionId.SubscriptionIdType)
	assert.Equal(t, datatype.UTF8String("12345"), profile.SubscriptionId.SubscriptionIdData)
}

// checkSAASuccessDeregistration ensures that a successful SAA with a deregistration command contains all the expected data
func checkSAASuccessDeregistration(t *testing.T, response *diam.Message) {
	saa := testUnmarshalSAA(t, response)
	assert.Equal(t, diam.Success, int(saa.ResultCode))
	assert.Equal(t, datatype.DiameterIdentity(""), saa.AAAServerName)
	assert.Equal(t, datatype.UTF8String("sub1"), saa.UserName)
	assert.Equal(t, int32(definitions.AuthSessionState_NO_STATE_MAINTAINED), saa.AuthSessionState)

}

// testUnmarshalSAA unmarshals an SAA message and checks that the SessionID,
// OriginHost, and OriginRealm fields are as expected.
func testUnmarshalSAA(t *testing.T, response *diam.Message) definitions.SAA {
	var saa definitions.SAA
	err := response.Unmarshal(&saa)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", saa.SessionID)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), saa.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), saa.OriginRealm)
	return saa
}
