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
	"net"
	"time"

	"magma/lte/cloud/go/services/nprobe"
	"magma/orc8r/cloud/go/services/eventd/obsidian/models"
)

func encodeGeneralizedTime(timestamp string) ([]byte, error) {
	ptime, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return []byte{}, nil
	}
	return ptime.MarshalBinary()
}

func encodePdnAddressAllocation(event *models.Event) []byte {
	ipAddr, addrType := getEventSourceIP(event)
	allocatedIP := []byte{byte(addrType)}
	switch addrType {
	case IPV4Type:
		allocatedIP = append(allocatedIP, net.ParseIP(ipAddr)[12:16]...)
		break
	case IPV6Type:
		allocatedIP = append(allocatedIP, net.ParseIP(ipAddr)...)
		break
	}
	return allocatedIP
}

func getEventSourceIP(event *models.Event) (string, asn1.Enumerated) {
	eventData := event.Value.(map[string]interface{})
	ipAddress := eventData["ip_addr"].(string)
	if ipAddress != "" {
		return ipAddress, IPV4Type
	}
	return eventData["ipv6_addr"].(string), IPV6Type
}

func getOriginNodeIP(event *models.Event) string {
	eventData := event.Value.(map[string]interface{})
	originIP := ""
	switch event.StreamName {
	// TBD support mme and hss
	case "sessiond":
		originIP = eventData["spgw_ip"].(string)
	}
	return originIP
}

func processEventSpecificData(event *models.Event) (
	EPSSpecificParameters,
	EventAdditionalParams,
) {
	tag, eventID := IRIReportRecord, UnsupportedEvent
	specificData := EPSSpecificParameters{}
	switch event.EventType {
	case nprobe.SessionCreated:
		eventID, tag, specificData = makeBearerActivationParams(event)
	case nprobe.SessionUpdated:
		eventID, tag, specificData = makeBearerModificationParams(event)
	case nprobe.SessionTerminated:
		eventID, tag, specificData = makeBearerDeactivationParams(event)
	case nprobe.AttachSuccess:
		eventID, tag = EutranAttach, IRIReportRecord
	case nprobe.DetachSuccess:
		eventID, tag = EutranDetach, IRIReportRecord
	}

	params := EventAdditionalParams{
		eventID:  eventID,
		asn1Tag:  tag,
		originIP: getOriginNodeIP(event),
	}
	return specificData, params
}

func makeTimestamp(timestamp string) (*Timestamp, error) {
	ptime, err := time.Parse(time.RFC3339Nano, timestamp)
	if err != nil {
		return nil, err
	}
	return &Timestamp{
		LocalTime: LocalTimestamp{
			GeneralizedTime:        ptime,
			WinterSummerIndication: IndicationNotAvailable,
		},
	}, nil
}

func makeNetworkIdentifier(ipAddr, operatorID string) NetworkIdentifier {
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

func makePartyInformation(event *models.Event) []PartyInformation {
	eventData := event.Value.(map[string]interface{})
	info := PartyInformation{
		PartyQualified: PartyQualifierTarget,
		PartyIdentity: PartyIdentity{
			IMSI: []byte(eventData["imsi"].(string)),
		},
	}

	if event.StreamName == "sessiond" {
		info.PartyIdentity.IMEI = []byte(eventData["imei"].(string))
		info.PartyIdentity.MSISDN = []byte(eventData["msisdn"].(string))
	}
	return []PartyInformation{info}
}

func makeBearerActivationParams(event *models.Event) (
	asn1.Enumerated,
	string,
	EPSSpecificParameters,
) {
	eventData := event.Value.(map[string]interface{})
	specificData := EPSSpecificParameters{
		PDNAddressAllocation: encodePdnAddressAllocation(event), // Should only be provided if default bearer
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
	return BearerActivation, IRIBeginRecord, specificData
}

func makeBearerModificationParams(event *models.Event) (
	asn1.Enumerated,
	string,
	EPSSpecificParameters,
) {
	eventData := event.Value.(map[string]interface{})
	specificData := EPSSpecificParameters{
		EPSBearerIdentity:      []byte(eventData["session_id"].(string)),
		BearerDeactivationType: DefaultBearer, // TBD
		EPSLocationOfTheTarget: EPSLocation{
			UserLocationInfo: []byte(eventData["user_location"].(string)),
		},
	}
	return BearerModification, IRIContinueRecord, specificData
}

func makeBearerDeactivationParams(event *models.Event) (
	asn1.Enumerated,
	string,
	EPSSpecificParameters,
) {
	eventData := event.Value.(map[string]interface{})
	specificData := EPSSpecificParameters{
		EPSBearerIdentity:      []byte(eventData["session_id"].(string)),
		BearerDeactivationType: DefaultBearer, // TBD
		EPSLocationOfTheTarget: EPSLocation{
			UserLocationInfo: []byte(eventData["user_location"].(string)),
		},
	}
	return BearerDeactivation, IRIEndRecord, specificData
}
