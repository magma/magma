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
	"errors"

	"github.com/gofrs/uuid"
)

var (
	HeaderFixLen        uint32 = 40
	HeaderVersion       uint16 = 2
	HeaderPduType       uint16 = 1  // X2 PDU
	HeaderPayloadFormat uint16 = 14 // ETSI TS 133 108 [B.9] Defined Payload

	AttributeDomainID  uint16 = 5
	AttributeNetworkFn uint16 = 6
	AttributeTimestamp uint16 = 9
	AttributeSeqNumber uint16 = 8
	AttributeTargetID  uint16 = 17

	PayloadDirectionUnkown     uint16 = 1
	PayloadDirectionToTarget   uint16 = 2
	PayloadDirectionFromTarget uint16 = 3
)

// Attribute represents an X2 IRI conditional attribute field
// as defined in ETSI TS 103 221-2.
type Attribute struct {
	Tag   uint16
	Len   uint16
	Value []byte
}

// NewAttribute creates and returns a new attribute
func NewAttribute(tag uint16, value []byte) Attribute {
	return Attribute{
		Tag:   tag,
		Len:   uint16(len(value)),
		Value: value,
	}
}

// marshal returns a byte sequence of the attribute in network byte order.
func (t *Attribute) marshal() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b[0:2], t.Tag)
	binary.BigEndian.PutUint16(b[2:4], t.Len)
	return append(b, t.Value...)
}

// marshalAttributes parses a slice of attributes and marshals them.
func marshalAttributes(attrs []Attribute) []byte {
	var b []byte
	for _, attr := range attrs {
		b = append(b, attr.marshal()...)
	}
	return b
}

// Unmarshal parses a byte array to reconstruct an attribute.
func (t *Attribute) unmarshal(b []byte) error {
	if len(b) < 4 {
		return errors.New("invalid input size")
	}

	t.Tag = binary.BigEndian.Uint16(b[0:2])
	t.Len = binary.BigEndian.Uint16(b[2:4])

	if int(t.Len)+4 > len(b) {
		return errors.New("truncated attribute or wrong length")
	}
	t.Value = make([]byte, t.Len)
	copy(t.Value, b[4:])
	return nil
}

func parseAttributes(b []byte) ([]Attribute, error) {
	var attrs []Attribute
	for i := 0; i < len(b); {
		var attr Attribute
		if err := attr.unmarshal(b[i:]); err != nil {
			return attrs, err
		}
		attrs = append(attrs, attr)
		i += int(attr.Len) + 4
	}
	return attrs, nil
}

// EpsIRIHeader represents an X2 IRI header as defined in ETSI TS 103 221-2.
type EpsIRIHeader struct {
	Version               uint16
	PduType               uint16
	HeaderLength          uint32
	PayloadLength         uint32
	PayloadFormat         uint16
	PayloadDirection      uint16
	XID                   uuid.UUID
	CorrelationID         uint64
	ConditionalAttributes []Attribute
}

// NewEpsIRIHeader creates and returns a new EpsIRIHeader
func NewEpsIRIHeader(uuid uuid.UUID, corrID uint64, attrs []Attribute, attrs_len uint32) EpsIRIHeader {
	return EpsIRIHeader{
		Version:               HeaderVersion,
		PduType:               HeaderPduType,
		HeaderLength:          attrs_len + HeaderFixLen,
		PayloadFormat:         HeaderPayloadFormat,
		PayloadDirection:      PayloadDirectionUnkown,
		XID:                   uuid,
		ConditionalAttributes: attrs,
		CorrelationID:         corrID,
	}
}

// Marshal returns a byte sequence of the header in network byte order.
func (h *EpsIRIHeader) Marshal() []byte {
	b := make([]byte, h.HeaderLength)
	h.marshalTo(b)
	return b
}

func (h *EpsIRIHeader) marshalTo(b []byte) {
	binary.BigEndian.PutUint16(b[0:2], h.Version)
	binary.BigEndian.PutUint16(b[2:4], h.PduType)
	binary.BigEndian.PutUint32(b[4:8], h.HeaderLength)
	binary.BigEndian.PutUint32(b[8:12], h.PayloadLength)
	binary.BigEndian.PutUint16(b[12:14], h.PayloadFormat)
	binary.BigEndian.PutUint16(b[14:16], h.PayloadDirection)
	copy(b[16:32], h.XID.Bytes()) // append UUID
	binary.BigEndian.PutUint64(b[32:40], h.CorrelationID)

	attrs := marshalAttributes(h.ConditionalAttributes)
	copy(b[40:], attrs) // append attributes
}

// Unmarshal parses the BER decoded ASN.1 as defined in ETSI TS 103 221-2.
func (h *EpsIRIHeader) Unmarshal(b []byte) error {
	if len(b) < int(HeaderFixLen) {
		return errors.New("invalid input size")
	}

	var err error
	h.XID, err = uuid.FromBytes(b[16:32])
	if err != nil {
		return err
	}

	h.ConditionalAttributes, err = parseAttributes(b[40:])
	if err != nil {
		return err
	}
	h.Version = binary.BigEndian.Uint16(b[0:2])
	h.PduType = binary.BigEndian.Uint16(b[2:4])
	h.HeaderLength = binary.BigEndian.Uint32(b[4:8])
	h.PayloadLength = binary.BigEndian.Uint32(b[8:12])
	h.PayloadFormat = binary.BigEndian.Uint16(b[12:14])
	h.PayloadDirection = binary.BigEndian.Uint16(b[14:16])
	h.CorrelationID = binary.BigEndian.Uint64(b[32:40])
	return nil
}
