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
	"testing"

	fegprotos "magma/feg/cloud/go/protos"
	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/diameter"
	definitions "magma/feg/gateway/services/s6a_proxy/servicers"
	hss "magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/servicers/test_utils"
	"magma/feg/gateway/services/testcore/hss/storage"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/stretchr/testify/assert"
)

func TestNewULA_MissingMandatoryAVP(t *testing.T) {
	server := test_utils.NewTestHomeSubscriberServer(t)
	m := diam.NewMessage(diam.UpdateLocation, diam.RequestFlag|diam.ProxiableFlag, diam.TGPP_S6A_APP_ID, 1, 1, dict.Default)
	response, err := hss.NewULA(server, m)
	assert.EqualError(t, err, "Missing IMSI in message")

	// Check that the ULA is a failure message.
	var ula definitions.ULA
	err = response.Unmarshal(&ula)
	assert.NoError(t, err)
	assert.Equal(t, diam.MissingAVP, int(ula.ResultCode))
}

func TestNewULA_UnknownSubscriber(t *testing.T) {
	ulr := createULR("sub_unknown")
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewULA(server, ulr)
	assert.Exactly(t, err, storage.NewUnknownSubscriberError("sub_unknown"))

	// Check that the ULA is a failure message.
	var ula definitions.ULA
	err = response.Unmarshal(&ula)
	assert.NoError(t, err)
	assert.Equal(t, uint32(fegprotos.ErrorCode_USER_UNKNOWN), ula.ExperimentalResult.ExperimentalResultCode)
}

func TestNewULA_SuccessfulResponse(t *testing.T) {
	ulr := createULR("sub1")
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewULA(server, ulr)
	assert.NoError(t, err)

	// Check that the ULA has all the expected data.
	var ula definitions.ULA
	err = response.Unmarshal(&ula)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", ula.SessionID)
	assert.Equal(t, diam.Success, int(ula.ResultCode))
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), ula.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), ula.OriginRealm)
}

func TestNewULA_NewSuccessfulULA(t *testing.T) {
	ulr := createULR("sub1")
	server := test_utils.NewTestHomeSubscriberServer(t)
	profile := &mconfig.HSSConfig_SubscriptionProfile{
		MaxUlBitRate: 123,
		MaxDlBitRate: 456,
	}
	response := server.NewSuccessfulULA(ulr, datatype.UTF8String("magma;123_1234"), profile)

	var ula definitions.ULA
	err := response.Unmarshal(&ula)
	assert.NoError(t, err)
	assert.Equal(t, diam.Success, int(ula.ResultCode))
	assert.Equal(t, "magma;123_1234", ula.SessionID)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), ula.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), ula.OriginRealm)
	assert.Equal(t, uint32(0), ula.ULAFlags)

	assert.Equal(t, uint32(0), ula.SubscriptionData.APNConfigurationProfile.ContextIdentifier)
	assert.Equal(t, int32(0), ula.SubscriptionData.APNConfigurationProfile.AllAPNConfigurationsIncludedIndicator)
	assert.Equal(t, uint32(123), ula.SubscriptionData.AMBR.MaxRequestedBandwidthUL)
	assert.Equal(t, uint32(456), ula.SubscriptionData.AMBR.MaxRequestedBandwidthDL)
	assert.Equal(t, int32(2), ula.SubscriptionData.NetworkAccessMode)
	assert.Equal(t, datatype.OctetString("12345"), ula.SubscriptionData.MSISDN)
	assert.Equal(t, uint32(47), ula.SubscriptionData.AccessRestrictionData)
	assert.Equal(t, int32(0), ula.SubscriptionData.SubscriberStatus)

	assert.Equal(t, 1, len(ula.SubscriptionData.APNConfigurationProfile.APNConfigs))
	config := ula.SubscriptionData.APNConfigurationProfile.APNConfigs[0]
	assert.Equal(t, uint32(123), config.AMBR.MaxRequestedBandwidthUL)
	assert.Equal(t, uint32(456), config.AMBR.MaxRequestedBandwidthDL)
	assert.Equal(t, uint32(0), config.ContextIdentifier)
	assert.Equal(t, uint32(fegprotos.UpdateLocationAnswer_APNConfiguration_IPV4), config.PDNType)
	assert.Equal(t, "oai.ipv4", config.ServiceSelection)

	eps := config.EPSSubscribedQoSProfile
	assert.Equal(t, int32(9), eps.QoSClassIdentifier)
	assert.Equal(t, int32(1), eps.AllocationRetentionPriority.PreemptionCapability)
	assert.Equal(t, int32(0), eps.AllocationRetentionPriority.PreemptionVulnerability)
	assert.Equal(t, uint32(15), eps.AllocationRetentionPriority.PriorityLevel)
}

func TestValidateULR_MissingUserName(t *testing.T) {
	ulr := createBaseULR()
	ulr.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	ulr.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	ulr.NewAVP(avp.ULRFlags, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	ulr.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))

	err := hss.ValidateULR(ulr)
	assert.EqualError(t, err, "Missing IMSI in message")
}

func TestValidateULR_MissingVisitedPLMNID(t *testing.T) {
	ulr := createBaseULR()
	ulr.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	ulr.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("sub1"))
	ulr.NewAVP(avp.ULRFlags, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	ulr.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))

	err := hss.ValidateULR(ulr)
	assert.EqualError(t, err, "Missing Visited PLMN ID in message")
}

func TestValidateULR_MissingULRFlags(t *testing.T) {
	ulr := createBaseULR()
	ulr.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	ulr.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("sub1"))
	ulr.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	ulr.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))

	err := hss.ValidateULR(ulr)
	assert.EqualError(t, err, "Missing ULR flags in message")
}

func TestValidateULR_MissingRATType(t *testing.T) {
	ulr := createBaseULR()
	ulr.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	ulr.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("sub1"))
	ulr.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	ulr.NewAVP(avp.ULRFlags, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))

	err := hss.ValidateULR(ulr)
	assert.EqualError(t, err, "Missing RAT type in message")
}

func TestValidateULR_MissingSessionID(t *testing.T) {
	ulr := createBaseULR()
	ulr.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("sub1"))
	ulr.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	ulr.NewAVP(avp.ULRFlags, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	ulr.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))

	err := hss.ValidateULR(ulr)
	assert.EqualError(t, err, "Missing SessionID in message")
}

func TestNewULA_RATTypeNotAllowed(t *testing.T) {
	ulr := createULRExtended("sub1", 10)
	server := test_utils.NewTestHomeSubscriberServer(t)
	response, err := hss.NewULA(server, ulr)
	assert.EqualError(t, err, "RAT-Type not allowed: 10")

	// Check that the ULA has all the expected data.
	var ula definitions.ULA
	err = response.Unmarshal(&ula)
	assert.NoError(t, err)
	assert.Equal(t, "magma;123_1234", ula.SessionID)
	assert.Equal(t, uint32(fegprotos.ErrorCode_RAT_NOT_ALLOWED), ula.ExperimentalResult.ExperimentalResultCode)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), ula.OriginHost)
	assert.Equal(t, datatype.DiameterIdentity("magma.com"), ula.OriginRealm)
}

func createBaseULR() *diam.Message {
	ulr := diameter.NewProxiableRequest(diam.UpdateLocation, diam.TGPP_S6A_APP_ID, dict.Default)
	ulr.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	ulr.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("magma.com"))
	return ulr
}

func createULR(userName string) *diam.Message {
	return createULRExtended(userName, definitions.RadioAccessTechnologyType_EUTRAN)
}

func createULRExtended(userName string, ratType uint32) *diam.Message {
	ulr := createBaseULR()
	ulr.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("magma;123_1234"))
	ulr.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(userName))
	ulr.NewAVP(avp.VisitedPLMNID, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	ulr.NewAVP(avp.ULRFlags, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(0))
	ulr.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(ratType))
	return ulr
}
