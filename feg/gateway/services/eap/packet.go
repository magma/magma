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

// Package eap (EAP Authenticator) provides interface to supported & registered EAP Authenticator Providers
package eap

import (
	"fmt"
	"io"
)

// Packet represents EAP Packet
type Packet []byte

// NewPacket creates an EAP Packet with initialized header and appends provided data
// if additionalCapacity is specified, NewPacket reserves extra additionalCapacity bytes capacity in the
// returned packet byte slice
func NewPacket(code, identifier uint8, data []byte, additionalCapacity ...uint) Packet {
	l := len(data) + EapHeaderLen
	packetCap := l
	if len(additionalCapacity) > 0 && l < int(EapMaxLen) {
		ac := additionalCapacity[0]
		packetCap = l + int(ac)
		if packetCap > int(EapMaxLen) {
			packetCap = int(EapMaxLen)
		}
	}
	p := make([]byte, EapHeaderLen, packetCap)
	if l > EapHeaderLen {
		p = append(p, data...)
	}
	p[EapMsgCode], p[EapMsgIdentifier], p[EapMsgLenLow], p[EapMsgLenHigh] = code, identifier, uint8(l), uint8(l>>8)
	return p
}

// NewPreallocatedPacket creates an EAP Packet from/in passed data slice & initializes its header
func NewPreallocatedPacket(identifier uint8, data []byte) (Packet, error) {
	l := len(data)
	if l < EapHeaderLen {
		return nil, fmt.Errorf("Data is too short: %d, must be at least %d bytes", l, EapHeaderLen)
	}
	p := Packet(data)
	p[EapMsgIdentifier], p[EapMsgLenLow], p[EapMsgLenHigh] = identifier, uint8(l), uint8(l>>8)
	return p, nil
}

// Validate verifies that the packet is not nil & it's length is correct
func (p Packet) Validate() error {
	lp := len(p)
	if lp < EapHeaderLen {
		return io.ErrShortBuffer
	}
	if p.Len() != lp {
		return fmt.Errorf("Invalid Packet Length: header => %d, actual => %d", p.Len(), lp)
	}
	return nil
}

// Len returns EAP Packet length derived from its header (vs. len of []byte)
func (p Packet) Len() int {
	return (int(p[EapMsgLenHigh]) << 8) + int(p[EapMsgLenLow])
}

// Identifier returns EAP Message Identifier
func (p Packet) Identifier() uint8 {
	return p[EapMsgIdentifier]
}

// Code returns EAP Packet Message Code
func (p Packet) Code() uint8 {
	if len(p) > EapMsgCode {
		return p[EapMsgCode]
	}
	return UndefinedCode
}

// IsSuccess returns if the EAP Packet Code is Success (3)
func (p Packet) IsSuccess() bool {
	return len(p) > EapMsgCode && p[EapMsgCode] == SuccessCode
}

// Type returns EAP Method Type or 0 - reserved if not available
func (p Packet) Type() uint8 {
	if len(p) <= EapMsgMethodType {
		return 0
	}
	return p[EapMsgMethodType]
}

// Truncate truncates EAP Packet to its header defined length
func (p Packet) Truncate() Packet {
	mLen := p.Len()
	if len(p) > mLen {
		p = p[:mLen]
	}
	return p
}

// TypeData returns EAP Packet Type data part
func (p Packet) TypeData() []byte {
	l := p.Len()
	if l < EapMsgData {
		return nil
	}
	return p[EapMsgData:l]
}

// TypeDataUnsafe - same as TypeData, but doesn't check the packet length
func (p Packet) TypeDataUnsafe() []byte {
	return p[EapMsgData : (int(p[EapMsgLenHigh])<<8)+int(p[EapMsgLenLow])]
}

// Failure returns RFC 3748, 4.2 Failure packet with Identifier set from p
func (p Packet) Failure() Packet {
	var identifier uint8
	if len(p) > EapMsgIdentifier {
		identifier = p[EapMsgIdentifier]
	}
	// Return RFC 3748 p4.2 EAP Failure packet
	//  0                   1                   2                   3
	//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	// |     Code      |  Identifier   |            Length             |
	// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	return []byte{
		FailureCode, // Code
		identifier,  // Identifier
		0, 4}        // Length
}
