package eapauth

import (
	"fmt"
	"io"
)

// attribute EAP Method's attribute implementation
// attribute provides the following interface:
//
// type Attribute interface {
// 	   Type() uint8
// 	   Value() []byte
// 	   String() string
// 	   Len() int
//     Marshal() []byte
// }
type attribute []byte

// NewAttribute creates and returns new attribute of given type (typ) & value
// the new attribute is padded with zeros to 4 byte boundary
func NewAttribute(typ uint8, data []byte) attribute {
	ld := len(data)
	l := 2 + ld
	pad := (4 - l&3) & 3
	l += pad
	res := make([]byte, 2, l)
	res[0], res[1] = typ, uint8(l<<2)
	if ld > 0 {
		res = append(res, data...)
	}
	if pad > 0 {
		res = append(res, make([]byte, pad)...)
	}
	return res
}

// String - implements stringer interface
func (a attribute) String() string {
	if a == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Type: %d, Len: %d, Value: %v", a.Type(), a.Len(), a.Value())
}

// Type returns attribute type byte
func (a attribute) Type() uint8 {
	return a[0]
}

// Type returns attribute type byte
func (a attribute) Value() []byte {
	return a[2:]
}

// Len returns total serialized length of the attribute (the size that the attribute will occupy in EAP packet)
func (a attribute) Len() int {
	return len(a)
}

// Marshal serialized attribute (the form that the attribute will occupy in EAP packet)
func (a attribute) Marshal() []byte {
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

func NewAttributeScanner(eapData []byte) (*attributeScanner, error) {
	el := len(eapData)
	if el <= EapFirstAttribute {
		return nil, fmt.Errorf("EAP Message is too short (%d bytes) to include attributes", el)
	}
	dataLen := el - EapFirstAttribute
	return &attributeScanner{
		subtype:       eapData[EapSubtype],
		len:           dataLen,
		lastAttrStart: dataLen - 2,
		data:          eapData[EapFirstAttribute:]}, nil
}

// Next Returns next available attribute (if any) and adjusts internal reference to point past it
func (sc *attributeScanner) Next() (attribute, error) {
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
	return sc.data[attrStart:sc.current], nil
}

// Reset resets the scanner to the beginning of the attributes
func (sc *attributeScanner) Reset() {
	sc.current = 0
}

// Append adds given Attribute to EAP Packet and amends the packet's header to reflect the new EAP packet size
func (eap Packet) Append(a attribute) (Packet, error) {
	if a == nil {
		return eap, fmt.Errorf("Nil Attribute")
	}
	eapLen := len(eap)
	mLen := uint(eap[EapMsgLenHigh])<<8 + uint(eap[EapMsgLenLow])
	if mLen > uint(eapLen) {
		return eap, fmt.Errorf("Invalid EAP Length, header: %d, actual: %d", mLen, eapLen)
	}
	if mLen < uint(EapFirstAttribute) {
		return eap, fmt.Errorf("Insufficient EAP length, header: %d, data: %d", mLen, eapLen)
	}
	if uint(eapLen) > mLen {
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
	mLen += uint(alen << 2)
	if mLen > EapMaxLen {
		return eap, fmt.Errorf("EAP Len would exceeds %d bytes: %d", EapMaxLen, mLen)
	}
	eap = append(eap, a...)
	if l > 0 { // attribute was not 4 bytes aligned
		eap = append(eap, make([]byte, l)...)
	}
	eap[eapLen+1] = uint8(alen)
	// Update EAP length
	eap[EapMsgLenLow] = uint8(mLen)
	eap[EapMsgLenHigh] = uint8(mLen >> 8)
	return eap, nil
}
