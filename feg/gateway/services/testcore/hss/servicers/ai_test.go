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

	"github.com/emakeev/milenage"
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/stretchr/testify/assert"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	definitions "magma/feg/gateway/services/s6a_proxy/servicers"
	hss "magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/servicers/test_utils"
	"magma/feg/gateway/services/testcore/hss/storage"
	lteprotos "magma/lte/cloud/go/protos"
)

func TestNewAIA_MissingSessionID(t *testing.T) {
	m := diameter.NewProxiableRequest(diam.AuthenticationInformation, diam.TGPP_S6A_APP_ID, dict.Default)
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewAIA(server, m)
	assert.Error(t, err)

	// Check that the AIA is a failure message.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, diam.MissingAVP, int(aia.ResultCode))
}

func TestNewAIA_UnknownIMSI(t *testing.T) {
	air := createAIR("sub_unknown")
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewAIA(server, air)
	assert.Exactly(t, storage.NewUnknownSubscriberError("sub_unknown"), err)

	// Check that the AIA is a failure message.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, uint32(protos.ErrorCode_USER_UNKNOWN), aia.ExperimentalResult.ExperimentalResultCode)
}

func TestNewAIA_SuccessfulResponse(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)
	amf := []byte("\x80\x00")
	rand := []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f")
	milenage, err := milenage.NewMockCipher(amf, rand)
	assert.NoError(t, err)
	server.Milenage = milenage

	air := createAIR("sub1")
	response, err := hss.NewAIA(server, air)
	assert.NoError(t, err)

	// Check that the AIA has all the expected data.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", aia.SessionID)
	assert.Equal(t, diam.Success, int(aia.ResultCode))
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginRealm)

	ai := &aia.AI
	assert.Equal(t, 1, len(ai.EUtranVectors))

	vec := ai.EUtranVectors[0]
	assert.Equal(t, datatype.OctetString(rand), vec.RAND)
	assert.Equal(t, datatype.OctetString([]byte("\x2d\xaf\x87\x3d\x73\xf3\x10\xc6")), vec.XRES)
	assert.Equal(t, datatype.OctetString([]byte("o\xbf\xa3\x83\x95 \x80\x00\xb1\x1f \xbd\xdc\xf5\xeeS")), vec.AUTN)
	assert.Equal(t, datatype.OctetString([]byte("Q\xd0g\xde?\x95\xecB\x94\xf8\xe7\xc4\x0f\x92\x81i\x8e\\Cu\xc1\xe5\xab\x1a\xc0\xe6z\x117\nkz")), vec.KASME)

	subscriber, err := server.GetSubscriberData(context.Background(), &lteprotos.SubscriberID{Id: "sub1"})
	assert.NoError(t, err)
	assert.Equal(t, uint64(7351), subscriber.State.LteAuthNextSeq)
}

func TestNewAIA_MultipleVectors(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)
	air := createAIRExtended("sub1", 3, 0)
	response, err := hss.NewAIA(server, air)
	assert.NoError(t, err)

	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(aia.AI.EUtranVectors))

	vector := aia.AI.EUtranVectors[0]
	assert.Equal(t, milenage.RandChallengeBytes, len(vector.RAND))
	assert.Equal(t, milenage.XresBytes, len(vector.XRES))
	assert.Equal(t, milenage.AutnBytes, len(vector.AUTN))
	assert.Equal(t, milenage.KasmeBytes, len(vector.KASME))

	for i := 0; i < len(aia.AI.EUtranVectors); i++ {
		for j := i + 1; j < len(aia.AI.EUtranVectors); j++ {
			assert.NotEqual(t, aia.AI.EUtranVectors[i], aia.AI.EUtranVectors[j])
		}
	}
}

func TestNewAIA_MissingAuthKey(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)

	air := createAIR("missing_auth_key")
	response, err := hss.NewAIA(server, air)
	assert.Exactly(t, hss.NewAuthRejectedError("incorrect key size. Expected 16 bytes, but got 0 bytes"), err)

	// Check that the AIA has the expected error.
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)
	assert.Equal(t, uint32(protos.ErrorCode_AUTHORIZATION_REJECTED), aia.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, 0, len(aia.AI.EUtranVectors))
	assert.Equal(t, "magma;123_1234", aia.SessionID)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), aia.OriginRealm)
}

func TestValidateAIR_MissingUserName(t *testing.T) {
	m := createBaseAIR()
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	m.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				diameter.Vendor3GPP,
				datatype.Unsigned32(1)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0)),
		},
	}
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, authInfo)

	assert.EqualError(t, hss.ValidateAIR(m), "Missing IMSI in message")
}

func TestValidateAIR_MissingVistedPLMNID(t *testing.T) {
	m := createBaseAIR()
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("magma"))
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				diameter.Vendor3GPP,
				datatype.Unsigned32(1)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0)),
		},
	}
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, authInfo)

	assert.EqualError(t, hss.ValidateAIR(m), "Missing Visited PLMN ID in message")
}

func TestValidateAIR_MissingEUTRANInfo(t *testing.T) {
	m := createBaseAIR()
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("magma"))
	m.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))

	assert.EqualError(t, hss.ValidateAIR(m), "Missing requested E-UTRAN and UTRAN/GERAN authentication info in message")
}

func TestValidateAIR_MissingSessionId(t *testing.T) {
	m := createBaseAIR()
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("magma"))
	m.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				diameter.Vendor3GPP,
				datatype.Unsigned32(1)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0)),
		},
	}
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, authInfo)

	assert.EqualError(t, hss.ValidateAIR(m), "Missing SessionID in message")
}

func TestValidateAIR_Success(t *testing.T) {
	air := createAIR("sub1")
	assert.NoError(t, hss.ValidateAIR(air))
	air = createAIRExtended("sub2", 0, 1)
	assert.NoError(t, hss.ValidateAIR(air))
	air = createAIRExtended("sub2", 1, 1)
	assert.NoError(t, hss.ValidateAIR(air))
}

// createBaseAIR outputs a mock authentication information request with only a
// few AVPs added.
func createBaseAIR() *diam.Message {
	air := diameter.NewProxiableRequest(diam.AuthenticationInformation, diam.TGPP_S6A_APP_ID, dict.Default)
	air.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	air.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	air.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	return air
}

// createAIR outputs a mock authentication information request.
func createAIR(userName string) *diam.Message {
	return createAIRExtended(userName, 1, 0)
}

// createAIRExtended outputs a mock authentication information request.
// It allows specifying more options than createAIR.
func createAIRExtended(userName string, numRequestedVectors, numRequestedUtranVectors uint32) *diam.Message {
	m := createBaseAIR()
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(userName))
	m.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	authInfo := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(
				avp.NumberOfRequestedVectors,
				avp.Vbit|avp.Mbit,
				diameter.Vendor3GPP,
				datatype.Unsigned32(numRequestedVectors)),
			diam.NewAVP(
				avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0)),
		},
	}
	m.NewAVP(avp.RequestedEUTRANAuthenticationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, authInfo)
	if numRequestedUtranVectors > 0 {
		authInfo := &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(
					avp.NumberOfRequestedVectors,
					avp.Vbit|avp.Mbit,
					diameter.Vendor3GPP,
					datatype.Unsigned32(numRequestedUtranVectors)),
				diam.NewAVP(
					avp.ImmediateResponsePreferred, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(0)),
			},
		}
		m.NewAVP(avp.RequestedUTRANGERANAuthenticationInfo, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, authInfo)
	}
	return m
}

func TestNewSuccessfulAIA(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)
	serverCfg := server.Config.Server

	msg := createAIRExtended("user1", 1, 1)
	var air definitions.AIR
	err := msg.Unmarshal(&air)
	assert.NoError(t, err)

	vector := &milenage.EutranVector{}
	copy(vector.Rand[:], []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f"))
	copy(vector.Xres[:], []byte("\x2d\xaf\x87\x3d\x73\xf3\x10\xc6"))
	copy(vector.Autn[:], []byte("o\xbf\xa3\x80\x1fW\x80\x00{\xdeY\x88n\x96\xe4\xfe"))
	copy(vector.Kasme[:],
		[]byte("\x87H\xc1\xc0\xa2\x82o\xa4\x05\xb1\xe2~\xa1\x04CJ\xe5V\xc7e\xe8\xf0a\xeb\xdb\x8a\xe2\x86\xc4F\x16\xc2"))

	utranVector := &milenage.UtranVector{}
	copy(utranVector.Rand[:], []byte("\xb0\x02\x7e\x16\x95\x59\xc4\x8a\x2d\xc8\x3d\xb6\xb7\x6b\x77\xb5"))
	copy(utranVector.Xres[:], []byte("\xc8\x8e\x71\xfa\x5a\x1b\x41\x50"))
	copy(utranVector.Autn[:], []byte("\x56\x63\xa4\x78\x7f\x61\x80\x00\x8d\x45\x64\x77\x96\x06\x39\x22"))
	copy(utranVector.ConfidentialityKey[:], []byte("\x5e\x8d\xc7\x0e\x0d\xfa\xa5\xde\x81\x90\xe9\xc5\x98\x6a\xc2\x23"))
	copy(utranVector.IntegrityKey[:], []byte("\xf2\xc1\x9c\x15\xf5\x5a\xf1\xe8\xbb\xdb\x76\x19\x30\xeb\xef\x7c"))

	response := server.NewSuccessfulAIA(
		msg, air.SessionID, []*milenage.EutranVector{vector}, []*milenage.UtranVector{utranVector})
	var aia definitions.AIA
	err = response.Unmarshal(&aia)
	assert.NoError(t, err)

	assert.Equal(t, uint32(diam.Success), aia.ResultCode)
	assert.Equal(t, air.SessionID, datatype.UTF8String(aia.SessionID))
	assert.Equal(t, datatype.DiameterIdentity(serverCfg.DestHost), aia.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity(serverCfg.DestRealm), aia.OriginRealm)

	ai := &aia.AI
	assert.Equal(t, 1, len(ai.EUtranVectors))

	vec := ai.EUtranVectors[0]
	assert.Equal(t, datatype.OctetString(vector.Rand[:]), vec.RAND)
	assert.Equal(t, datatype.OctetString(vector.Xres[:]), vec.XRES)
	assert.Equal(t, datatype.OctetString(vector.Autn[:]), vec.AUTN)
	assert.Equal(t, datatype.OctetString(vector.Kasme[:]), vec.KASME)

	uvec := ai.UtranVectors[0]
	assert.Equal(t, datatype.OctetString(utranVector.Rand[:]), uvec.RAND)
	assert.Equal(t, datatype.OctetString(utranVector.Xres[:]), uvec.XRES)
	assert.Equal(t, datatype.OctetString(utranVector.Autn[:]), uvec.AUTN)
	assert.Equal(t, datatype.OctetString(utranVector.ConfidentialityKey[:]), uvec.CK)
	assert.Equal(t, datatype.OctetString(utranVector.IntegrityKey[:]), uvec.IK)
}
