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
	"encoding/binary"
	"fmt"
	"net"
	"time"
	"unsafe"
	"encoding/asn1"

	"magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/golang/glog"
)

// ETSI TS 133 108 [B9] ASN.1 OID
func getAsn1ObjectIdentifier() asn1.ObjectIdentifier {
	return []int{0, 4, 0, 2, 2, 4, 8, 15, 4}
}

func getPdnAddressAllocation(eventData map[string]interface{}) []byte {
	var allocatedIP []byte
	ipAddress := eventData["ip_addr"].(string)
	ip6Address := eventData["ipv6_addr"].(string)

	if len(ipAddress) == 0 && len(ip6Address) == 0 {
		glog.Errorf("PDN allocated address not found %v", eventData)
		return allocatedIP
	}

	if len(ipAddress) > 0 {
		allocatedIP = append(allocatedIP, byte(IPV4Type))
		allocatedIP = append(allocatedIP, net.ParseIP(ipAddress)[12:16]...)
	} else {
		allocatedIP = append(allocatedIP, byte(IPV6Type))
		allocatedIP = append(allocatedIP, net.ParseIP(ip6Address)...)
	}
	return allocatedIP
}

func getNetworkIdentifier(eventData map[string]interface{}, operatorID string) NetworkIdentifier {
	spgwIP := eventData["spgw_ip"].(string)
	ipAddress := IPAddress{
		IPType:  IPV6Type,
		IPValue: IPValue{IPBinaryAddress: []byte(spgwIP)},
	}

	return NetworkIdentifier{
		OperatorIdentifier: []byte(operatorID),
		NetworkElementIdentifier: NetworkElementIdentifier{
			IPAddress: ipAddress,
		},
	}
}

func getPartyInformation(eventData map[string]interface{}) []PartyInformation {
	imei := eventData["imei"].(string)
	imsi := eventData["imsi"].(string)
	msisdn := eventData["msisdn"].(string)

	partyInformation := PartyInformation{
		PartyQualified: PartyQualifierTarget,
		PartyIdentity: PartyIdentity{
			IMEI:   []byte(imei),
			IMSI:   []byte(imsi),
			MSISDN: []byte(msisdn),
		},
	}
	return []PartyInformation{partyInformation}
}

func getTimestamp(timestamp string) (*Timestamp, error) {
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

func processEventSpecificData(event_type string, eventData map[string]interface{}) (
	EPSSpecificParameters,
	asn1.Enumerated,
	string,
) {
	switch event_type {
	case "session_created":
		return getBearerActivationSpecificParams(eventData), BearerActivation, IRIBeginRecord
	case "session_modified":
		return getBearerModificationSpecificParams(eventData), BearerModification, IRIContinueRecord
	case "session_terminated":
		return getBearerDeactivationSpecificParams(eventData), BearerDeactivation, IRIEndRecord
	case "session_create_failure":
		return getBearerActivationFailureSpecificParams(eventData), BearerActivation, IRIReportRecord
	case "attach_success":
		return getEutranAttachSpecificParams(eventData), BearerActivation, IRIReportRecord
	case "detach_success":
		return getEutranADettachSpecificParams(eventData), BearerActivation, IRIReportRecord
	case "s1_setup_success":
		return getUERequestedPDNConnectivitySpecificParams(eventData), UERequestedPDNConnectivity, IRIReportRecord
	default:
		return EPSSpecificParameters{}, UnsupportedEvent, IRIReportRecord
	}
}

func buildInterceptRelatedInformationContent(event models.Event, opID string, lawfulInterceptionID, correlationID []byte) (
	IRIParameter,
	string,
	error,
) {
	eventData := event.Value.(map[string]interface{})
	specificData, eventID, tags := processEventSpecificData(event.EventType, eventData)
	if eventID == UnsupportedEvent {
		return IRIParameter{}, "", fmt.Errorf("Unsupported Event %s", event.EventType)
	}

	timestamp, err := getTimestamp(event.Timestamp)
	if err != nil {
		return IRIParameter{}, "", fmt.Errorf("Failed to parse timestamp %s. Skip this event %s", event.Timestamp, event.EventType)
	}

	return IRIParameter{
		Hi2epsDomainID:        getAsn1ObjectIdentifier(),
		LawfulInterceptionID:  lawfulInterceptionID,
		TimeStamp:             *timestamp,
		Initiator:             InitiatorNotAvailable,
		PartyInformation:      getPartyInformation(eventData),
		EPSCorrelationNumber:  correlationID,
		EPSevent:              eventID,
		NetworkIdentifier:     getNetworkIdentifier(eventData, opID),
		EPSSpecificParameters: specificData,
	}, tags, nil
}

func BuildInterceptRelatedInformationPDU(event models.Event, opID string, lawfulInterceptionID, correlationID []byte) (EpsIRIPDU, error) {
	content, tags, err := buildInterceptRelatedInformationContent(event, opID, lawfulInterceptionID, lawfulInterceptionID)
	if err != nil {
		glog.Errorf("Failed to build IRI Content for this event %v\n", event)
		return EpsIRIPDU{}, err
	}

	encodedContent, err := asn1.MarshalWithParams(content, tags)
	if err != nil {
		glog.Errorf("Failed to ASN.1 encode this content %v\n", content)
		return EpsIRIPDU{}, err
	}
	pdu := EpsIRIPDU{
		Version:               X2HeaderVersion,
		PduType:               X2HeaderPduType,
		PayloadLength:         uint32(len(encodedContent)),
		PayloadDirection:      X2PayloadDirectionUnkown,
		CorrelationID:         binary.BigEndian.Uint64(correlationID),
		ConditionalAttrFields: []Attribute{},
		Payload:               encodedContent,
	}
	pdu.HeaderLength = uint32(unsafe.Sizeof(pdu)) - pdu.PayloadLength
	return pdu, nil
}
