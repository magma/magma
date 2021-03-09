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

	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	eventd_models "magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/gofrs/uuid"
	"github.com/golang/glog"
)

// EpsIRIRecord
type EpsIRIRecord struct {
	Header  EpsIRIHeader
	Payload []byte
}

// MakeRecord build a EpsIRIRecord record
func MakeRecord(event *eventd_models.Event, operatorID string, task *models.NetworkProbeTask, seqNbr uint64) (EpsIRIRecord, error) {
	// build payload content
	details := task.TaskDetails
	payload, err := makePayload(event, operatorID, details.CorrelationID)
	if err != nil {
		glog.Errorf("Failed to build IRI content\n")
		return EpsIRIRecord{}, err
	}

	// build record header
	header, err := makeHeader(event, string(task.TaskID), details.TargetID, details.CorrelationID, seqNbr, uint32(len(payload)))
	if err != nil {
		glog.Errorf("Failed to build IRI header\n")
		return EpsIRIRecord{}, err
	}
	return EpsIRIRecord{Header: header, Payload: payload}, nil
}

func makeHeader(event *eventd_models.Event, xID, targetID string, correlationID, seqNbr uint64, pLen uint32) (EpsIRIHeader, error) {
	uuid, err := uuid.FromString(xID)
	if err != nil {
		glog.Errorf("Failed to parse xID - type 4 %s\n", xID)
		return EpsIRIHeader{}, err
	}

	bTime, err := encodeGeneralizedTime(event.Timestamp)
	if err != nil {
		glog.Errorf("Failed to encode timestamp %s\n", event.Timestamp)
		return EpsIRIHeader{}, err
	}

	attrs, aLen := addConditionalAttributes(targetID, event.StreamName, bTime, seqNbr)
	header := makeEpsIRIHeader(uuid, attrs, correlationID, aLen+FixHeaderLength, pLen)
	return header, nil
}

func addConditionalAttributes(targetID, networkFunc string, bTime []byte, seqNbr uint64) ([]ConditionalAttribute, uint32) {

	bSeqNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(bSeqNumber, seqNbr)

	attrs := []ConditionalAttribute{}
	attrs = append(attrs, makeStrConditionalAttribute(targetID, X2AttrIDTargetID))
	attrs = append(attrs, makeStrConditionalAttribute(networkFunc, X2AttrIDNetworkFunc))
	attrs = append(attrs, makeBinConditionalAttribute(bTime, X2AttrIDTimestamp))
	attrs = append(attrs, makeBinConditionalAttribute(bSeqNumber, X2AttrIDSeqNumber))

	aLen := uint32(0)
	for _, attr := range attrs {
		aLen += uint32(attr.Length) + 4 // Tag(2 Bytes) + Len(2 Bytes)
	}
	return attrs, aLen
}

func makePayload(event *eventd_models.Event, operatorID string, correlationID uint64) ([]byte, error) {
	specificData, params := processEventSpecificData(event)
	if params.eventID == UnsupportedEvent {
		return []byte{}, fmt.Errorf("Unsupported event %s\n", event.EventType)
	}

	timestamp, err := makeTimestamp(event.Timestamp)
	if err != nil {
		return []byte{}, err
	}

	corrID := make([]byte, 8)
	binary.BigEndian.PutUint64(corrID, correlationID)
	content := EpsIRIContent{
		Hi2epsDomainID:        GetOID(),
		TimeStamp:             *timestamp,
		Initiator:             InitiatorNotAvailable,
		PartyInformation:      makePartyInformation(event),
		EPSCorrelationNumber:  corrID,
		EPSEvent:              params.eventID,
		NetworkIdentifier:     makeNetworkIdentifier(params.originIP, operatorID),
		EPSSpecificParameters: specificData,
	}
	return asn1.MarshalWithParams(content, params.asn1Tag)
}
