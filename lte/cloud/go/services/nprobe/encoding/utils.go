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

package encoding

import (
	"encoding/asn1"
	"encoding/binary"
	"encoding/hex"
	"net"
	"strconv"
	"strings"
	"time"

	"magma/lte/cloud/go/services/nprobe"
	"magma/orc8r/cloud/go/services/eventd/obsidian/models"
)

func convertUint64ToBytes(v uint64) []byte {
	o := make([]byte, 8)
	binary.BigEndian.PutUint64(o, v)
	return o
}

func convertUint32ToBytes(v uint32) []byte {
	o := make([]byte, 4)
	binary.BigEndian.PutUint32(o, v)
	return o
}

// encodeAPN encodes APN to a byte sequence as specified in TS TS 29.274
func encodeAPN(apn string) []byte {
	encodedAPN := []byte{}
	if alen := len(apn); alen > 0 {
		encodedAPN = append(encodedAPN, byte(alen))
		encodedAPN = append(encodedAPN, []byte(apn)...)
	}
	return encodedAPN
}

// encodeUserLocation converts user location encodingfrom TS 29.061 to TS 29.274
func encodeUserLocation(s string) []byte {
	ulStr := strings.Join(strings.Fields(s), "")
	userLocation, err := hex.DecodeString(ulStr)
	if err != nil || len(userLocation) == 0 {
		return []byte{}
	}

	encodedUL := []byte{}
	if userLocation[0] == 0x82 {
		encodedUL = append(encodedUL, byte(0x18)) // TAI(8)+ECGI(16)
		encodedUL = append(encodedUL, userLocation[1:]...)
	}
	return encodedUL
}

// encodeIMSI encodes IMSI to a byte sequence as specified in TS 29.274
func encodeIMSI(s string) []byte {
	encodedIMSI := []byte{}
	imsi := strings.TrimPrefix(s, "IMSI")
	if len(imsi) < 6 || len(imsi) > 16 {
		return encodedIMSI
	}
	for idx := 0; idx+1 < len(imsi); idx += 2 {
		firstDigit, err := strconv.Atoi(string(imsi[idx]))
		if err != nil {
			return []byte{}
		}
		secondDigit, err := strconv.Atoi(string(imsi[idx+1]))
		if err != nil {
			return []byte{}
		}
		encodedIMSI = append(encodedIMSI, byte((secondDigit<<4)+firstDigit))
	}
	if len(imsi)%2 == 1 {
		digit, err := strconv.Atoi(imsi[len(imsi)-1:])
		if err != nil {
			return []byte{}
		}
		encodedIMSI = append(encodedIMSI, byte(0xF0+digit))
	}
	return encodedIMSI
}

// encodeIMEI encodes IMEI to a byte sequence as specified in TS 29.274
func encodeIMEI(imei []rune) []byte {
	iLen := len(imei)
	if iLen != 15 && iLen != 16 {
		return []byte{}
	}

	// Encoding snr and tac
	encodedIMEI := make([]byte, 8)
	for idx := 0; idx+1 < len(imei)-2; idx += 2 {
		firstDigit := int(imei[idx])
		secondDigit := int(imei[idx+1])
		i := 3 - idx/2
		if idx >= 8 {
			i = 7 - idx/2 + 3
		}
		encodedIMEI[i] = byte((firstDigit << 4) + secondDigit)
	}

	// Encoding spare
	lastDigit := int(imei[iLen-1])
	if iLen%2 == 1 {
		encodedIMEI[7] = byte(0xF0 + lastDigit)
	} else {
		encodedIMEI[7] = byte((int(imei[iLen-2]) << 4) + lastDigit)
	}
	return encodedIMEI
}

// encodeUnixTime encodes unix time to a byte sequence
func encodeUnixTime(timestamp time.Time) []byte {
	unix := convertUint32ToBytes(uint32(timestamp.Unix()))
	return append(
		unix,
		convertUint32ToBytes(uint32(timestamp.Nanosecond()))...,
	)
}

// encodeGeneralizedTime encodes timestamp to a byte sequence as specified in TS 33.108
func encodeGeneralizedTime(timestamp time.Time) []byte {
	formatStr := "20060102150405.000Z"
	return []byte(timestamp.Format(formatStr))
}

// decodeRecordType decodes record type (parent structure tag)
func decodeRecordType(b byte) string {
	switch b {
	case 0xa1:
		return IRIBeginRecord
	case 0xa2:
		return IRIEndRecord
	case 0xa3:
		return IRIContinueRecord
	default:
		return IRIReportRecord
	}
}

// getEPSEventID maps eventd event types to 3GPP event IDs
func getEPSEventID(eventType string) asn1.Enumerated {
	switch eventType {
	case nprobe.SessionCreated:
		return BearerActivation
	case nprobe.SessionUpdated:
		return BearerModification
	case nprobe.SessionTerminated:
		return BearerDeactivation
	case nprobe.AttachSuccess:
		return EutranAttach
	case nprobe.DetachSuccess:
		return EutranDetach
	}
	return UnsupportedEvent
}

// getRecordType maps 3GPP event IDs to corresponding record type
func getRecordType(eventID asn1.Enumerated) string {
	switch eventID {
	case BearerActivation:
		return IRIBeginRecord
	case BearerDeactivation:
		return IRIEndRecord
	case BearerModification:
		return IRIContinueRecord
	default:
		return IRIReportRecord
	}
}

// processEventSpecificData process specific data from each events to from
// 3GPP EPSSpecificParameters structure.
// Some events does not require this processing.
func processEventSpecificData(event *models.Event) EPSSpecificParameters {
	switch event.EventType {
	case nprobe.SessionCreated:
		return makeBearerActivationParams(event)
	case nprobe.SessionUpdated:
		return makeBearerModificationParams(event)
	case nprobe.SessionTerminated:
		return makeBearerDeactivationParams(event)
	}
	return EPSSpecificParameters{}
}

// makeTimestamp returns a Timestamp object as defined in the asn1 schema
func makeTimestamp(timestamp time.Time) Timestamp {
	return Timestamp{
		LocalTime: LocalTimestamp{
			GeneralizedTime:        encodeGeneralizedTime(timestamp),
			WinterSummerIndication: IndicationNotAvailable,
		},
	}
}

// makePdnAddressAllocation returns an encoded PDN address allocation
func makePdnAddressAllocation(event *models.Event) []byte {
	eventData := event.Value.(map[string]interface{})
	if ipAddr, ok := eventData["ip_addr"]; ok {
		allocatedIP := []byte{byte(IPV4PdnType)}
		bIP := net.ParseIP(ipAddr.(string))
		return append(allocatedIP, net.IP.To4(bIP)...)
	}

	if ipv6Addr, ok := eventData["ipv6_addr"]; ok {
		allocatedIPv6 := []byte{byte(IPV6PdnType)}
		bIP6 := net.ParseIP(ipv6Addr.(string))
		return append(allocatedIPv6, net.IP.To16(bIP6)...)
	}
	return []byte{}
}

// makeEPSLocation returns an EPSLocation object as definied in the asn1 schema
func makeEPSLocation(event *models.Event) EPSLocation {
	eventData := event.Value.(map[string]interface{})
	if v, ok := eventData["user_location"]; ok {
		return EPSLocation{
			UserLocationInfo: encodeUserLocation(v.(string)),
		}
	}
	return EPSLocation{}
}

// makePartyInformation returns a PartyInformation slice as defined in the asn1 schema
func makePartyInformation(event *models.Event) []PartyInformation {
	partyID := PartyIdentity{}
	eventData := event.Value.(map[string]interface{})
	if v, ok := eventData["imsi"]; ok {
		partyID.IMSI = encodeIMSI(v.(string))
	}
	if v, ok := eventData["imei"]; ok {
		partyID.IMEI = encodeIMEI([]rune(v.(string)))
	}
	if v, ok := eventData["msisdn"]; ok {
		msisdn := v.(string)
		partyID.MSISDN = []byte(msisdn)
	}
	return []PartyInformation{
		{
			PartyQualified: PartyQualifierTarget,
			PartyIdentity:  partyID,
		},
	}
}

// makeNetworkIdentifier returns a NetworkIdentifier object as defined in the asn1 schema
func makeNetworkIdentifier(event *models.Event, operatorID uint32) NetworkIdentifier {
	eventData := event.Value.(map[string]interface{})
	ipAddr := IPAddress{}
	if originIP, ok := eventData["spgw_ip"]; ok {
		bIP := net.ParseIP(originIP.(string))
		ipAddr.IPType = IPV4Type
		ipAddr.IPValue = IPValue{
			IPBinaryAddress: net.IP.To4(bIP),
		}
	}
	return NetworkIdentifier{
		OperatorIdentifier: []byte(strconv.Itoa(int(operatorID))),
		NetworkElementIdentifier: NetworkElementIdentifier{
			IPAddress: ipAddr,
		},
	}
}

// makeBearerActivationParams returns the corresponding EPSSpecificParameters
// for bearer activation as defined in the asn1 schema
func makeBearerActivationParams(event *models.Event) EPSSpecificParameters {
	eventData := event.Value.(map[string]interface{})
	return EPSSpecificParameters{
		EPSBearerIdentity:      []byte{BearerID},
		PDNAddressAllocation:   makePdnAddressAllocation(event),
		APN:                    encodeAPN(eventData["apn"].(string)),
		RATType:                []byte{RatTypeEutran},
		BearerActivationType:   DefaultBearer,
		EPSLocationOfTheTarget: makeEPSLocation(event),
	}
}

// makeBearerModificationParams returns the corresponding EPSSpecificParameters
// for bearer modification as defined in the asn1 schema
func makeBearerModificationParams(event *models.Event) EPSSpecificParameters {
	return EPSSpecificParameters{
		EPSBearerIdentity:      []byte{BearerID},
		EPSLocationOfTheTarget: makeEPSLocation(event),
	}
}

// makeBearerDeactivationParams returns the corresponding EPSSpecificParameters
// for bearer deactivation as defined in the asn1 schema
func makeBearerDeactivationParams(event *models.Event) EPSSpecificParameters {
	return EPSSpecificParameters{
		EPSBearerIdentity:      []byte{BearerID},
		BearerDeactivationType: DefaultBearer,
		EPSLocationOfTheTarget: makeEPSLocation(event),
	}
}
