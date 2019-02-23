/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package message_test

import (
	"fmt"
	"testing"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/decode/message"
	"magma/feg/gateway/services/csfb/servicers/decode/test_utils"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

func TestSGsMessageDecoder(t *testing.T) {
	chunk := []byte{byte(0xFF)}
	msgType, msg, err := message.SGsMessageDecoder(chunk)
	assert.EqualError(t, err, "unknown message type")
	assert.Equal(t, decode.SGsMessageType(byte(0xFF)), msgType)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPIMSIDetachAck(t *testing.T) {
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPIMSIDetachAck)}, imsi...)
	msg, err := message.DecodeSGsAPIMSIDetachAck(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.IMSIDetachAck{
		Imsi: "111111",
	})
	assert.Equal(t, expectedMsg, msg)
}

func TestDecodeSGsAPLocationUpdateAccept(t *testing.T) {
	// without mobile identity
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPLocationUpdateAccept)}, imsi...)
	LAI := test_utils.ConstructDefaultLocationAreaIdentifier()
	chunk = append(chunk, LAI...)

	msg, err := message.DecodeSGsAPLocationUpdateAccept(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.LocationUpdateAccept{
		Imsi:                   "111111",
		LocationAreaIdentifier: LAI[2:],
	})
	assert.Equal(t, expectedMsg, msg)

	// with new IMSI
	chunk = append([]byte{byte(decode.SGsAPLocationUpdateAccept)}, imsi...)
	chunk = append(chunk, LAI...)
	newIMSI, _ := test_utils.ConstructMobileIdentity(
		"222222",
		[]byte{},
	)
	chunk = append(chunk, newIMSI...)

	msg, err = message.DecodeSGsAPLocationUpdateAccept(chunk)
	assert.NoError(t, err)
	expectedMsg, _ = ptypes.MarshalAny(&protos.LocationUpdateAccept{
		Imsi:                   "111111",
		LocationAreaIdentifier: LAI[2:],
		NewIMSITMSI:            &protos.LocationUpdateAccept_NewImsi{NewImsi: "222222"},
	})
	assert.Equal(t, expectedMsg, msg)

	// with new TMSI
	chunk = append([]byte{byte(decode.SGsAPLocationUpdateAccept)}, imsi...)
	chunk = append(chunk, LAI...)
	newTMSI, _ := test_utils.ConstructMobileIdentity(
		"",
		[]byte{byte(0x11), byte(0x12), byte(0x13), byte(0x14)},
	)
	chunk = append(chunk, newTMSI...)

	msg, err = message.DecodeSGsAPLocationUpdateAccept(chunk)
	assert.NoError(t, err)
	expectedMsg, _ = ptypes.MarshalAny(&protos.LocationUpdateAccept{
		Imsi:                   "111111",
		LocationAreaIdentifier: LAI[2:],
		NewIMSITMSI: &protos.LocationUpdateAccept_NewTmsi{
			NewTmsi: []byte{byte(0x11), byte(0x12), byte(0x13), byte(0x14)},
		},
	})
	assert.Equal(t, expectedMsg, msg)

	// with wrong byte 3 for TMSI
	chunk = append([]byte{byte(decode.SGsAPLocationUpdateAccept)}, imsi...)
	chunk = append(chunk, LAI...)
	newTMSI[2] = byte(0xE4)
	chunk = append(chunk, newTMSI...)

	msg, err = message.DecodeSGsAPLocationUpdateAccept(chunk)
	assert.EqualError(t, err, "byte 3 of mobile identity field as TMSI should be 0xF4, not 0xE4")
	assert.Equal(t, &any.Any{}, msg)

	// with wrong identity type
	chunk = append([]byte{byte(decode.SGsAPLocationUpdateAccept)}, imsi...)
	chunk = append(chunk, LAI...)
	newTMSI[2] = byte(0xF2)
	chunk = append(chunk, newTMSI...)

	msg, err = message.DecodeSGsAPLocationUpdateAccept(chunk)
	assert.EqualError(t, err, "cannot recognize the identity type 2 for mobile identity field")
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPLocationUpdateReject(t *testing.T) {
	// with LAI
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPLocationUpdateReject)}, imsi...)
	rejectCause := []byte{byte(0x11), byte(0x11), byte(0x11)}
	chunk = append(chunk, rejectCause...)
	LAI := test_utils.ConstructDefaultLocationAreaIdentifier()
	chunk = append(chunk, LAI...)
	msg, err := message.DecodeSGsAPLocationUpdateReject(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.LocationUpdateReject{
		Imsi:                   "111111",
		LocationAreaIdentifier: LAI[2:],
		RejectCause:            rejectCause,
	})
	assert.Equal(t, expectedMsg, msg)

	// without LAI
	chunk = append([]byte{byte(decode.SGsAPLocationUpdateReject)}, imsi...)
	chunk = append(chunk, rejectCause...)
	msg, err = message.DecodeSGsAPLocationUpdateReject(chunk)
	assert.NoError(t, err)
	expectedMsg, _ = ptypes.MarshalAny(&protos.LocationUpdateReject{
		Imsi:                   "111111",
		LocationAreaIdentifier: nil,
		RejectCause:            rejectCause,
	})
	assert.Equal(t, expectedMsg, msg)
}

func TestDecodeSGsAPMMInformationRequest(t *testing.T) {
	// Successfully decode message
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPMMInformationRequest)}, imsi...)
	mmInfo := test_utils.ConstructDefaultMMInformation()
	chunk = append(chunk, mmInfo...)
	msg, err := message.DecodeSGsAPMMInformationRequest(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.MMInformationRequest{
		Imsi:          "111111",
		MmInformation: mmInfo[2:],
	})
	assert.Equal(t, expectedMsg, msg)

	// chunk too short
	imsi, _ = test_utils.ConstructIMSI("111111")
	chunk = append([]byte{byte(decode.SGsAPMMInformationRequest)}, imsi...)
	msg, err = message.DecodeSGsAPMMInformationRequest(chunk)
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin+decode.IELengthMMInformationMin,
		len(chunk),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPPagingRequest(t *testing.T) {
	// All fields present
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPPagingRequest)}, imsi...)
	vlrName := test_utils.ConstructDefaultVLRName()
	chunk = append(chunk, vlrName...)
	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	serviceIndicator := test_utils.ConstructDefaultIE(
		decode.IEIServiceIndicator,
		decode.IELengthServiceIndicator-mandatoryFieldLength,
	)
	chunk = append(chunk, serviceIndicator...)
	tmsi := test_utils.ConstructDefaultIE(
		decode.IEITMSI,
		decode.IELengthTMSI-mandatoryFieldLength,
	)
	chunk = append(chunk, tmsi...)
	cli := test_utils.ConstructDefaultIE(
		decode.IEICLI,
		decode.IELengthCLIMin-mandatoryFieldLength,
	)
	chunk = append(chunk, cli...)
	lai := test_utils.ConstructDefaultIE(
		decode.IEILocationAreaIdentifier,
		decode.IELengthLocationAreaIdentifier-mandatoryFieldLength,
	)
	chunk = append(chunk, lai...)
	globalCNId := test_utils.ConstructDefaultIE(
		decode.IEIGlobalCNId,
		decode.IELengthGlobalCNId-mandatoryFieldLength,
	)
	chunk = append(chunk, globalCNId...)
	ssCode := test_utils.ConstructDefaultIE(
		decode.IEISSCode,
		decode.IELengthSSCode-mandatoryFieldLength,
	)
	chunk = append(chunk, ssCode...)
	lcsIndicator := test_utils.ConstructDefaultIE(
		decode.IEILCSIndicator,
		decode.IELengthLCSIndicator-mandatoryFieldLength,
	)
	chunk = append(chunk, lcsIndicator...)
	lcsClientIdentity := test_utils.ConstructDefaultIE(
		decode.IEILCSClientIdentity,
		decode.IELengthLCSClientIdentityMin-mandatoryFieldLength,
	)
	chunk = append(chunk, lcsClientIdentity...)
	channelNeeded := test_utils.ConstructDefaultIE(
		decode.IEIChannelNeeded,
		decode.IELengthChannelNeeded-mandatoryFieldLength,
	)
	chunk = append(chunk, channelNeeded...)
	eMLPPPriprity := test_utils.ConstructDefaultIE(
		decode.IEIeMLPPPriority,
		decode.IELengthEMLPPPriority-mandatoryFieldLength,
	)
	chunk = append(chunk, eMLPPPriprity...)
	msg, err := message.DecodeSGsAPPagingRequest(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.PagingRequest{
		Imsi:                   "111111",
		VlrName:                "www.facebook.com",
		ServiceIndicator:       serviceIndicator[mandatoryFieldLength:],
		Tmsi:                   tmsi[mandatoryFieldLength:],
		Cli:                    cli[mandatoryFieldLength:],
		LocationAreaIdentifier: lai[mandatoryFieldLength:],
		GlobalCnId:             globalCNId[mandatoryFieldLength:],
		SsCode:                 ssCode[mandatoryFieldLength:],
		LcsIndicator:           lcsIndicator[mandatoryFieldLength:],
		LcsClientIdentity:      lcsClientIdentity[mandatoryFieldLength:],
		ChannelNeeded:          channelNeeded[mandatoryFieldLength:],
		EmlppPriority:          eMLPPPriprity[mandatoryFieldLength:],
	})
	assert.Equal(t, expectedMsg, msg)

	// TMSI not present
	chunk = append([]byte{byte(decode.SGsAPPagingRequest)}, imsi...)
	vlrName = test_utils.ConstructDefaultVLRName()
	chunk = append(chunk, vlrName...)
	chunk = append(chunk, serviceIndicator...)
	chunk = append(chunk, cli...)
	chunk = append(chunk, lai...)
	chunk = append(chunk, globalCNId...)
	chunk = append(chunk, ssCode...)
	chunk = append(chunk, lcsIndicator...)
	chunk = append(chunk, lcsClientIdentity...)
	chunk = append(chunk, channelNeeded...)
	chunk = append(chunk, eMLPPPriprity...)
	msg, err = message.DecodeSGsAPPagingRequest(chunk)
	assert.NoError(t, err)
	expectedMsg, _ = ptypes.MarshalAny(&protos.PagingRequest{
		Imsi:             "111111",
		VlrName:          "www.facebook.com",
		ServiceIndicator: serviceIndicator[mandatoryFieldLength:],
		// Tmsi: tmsi[mandatoryFieldLength:],
		Cli:                    cli[mandatoryFieldLength:],
		LocationAreaIdentifier: lai[mandatoryFieldLength:],
		GlobalCnId:             globalCNId[mandatoryFieldLength:],
		SsCode:                 ssCode[mandatoryFieldLength:],
		LcsIndicator:           lcsIndicator[mandatoryFieldLength:],
		LcsClientIdentity:      lcsClientIdentity[mandatoryFieldLength:],
		ChannelNeeded:          channelNeeded[mandatoryFieldLength:],
		EmlppPriority:          eMLPPPriprity[mandatoryFieldLength:],
	})
	assert.Equal(t, expectedMsg, msg)

	// CLI, SS Code, Channel Needed, eMLPP Priority not present
	chunk = append([]byte{byte(decode.SGsAPPagingRequest)}, imsi...)
	vlrName = test_utils.ConstructDefaultVLRName()
	chunk = append(chunk, vlrName...)
	chunk = append(chunk, serviceIndicator...)
	chunk = append(chunk, tmsi...)
	chunk = append(chunk, lai...)
	chunk = append(chunk, globalCNId...)
	chunk = append(chunk, lcsIndicator...)
	chunk = append(chunk, lcsClientIdentity...)
	msg, err = message.DecodeSGsAPPagingRequest(chunk)
	assert.NoError(t, err)
	expectedMsg, _ = ptypes.MarshalAny(&protos.PagingRequest{
		Imsi:                   "111111",
		VlrName:                "www.facebook.com",
		ServiceIndicator:       serviceIndicator[mandatoryFieldLength:],
		Tmsi:                   tmsi[mandatoryFieldLength:],
		LocationAreaIdentifier: lai[mandatoryFieldLength:],
		GlobalCnId:             globalCNId[mandatoryFieldLength:],
		LcsIndicator:           lcsIndicator[mandatoryFieldLength:],
		LcsClientIdentity:      lcsClientIdentity[mandatoryFieldLength:],
	})
	assert.Equal(t, expectedMsg, msg)

	// All optional fields not present
	chunk = append([]byte{byte(decode.SGsAPPagingRequest)}, imsi...)
	vlrName = test_utils.ConstructDefaultVLRName()
	chunk = append(chunk, vlrName...)
	chunk = append(chunk, serviceIndicator...)
	msg, err = message.DecodeSGsAPPagingRequest(chunk)
	assert.NoError(t, err)
	expectedMsg, _ = ptypes.MarshalAny(&protos.PagingRequest{
		Imsi:             "111111",
		VlrName:          "www.facebook.com",
		ServiceIndicator: serviceIndicator[mandatoryFieldLength:],
	})
	assert.Equal(t, expectedMsg, msg)

	// Optional fields in wrong order (SS Code <-> LCS Indicator)
	chunk = append([]byte{byte(decode.SGsAPPagingRequest)}, imsi...)
	vlrName = test_utils.ConstructDefaultVLRName()
	chunk = append(chunk, vlrName...)
	chunk = append(chunk, serviceIndicator...)
	chunk = append(chunk, tmsi...)
	chunk = append(chunk, cli...)
	chunk = append(chunk, lai...)
	chunk = append(chunk, globalCNId...)
	chunk = append(chunk, lcsIndicator...)
	chunk = append(chunk, ssCode...)
	chunk = append(chunk, lcsClientIdentity...)
	chunk = append(chunk, channelNeeded...)
	chunk = append(chunk, eMLPPPriprity...)
	msg, err = message.DecodeSGsAPPagingRequest(chunk)
	assert.EqualError(t, err, "tried all possible IE but still some bytes undecoded")
	assert.Equal(t, &any.Any{}, msg)

	// Wrong IE present in optional fields (MM Information)
	chunk = append([]byte{byte(decode.SGsAPPagingRequest)}, imsi...)
	vlrName = test_utils.ConstructDefaultVLRName()
	chunk = append(chunk, vlrName...)
	chunk = append(chunk, serviceIndicator...)
	chunk = append(chunk, tmsi...)
	chunk = append(chunk, cli...)
	mminfo := test_utils.ConstructDefaultIE(
		decode.IEIMMInformation,
		decode.IELengthMMInformationMin-mandatoryFieldLength,
	)
	chunk = append(chunk, mminfo...)
	chunk = append(chunk, lai...)
	chunk = append(chunk, globalCNId...)
	chunk = append(chunk, lcsIndicator...)
	chunk = append(chunk, ssCode...)
	chunk = append(chunk, lcsClientIdentity...)
	chunk = append(chunk, channelNeeded...)
	chunk = append(chunk, eMLPPPriprity...)
	msg, err = message.DecodeSGsAPPagingRequest(chunk)
	assert.EqualError(t, err, "tried all possible IE but still some bytes undecoded")
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPEPSDetachAck(t *testing.T) {
	// Successfully decode message
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPEPSDetachAck)}, imsi...)
	msg, err := message.DecodeSGsAPEPSDetachAck(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.EPSDetachAck{
		Imsi: "111111",
	})
	assert.Equal(t, expectedMsg, msg)

	// Chunk too short
	msg, err = message.DecodeSGsAPEPSDetachAck(chunk[:2])
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin,
		decode.IELengthMessageType+decode.IELengthIMSIMax,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPAlertRequest(t *testing.T) {
	// Successfully decode message
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPAlertRequest)}, imsi...)
	msg, err := message.DecodeSGsAPAlertRequest(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.AlertRequest{
		Imsi: "111111",
	})
	assert.Equal(t, expectedMsg, msg)

	// Chunk too short
	msg, err = message.DecodeSGsAPAlertRequest(chunk[:2])
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin,
		decode.IELengthMessageType+decode.IELengthIMSIMax,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPDownlinkUnitdata(t *testing.T) {
	// Successfully decode message
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPDownlinkUnitdata)}, imsi...)
	nasMessageContainer := test_utils.ConstructDefaultIE(decode.IEINASMessageContainer, 5)
	chunk = append(chunk, nasMessageContainer...)
	msg, err := message.DecodeSGsAPDownlinkUnitdata(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.DownlinkUnitdata{
		Imsi:                "111111",
		NasMessageContainer: nasMessageContainer[2:],
	})
	assert.Equal(t, expectedMsg, msg)

	// Chunk too short
	msg, err = message.DecodeSGsAPDownlinkUnitdata(chunk[:2])
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin+decode.IELengthNASMessageContainerMin,
		decode.IELengthMessageType+decode.IELengthIMSIMax+decode.IELengthNASMessageContainerMax,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)

	// Chunk too long
	wrongChunk := append([]byte{byte(decode.SGsAPDownlinkUnitdata)}, imsi...)
	nasMessageContainer = test_utils.ConstructDefaultIE(decode.IEINASMessageContainer, 500)
	wrongChunk = append(wrongChunk, nasMessageContainer...)
	msg, err = message.DecodeSGsAPDownlinkUnitdata(wrongChunk)
	errorMsg = fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin+decode.IELengthNASMessageContainerMin,
		decode.IELengthMessageType+decode.IELengthIMSIMax+decode.IELengthNASMessageContainerMax,
		len(wrongChunk),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPReleaseRequest(t *testing.T) {
	// With SGs Cause
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPServiceRequest)}, imsi...)
	sgsCause := test_utils.ConstructDefaultIE(decode.IEISGsCause, 1)
	chunk = append(chunk, sgsCause...)
	msg, err := message.DecodeSGsAPReleaseRequest(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.ReleaseRequest{
		Imsi:     "111111",
		SgsCause: sgsCause[2:],
	})
	assert.Equal(t, expectedMsg, msg)

	// Without SGs Cause
	chunk = append([]byte{byte(decode.SGsAPServiceRequest)}, imsi...)
	msg, err = message.DecodeSGsAPReleaseRequest(chunk)
	assert.NoError(t, err)
	expectedMsg, _ = ptypes.MarshalAny(&protos.ReleaseRequest{
		Imsi: "111111",
	})
	assert.Equal(t, expectedMsg, msg)

	// Chunk too short
	msg, err = message.DecodeSGsAPReleaseRequest(chunk[:2])
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin,
		decode.IELengthMessageType+decode.IELengthIMSIMax+decode.IELengthSGsCause,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)

	// Chunk too long
	imsi, _ = test_utils.ConstructIMSI("111111111111111111111111")
	chunk = append([]byte{byte(decode.SGsAPServiceRequest)}, imsi...)
	chunk = append(chunk, sgsCause...)
	msg, err = message.DecodeSGsAPReleaseRequest(chunk)
	errorMsg = fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin,
		decode.IELengthMessageType+decode.IELengthIMSIMax+decode.IELengthSGsCause,
		len(chunk),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPServiceAbortRequest(t *testing.T) {
	// Successfully decode message
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPServiceAbortRequest)}, imsi...)
	msg, err := message.DecodeSGsAPServiceAbortRequest(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.ServiceAbortRequest{
		Imsi: "111111",
	})
	assert.Equal(t, expectedMsg, msg)

	// Chunk too short
	msg, err = message.DecodeSGsAPServiceAbortRequest(chunk[:2])
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin,
		decode.IELengthMessageType+decode.IELengthIMSIMax,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)

	// Chunk too long
	imsi, _ = test_utils.ConstructIMSI("111111111111111111111111")
	chunk = append([]byte{byte(decode.SGsAPServiceAbortRequest)}, imsi...)
	msg, err = message.DecodeSGsAPServiceAbortRequest(chunk)
	errorMsg = fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin,
		decode.IELengthMessageType+decode.IELengthIMSIMax,
		len(chunk),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPStatus(t *testing.T) {
	// With IMSI
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPStatus)}, imsi...)
	sgsCause := test_utils.ConstructDefaultIE(decode.IEISGsCause, 1)
	chunk = append(chunk, sgsCause...)
	erroneousMsg := test_utils.ConstructDefaultIE(decode.IEIErroneousMessage, 10)
	chunk = append(chunk, erroneousMsg...)
	msg, err := message.DecodeSGsAPStatus(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.Status{
		Imsi:             "111111",
		SgsCause:         sgsCause[2:],
		ErroneousMessage: erroneousMsg[2:],
	})
	assert.Equal(t, expectedMsg, msg)

	// Without IMSI
	chunk = append([]byte{byte(decode.SGsAPStatus)}, sgsCause...)
	chunk = append(chunk, erroneousMsg...)
	msg, err = message.DecodeSGsAPStatus(chunk)
	assert.NoError(t, err)
	expectedMsg, _ = ptypes.MarshalAny(&protos.Status{
		SgsCause:         sgsCause[2:],
		ErroneousMessage: erroneousMsg[2:],
	})
	assert.Equal(t, expectedMsg, msg)

	// Chunk too short
	msg, err = message.DecodeSGsAPStatus(chunk[:3])
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthSGsCause+decode.IELengthErroneousMessageMin,
		3,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPResetAck(t *testing.T) {
	// Successfully decode message
	mmeName := test_utils.ConstructDefaultMMEName()
	chunk := append([]byte{byte(decode.SGsAPResetAck)}, mmeName...)
	msg, err := message.DecodeSGsAPResetAck(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.ResetAck{
		MmeName: "abcdefghijabcdefghijabcdefghijabcdefghijabcdefghijabcde",
	})
	assert.Equal(t, expectedMsg, msg)

	// chunk too short
	msg, err = message.DecodeSGsAPResetAck(chunk[:2])
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, "+
			"actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthMMEName,
		decode.IELengthMessageType+decode.IELengthMMEName,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPResetIndication(t *testing.T) {
	// Successfully decode message
	vlrName := test_utils.ConstructDefaultVLRName()
	chunk := append([]byte{byte(decode.SGsAPResetIndication)}, vlrName...)
	msg, err := message.DecodeSGsAPResetIndication(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.ResetIndication{
		VlrName: "www.facebook.com",
	})
	assert.Equal(t, expectedMsg, msg)

	// chunk too short
	msg, err = message.DecodeSGsAPResetIndication(chunk[:2])
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, "+
			"actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthVLRNameMin,
		2,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPAlertAck(t *testing.T) {
	// Successfully decode message
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPAlertAck)}, imsi...)
	msg, err := message.DecodeSGsAPAlertAck(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.AlertAck{
		Imsi: "111111",
	})
	assert.Equal(t, expectedMsg, msg)

	// Chunk too short
	msg, err = message.DecodeSGsAPAlertAck(chunk[:3])
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin,
		decode.IELengthMessageType+decode.IELengthIMSIMax,
		len(chunk[:3]),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)

	// Chunk too long
	imsi, _ = test_utils.ConstructIMSI("111111111111111111111111")
	chunk = append([]byte{byte(decode.SGsAPAlertAck)}, imsi...)
	msg, err = message.DecodeSGsAPAlertAck(chunk)
	errorMsg = fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin,
		decode.IELengthMessageType+decode.IELengthIMSIMax,
		len(chunk),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)
}

func TestDecodeSGsAPAlertReject(t *testing.T) {
	// Successfully decode message
	imsi, _ := test_utils.ConstructIMSI("111111")
	chunk := append([]byte{byte(decode.SGsAPAlertReject)}, imsi...)
	sgsCause := test_utils.ConstructDefaultIE(decode.IEISGsCause, 1)
	chunk = append(chunk, sgsCause...)
	msg, err := message.DecodeSGsAPAlertReject(chunk)
	assert.NoError(t, err)
	expectedMsg, _ := ptypes.MarshalAny(&protos.AlertReject{
		Imsi:     "111111",
		SgsCause: sgsCause[2:],
	})
	assert.Equal(t, expectedMsg, msg)

	// Chunk too Short
	chunk = append([]byte{byte(decode.SGsAPAlertReject)}, imsi[:3]...)
	chunk = append(chunk, sgsCause...)
	msg, err = message.DecodeSGsAPAlertReject(chunk)
	errorMsg := fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin+decode.IELengthSGsCause,
		decode.IELengthMessageType+decode.IELengthIMSIMax+decode.IELengthSGsCause,
		len(chunk),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)

	// Chunk too long
	imsi, _ = test_utils.ConstructIMSI("11111111111111111111111111")
	chunk = append([]byte{byte(decode.SGsAPAlertReject)}, imsi...)
	chunk = append(chunk, sgsCause...)
	msg, err = message.DecodeSGsAPAlertReject(chunk)
	errorMsg = fmt.Sprintf(
		"wrong chunk size, min length: %d, max length: %d, actual length of the chunk: %d",
		decode.IELengthMessageType+decode.IELengthIMSIMin+decode.IELengthSGsCause,
		decode.IELengthMessageType+decode.IELengthIMSIMax+decode.IELengthSGsCause,
		len(chunk),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)

	// IMSI too short
	imsi, _ = test_utils.ConstructIMSI("11111")
	chunk = append([]byte{byte(decode.SGsAPAlertReject)}, imsi...)
	chunk = append(chunk, sgsCause...)
	chunk = append(chunk, sgsCause...)

	mandatoryFieldLength := decode.LengthIEI + decode.LengthLengthIndicator
	msg, err = message.DecodeSGsAPAlertReject(chunk)
	errorMsg = fmt.Sprintf(
		"failed to decode IMSI: wrong length indicator, \nmin value: %d, max value: %d, length indicator: %d",
		decode.IELengthIMSIMin-mandatoryFieldLength,
		decode.IELengthIMSIMax-mandatoryFieldLength,
		len(imsi)-mandatoryFieldLength,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)

	// IMSI too long
	imsi, _ = test_utils.ConstructIMSI("111111111111111111111")
	chunk = append([]byte{byte(decode.SGsAPAlertReject)}, imsi...)
	msg, err = message.DecodeSGsAPAlertReject(chunk)
	errorMsg = fmt.Sprintf(
		"failed to decode IMSI: wrong length indicator, \nmin value: %d, max value: %d, length indicator: %d",
		decode.IELengthIMSIMin-mandatoryFieldLength,
		decode.IELengthIMSIMax-mandatoryFieldLength,
		len(imsi)-mandatoryFieldLength,
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)

	// sgsCause too short
	imsi, _ = test_utils.ConstructIMSI("1111111111")
	chunk = append([]byte{byte(decode.SGsAPAlertReject)}, imsi...)
	chunk = append(chunk, sgsCause[:1]...)
	msg, err = message.DecodeSGsAPAlertReject(chunk)
	errorMsg = fmt.Sprintf(
		"failed to decode SGsCause: chunk too short, \nmin length of information element: %d, number of undecoded bytes: %d",
		decode.IELengthSGsCause,
		len(sgsCause[:1]),
	)
	assert.EqualError(t, err, errorMsg)
	assert.Equal(t, &any.Any{}, msg)

}
