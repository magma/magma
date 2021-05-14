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

package eap

import (
	"fmt"
	"io"
)

type AttrType uint8

// attribute EAP Method's attribute implementation
// attribute provides the following interface:
//
type Attribute interface {
	Type() AttrType
	Value() []byte
	String() string
	Len() int
	AttrLen() uint8
	Marshaled() []byte
}

type attribute []byte

// NewAttribute creates and returns new attribute of given type (typ) & value
// the new attribute is padded with zeros to 4 byte boundary
func NewAttribute(typ AttrType, data []byte) attribute {
	ld := len(data)
	l := 2 + ld
	pad := (4 - l&3) & 3
	l += pad
	res := make([]byte, 2, l)
	res[0], res[1] = byte(typ), byte(l<<2)
	if ld > 0 {
		res = append(res, data...)
	}
	if pad > 0 {
		res = append(res, make([]byte, pad)...)
	}
	return res
}

// NewRawAttribute 'casts' given slice to Attribute type without any validations
func NewRawAttribute(data []byte) attribute {
	return data
}

// String - implements stringer interface
func (a attribute) String() string {
	if a == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Type: %d, Len: %d, Value: %v", a.Type(), a.Len(), a.Value())
}

// Type returns attribute type byte
func (a attribute) Type() AttrType {
	return AttrType(a[0])
}

// Value returns attribute type bytes
func (a attribute) Value() []byte {
	return a[2:]
}

// AttrLen returns the length byte from the attribute header
func (a attribute) AttrLen() uint8 {
	if a == nil || len(a) < 1 {
		return 0
	}
	return a[1]
}

// Len returns total serialized length of the attribute (the size that the attribute will occupy in EAP packet)
func (a attribute) Len() int {
	return len(a)
}

// Marshal serialized attribute (the form that the attribute will occupy in EAP packet)
func (a attribute) Marshaled() []byte {
	return a
}

// attributeScanner holds the internal state while processing
// a given EAP Provider payload.
type attributeScanner struct {
	subtype       uint8
	len           int // total data len
	lastAttrStart int // last viable attribute start offset
	data          []byte
	current       int
}

func NewAttributeScanner(eapData Packet) (*attributeScanner, error) {
	el := len(eapData)
	if el <= EapFirstAttribute {
		return nil, fmt.Errorf("EAP Message is too short (%d bytes) to include attributes", el)
	}
	if err := eapData.Validate(); err != nil {
		return nil, err
	}
	dataLen := el - EapFirstAttribute
	return &attributeScanner{
		subtype:       eapData[EapSubtype],
		len:           dataLen,
		lastAttrStart: dataLen - 2,
		data:          eapData[EapFirstAttribute:]}, nil
}

// Next Returns next available attribute (if any) and adjusts internal reference to point past it
func (sc *attributeScanner) Next() (Attribute, error) {
	attrStart := sc.current
	if attrStart >= sc.len {
		return nil, io.EOF
	}
	if attrStart > sc.lastAttrStart {
		return nil, fmt.Errorf(
			"Attribute Offset %d is past Last Viable Attribute position %d",
			attrStart, sc.lastAttrStart)
	}
	attrLen := int(sc.data[attrStart+1]) << 2
	if attrLen == 0 {
		attrLen = 1024 // max attr len according to RFC
	}
	if attrStart+attrLen > sc.len {
		return nil, fmt.Errorf(
			"Expected Attribute (%d:%d, len:%d) is past the EAP end %d",
			attrStart, attrStart+attrLen, attrLen, sc.len)
	}
	sc.current += attrLen
	return attribute(sc.data[attrStart:sc.current]), nil
}

// Reset resets the scanner to the beginning of the attributes
func (sc *attributeScanner) Reset() {
	sc.current = 0
}

// Append adds given Attribute to EAP Packet and amends the packet's header to reflect the new EAP packet size
func (eap Packet) Append(a Attribute) (Packet, error) {
	if a == nil {
		return eap, fmt.Errorf("Nil Attribute")
	}
	eapLen := len(eap)
	mLen := eap.Len()
	if mLen > eapLen {
		return eap, fmt.Errorf("Invalid EAP Length, header: %d, actual: %d", mLen, eapLen)
	}
	if mLen < EapFirstAttribute {
		return eap, fmt.Errorf("Insufficient EAP length, header: %d, data: %d", mLen, eapLen)
	}
	if eapLen > mLen {
		eap = eap[:mLen]
	}
	l := a.Len()
	if l < 2 {
		return eap, fmt.Errorf("Insufficient EAP length: %d", l)
	}
	if l > 1024 {
		return eap, fmt.Errorf("Attribute exceeds 1024 bytes: %d", l)
	}
	alen := l >> 2
	l &= 3
	if l != 0 {
		alen += 1
		l = 4 - l
	}
	mLen += alen << 2
	if mLen > int(EapMaxLen) {
		return eap, fmt.Errorf("EAP Len would exceeds %d bytes: %d", EapMaxLen, mLen)
	}
	eap = append(eap, a.Marshaled()...)
	if l > 0 { // attribute was not 4 bytes aligned
		eap = append(eap, make([]byte, l)...)
	}
	eap[eapLen+1] = uint8(alen)
	// Update EAP length
	eap[EapMsgLenLow] = uint8(mLen)
	eap[EapMsgLenHigh] = uint8(mLen >> 8)
	return eap, nil
}
