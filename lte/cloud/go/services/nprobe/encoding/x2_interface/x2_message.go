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

package x2_interface

import (
	"encoding/asn1"
	"encoding/binary"
	"fmt"
	"net"
	"time"
	"unsafe"

	"magma/lte/cloud/go/services/nprobe"
	"magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/gofrs/uuid"
	"github.com/golang/glog"
)

func encodeTimestamp(timestamp string) (*Timestamp, error) {
	generalizedTime, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return nil, err
	}
	return &Timestamp{
		LocalTime: LocalTimestamp{
			GeneralizedTime:        generalizedTime,
			WinterSummerIndication: IndicationNotAvailable,
		},
	}, nil
}

func encodePdnAddressAllocation(eventData map[string]interface{}) []byte {
	ipAddress := eventData["ip_addr"].(string)
	addrType := IPV4Type

	if len(ipAddress) == 0 {
		ipAddress = eventData["ipv6_addr"].(string)
		addrType = IPV6Type
	}

	allocatedIP := []byte{byte(addrType)}
	switch addrType {
	case IPV4Type:
		allocatedIP = append(allocatedIP, net.ParseIP(ipAddress)[12:16]...)
		break
	case IPV6Type:
		allocatedIP = append(allocatedIP, net.ParseIP(ipAddress)...)
		break
	}
	return allocatedIP
}

func constructNetworkIdentifier(ipAddr, operatorID string) NetworkIdentifier {
	return NetworkIdentifier{
		OperatorIdentifier: []byte(operatorID),
		NetworkElementIdentifier: NetworkElementIdentifier{
			IPAddress: IPAddress{
				IPType:  IPV6Type,
				IPValue: IPValue{IPBinaryAddress: []byte(ipAddr)},
			},
		},
	}
}

func constructPartyInformation(eventData map[string]interface{}) []PartyInformation {
	return []PartyInformation{
		{
			PartyQualified: PartyQualifierTarget,
			PartyIdentity: PartyIdentity{
				IMEI:   []byte(eventData["imei"].(string)),
				IMSI:   []byte(eventData["imsi"].(string)),
				MSISDN: []byte(eventData["msisdn"].(string)),
			},
		},
	}
}

func getBearerActivationParams(eventData map[string]interface{}) EPSSpecificParameters {
	return EPSSpecificParameters{
		PDNAddressAllocation: encodePdnAddressAllocation(eventData), // Should only be provided if default bearer
		APN:                  []byte(eventData["apn"].(string)),
		EPSBearerIdentity:    []byte(eventData["session_id"].(string)),
		RATType:              []byte{byte(RatTypeEutran)},
		EPSBearerQoS:         []byte(""),    // TBD
		BearerActivationType: DefaultBearer, // TBD
		ApnAmbr:              []byte(""),    // TBD
		EPSLocationOfTheTarget: EPSLocation{
			UserLocationInfo: []byte(eventData["user_location"].(string)),
		},
	}
}

func getBearerDeactivationParams(eventData map[string]interface{}) EPSSpecificParameters {
	return EPSSpecificParameters{
		EPSBearerIdentity:      []byte(eventData["session_id"].(string)),
		BearerDeactivationType: DefaultBearer, // TBD
		EPSLocationOfTheTarget: EPSLocation{
			UserLocationInfo: []byte(eventData["user_location"].(string)),
		},
	}
}

func processEventSpecificData(event *models.Event) (
	EPSSpecificParameters,
	EventAdditionalParams,
) {
	var params = EventAdditionalParams{
		eventID: UnsupportedEvent,
		asn1Tag: IRIReportRecord,
	}

	eventData := event.Value.(map[string]interface{})
	switch event.StreamName {
	case "sessiond":
		params.originIP = eventData["spgw_ip"].(string)
	}

	switch event.EventType {
	case nprobe.SessionCreated:
		params.eventID = BearerActivation
		params.asn1Tag = IRIBeginRecord
		return getBearerActivationParams(eventData), params
	case nprobe.SessionTerminated:
		params.eventID = BearerModification
		params.asn1Tag = IRIContinueRecord
		return getBearerDeactivationParams(eventData), params
	default:
		return EPSSpecificParameters{}, params
	}
}

func buildAsn1EpsIRIContent(event *models.Event, operatorID string, correlationID uint64) ([]byte, error) {
	eventData := event.Value.(map[string]interface{})
	specificData, params := processEventSpecificData(event)
	if params.eventID == UnsupportedEvent {
		return []byte{}, fmt.Errorf("Unsupported event %s\n", event.EventType)
	}

	timestamp, err := encodeTimestamp(event.Timestamp)
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to parse timestamp %s.", event.Timestamp)
	}

	corrID := make([]byte, 8)
	binary.BigEndian.PutUint64(corrID, correlationID)

	content := EpsIRIContent{
		Hi2epsDomainID:        GetOID(),
		TimeStamp:             *timestamp,
		Initiator:             InitiatorNotAvailable,
		PartyInformation:      constructPartyInformation(eventData),
		EPSCorrelationNumber:  corrID,
		EPSEvent:              params.eventID,
		NetworkIdentifier:     constructNetworkIdentifier(params.originIP, operatorID),
		EPSSpecificParameters: specificData,
	}

	encodedContent, err := asn1.MarshalWithParams(content, params.asn1Tag)
	if err != nil {
		glog.Errorf("Failed to encode EpsIRIContent %v\n", content)
		return []byte{}, err
	}
	return encodedContent, nil
}

func buildEpsIRIConditionalAttributes() ([]ConditionalAttribute, error) {
	// TBD
	return []ConditionalAttribute{}, nil
}

func ConstructEpsIRIMessage(event *models.Event, operatorID, xID string, correlationID uint64) (EpsIRIMessage, error) {
	uuid, err := uuid.FromString(xID)
	if err != nil {
		glog.Errorf("Failed to parse xID - type 4 %s\n", xID)

	}
	content, err := buildAsn1EpsIRIContent(event, operatorID, correlationID)
	if err != nil {
		glog.Errorf("Failed to build IRI Content for this event %s\n", event.EventType)
		return EpsIRIMessage{}, err
	}

	attrs, err := buildEpsIRIConditionalAttributes()
	if err != nil {
		glog.Errorf("Failed to build IRI Content for this event %s\n", event.EventType)
		return EpsIRIMessage{}, err
	}

	record := EpsIRIMessage{
		Version:               X2HeaderVersion,
		PduType:               X2HeaderPduType,
		PayloadLength:         uint32(len(content)),
		PayloadDirection:      X2PayloadDirectionUnkown,
		XID:                   uuid,
		CorrelationID:         correlationID,
		ConditionalAttrFields: attrs,
		Payload:               content,
	}
	record.HeaderLength = uint32(unsafe.Sizeof(record)) - record.PayloadLength
	return record, nil
}
