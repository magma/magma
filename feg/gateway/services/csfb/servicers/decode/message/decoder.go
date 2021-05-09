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

package message

import (
	"errors"
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/decode/ie"

	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

type decoderImpl func([]byte) (*any.Any, error)

var decoderMap = map[decode.SGsMessageType]decoderImpl{
	decode.SGsAPLocationUpdateAccept:     DecodeSGsAPLocationUpdateAccept,
	decode.SGsAPLocationUpdateReject:     DecodeSGsAPLocationUpdateReject,
	decode.SGsAPIMSIDetachAck:            DecodeSGsAPIMSIDetachAck,
	decode.SGsAPPagingRequest:            DecodeSGsAPPagingRequest,
	decode.SGsAPEPSDetachAck:             DecodeSGsAPEPSDetachAck,
	decode.SGsAPAlertRequest:             DecodeSGsAPAlertRequest,
	decode.SGsAPDownlinkUnitdata:         DecodeSGsAPDownlinkUnitdata,
	decode.SGsAPMMInformationRequest:     DecodeSGsAPMMInformationRequest,
	decode.SGsAPReleaseRequest:           DecodeSGsAPReleaseRequest,
	decode.SGsAPServiceAbortRequest:      DecodeSGsAPServiceAbortRequest,
	decode.SGsAPStatus:                   DecodeSGsAPStatus,
	decode.SGsAPResetAck:                 DecodeSGsAPResetAck,
	decode.SGsAPResetIndication:          DecodeSGsAPResetIndication,
	decode.SGsAPAlertAck:                 DecodeSGsAPAlertAck,
	decode.SGsAPAlertReject:              DecodeSGsAPAlertReject,
	decode.SGsAPEPSDetachIndication:      DecodeSGsAPEPSDetachIndication,
	decode.SGsAPIMSIDetachIndication:     DecodeSGsAPIMSIDetachIndication,
	decode.SGsAPLocationUpdateRequest:    DecodeSGsAPLocationUpdateRequest,
	decode.SGsAPPagingReject:             DecodeSGsAPPagingReject,
	decode.SGsAPServiceRequest:           DecodeSGsAPServiceRequest,
	decode.SGsAPTMSIReallocationComplete: DecodeSGsAPTMSIReallocationComplete,
	decode.SGsAPUEActivityIndication:     DecodeSGsAPUEActivityIndication,
	decode.SGsAPUEUnreachable:            DecodeSGsAPUEUnreachable,
	decode.SGsAPUplinkUnitdata:           DecodeSGsAPUplinkUnitdata,
}

func SGsMessageDecoder(chunk []byte) (decode.SGsMessageType, *any.Any, error) {
	if decoderFunc, ok := decoderMap[decode.SGsMessageType(chunk[0])]; ok {
		glog.V(2).Infof(
			"Received message in bytes (hex format): % x",
			chunk,
		)
		glog.V(2).Infof(
			"Received message is being decoded as %s. ",
			decode.MsgTypeNameByCode[decode.SGsMessageType(chunk[0])],
		)
		marshalledMsg, err := decoderFunc(chunk)
		return decode.SGsMessageType(chunk[0]), marshalledMsg, err
	}
	return decode.SGsMessageType(chunk[0]), &any.Any{}, errors.New("unknown message type")
}

// DecodeSGsAPLocationUpdateAccept decodes the SGsAPLocationUpdateAccept message byte-by-byte
func DecodeSGsAPLocationUpdateAccept(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthLocationAreaIdentifier
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax + decode.IELengthLocationAreaIdentifier + decode.IELengthIMSIMax
	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType

	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += imsiLength

	minLength = decode.IELengthMessageType + imsiLength + decode.IELengthLocationAreaIdentifier
	err = validateMessageMinLength(chunk, minLength)
	if err != nil {
		return &any.Any{}, err
	}

	lai, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthLocationAreaIdentifier, decode.IEILocationAreaIdentifier)
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += decode.IELengthLocationAreaIdentifier

	mobileIdentity, err := readMobileIdentity(chunk, imsiLength, readIdx)
	if err != nil {
		return &any.Any{}, err
	}

	msg := protos.LocationUpdateAccept{
		Imsi:                   imsi,
		LocationAreaIdentifier: lai,
	}
	if mobileIdentity.IMSI != "" {
		msg.NewIMSITMSI = &protos.LocationUpdateAccept_NewImsi{NewImsi: mobileIdentity.IMSI}
	} else if mobileIdentity.TMSI != nil {
		msg.NewIMSITMSI = &protos.LocationUpdateAccept_NewTmsi{NewTmsi: mobileIdentity.TMSI}
	}

	marshalledMsg, err := ptypes.MarshalAny(&msg)
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

// DecodeSGsAPLocationUpdateReject decodes the SGsAPLocationUpdateReject message byte-by-byte
func DecodeSGsAPLocationUpdateReject(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.LengthRejectCause
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax + decode.LengthRejectCause + decode.IELengthLocationAreaIdentifier
	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType

	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += imsiLength

	lai, err := readLAI(chunk, imsiLength, readIdx+decode.LengthRejectCause)
	if err != nil {
		return &any.Any{}, err
	}

	rejectCause := chunk[readIdx : readIdx+decode.LengthRejectCause]

	marshalledMsg, err := ptypes.MarshalAny(&protos.LocationUpdateReject{
		Imsi:                   imsi,
		RejectCause:            rejectCause,
		LocationAreaIdentifier: lai,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

// DecodeSGsAPIMSIDetachAck decodes the SGsAPIMSIDetachAck message byte-by-byte
func DecodeSGsAPIMSIDetachAck(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax
	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, _, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.IMSIDetachAck{
		Imsi: imsi,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

// DecodeSGsAPMMInformationRequest decodes the SGsAPMMInformationRequest message byte-by-byte
func DecodeSGsAPMMInformationRequest(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthMMInformationMin
	err := validateMessageMinLength(chunk, minLength)
	if err != nil {
		return &any.Any{}, err
	}

	bytePtr := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[bytePtr:])
	if err != nil {
		return &any.Any{}, err
	}

	bytePtr += imsiLength
	mmInfo, _, err := ie.DecodeVariableLengthIE(chunk[bytePtr:], decode.IELengthMMInformationMin, decode.IEIMMInformation)
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.MMInformationRequest{
		Imsi:          imsi,
		MmInformation: mmInfo,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPPagingRequest(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthVLRNameMin + decode.IELengthServiceIndicator
	err := validateMessageMinLength(chunk, minLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}

	readIdx += imsiLength
	vlrName, vlrLength, err := ie.DecodeVariableLengthIE(chunk[readIdx:], decode.IELengthVLRNameMin, decode.IEIVLRName)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx += vlrLength
	serviceIndicator, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthServiceIndicator, decode.IEIServiceIndicator)
	if err != nil {
		return &any.Any{}, err
	}

	decodedMsg := &protos.PagingRequest{
		Imsi:             imsi,
		VlrName:          string(vlrName),
		ServiceIndicator: serviceIndicator,
	}

	// rest of fields are optional
	optionalFieldIEIOrder := []decode.InformationElementIdentifier{
		decode.IEITMSI,
		decode.IEICLI,
		decode.IEILocationAreaIdentifier,
		decode.IEIGlobalCNId,
		decode.IEISSCode,
		decode.IEILCSIndicator,
		decode.IEILCSClientIdentity,
		decode.IEIChannelNeeded,
		decode.IEIeMLPPPriority,
	}
	optionalFieldIdx := 0
	readIdx += decode.IELengthServiceIndicator
	for optionalFieldIdx < len(optionalFieldIEIOrder) && readIdx < len(chunk) {
		if optionalFieldIEIOrder[optionalFieldIdx] == decode.InformationElementIdentifier(chunk[readIdx]) {
			switch decode.InformationElementIdentifier(chunk[readIdx]) {
			case decode.IEITMSI:
				tmsi, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthTMSI, decode.IEITMSI)
				if err != nil {
					return &any.Any{}, err
				}
				decodedMsg.Tmsi = tmsi
				readIdx += decode.IELengthTMSI
			case decode.IEICLI:
				cli, cliLength, err := ie.DecodeLimitedLengthIE(chunk[readIdx:], decode.IELengthCLIMin, decode.IELengthCLIMax, decode.IEICLI)
				if err != nil {
					return &any.Any{}, err
				}
				decodedMsg.Cli = cli
				readIdx += cliLength
			case decode.IEILocationAreaIdentifier:
				lai, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthLocationAreaIdentifier, decode.IEILocationAreaIdentifier)
				if err != nil {
					return &any.Any{}, err
				}
				decodedMsg.LocationAreaIdentifier = lai
				readIdx += decode.IELengthLocationAreaIdentifier
			case decode.IEIGlobalCNId:
				globalCNId, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthGlobalCNId, decode.IEIGlobalCNId)
				if err != nil {
					return &any.Any{}, err
				}
				decodedMsg.GlobalCnId = globalCNId
				readIdx += decode.IELengthGlobalCNId
			case decode.IEISSCode:
				ssCode, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthSSCode, decode.IEISSCode)
				if err != nil {
					return &any.Any{}, err
				}
				decodedMsg.SsCode = ssCode
				readIdx += decode.IELengthSSCode
			case decode.IEILCSIndicator:
				lcsIndicator, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthLCSIndicator, decode.IEILCSIndicator)
				if err != nil {
					return &any.Any{}, err
				}
				decodedMsg.LcsIndicator = lcsIndicator
				readIdx += decode.IELengthLCSIndicator
			case decode.IEILCSClientIdentity:
				lcsCLientIdentity, lengthLCSClientIdentity, err := ie.DecodeVariableLengthIE(
					chunk[readIdx:],
					decode.IELengthLCSClientIdentityMin,
					decode.IEILCSClientIdentity,
				)
				if err != nil {
					return &any.Any{}, err
				}
				decodedMsg.LcsClientIdentity = lcsCLientIdentity
				readIdx += lengthLCSClientIdentity
			case decode.IEIChannelNeeded:
				channelNeeded, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthChannelNeeded, decode.IEIChannelNeeded)
				if err != nil {
					return &any.Any{}, err
				}
				decodedMsg.ChannelNeeded = channelNeeded
				readIdx += decode.IELengthChannelNeeded
			case decode.IEIeMLPPPriority:
				emlppPriority, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthEMLPPPriority, decode.IEIeMLPPPriority)
				if err != nil {
					return &any.Any{}, err
				}
				decodedMsg.EmlppPriority = emlppPriority
				readIdx += decode.IELengthEMLPPPriority
			}
		}
		optionalFieldIdx += 1
	}

	if readIdx < len(chunk) {
		return &any.Any{}, errors.New("tried all possible IE but still some bytes undecoded")
	}

	marshalledMsg, err := ptypes.MarshalAny(decodedMsg)
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPEPSDetachAck(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax
	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, _, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.EPSDetachAck{
		Imsi: imsi,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPAlertRequest(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax
	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, _, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.AlertRequest{
		Imsi: imsi,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPDownlinkUnitdata(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthNASMessageContainerMin
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax + decode.IELengthNASMessageContainerMax
	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}

	readIdx += imsiLength
	nasMessageContainer, _, err := ie.DecodeLimitedLengthIE(
		chunk[readIdx:],
		decode.IELengthNASMessageContainerMin,
		decode.IELengthNASMessageContainerMax,
		decode.IEINASMessageContainer,
	)
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.DownlinkUnitdata{
		Imsi:                imsi,
		NasMessageContainer: nasMessageContainer,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPReleaseRequest(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax + decode.IELengthSGsCause
	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}

	decodedMsg := &protos.ReleaseRequest{
		Imsi: imsi,
	}

	readIdx += imsiLength
	if readIdx < len(chunk) {
		sgsCause, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthSGsCause, decode.IEISGsCause)
		if err != nil {
			return &any.Any{}, err
		}
		decodedMsg.SgsCause = sgsCause
	}

	marshalledMsg, err := ptypes.MarshalAny(decodedMsg)
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPServiceAbortRequest(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax
	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, _, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.ServiceAbortRequest{
		Imsi: imsi,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPStatus(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthSGsCause + decode.IELengthErroneousMessageMin
	err := validateMessageMinLength(chunk, minLength)
	if err != nil {
		return &any.Any{}, err
	}

	decodedMsg := &protos.Status{}

	readIdx := decode.IELengthMessageType
	if decode.InformationElementIdentifier(chunk[readIdx]) == decode.IEIIMSI {
		imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
		if err != nil {
			return &any.Any{}, err
		}
		decodedMsg.Imsi = imsi
		readIdx += imsiLength
	}

	sgsCause, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthSGsCause, decode.IEISGsCause)
	if err != nil {
		return &any.Any{}, err
	}
	decodedMsg.SgsCause = sgsCause

	readIdx += decode.IELengthSGsCause
	erroneousMsg, _, err := ie.DecodeVariableLengthIE(
		chunk[readIdx:],
		decode.IELengthErroneousMessageMin,
		decode.IEIErroneousMessage,
	)
	if err != nil {
		return &any.Any{}, err
	}
	decodedMsg.ErroneousMessage = erroneousMsg

	marshalledMsg, err := ptypes.MarshalAny(decodedMsg)
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPResetAck(chunk []byte) (*any.Any, error) {
	length := decode.IELengthMessageType + decode.IELengthMMEName
	err := validateMessageLength(chunk, length, length)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	mmeName, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthMMEName, decode.IEIMMEName)
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.ResetAck{
		MmeName: string(mmeName),
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}
	return marshalledMsg, nil
}

func DecodeSGsAPResetIndication(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthVLRNameMin
	err := validateMessageMinLength(chunk, minLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	vlrName, _, err := ie.DecodeVariableLengthIE(chunk[readIdx:], decode.IELengthVLRNameMin, decode.IEIVLRName)
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.ResetIndication{
		VlrName: string(vlrName),
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}
	return marshalledMsg, nil
}

// DecodeSGsAPAlertAck decodes the SGsAPAlertAck message byte-by-byte
func DecodeSGsAPAlertAck(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax

	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, _, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.AlertAck{
		Imsi: imsi,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

// DecodeSGsAPAlertReject decodes the SGsAPAlertReject message byte-by-byte
func DecodeSGsAPAlertReject(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthSGsCause
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax + decode.IELengthSGsCause

	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += imsiLength

	sgsCause, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthSGsCause, decode.IEISGsCause)
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.AlertReject{
		Imsi:     imsi,
		SgsCause: sgsCause,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPEPSDetachIndication(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthMMEName + decode.IELengthIMSIDetachFromEPSServiceType
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax + decode.IELengthMMEName + decode.IELengthIMSIDetachFromEPSServiceType

	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += imsiLength

	mmeName, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthMMEName, decode.IEIMMEName)
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += decode.IELengthMMEName

	imsDetachFromEPSServiceType, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthIMSIDetachFromEPSServiceType,
		decode.IEIIMSIDetachFromEPSServiceType)
	if err != nil {
		return &any.Any{}, err
	}
	marshalledMsg, err := ptypes.MarshalAny(&protos.EPSDetachIndication{
		Imsi:                         imsi,
		MmeName:                      string(mmeName),
		ImsiDetachFromEpsServiceType: imsDetachFromEPSServiceType,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPIMSIDetachIndication(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthMMEName + decode.IELengthIMSIDetachFromNonEPSServiceType
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax + decode.IELengthMMEName + decode.IELengthIMSIDetachFromNonEPSServiceType

	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += imsiLength

	mmeName, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthMMEName, decode.IEIMMEName)
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += decode.IELengthMMEName

	imsDetachFromNonEPSServiceType, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthIMSIDetachFromNonEPSServiceType,
		decode.IEIIMSIDetachFromNonEPSServiceType)
	if err != nil {
		return &any.Any{}, err
	}
	marshalledMsg, err := ptypes.MarshalAny(&protos.IMSIDetachIndication{
		Imsi:                            imsi,
		MmeName:                         string(mmeName),
		ImsiDetachFromNonEpsServiceType: imsDetachFromNonEPSServiceType,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPLocationUpdateRequest(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthMMEName +
		decode.IELengthEPSLocationUpdateType + decode.IELengthLocationAreaIdentifier
	err := validateMessageMinLength(chunk, minLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += imsiLength

	mmeName, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthMMEName,
		decode.IEIMMEName)
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += decode.IELengthMMEName

	epsLocationUpdateType, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthEPSLocationUpdateType,
		decode.IEIEPSLocationUpdateType)
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += decode.IELengthEPSLocationUpdateType

	newLocationAreaIdentifier, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthLocationAreaIdentifier,
		decode.IEILocationAreaIdentifier)
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += decode.IELengthLocationAreaIdentifier

	locationUpdateRequest := &protos.LocationUpdateRequest{
		Imsi:                      imsi,
		MmeName:                   string(mmeName),
		EpsLocationUpdateType:     epsLocationUpdateType,
		NewLocationAreaIdentifier: newLocationAreaIdentifier,
	}

	// rest of fields are optional
	optionalFieldIEIOrder := []decode.InformationElementIdentifier{
		decode.IEILocationAreaIdentifier,
		decode.IEITMSIStatus,
		decode.IEIIMEISV,
		decode.IEITAI,
		decode.IEIEUTRANCellGlobalIdentity,
	}
	optionalFieldIdx := 0
	for optionalFieldIdx < len(optionalFieldIEIOrder) && readIdx < len(chunk) {
		if optionalFieldIEIOrder[optionalFieldIdx] == decode.InformationElementIdentifier(chunk[readIdx]) {
			switch decode.InformationElementIdentifier(chunk[readIdx]) {
			case decode.IEILocationAreaIdentifier:
				{
					oldLocationAreaIdentifier, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthLocationAreaIdentifier,
						decode.IEILocationAreaIdentifier)
					if err != nil {
						return &any.Any{}, err
					}
					locationUpdateRequest.OldLocationAreaIdentifier = oldLocationAreaIdentifier
					readIdx += decode.IELengthLocationAreaIdentifier
				}
			case decode.IEITMSIStatus:
				{
					tmsiStatus, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthTMSIStatus,
						decode.IEITMSIStatus)
					if err != nil {
						return &any.Any{}, err
					}
					locationUpdateRequest.TmsiStatus = tmsiStatus
					readIdx += decode.IELengthTMSIStatus
				}
			case decode.IEIIMEISV:
				{
					imeisv, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthIMEISV,
						decode.IEIIMEISV)
					if err != nil {
						return &any.Any{}, err
					}
					locationUpdateRequest.Imeisv = imeisv
					readIdx += decode.IELengthIMEISV
				}
			case decode.IEITAI:
				{
					tai, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthTAI, decode.IEITAI)
					if err != nil {
						return &any.Any{}, err
					}
					locationUpdateRequest.Tai = tai
					readIdx += decode.IELengthTAI
				}
			case decode.IEIEUTRANCellGlobalIdentity:
				{
					eCgi, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthEUTRANCellGlobalIdentity,
						decode.IEIEUTRANCellGlobalIdentity)
					if err != nil {
						return &any.Any{}, err
					}
					locationUpdateRequest.ECgi = eCgi
					readIdx += decode.IELengthEUTRANCellGlobalIdentity
				}
			}
		}
		optionalFieldIdx += 1
	}
	if readIdx < len(chunk) {
		return &any.Any{}, errors.New("tried all possible IE but still some bytes undecoded")
	}

	marshalledMsg, err := ptypes.MarshalAny(locationUpdateRequest)
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPPagingReject(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthSGsCause
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax + decode.IELengthSGsCause

	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += imsiLength

	sgsCause, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthSGsCause, decode.IEISGsCause)
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.PagingReject{
		Imsi:     imsi,
		SgsCause: sgsCause,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPServiceRequest(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthServiceIndicator
	err := validateMessageMinLength(chunk, minLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += imsiLength

	serviceIndicator, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthServiceIndicator,
		decode.IEIServiceIndicator)
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += decode.IELengthServiceIndicator

	serviceRequest := &protos.ServiceRequest{
		Imsi:             imsi,
		ServiceIndicator: serviceIndicator,
	}

	// rest of fields are optional
	optionalFieldIEIOrder := []decode.InformationElementIdentifier{
		decode.IEIIMEISV,
		decode.IEIUETimeZone,
		decode.IEIMobileStationClassmark2,
		decode.IEITAI,
		decode.IEIEUTRANCellGlobalIdentity,
		decode.IEIUEEMMMode,
	}
	optionalFieldIdx := 0
	for optionalFieldIdx < len(optionalFieldIEIOrder) && readIdx < len(chunk) {
		if optionalFieldIEIOrder[optionalFieldIdx] == decode.InformationElementIdentifier(chunk[readIdx]) {
			switch decode.InformationElementIdentifier(chunk[readIdx]) {
			case decode.IEIIMEISV:
				{
					imeIsv, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthIMEISV,
						decode.IEIIMEISV)
					if err != nil {
						return &any.Any{}, err
					}
					serviceRequest.Imeisv = imeIsv
					readIdx += decode.IELengthIMEISV
				}
			case decode.IEIUETimeZone:
				{
					ueTimeZone, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthUETimeZone,
						decode.IEIUETimeZone)
					if err != nil {
						return &any.Any{}, err
					}
					serviceRequest.UeTimeZone = ueTimeZone
					readIdx += decode.IELengthUETimeZone
				}
			case decode.IEIMobileStationClassmark2:
				{
					mobileStationClassmark2, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthMobileStationClassmark2,
						decode.IEIMobileStationClassmark2)
					if err != nil {
						return &any.Any{}, err
					}
					serviceRequest.MobileStationClassmark2 = mobileStationClassmark2
					readIdx += decode.IELengthMobileStationClassmark2
				}
			case decode.IEITAI:
				{
					tai, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthTAI, decode.IEITAI)
					if err != nil {
						return &any.Any{}, err
					}
					serviceRequest.Tai = tai
					readIdx += decode.IELengthTAI
				}
			case decode.IEIEUTRANCellGlobalIdentity:
				{
					eCgi, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthEUTRANCellGlobalIdentity,
						decode.IEIEUTRANCellGlobalIdentity)
					if err != nil {
						return &any.Any{}, err
					}
					serviceRequest.ECgi = eCgi
					readIdx += decode.IELengthEUTRANCellGlobalIdentity
				}
			case decode.IEIUEEMMMode:
				{
					ueEmmMode, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthUEEMMMode,
						decode.IEIUEEMMMode)
					if err != nil {
						return &any.Any{}, err
					}
					serviceRequest.UeEmmMode = ueEmmMode
					readIdx += decode.IELengthUEEMMMode
				}
			}
		}
		optionalFieldIdx += 1
	}
	if readIdx < len(chunk) {
		return &any.Any{}, errors.New("tried all possible IE but still some bytes undecoded")
	}

	marshalledMsg, err := ptypes.MarshalAny(serviceRequest)
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPTMSIReallocationComplete(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax

	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, _, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.TMSIReallocationComplete{
		Imsi: imsi,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPUEActivityIndication(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax

	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, _, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.UEActivityIndication{
		Imsi: imsi,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPUEUnreachable(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthSGsCause
	maxLength := decode.IELengthMessageType + decode.IELengthIMSIMax + decode.IELengthSGsCause

	err := validateMessageLength(chunk, minLength, maxLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += imsiLength

	sgsCause, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthSGsCause, decode.IEISGsCause)
	if err != nil {
		return &any.Any{}, err
	}

	marshalledMsg, err := ptypes.MarshalAny(&protos.UEUnreachable{
		Imsi:     imsi,
		SgsCause: sgsCause,
	})
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}

func DecodeSGsAPUplinkUnitdata(chunk []byte) (*any.Any, error) {
	minLength := decode.IELengthMessageType + decode.IELengthIMSIMin + decode.IELengthNASMessageContainerMin
	err := validateMessageMinLength(chunk, minLength)
	if err != nil {
		return &any.Any{}, err
	}

	readIdx := decode.IELengthMessageType
	imsi, imsiLength, err := ie.DecodeIMSI(chunk[readIdx:])
	if err != nil {
		return &any.Any{}, err
	}
	readIdx += imsiLength
	nasMessageContainer, nasMessageContainerLength, err := ie.DecodeLimitedLengthIE(chunk[readIdx:],
		decode.IELengthNASMessageContainerMin, decode.IELengthNASMessageContainerMax,
		decode.IEINASMessageContainer)
	if err != nil {
		return &any.Any{}, err
	}

	sgUplinkUnitData := &protos.UplinkUnitdata{
		Imsi:                imsi,
		NasMessageContainer: nasMessageContainer,
	}
	readIdx += nasMessageContainerLength

	// rest of fields are optional
	optionalFieldIEIOrder := []decode.InformationElementIdentifier{
		decode.IEIIMEISV,
		decode.IEIUETimeZone,
		decode.IEIMobileStationClassmark2,
		decode.IEITAI,
		decode.IEIEUTRANCellGlobalIdentity,
	}
	optionalFieldIdx := 0
	for optionalFieldIdx < len(optionalFieldIEIOrder) && readIdx < len(chunk) {
		if optionalFieldIEIOrder[optionalFieldIdx] == decode.InformationElementIdentifier(chunk[readIdx]) {
			switch decode.InformationElementIdentifier(chunk[readIdx]) {
			case decode.IEIIMEISV:
				{
					imeIsv, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthIMEISV,
						decode.IEIIMEISV)
					if err != nil {
						return &any.Any{}, err
					}
					sgUplinkUnitData.Imeisv = imeIsv
					readIdx += decode.IELengthIMEISV
				}
			case decode.IEIUETimeZone:
				{
					ueTimeZone, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthUETimeZone,
						decode.IEIUETimeZone)
					if err != nil {
						return &any.Any{}, err
					}
					sgUplinkUnitData.UeTimeZone = ueTimeZone
					readIdx += decode.IELengthUETimeZone
				}
			case decode.IEIMobileStationClassmark2:
				{
					mobileStationClassmark2, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthMobileStationClassmark2,
						decode.IEIMobileStationClassmark2)
					if err != nil {
						return &any.Any{}, err
					}
					sgUplinkUnitData.MobileStationClassmark2 = mobileStationClassmark2
					readIdx += decode.IELengthMobileStationClassmark2
				}
			case decode.IEITAI:
				{
					tai, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthTAI, decode.IEITAI)
					if err != nil {
						return &any.Any{}, err
					}
					sgUplinkUnitData.Tai = tai
					readIdx += decode.IELengthTAI
				}
			case decode.IEIEUTRANCellGlobalIdentity:
				{
					eCgi, err := ie.DecodeFixedLengthIE(chunk[readIdx:], decode.IELengthEUTRANCellGlobalIdentity,
						decode.IEIEUTRANCellGlobalIdentity)
					if err != nil {
						return &any.Any{}, err
					}
					sgUplinkUnitData.ECgi = eCgi
					readIdx += decode.IELengthEUTRANCellGlobalIdentity
				}
			}
		}
		optionalFieldIdx += 1
	}
	if readIdx < len(chunk) {
		return &any.Any{}, errors.New("tried all possible IE but still some bytes undecoded")
	}

	marshalledMsg, err := ptypes.MarshalAny(sgUplinkUnitData)
	if err != nil {
		return &any.Any{}, fmt.Errorf("Error marshaling SGs message to Any: %s", err)
	}

	return marshalledMsg, nil
}
