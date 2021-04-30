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
	"errors"
	"fmt"

	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	eventdM "magma/orc8r/cloud/go/services/eventd/obsidian/models"

	"github.com/gofrs/uuid"
)

// EpsIRIRecord represents a full IRI record combining header and payload
type EpsIRIRecord struct {
	Header  EpsIRIHeader
	Payload EpsIRIContent
}

// Encode returns a byte sequence of the EpsIRIRecord in network byte order
func (r *EpsIRIRecord) Encode() ([]byte, error) {
	recordType := getRecordType(r.Payload.EPSEvent)
	content, err := asn1.MarshalWithParams(r.Payload, recordType)
	if err != nil {
		return []byte{}, err
	}

	// Update payload length before marshaling the header
	r.Header.PayloadLength = uint32(len(content))
	msg := r.Header.Marshal()
	return append(msg, content...), nil
}

// Decode constructs an IRI record from a byte sequence
func (r *EpsIRIRecord) Decode(b []byte) error {
	if len(b) < int(HeaderFixLen) {
		return errors.New("input too small")
	}

	hdr_len := binary.BigEndian.Uint32(b[4:8])
	pld_len := binary.BigEndian.Uint32(b[8:12])
	if int(hdr_len+pld_len) > len(b) {
		return errors.New("invalid input size")
	}

	if err := r.Header.Unmarshal(b[:hdr_len]); err != nil {
		return err
	}

	content := b[hdr_len : hdr_len+pld_len]
	recordType := decodeRecordType(b[hdr_len])
	if _, err := asn1.UnmarshalWithParams(content, &r.Payload, recordType); err != nil {
		return err
	}
	return nil
}

// makeConditionalAttributes builds the mandatory conditional attributes defined in
// ETSI TS 103 221-2
func makeConditionalAttributes(
	streamName, targetID string,
	timestamp []byte,
	seqNbr uint32,
) ([]Attribute, uint32) {
	attrs := []Attribute{}
	attrs = append(attrs, NewAttribute(AttributeNetworkFn, []byte(streamName)))
	attrs = append(attrs, NewAttribute(AttributeTargetID, []byte(targetID)))
	attrs = append(attrs, NewAttribute(AttributeTimestamp, timestamp))
	attrs = append(attrs, NewAttribute(AttributeSeqNumber, convertUint32ToBytes(seqNbr)))

	attrs_len := uint32(0)
	for _, attr := range attrs {
		attrs_len += uint32(attr.Len) + 4 // TAG(2B) + LEN(2B)
	}
	return attrs, attrs_len
}

// makeEpsIRIContent builds the IRI Content structure with the available information
// for each event
func makeEpsIRIContent(
	event *eventdM.Event,
	eventID asn1.Enumerated,
	correlationID uint64,
	operatorID uint32,
	timestamp []byte,
) EpsIRIContent {
	// build iri content
	return EpsIRIContent{
		Hi2epsDomainID:        GetOID(),
		TimeStamp:             makeTimestamp(timestamp),
		Initiator:             InitiatorNotAvailable,
		PartyInformation:      makePartyInformation(event),
		EPSCorrelationNumber:  convertUint64ToBytes(correlationID),
		EPSEvent:              eventID,
		NetworkIdentifier:     makeNetworkIdentifier(event, operatorID),
		EPSSpecificParameters: processEventSpecificData(event),
	}
}

// MakeEpsIRIRecord build a new record and encode it to a byte sequence
func MakeRecord(
	event *eventdM.Event,
	task *models.NetworkProbeTask,
	operatorID, sequenceNbr uint32,
) ([]byte, error) {

	// map event type to 3gpp event id
	eventID := getEPSEventID(event.EventType)
	if eventID == UnsupportedEvent {
		return []byte{}, fmt.Errorf("Unsupported event type %s\n", event.EventType)
	}

	bTimestamp, err := encodeGeneralizedTime(event.Timestamp)
	if err != nil {
		return []byte{}, err
	}

	attrs, attrs_len := makeConditionalAttributes(
		event.StreamName,
		task.TaskDetails.TargetID,
		bTimestamp,
		sequenceNbr,
	)

	uuid, err := uuid.FromString(string(task.TaskID))
	if err != nil {
		return []byte{}, err
	}

	correlationID := task.TaskDetails.CorrelationID
	record := EpsIRIRecord{
		Header:  NewEpsIRIHeader(uuid, correlationID, attrs, attrs_len),
		Payload: makeEpsIRIContent(event, eventID, correlationID, operatorID, bTimestamp),
	}
	return record.Encode()
}
