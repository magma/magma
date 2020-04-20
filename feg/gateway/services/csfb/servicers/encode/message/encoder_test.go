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

package message_test

import (
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/encode/ie"
	"magma/feg/gateway/services/csfb/servicers/encode/message"

	"github.com/stretchr/testify/assert"
)

func TestEncodeSGsAPAlertAck(t *testing.T) {
	msg := &protos.AlertAck{
		Imsi: "001010000000001",
	}
	encodedMsg, err := message.EncodeSGsAPAlertAck(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPAlertAck), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	assert.Equal(t, encodedIMSI, encodedMsg[1:])
}

func TestEncodeSGsAPAlertReject(t *testing.T) {
	msg := &protos.AlertReject{
		Imsi:     "001010000000001",
		SgsCause: []byte{byte(0x11)},
	}
	encodedMsg, err := message.EncodeSGsAPAlertReject(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPAlertReject), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedSGsCause, _ := ie.EncodeFixedLengthIE(decode.IEISGsCause, decode.IELengthSGsCause, []byte{byte(0x11)})
	restOfFields := append(encodedIMSI, encodedSGsCause...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPAlertRequest(t *testing.T) {
	msg := &protos.AlertRequest{
		Imsi: "001010000000001",
	}
	encodedMsg, err := message.EncodeSGsAPAlertRequest(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPAlertRequest), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	assert.Equal(t, encodedIMSI, encodedMsg[1:])
}

func TestEncodeSGsAPDownlinkUnitdata(t *testing.T) {
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	msg := &protos.DownlinkUnitdata{
		Imsi:                "001010000000001",
		NasMessageContainer: make([]byte, decode.IELengthNASMessageContainerMax-mandatoryFieldLength),
	}
	encodedMsg, err := message.EncodeSGsAPDownlinkUnitdata(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPDownlinkUnitdata), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedNasMessageContainer, _ := ie.EncodeLimitedLengthIE(
		decode.IEINASMessageContainer,
		decode.IELengthNASMessageContainerMin,
		decode.IELengthNASMessageContainerMax,
		make([]byte, decode.IELengthNASMessageContainerMax-mandatoryFieldLength),
	)
	restOfFields := append(encodedIMSI, encodedNasMessageContainer...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPEPSDetachAck(t *testing.T) {
	msg := &protos.EPSDetachAck{
		Imsi: "001010000000001",
	}
	encodedMsg, err := message.EncodeSGsAPEPSDetachAck(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPEPSDetachAck), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	assert.Equal(t, encodedIMSI, encodedMsg[1:])
}

func TestEncodeSGsAPEPSDetachIndication(t *testing.T) {
	msg := &protos.EPSDetachIndication{
		Imsi:                         "001010000000001",
		MmeName:                      ".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
		ImsiDetachFromEpsServiceType: []byte{byte(0x11)},
	}
	encodedMsg, err := message.EncodeSGsAPEPSDetachIndication(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPEPSDetachIndication), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedMMEName, _ := ie.EncodeMMEName(
		".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
	)
	encodedServiceType, _ := ie.EncodeFixedLengthIE(
		decode.IEIIMSIDetachFromEPSServiceType,
		decode.IELengthIMSIDetachFromEPSServiceType,
		[]byte{byte(0x11)},
	)
	restOfFields := append(encodedIMSI, encodedMMEName...)
	restOfFields = append(restOfFields, encodedServiceType...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPIMSIDetachAck(t *testing.T) {
	msg := &protos.IMSIDetachAck{
		Imsi: "001010000000001",
	}
	encodedMsg, err := message.EncodeSGsAPIMSIDetachAck(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPIMSIDetachAck), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	assert.Equal(t, encodedIMSI, encodedMsg[1:])
}

func TestEncodeSGsAPIMSIDetachIndication(t *testing.T) {
	msg := &protos.IMSIDetachIndication{
		Imsi:                            "001010000000001",
		MmeName:                         ".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
		ImsiDetachFromNonEpsServiceType: []byte{byte(0x11)},
	}
	encodedMsg, err := message.EncodeSGsAPIMSIDetachIndication(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPIMSIDetachIndication), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedMMEName, _ := ie.EncodeMMEName(
		".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
	)
	encodedServiceType, _ := ie.EncodeFixedLengthIE(
		decode.IEIIMSIDetachFromNonEPSServiceType,
		decode.IELengthIMSIDetachFromNonEPSServiceType,
		[]byte{byte(0x11)},
	)
	restOfFields := append(encodedIMSI, encodedMMEName...)
	restOfFields = append(restOfFields, encodedServiceType...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPLocationUpdateAccept(t *testing.T) {
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator

	// without IMSI and TMSI
	msg := &protos.LocationUpdateAccept{
		Imsi:                   "001010000000001",
		LocationAreaIdentifier: make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
	}
	encodedMsg, err := message.EncodeSGsAPLocationUpdateAccept(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPLocationUpdateAccept), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedLAI, _ := ie.EncodeFixedLengthIE(
		decode.IEILocationAreaIdentifier,
		decode.IELengthLocationAreaIdentifier,
		make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
	)
	mandatoryFields := append(encodedIMSI, encodedLAI...)
	assert.Equal(t, mandatoryFields, encodedMsg[1:])

	// with IMSI
	msg = &protos.LocationUpdateAccept{
		Imsi:                   "001010000000001",
		LocationAreaIdentifier: make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
		NewIMSITMSI:            &protos.LocationUpdateAccept_NewImsi{NewImsi: "22222222"},
	}
	encodedMsg, err = message.EncodeSGsAPLocationUpdateAccept(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPLocationUpdateAccept), encodedMsg[0])
	encodedNewIMSI, _ := ie.EncodeIMSI("22222222")
	restOfFields := append(mandatoryFields, encodedNewIMSI...)
	assert.Equal(t, restOfFields, encodedMsg[1:])

	// with TMSI
	msg = &protos.LocationUpdateAccept{
		Imsi:                   "001010000000001",
		LocationAreaIdentifier: make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
		NewIMSITMSI: &protos.LocationUpdateAccept_NewTmsi{
			NewTmsi: make([]byte, decode.IELengthTMSI-mandatoryFieldLength),
		},
	}
	encodedMsg, err = message.EncodeSGsAPLocationUpdateAccept(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPLocationUpdateAccept), encodedMsg[0])
	encodedNewTMSI, _ := ie.EncodeFixedLengthIE(
		decode.IEITMSI,
		decode.IELengthTMSI,
		make([]byte, decode.IELengthTMSI-mandatoryFieldLength),
	)
	restOfFields = append(mandatoryFields, encodedNewTMSI...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPLocationUpdateReject(t *testing.T) {
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator

	// without LAI
	msg := &protos.LocationUpdateReject{
		Imsi:        "001010000000001",
		RejectCause: make([]byte, decode.IELengthRejectCause-mandatoryFieldLength),
	}
	encodedMsg, err := message.EncodeSGsAPLocationUpdateReject(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPLocationUpdateReject), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedRejectCause, _ := ie.EncodeFixedLengthIE(
		decode.IEIRejectCause,
		decode.IELengthRejectCause,
		make([]byte, decode.IELengthRejectCause-mandatoryFieldLength),
	)
	restOfFields := append(encodedIMSI, encodedRejectCause...)
	assert.Equal(t, restOfFields, encodedMsg[1:])

	// with LAI
	msg = &protos.LocationUpdateReject{
		Imsi:                   "001010000000001",
		RejectCause:            make([]byte, decode.IELengthRejectCause-mandatoryFieldLength),
		LocationAreaIdentifier: make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
	}
	encodedMsg, err = message.EncodeSGsAPLocationUpdateReject(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPLocationUpdateReject), encodedMsg[0])
	encodedLAI, _ := ie.EncodeFixedLengthIE(
		decode.IEILocationAreaIdentifier,
		decode.IELengthLocationAreaIdentifier,
		make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
	)
	restOfFields = append(restOfFields, encodedLAI...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPLocationUpdateRequest(t *testing.T) {
	// without optional fields
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	msg := &protos.LocationUpdateRequest{
		Imsi:                      "001010000000001",
		MmeName:                   ".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
		EpsLocationUpdateType:     []byte{byte(0x11)},
		NewLocationAreaIdentifier: make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
	}
	encodedMsg, err := message.EncodeSGsAPLocationUpdateRequest(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPLocationUpdateRequest), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedMMEName, _ := ie.EncodeMMEName(
		".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
	)
	encodedLocationUpdateType, _ := ie.EncodeFixedLengthIE(
		decode.IEIEPSLocationUpdateType,
		decode.IELengthEPSLocationUpdateType,
		[]byte{byte(0x11)},
	)
	encodedNewLAI, _ := ie.EncodeFixedLengthIE(
		decode.IEILocationAreaIdentifier,
		decode.IELengthLocationAreaIdentifier,
		make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
	)
	restOfFields := append(encodedIMSI, encodedMMEName...)
	restOfFields = append(restOfFields, encodedLocationUpdateType...)
	restOfFields = append(restOfFields, encodedNewLAI...)
	assert.Equal(t, restOfFields, encodedMsg[1:])

	// with optional fields
	msg = &protos.LocationUpdateRequest{
		Imsi:                      "001010000000001",
		MmeName:                   ".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
		EpsLocationUpdateType:     []byte{byte(0x11)},
		NewLocationAreaIdentifier: make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
		OldLocationAreaIdentifier: make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
		TmsiStatus:                []byte{byte(0x11)},
		Imeisv:                    make([]byte, decode.IELengthIMEISV-mandatoryFieldLength),
		Tai:                       make([]byte, decode.IELengthTAI-mandatoryFieldLength),
		ECgi:                      make([]byte, decode.IELengthEUTRANCellGlobalIdentity-mandatoryFieldLength),
	}
	encodedMsg, err = message.EncodeSGsAPLocationUpdateRequest(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPLocationUpdateRequest), encodedMsg[0])
	encodedOldLAI, _ := ie.EncodeFixedLengthIE(
		decode.IEILocationAreaIdentifier,
		decode.IELengthLocationAreaIdentifier,
		make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
	)
	encodedTMSIStatus, _ := ie.EncodeFixedLengthIE(
		decode.IEITMSIStatus,
		decode.IELengthTMSIStatus,
		[]byte{byte(0x11)},
	)
	encodedIMEISV, _ := ie.EncodeFixedLengthIE(
		decode.IEIIMEISV,
		decode.IELengthIMEISV,
		make([]byte, decode.IELengthIMEISV-mandatoryFieldLength),
	)
	encodedTAI, _ := ie.EncodeFixedLengthIE(
		decode.IEITAI,
		decode.IELengthTAI,
		make([]byte, decode.IELengthTAI-mandatoryFieldLength),
	)
	encodedECGI, _ := ie.EncodeFixedLengthIE(
		decode.IEIEUTRANCellGlobalIdentity,
		decode.IELengthEUTRANCellGlobalIdentity,
		make([]byte, decode.IELengthEUTRANCellGlobalIdentity-mandatoryFieldLength),
	)
	restOfFields = append(restOfFields, encodedOldLAI...)
	restOfFields = append(restOfFields, encodedTMSIStatus...)
	restOfFields = append(restOfFields, encodedIMEISV...)
	restOfFields = append(restOfFields, encodedTAI...)
	restOfFields = append(restOfFields, encodedECGI...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPMMInformationRequest(t *testing.T) {
	msg := &protos.MMInformationRequest{
		Imsi:          "001010000000001",
		MmInformation: []byte{byte(0x11)},
	}
	encodedMsg, err := message.EncodeSGsAPMMInformationRequest(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPMMInformationRequest), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedMMInfo, _ := ie.EncodeVariableLengthIE(
		decode.IEIMMInformation,
		decode.IELengthMMInformationMin,
		[]byte{byte(0x11)},
	)
	restOfFields := append(encodedIMSI, encodedMMInfo...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPPagingReject(t *testing.T) {
	msg := &protos.PagingReject{
		Imsi:     "001010000000001",
		SgsCause: []byte{byte(0x11)},
	}
	encodedMsg, err := message.EncodeSGsAPPagingReject(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPPagingReject), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedSGsCause, _ := ie.EncodeFixedLengthIE(
		decode.IEISGsCause,
		decode.IELengthSGsCause,
		[]byte{byte(0x11)},
	)
	restOfFields := append(encodedIMSI, encodedSGsCause...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPPagingRequest(t *testing.T) {
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator

	// without optional fields
	msg := &protos.PagingRequest{
		Imsi:             "001010000000001",
		VlrName:          "aaaaaaaaa",
		ServiceIndicator: make([]byte, decode.IELengthServiceIndicator-mandatoryFieldLength),
	}
	encodedMsg, err := message.EncodeSGsAPPagingRequest(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPPagingRequest), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedVLRName, _ := ie.EncodeVariableLengthIE(
		decode.IEIVLRName,
		decode.IELengthVLRNameMin,
		[]byte("aaaaaaaaa"),
	)
	encodedServiceIndicator, _ := ie.EncodeFixedLengthIE(
		decode.IEIServiceIndicator,
		decode.IELengthServiceIndicator,
		make([]byte, decode.IELengthServiceIndicator-mandatoryFieldLength),
	)
	restOfFields := append(encodedIMSI, encodedVLRName...)
	restOfFields = append(restOfFields, encodedServiceIndicator...)
	assert.Equal(t, restOfFields, encodedMsg[1:])

	// with optional fields
	msg = &protos.PagingRequest{
		Imsi:                   "001010000000001",
		VlrName:                "aaaaaaaaa",
		ServiceIndicator:       make([]byte, decode.IELengthServiceIndicator-mandatoryFieldLength),
		Tmsi:                   make([]byte, decode.IELengthTMSI-mandatoryFieldLength),
		Cli:                    make([]byte, decode.IELengthCLIMin-mandatoryFieldLength),
		LocationAreaIdentifier: make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
		GlobalCnId:             make([]byte, decode.IELengthGlobalCNId-mandatoryFieldLength),
		SsCode:                 make([]byte, decode.IELengthSSCode-mandatoryFieldLength),
		LcsIndicator:           make([]byte, decode.IELengthLCSIndicator-mandatoryFieldLength),
		LcsClientIdentity:      make([]byte, decode.IELengthLCSClientIdentityMin-mandatoryFieldLength),
		ChannelNeeded:          make([]byte, decode.IELengthChannelNeeded-mandatoryFieldLength),
		EmlppPriority:          make([]byte, decode.IELengthEMLPPPriority-mandatoryFieldLength),
	}
	encodedMsg, err = message.EncodeSGsAPPagingRequest(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPPagingRequest), encodedMsg[0])
	encodedTMSI, _ := ie.EncodeFixedLengthIE(
		decode.IEITMSI,
		decode.IELengthTMSI,
		make([]byte, decode.IELengthTMSI-mandatoryFieldLength),
	)
	encodedCLI, _ := ie.EncodeLimitedLengthIE(
		decode.IEICLI,
		decode.IELengthCLIMin,
		decode.IELengthCLIMax,
		make([]byte, decode.IELengthCLIMin-mandatoryFieldLength),
	)
	encodedLAI, _ := ie.EncodeFixedLengthIE(
		decode.IEILocationAreaIdentifier,
		decode.IELengthLocationAreaIdentifier,
		make([]byte, decode.IELengthLocationAreaIdentifier-mandatoryFieldLength),
	)
	encodedGlobalCNID, _ := ie.EncodeFixedLengthIE(
		decode.IEIGlobalCNId,
		decode.IELengthGlobalCNId,
		make([]byte, decode.IELengthGlobalCNId-mandatoryFieldLength),
	)
	encodedSSCode, _ := ie.EncodeFixedLengthIE(
		decode.IEISSCode,
		decode.IELengthSSCode,
		make([]byte, decode.IELengthSSCode-mandatoryFieldLength),
	)
	encodedLCSIndicator, _ := ie.EncodeFixedLengthIE(
		decode.IEILCSIndicator,
		decode.IELengthLCSIndicator,
		make([]byte, decode.IELengthLCSIndicator-mandatoryFieldLength),
	)
	encodedLCSClientID, _ := ie.EncodeVariableLengthIE(
		decode.IEILCSClientIdentity,
		decode.IELengthLCSClientIdentityMin,
		make([]byte, decode.IELengthLCSClientIdentityMin-mandatoryFieldLength),
	)
	encodedChannelNeeded, _ := ie.EncodeFixedLengthIE(
		decode.IEIChannelNeeded,
		decode.IELengthChannelNeeded,
		make([]byte, decode.IELengthChannelNeeded-mandatoryFieldLength),
	)
	encodedEMLPPPriority, _ := ie.EncodeFixedLengthIE(
		decode.IEIeMLPPPriority,
		decode.IELengthEMLPPPriority,
		make([]byte, decode.IELengthEMLPPPriority-mandatoryFieldLength),
	)
	restOfFields = append(restOfFields, encodedTMSI...)
	restOfFields = append(restOfFields, encodedCLI...)
	restOfFields = append(restOfFields, encodedLAI...)
	restOfFields = append(restOfFields, encodedGlobalCNID...)
	restOfFields = append(restOfFields, encodedSSCode...)
	restOfFields = append(restOfFields, encodedLCSIndicator...)
	restOfFields = append(restOfFields, encodedLCSClientID...)
	restOfFields = append(restOfFields, encodedChannelNeeded...)
	restOfFields = append(restOfFields, encodedEMLPPPriority...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPReleaseRequest(t *testing.T) {
	msg := &protos.ReleaseRequest{
		Imsi:     "001010000000001",
		SgsCause: []byte{byte(0x11)},
	}
	encodedMsg, err := message.EncodeSGsAPReleaseRequest(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPReleaseRequest), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedSGsCause, _ := ie.EncodeFixedLengthIE(
		decode.IEISGsCause,
		decode.IELengthSGsCause,
		[]byte{byte(0x11)},
	)
	restOfFields := append(encodedIMSI, encodedSGsCause...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPResetAck(t *testing.T) {
	msg := &protos.ResetAck{
		MmeName: ".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
	}
	encodedMsg, err := message.EncodeSGsAPResetAck(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPResetAck), encodedMsg[0])
	encodedMMEName, _ := ie.EncodeMMEName(
		".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
	)
	assert.Equal(t, encodedMMEName, encodedMsg[1:])
}

func TestEncodeSGsAPResetIndication(t *testing.T) {
	msg := &protos.ResetIndication{
		MmeName: ".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
	}
	encodedMsg, err := message.EncodeSGsAPResetIndication(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPResetIndication), encodedMsg[0])
	encodedMMEName, _ := ie.EncodeMMEName(
		".mmec01.mmegi0001.mme.EPC.mnc001.mcc001.3gppnetwork.org",
	)
	assert.Equal(t, encodedMMEName, encodedMsg[1:])
}

func TestEncodeSGsAPServiceAbortRequest(t *testing.T) {
	msg := &protos.ServiceAbortRequest{
		Imsi: "001010000000001",
	}
	encodedMsg, err := message.EncodeSGsAPServiceAbortRequest(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPServiceAbortRequest), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	assert.Equal(t, encodedIMSI, encodedMsg[1:])
}

func TestEncodeSGsAPServiceRequest(t *testing.T) {
	// without optional fields
	msg := &protos.ServiceRequest{
		Imsi:             "001010000000001",
		ServiceIndicator: []byte{byte(0x11)},
	}
	encodedMsg, err := message.EncodeSGsAPServiceRequest(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPServiceRequest), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedServiceIndicator, _ := ie.EncodeFixedLengthIE(
		decode.IEIServiceIndicator,
		decode.IELengthServiceIndicator,
		[]byte{byte(0x11)},
	)
	restOfFields := append(encodedIMSI, encodedServiceIndicator...)
	assert.Equal(t, restOfFields, encodedMsg[1:])

	// with optional fields
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	msg = &protos.ServiceRequest{
		Imsi:                    "001010000000001",
		ServiceIndicator:        []byte{byte(0x11)},
		Imeisv:                  make([]byte, decode.IELengthIMEISV-mandatoryFieldLength),
		UeTimeZone:              []byte{byte(0x12)},
		MobileStationClassmark2: make([]byte, decode.IELengthMobileStationClassmark2-mandatoryFieldLength),
		Tai:                     make([]byte, decode.IELengthTAI-mandatoryFieldLength),
		ECgi:                    make([]byte, decode.IELengthEUTRANCellGlobalIdentity-mandatoryFieldLength),
		UeEmmMode:               []byte{byte(0x13)},
	}
	encodedMsg, err = message.EncodeSGsAPServiceRequest(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPServiceRequest), encodedMsg[0])
	encodedIMEISV, _ := ie.EncodeFixedLengthIE(
		decode.IEIIMEISV,
		decode.IELengthIMEISV,
		make([]byte, decode.IELengthIMEISV-mandatoryFieldLength),
	)
	encodedUETimeZone, _ := ie.EncodeFixedLengthIE(
		decode.IEIUETimeZone,
		decode.IELengthUETimeZone,
		[]byte{byte(0x12)},
	)
	encodedMSC2, _ := ie.EncodeFixedLengthIE(
		decode.IEIMobileStationClassmark2,
		decode.IELengthMobileStationClassmark2,
		make([]byte, decode.IELengthMobileStationClassmark2-mandatoryFieldLength),
	)
	encodedTAI, _ := ie.EncodeFixedLengthIE(
		decode.IEITAI,
		decode.IELengthTAI,
		make([]byte, decode.IELengthTAI-mandatoryFieldLength),
	)
	encodedECGI, _ := ie.EncodeFixedLengthIE(
		decode.IEIEUTRANCellGlobalIdentity,
		decode.IELengthEUTRANCellGlobalIdentity,
		make([]byte, decode.IELengthEUTRANCellGlobalIdentity-mandatoryFieldLength),
	)
	encodedUEEMMMode, _ := ie.EncodeFixedLengthIE(
		decode.IEIUEEMMMode,
		decode.IELengthUEEMMMode,
		[]byte{byte(0x13)},
	)
	restOfFields = append(restOfFields, encodedIMEISV...)
	restOfFields = append(restOfFields, encodedUETimeZone...)
	restOfFields = append(restOfFields, encodedMSC2...)
	restOfFields = append(restOfFields, encodedTAI...)
	restOfFields = append(restOfFields, encodedECGI...)
	restOfFields = append(restOfFields, encodedUEEMMMode...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPStatus(t *testing.T) {
	// without IMSI
	msg := &protos.Status{
		SgsCause:         []byte{byte(0x13)},
		ErroneousMessage: []byte{byte(0x14), byte(0x15), byte(0x16), byte(0x16)},
	}
	encodedMsg, err := message.EncodeSGsAPStatus(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPStatus), encodedMsg[0])
	encodedSGsCause, _ := ie.EncodeFixedLengthIE(
		decode.IEISGsCause,
		decode.IELengthSGsCause,
		[]byte{byte(0x13)},
	)
	encodedErroneousMessage, _ := ie.EncodeVariableLengthIE(
		decode.IEIErroneousMessage,
		decode.IELengthErroneousMessageMin,
		[]byte{byte(0x14), byte(0x15), byte(0x16), byte(0x16)},
	)
	restOfFields := append(encodedSGsCause, encodedErroneousMessage...)
	assert.Equal(t, restOfFields, encodedMsg[1:])

	// with IMSI
	msg = &protos.Status{
		Imsi:             "001010000000001",
		SgsCause:         []byte{byte(0x13)},
		ErroneousMessage: []byte{byte(0x14), byte(0x15), byte(0x16), byte(0x16)},
	}
	encodedMsg, err = message.EncodeSGsAPStatus(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPStatus), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	restOfFields = append(encodedIMSI, restOfFields...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPTMSIReallocationComplete(t *testing.T) {
	msg := &protos.TMSIReallocationComplete{
		Imsi: "001010000000001",
	}
	encodedMsg, err := message.EncodeSGsAPTMSIReallocationComplete(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPTMSIReallocationComplete), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	assert.Equal(t, encodedIMSI, encodedMsg[1:])
}

func TestEncodeSGsAPUEActivityIndication(t *testing.T) {
	msg := &protos.UEActivityIndication{
		Imsi: "001010000000001",
	}
	encodedMsg, err := message.EncodeSGsAPUEActivityIndication(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPUEActivityIndication), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	assert.Equal(t, encodedIMSI, encodedMsg[1:])
}

func TestEncodeSGsAPUEUnreachable(t *testing.T) {
	msg := &protos.UEUnreachable{
		Imsi:     "001010000000001",
		SgsCause: []byte{byte(0x13)},
	}
	encodedMsg, err := message.EncodeSGsAPUEUnreachable(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPUEUnreachable), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedSGsCause, _ := ie.EncodeFixedLengthIE(
		decode.IEISGsCause,
		decode.IELengthSGsCause,
		[]byte{byte(0x13)},
	)
	restOfFields := append(encodedIMSI, encodedSGsCause...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}

func TestEncodeSGsAPUplinkUnitdata(t *testing.T) {
	// without optional fields
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	msg := &protos.UplinkUnitdata{
		Imsi:                "001010000000001",
		NasMessageContainer: make([]byte, decode.IELengthNASMessageContainerMax-mandatoryFieldLength),
	}
	encodedMsg, err := message.EncodeSGsAPUplinkUnitdata(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPUplinkUnitdata), encodedMsg[0])
	encodedIMSI, _ := ie.EncodeIMSI("001010000000001")
	encodedNasMessageContainer, _ := ie.EncodeLimitedLengthIE(
		decode.IEINASMessageContainer,
		decode.IELengthNASMessageContainerMin,
		decode.IELengthNASMessageContainerMax,
		make([]byte, decode.IELengthNASMessageContainerMax-mandatoryFieldLength),
	)
	restOfFields := append(encodedIMSI, encodedNasMessageContainer...)
	assert.Equal(t, restOfFields, encodedMsg[1:])

	// with optional fields
	msg = &protos.UplinkUnitdata{
		Imsi:                    "001010000000001",
		NasMessageContainer:     make([]byte, decode.IELengthNASMessageContainerMax-mandatoryFieldLength),
		Imeisv:                  make([]byte, decode.IELengthIMEISV-mandatoryFieldLength),
		UeTimeZone:              []byte{byte(0x11)},
		MobileStationClassmark2: []byte{byte(0x12), byte(0x13), byte(0x14)},
		Tai:                     make([]byte, decode.IELengthTAI-mandatoryFieldLength),
		ECgi:                    make([]byte, decode.IELengthEUTRANCellGlobalIdentity-mandatoryFieldLength),
	}
	encodedMsg, err = message.EncodeSGsAPUplinkUnitdata(msg)
	assert.NoError(t, err)
	assert.Equal(t, byte(decode.SGsAPUplinkUnitdata), encodedMsg[0])
	encodedIMEISV, _ := ie.EncodeFixedLengthIE(
		decode.IEIIMEISV,
		decode.IELengthIMEISV,
		make([]byte, decode.IELengthIMEISV-mandatoryFieldLength),
	)
	encodedUETimeZone, _ := ie.EncodeFixedLengthIE(
		decode.IEIUETimeZone,
		decode.IELengthUETimeZone,
		[]byte{byte(0x11)},
	)
	encodedMSC2, _ := ie.EncodeFixedLengthIE(
		decode.IEIMobileStationClassmark2,
		decode.IELengthMobileStationClassmark2,
		[]byte{byte(0x12), byte(0x13), byte(0x14)},
	)
	encodedTAI, _ := ie.EncodeFixedLengthIE(
		decode.IEITAI,
		decode.IELengthTAI,
		make([]byte, decode.IELengthTAI-mandatoryFieldLength),
	)
	encodedECGI, _ := ie.EncodeFixedLengthIE(
		decode.IEIEUTRANCellGlobalIdentity,
		decode.IELengthEUTRANCellGlobalIdentity,
		make([]byte, decode.IELengthEUTRANCellGlobalIdentity-mandatoryFieldLength),
	)
	restOfFields = append(restOfFields, encodedIMEISV...)
	restOfFields = append(restOfFields, encodedUETimeZone...)
	restOfFields = append(restOfFields, encodedMSC2...)
	restOfFields = append(restOfFields, encodedTAI...)
	restOfFields = append(restOfFields, encodedECGI...)
	assert.Equal(t, restOfFields, encodedMsg[1:])
}
