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
	"net"
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

// encodeGeneralizedTime parses timestamp and marshals it to
// a byte sequence
func encodeGeneralizedTime(timestamp string) ([]byte, error) {
	ptime, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return []byte{}, err
	}
	return ptime.MarshalBinary()
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
func makeTimestamp(timestamp []byte) Timestamp {
	return Timestamp{
		LocalTime: LocalTimestamp{
			GeneralizedTime:        timestamp,
			WinterSummerIndication: IndicationNotAvailable,
		},
	}
}

// makePdnAddressAllocation returns an encoded PDN address allocation
func makePdnAddressAllocation(event *models.Event) []byte {
	eventData := event.Value.(map[string]interface{})
	if ipAddr, ok := eventData["ip_addr"]; ok {
		allocatedIP := []byte{byte(IPV4Type)}
		return append(allocatedIP, net.ParseIP(ipAddr.(string))[12:16]...)
	}

	if ipv6Addr, ok := eventData["ipv6_addr"]; ok {
		allocatedIPv6 := []byte{byte(IPV6Type)}
		return append(allocatedIPv6, net.ParseIP(ipv6Addr.(string))...)
	}
	return []byte{}
}

// makeEPSLocation returns an EPSLocation object as definied in the asn1 schema
func makeEPSLocation(event *models.Event) EPSLocation {
	eventData := event.Value.(map[string]interface{})
	if userLocation, ok := eventData["user_location"]; ok {
		return EPSLocation{
			UserLocationInfo: []byte(userLocation.(string)),
		}
	}
	return EPSLocation{}
}

// makePartyInformation returns a PartyInformation slice as defined in the asn1 schema
func makePartyInformation(event *models.Event) []PartyInformation {
	var infoIdentity PartyIdentity
	eventData := event.Value.(map[string]interface{})
	if imsi, ok := eventData["imsi"]; ok {
		infoIdentity.IMSI = []byte(imsi.(string))
	}
	if imei, ok := eventData["imei"]; ok {
		infoIdentity.IMEI = []byte(imei.(string))
	}
	if msisdn, ok := eventData["msisdn"]; ok {
		infoIdentity.MSISDN = []byte(msisdn.(string))
	}
	return []PartyInformation{
		{
			PartyQualified: PartyQualifierTarget,
			PartyIdentity:  infoIdentity,
		},
	}
}

// makeNetworkIdentifier returns a NetworkIdentifier object as defined in the asn1 schema
func makeNetworkIdentifier(event *models.Event, operatorID uint32) NetworkIdentifier {
	eventData := event.Value.(map[string]interface{})
	var ipAddr IPAddress
	if originIP, ok := eventData["spgw_ip"]; ok {
		ipAddr.IPType = IPV4Type
		ipAddr.IPValue = IPValue{
			IPBinaryAddress: []byte(originIP.(string)),
		}
	}
	return NetworkIdentifier{
		OperatorIdentifier: convertUint32ToBytes(operatorID),
		NetworkElementIdentifier: NetworkElementIdentifier{
			IPAddress: ipAddr,
		},
	}
}

// makeBearerActivationParams returns the corresponding EPSSpecificParameters
// for bearer activation as defined in the asn1 schema
func makeBearerActivationParams(event *models.Event) EPSSpecificParameters {
	eventData := event.Value.(map[string]interface{})
	if sessionID, ok := eventData["session_id"]; ok {
		apn := []byte{}
		if v, ok := eventData["apn"]; ok {
			apn = append(apn, []byte(v.(string))...)

		}
		return EPSSpecificParameters{
			EPSBearerIdentity:      []byte(sessionID.(string)),
			PDNAddressAllocation:   makePdnAddressAllocation(event),
			APN:                    apn,
			RATType:                []byte{RatTypeEutran},
			BearerActivationType:   DefaultBearer,
			EPSLocationOfTheTarget: makeEPSLocation(event),
		}
	}
	return EPSSpecificParameters{}
}

// makeBearerModificationParams returns the corresponding EPSSpecificParameters
// for bearer modification as defined in the asn1 schema
func makeBearerModificationParams(event *models.Event) EPSSpecificParameters {
	eventData := event.Value.(map[string]interface{})
	if sessionID, ok := eventData["session_id"]; ok {
		return EPSSpecificParameters{
			EPSBearerIdentity:      []byte(sessionID.(string)),
			EPSLocationOfTheTarget: makeEPSLocation(event),
		}
	}
	return EPSSpecificParameters{}
}

// makeBearerDeactivationParams returns the corresponding EPSSpecificParameters
// for bearer deactivation as defined in the asn1 schema
func makeBearerDeactivationParams(event *models.Event) EPSSpecificParameters {
	eventData := event.Value.(map[string]interface{})
	if sessionID, ok := eventData["session_id"]; ok {
		return EPSSpecificParameters{
			EPSBearerIdentity:      []byte(sessionID.(string)),
			BearerDeactivationType: DefaultBearer,
			EPSLocationOfTheTarget: makeEPSLocation(event),
		}
	}
	return EPSSpecificParameters{}
}
