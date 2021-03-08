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
	"encoding/binary"

	"github.com/gofrs/uuid"
)

type ConditionalAttribute struct {
	Tag    uint16
	Length uint16
	Value  []byte
}

// Serialize returns a byte sequence of the header in network byte order.
func (t ConditionalAttribute) serialize() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint16(b[0:2], t.Tag)
	binary.BigEndian.PutUint16(b[2:4], t.Length)
	return append(b, t.Value...)
}

func makeStrConditionalAttribute(value string, tag uint16) ConditionalAttribute {
	encodedValue := []byte(value)
	return ConditionalAttribute{
		Tag:    tag,
		Value:  encodedValue,
		Length: uint16(len(encodedValue)),
	}
}

func makeBinConditionalAttribute(value []byte, tag uint16) ConditionalAttribute {
	return ConditionalAttribute{
		Tag:    tag,
		Value:  value,
		Length: uint16(len(value)),
	}
}

// EpsIRIHeader represents an X2 IRI header
type EpsIRIHeader struct {
	Version               uint16
	PduType               uint16
	HeaderLength          uint32
	PayloadLength         uint32
	PayloadFormat         uint16
	PayloadDirection      uint16
	XID                   uuid.UUID
	CorrelationID         uint64
	ConditionalAttrFields []ConditionalAttribute
}

// Serialize returns a byte sequence of the header in network byte order.
func (h *EpsIRIHeader) Serialize() []byte {
	b := make([]byte, FixHeaderLength)
	h.SerializeTo(b)
	return b
}

// SerializeTo serializes the header to a byte sequence in network byte order.
func (h *EpsIRIHeader) SerializeTo(b []byte) {
	binary.BigEndian.PutUint16(b[0:2], h.Version)
	binary.BigEndian.PutUint16(b[2:4], h.PduType)
	binary.BigEndian.PutUint32(b[4:8], h.HeaderLength)
	binary.BigEndian.PutUint32(b[8:12], h.PayloadLength)
	binary.BigEndian.PutUint16(b[12:14], h.PayloadFormat)
	binary.BigEndian.PutUint16(b[14:16], h.PayloadDirection)
	copy(b[16:34], h.XID.Bytes()) // append UUID
	binary.BigEndian.PutUint64(b[32:40], h.CorrelationID)
	// serialize conditional attributes
	for _, tlv := range h.ConditionalAttrFields {
		b = append(b, tlv.serialize()...)
	}
}

func makeEpsIRIHeader(
	uuid uuid.UUID,
	attrs []ConditionalAttribute,
	correlationID uint64,
	hLen, pLen uint32,
) EpsIRIHeader {
	return EpsIRIHeader{
		Version:               X2HeaderVersion,
		PduType:               X2HeaderPduType,
		HeaderLength:          hLen,
		PayloadLength:         pLen,
		PayloadDirection:      X2PayloadDirectionUnkown,
		XID:                   uuid,
		CorrelationID:         correlationID,
		ConditionalAttrFields: attrs,
	}
}
