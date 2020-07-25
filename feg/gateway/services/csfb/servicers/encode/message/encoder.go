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
	"fmt"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/csfb/servicers/decode"
	"magma/feg/gateway/services/csfb/servicers/encode/ie"
)

func EncodeSGsAPAlertAck(message *protos.AlertAck) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPAlertAck, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	return encodedMsg, nil
}

func EncodeSGsAPAlertReject(message *protos.AlertReject) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPAlertReject, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	encodedSGsCause, err := ie.EncodeFixedLengthIE(decode.IEISGsCause, decode.IELengthSGsCause, message.SgsCause)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedSGsCause...)

	return encodedMsg, nil
}

func EncodeSGsAPAlertRequest(message *protos.AlertRequest) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPAlertRequest, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	return encodedMsg, nil
}

func EncodeSGsAPDownlinkUnitdata(message *protos.DownlinkUnitdata) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPDownlinkUnitdata, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	encodedNasMessageContainer, err := ie.EncodeLimitedLengthIE(
		decode.IEINASMessageContainer,
		decode.IELengthNASMessageContainerMin,
		decode.IELengthNASMessageContainerMax,
		message.NasMessageContainer,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedNasMessageContainer...)

	return encodedMsg, nil
}

func EncodeSGsAPEPSDetachAck(message *protos.EPSDetachAck) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(
		decode.SGsAPEPSDetachAck,
		message.Imsi,
	)
	if err != nil {
		return []byte{}, err
	}

	return encodedMsg, nil
}

func EncodeSGsAPEPSDetachIndication(message *protos.EPSDetachIndication) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPEPSDetachIndication, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	encodedMMEName, err := ie.EncodeMMEName(message.MmeName)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedMMEName...)

	encodedServiceType, err := ie.EncodeFixedLengthIE(
		decode.IEIIMSIDetachFromEPSServiceType,
		decode.IELengthIMSIDetachFromEPSServiceType,
		message.ImsiDetachFromEpsServiceType,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedServiceType...)

	return encodedMsg, nil
}

func EncodeSGsAPIMSIDetachAck(message *protos.IMSIDetachAck) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(
		decode.SGsAPIMSIDetachAck,
		message.Imsi,
	)
	if err != nil {
		return []byte{}, err
	}

	return encodedMsg, nil
}

func EncodeSGsAPIMSIDetachIndication(message *protos.IMSIDetachIndication) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPIMSIDetachIndication, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	encodedMMEName, err := ie.EncodeMMEName(message.MmeName)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedMMEName...)

	encodedServiceType, err := ie.EncodeFixedLengthIE(
		decode.IEIIMSIDetachFromNonEPSServiceType,
		decode.IELengthIMSIDetachFromNonEPSServiceType,
		message.ImsiDetachFromNonEpsServiceType,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedServiceType...)

	return encodedMsg, nil
}

func EncodeSGsAPLocationUpdateAccept(message *protos.LocationUpdateAccept) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(
		decode.SGsAPLocationUpdateAccept,
		message.Imsi,
	)
	if err != nil {
		return []byte{}, err
	}

	encodedLAI, err := ie.EncodeFixedLengthIE(
		decode.IEILocationAreaIdentifier,
		decode.IELengthLocationAreaIdentifier,
		message.LocationAreaIdentifier,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedLAI...)

	var encodedNewIMSITMSI []byte
	switch t := message.NewIMSITMSI.(type) {
	case *protos.LocationUpdateAccept_NewImsi:
		encodedNewIMSITMSI, err = ie.EncodeIMSI(message.GetNewImsi())
	case *protos.LocationUpdateAccept_NewTmsi:
		encodedNewIMSITMSI, err = ie.EncodeFixedLengthIE(
			decode.IEITMSI,
			decode.IELengthTMSI,
			message.GetNewTmsi(),
		)
	case nil:
		err = nil
	default:
		err = fmt.Errorf("Profile.Avatar has unexpected type %T", t)
	}
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedNewIMSITMSI...)

	return encodedMsg, nil
}

func EncodeSGsAPLocationUpdateReject(message *protos.LocationUpdateReject) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(
		decode.SGsAPLocationUpdateReject,
		message.Imsi,
	)
	if err != nil {
		return []byte{}, err
	}

	encodedCause, err := ie.EncodeFixedLengthIE(
		decode.IEIRejectCause,
		decode.IELengthRejectCause,
		message.RejectCause,
	)
	if err != nil {
		return []byte{}, nil
	}
	encodedMsg = append(encodedMsg, encodedCause...)

	if len(message.LocationAreaIdentifier) != 0 {
		encodedLAI, err := ie.EncodeFixedLengthIE(
			decode.IEILocationAreaIdentifier,
			decode.IELengthLocationAreaIdentifier,
			message.LocationAreaIdentifier,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedLAI...)
	}

	return encodedMsg, nil
}

func EncodeSGsAPLocationUpdateRequest(message *protos.LocationUpdateRequest) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPLocationUpdateRequest, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	encodedMMEName, err := ie.EncodeMMEName(message.MmeName)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedMMEName...)

	encodedLocationUpdateType, err := ie.EncodeFixedLengthIE(
		decode.IEIEPSLocationUpdateType,
		decode.IELengthEPSLocationUpdateType,
		message.EpsLocationUpdateType,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedLocationUpdateType...)

	encodedNewLAI, err := ie.EncodeFixedLengthIE(
		decode.IEILocationAreaIdentifier,
		decode.IELengthLocationAreaIdentifier,
		message.NewLocationAreaIdentifier,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedNewLAI...)

	if len(message.OldLocationAreaIdentifier) != 0 {
		encodedOldLAI, err := ie.EncodeFixedLengthIE(
			decode.IEILocationAreaIdentifier,
			decode.IELengthLocationAreaIdentifier,
			message.NewLocationAreaIdentifier,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedOldLAI...)
	}

	if len(message.TmsiStatus) != 0 {
		encodedTMSIStatus, err := ie.EncodeFixedLengthIE(
			decode.IEITMSIStatus,
			decode.IELengthTMSIStatus,
			message.TmsiStatus,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedTMSIStatus...)
	}

	if len(message.Imeisv) != 0 {
		encodedIMEISV, err := ie.EncodeFixedLengthIE(
			decode.IEIIMEISV,
			decode.IELengthIMEISV,
			message.Imeisv,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedIMEISV...)
	}

	if len(message.Tai) != 0 {
		encodedTAI, err := ie.EncodeFixedLengthIE(
			decode.IEITAI,
			decode.IELengthTAI,
			message.Tai,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedTAI...)
	}

	if len(message.ECgi) != 0 {
		encodedECGI, err := ie.EncodeFixedLengthIE(
			decode.IEIEUTRANCellGlobalIdentity,
			decode.IELengthEUTRANCellGlobalIdentity,
			message.ECgi,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedECGI...)
	}

	return encodedMsg, nil
}

func EncodeSGsAPMMInformationRequest(message *protos.MMInformationRequest) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPMMInformationRequest, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	mmInfo, err := ie.EncodeVariableLengthIE(
		decode.IEIMMInformation,
		decode.IELengthMMInformationMin,
		message.MmInformation,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, mmInfo...)

	return encodedMsg, nil
}

func EncodeSGsAPPagingReject(message *protos.PagingReject) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPPagingReject, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	encodedSGsCause, err := ie.EncodeFixedLengthIE(decode.IEISGsCause, decode.IELengthSGsCause, message.SgsCause)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedSGsCause...)

	return encodedMsg, nil
}

func EncodeSGsAPPagingRequest(message *protos.PagingRequest) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(
		decode.SGsAPPagingRequest,
		message.Imsi,
	)
	if err != nil {
		return []byte{}, err
	}

	encodedVLRName, err := ie.EncodeVariableLengthIE(
		decode.IEIVLRName,
		decode.IELengthVLRNameMin,
		[]byte(message.VlrName),
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedVLRName...)

	encodedServiceIndicator, err := ie.EncodeFixedLengthIE(
		decode.IEIServiceIndicator,
		decode.IELengthServiceIndicator,
		message.ServiceIndicator,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedServiceIndicator...)

	if len(message.Tmsi) != 0 {
		encodedTMSI, err := ie.EncodeFixedLengthIE(
			decode.IEITMSI,
			decode.IELengthTMSI,
			message.Tmsi,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedTMSI...)
	}

	if len(message.Cli) != 0 {
		encodedCLI, err := ie.EncodeLimitedLengthIE(
			decode.IEICLI,
			decode.IELengthCLIMin,
			decode.IELengthCLIMax,
			message.Cli,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedCLI...)
	}

	if len(message.LocationAreaIdentifier) != 0 {
		encodedLAI, err := ie.EncodeFixedLengthIE(
			decode.IEILocationAreaIdentifier,
			decode.IELengthLocationAreaIdentifier,
			message.LocationAreaIdentifier,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedLAI...)
	}

	if len(message.GlobalCnId) != 0 {
		encodedGlobalCnId, err := ie.EncodeFixedLengthIE(
			decode.IEIGlobalCNId,
			decode.IELengthGlobalCNId,
			message.GlobalCnId,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedGlobalCnId...)
	}

	if len(message.SsCode) != 0 {
		encodedSSCode, err := ie.EncodeFixedLengthIE(
			decode.IEISSCode,
			decode.IELengthSSCode,
			message.SsCode,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedSSCode...)
	}

	if len(message.LcsIndicator) != 0 {
		encodedLCSIndicator, err := ie.EncodeFixedLengthIE(
			decode.IEILCSIndicator,
			decode.IELengthLCSIndicator,
			message.LcsIndicator,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedLCSIndicator...)
	}

	if len(message.LcsClientIdentity) != 0 {
		encodedLCSClientID, err := ie.EncodeVariableLengthIE(
			decode.IEILCSClientIdentity,
			decode.IELengthLCSClientIdentityMin,
			message.LcsClientIdentity,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedLCSClientID...)
	}

	if len(message.ChannelNeeded) != 0 {
		encodedChannelNeeded, err := ie.EncodeFixedLengthIE(
			decode.IEIChannelNeeded,
			decode.IELengthChannelNeeded,
			message.ChannelNeeded,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedChannelNeeded...)
	}

	if len(message.EmlppPriority) != 0 {
		encodedEMLPPPriority, err := ie.EncodeFixedLengthIE(
			decode.IEIeMLPPPriority,
			decode.IELengthEMLPPPriority,
			message.EmlppPriority,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedEMLPPPriority...)
	}

	return encodedMsg, nil
}

func EncodeSGsAPReleaseRequest(message *protos.ReleaseRequest) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPReleaseRequest, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	encodedSGsCause, err := ie.EncodeFixedLengthIE(
		decode.IEISGsCause,
		decode.IELengthSGsCause,
		message.SgsCause,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedSGsCause...)

	return encodedMsg, nil
}

func EncodeSGsAPResetAck(message *protos.ResetAck) ([]byte, error) {
	encodedMsg := []byte{byte(decode.SGsAPResetAck)}
	encodedMMEName, err := ie.EncodeMMEName(message.MmeName)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedMMEName...)

	return encodedMsg, nil
}

func EncodeSGsAPResetIndication(message *protos.ResetIndication) ([]byte, error) {
	encodedMsg := []byte{byte(decode.SGsAPResetIndication)}
	encodedMMEName, err := ie.EncodeMMEName(message.MmeName)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedMMEName...)

	return encodedMsg, nil
}

func EncodeSGsAPServiceAbortRequest(message *protos.ServiceAbortRequest) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPServiceAbortRequest, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	return encodedMsg, nil
}

func EncodeSGsAPServiceRequest(message *protos.ServiceRequest) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPServiceRequest, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	encodedServiceIndicator, err := ie.EncodeFixedLengthIE(
		decode.IEIServiceIndicator,
		decode.IELengthServiceIndicator,
		message.ServiceIndicator,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedServiceIndicator...)

	if len(message.Imeisv) != 0 {
		encodedIMEISV, err := ie.EncodeFixedLengthIE(
			decode.IEIIMEISV,
			decode.IELengthIMEISV,
			message.Imeisv,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedIMEISV...)
	}

	if len(message.UeTimeZone) != 0 {
		encodedUETimeZone, err := ie.EncodeFixedLengthIE(
			decode.IEIUETimeZone,
			decode.IELengthUETimeZone,
			message.UeTimeZone,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedUETimeZone...)
	}

	if len(message.MobileStationClassmark2) != 0 {
		encodedMSC2, err := ie.EncodeFixedLengthIE(
			decode.IEIMobileStationClassmark2,
			decode.IELengthMobileStationClassmark2,
			message.MobileStationClassmark2,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedMSC2...)
	}

	if len(message.Tai) != 0 {
		encodedTAI, err := ie.EncodeFixedLengthIE(
			decode.IEITAI,
			decode.IELengthTAI,
			message.Tai,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedTAI...)
	}

	if len(message.ECgi) != 0 {
		encodedECGI, err := ie.EncodeFixedLengthIE(
			decode.IEIEUTRANCellGlobalIdentity,
			decode.IELengthEUTRANCellGlobalIdentity,
			message.ECgi,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedECGI...)
	}

	if len(message.UeEmmMode) != 0 {
		encodedUEEMMMode, err := ie.EncodeFixedLengthIE(
			decode.IEIUEEMMMode,
			decode.IELengthUEEMMMode,
			message.UeEmmMode,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedUEEMMMode...)
	}

	return encodedMsg, nil
}

func EncodeSGsAPStatus(message *protos.Status) ([]byte, error) {
	encodedMsg := []byte{byte(decode.SGsAPStatus)}

	if message.Imsi != "" {
		encodedIMSI, err := ie.EncodeIMSI(message.Imsi)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedIMSI...)
	}

	encodedSGsCause, err := ie.EncodeFixedLengthIE(decode.IEISGsCause, decode.IELengthSGsCause, message.SgsCause)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedSGsCause...)

	encodedErroneousMsg, err := ie.EncodeVariableLengthIE(
		decode.IEIErroneousMessage,
		decode.IELengthErroneousMessageMin,
		message.ErroneousMessage,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedErroneousMsg...)

	return encodedMsg, nil
}

func EncodeSGsAPTMSIReallocationComplete(message *protos.TMSIReallocationComplete) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPTMSIReallocationComplete, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	return encodedMsg, nil
}

func EncodeSGsAPUEActivityIndication(message *protos.UEActivityIndication) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPUEActivityIndication, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	return encodedMsg, nil
}

func EncodeSGsAPUEUnreachable(message *protos.UEUnreachable) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPUEUnreachable, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	encodedSGsCause, err := ie.EncodeFixedLengthIE(decode.IEISGsCause, decode.IELengthSGsCause, message.SgsCause)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedSGsCause...)

	return encodedMsg, nil
}

func EncodeSGsAPUplinkUnitdata(message *protos.UplinkUnitdata) ([]byte, error) {
	encodedMsg, err := encodeMessageTypeAndIMSI(decode.SGsAPUplinkUnitdata, message.Imsi)
	if err != nil {
		return []byte{}, err
	}

	encodedNASMessageContainer, err := ie.EncodeLimitedLengthIE(
		decode.IEINASMessageContainer,
		decode.IELengthNASMessageContainerMin,
		decode.IELengthNASMessageContainerMax,
		message.NasMessageContainer,
	)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedNASMessageContainer...)

	if len(message.Imeisv) != 0 {
		encodedIMEISV, err := ie.EncodeFixedLengthIE(
			decode.IEIIMEISV,
			decode.IELengthIMEISV,
			message.Imeisv,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedIMEISV...)
	}

	if len(message.UeTimeZone) != 0 {
		encodedUETimeZone, err := ie.EncodeFixedLengthIE(
			decode.IEIUETimeZone,
			decode.IELengthUETimeZone,
			message.UeTimeZone,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedUETimeZone...)
	}

	if len(message.MobileStationClassmark2) != 0 {
		encodedMSC2, err := ie.EncodeFixedLengthIE(
			decode.IEIMobileStationClassmark2,
			decode.IELengthMobileStationClassmark2,
			message.MobileStationClassmark2,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedMSC2...)
	}

	if len(message.Tai) != 0 {
		encodedTAI, err := ie.EncodeFixedLengthIE(
			decode.IEITAI,
			decode.IELengthTAI,
			message.Tai,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedTAI...)
	}

	if len(message.ECgi) != 0 {
		encodedECGI, err := ie.EncodeFixedLengthIE(
			decode.IEIEUTRANCellGlobalIdentity,
			decode.IELengthEUTRANCellGlobalIdentity,
			message.ECgi,
		)
		if err != nil {
			return []byte{}, err
		}
		encodedMsg = append(encodedMsg, encodedECGI...)
	}

	return encodedMsg, nil
}

func encodeMessageTypeAndIMSI(msgType decode.SGsMessageType, imsi string) ([]byte, error) {
	encodedMsg := []byte{byte(msgType)}
	encodedIMSI, err := ie.EncodeIMSI(imsi)
	if err != nil {
		return []byte{}, err
	}
	encodedMsg = append(encodedMsg, encodedIMSI...)

	return encodedMsg, nil
}
