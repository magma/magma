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

package packet

import (
	"errors"
	"fmt"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2869"
)

// EAP constants
const (
	EapMinLegalPacketLength  int = 4
	EapPacketHeaderLength    int = 4
	EapPacketTypeFieldLength int = 1
)

// EAP packet header offsets
const (
	EapCodeOffset       int = 0
	EapIdentifierOffset int = 1
	EapLengthMsbOffset  int = 2
	EapLengthLsbOffset  int = 3
)

// EapTypeOffset EAP type is not considered part of the header,
// as it is present only on EAP-Request and EAP-Response
// (and not in EAP-Success/Failure)
const EapTypeOffset int = 4

// Packet represents a single EAP Packet
type Packet struct {
	Code       Code
	EAPType    EAPType
	Identifier int
	Data       []byte
}

// NewPacket Creates an packet.Packet with the given parameters
// Note, If `code` is not EAPRequest or EAPResponse, `eapType`
// must have value of (verified at runtime)
func NewPacket(code Code, eapType EAPType, identifier int, data []byte) (*Packet, error) {
	// Enforce comment in section 5 of RFC 3748
	if (code != CodeRESPONSE) && (eapType == EAPTypeNAK) {
		return nil, fmt.Errorf(
			"EAP Type NAK in a packet with code other then Response is invalid (code=%d, eap_type=%d, id=%d)",
			int(code), int(eapType), identifier,
		)
	}

	// Assert type is applied on REQ and RES only
	if !code.IsRequestOrResponse() && eapType != EAPTypeNONE {
		return nil, fmt.Errorf(
			"EAP Packet Type field exists only for EAP Request and Response messages (code=%d, eap_type=%d, id=%d)",
			int(code), int(eapType), identifier,
		)
	}

	return &Packet{
		code,
		eapType,
		identifier,
		data,
	}, nil
}

// NewPacketFromRadius creates an eap.Packet from the given RADIUS packet
func NewPacketFromRadius(r *radius.Packet) (*Packet, error) {
	if r == nil {
		return nil, errors.New("got nil radius packet")
	}

	eapMessage := r.Get(rfc2869.EAPMessage_Type)
	if eapMessage == nil {
		return nil, errors.New("no EAP-Message attribute found")
	}

	return NewPacketFromRaw(eapMessage)
}

// NewPacketFromRaw Parses the given EAP packet and creates an object to represent this packet
func NewPacketFromRaw(b []byte) (*Packet, error) {
	if len(b) < EapMinLegalPacketLength {
		return nil, fmt.Errorf(
			"packet length must be at least %d bytes, got %d bytes",
			EapMinLegalPacketLength,
			len(b),
		)
	}

	// Extract code and type
	eapCode := Code(b[EapCodeOffset])
	if !eapCode.IsValid() {
		return nil, fmt.Errorf("invalid eap packet code '%d'", b[EapCodeOffset])
	}

	var eapType = EAPTypeNONE
	if eapCode.IsRequestOrResponse() {
		eapType = EAPType(b[EapTypeOffset])
		if !eapType.IsValid() {
			return nil, fmt.Errorf("invalid eap packet type '%d'", b[EapTypeOffset])
		}
	}

	// Extract and verify the length
	length := (int(b[EapLengthMsbOffset]) << 8) | int(b[EapLengthLsbOffset])
	if length != len(b) {
		return nil, fmt.Errorf(
			"length mismatch (packet header indicates %d, but packet contains %d data bytes)",
			length,
			len(b),
		)
	}

	identifier := int(b[EapIdentifierOffset])

	lenToRemove := EapPacketHeaderLength
	if eapCode.IsRequestOrResponse() {
		lenToRemove += EapPacketTypeFieldLength
	}

	return &Packet{
		eapCode,
		eapType,
		identifier,
		b[lenToRemove:],
	}, nil
}

// Bytes returns the bytes representation of the EAP packet
func (p Packet) Bytes() ([]byte, error) {
	eapTypeRequired := p.Code.IsRequestOrResponse()

	// Calculate the designated
	length := p.Length()

	// Build header
	raw := []byte{
		byte(p.Code),
		byte(p.Identifier),
		byte((length & 0xFF00) >> 8),
		byte(length & 0x00FF),
	}

	// Add type if exists
	if eapTypeRequired {
		if p.EAPType == EAPTypeNONE {
			return nil, errors.New("eap response/request must have a type, but NONE was set")
		}
		raw = append(raw, byte(p.EAPType))
	}

	// Append the data
	raw = append(raw, p.Data...)
	return raw, nil
}

// Length Returns the (calculated) length filed of the packet
func (p Packet) Length() int {
	result := EapPacketHeaderLength + len(p.Data)
	if p.Code.IsRequestOrResponse() {
		result += EapPacketTypeFieldLength
	}
	return result
}
