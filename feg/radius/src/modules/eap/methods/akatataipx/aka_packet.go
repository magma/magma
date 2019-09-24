package akatataipx

import (
	"errors"
	eap "fbc/cwf/radius/modules/eap/packet"
	"fmt"
)

// AkaSubtype the subtypes of AKA packet as defined in RFC4187 section 11
type AkaSubtype int

//@no-lint
const (
	AkaChallenge              AkaSubtype = 1
	AkaAuthenticationReject   AkaSubtype = 2
	AkaSynchronizationFailure AkaSubtype = 4
	AkaIdentity               AkaSubtype = 5
	SimStart                  AkaSubtype = 10
	SimChallenge              AkaSubtype = 11
	AkaNotification           AkaSubtype = 12
	AkaReauthentication       AkaSubtype = 13
	AkaClientError            AkaSubtype = 14
)

// AttributeType as enumerated in RFC4187 section 11
type AttributeType uint8

//nolint
const (
	AT_RAND              AttributeType = 1
	AT_AUTN              AttributeType = 2
	AT_RES               AttributeType = 3
	AT_AUTS              AttributeType = 4
	AT_PADDING           AttributeType = 6
	AT_NONCE_MT          AttributeType = 7
	AT_PERMANENT_ID_REQ  AttributeType = 10
	AT_MAC               AttributeType = 11
	AT_NOTIFICATION      AttributeType = 12
	AT_ANY_ID_REQ        AttributeType = 13
	AT_IDENTITY          AttributeType = 14
	AT_VERSION_LIST      AttributeType = 15
	AT_SELECTED_VERSION  AttributeType = 16
	AT_FULLAUTH_ID_REQ   AttributeType = 17
	AT_COUNTER           AttributeType = 19
	AT_COUNTER_TOO_SMALL AttributeType = 20
	AT_NONCE_S           AttributeType = 21
	AT_CLIENT_ERROR_CODE AttributeType = 22
	AT_IV                AttributeType = 129
	AT_ENCR_DATA         AttributeType = 130
	AT_NEXT_PSEUDONYM    AttributeType = 132
	AT_NEXT_REAUTH_ID    AttributeType = 133
	AT_CHECKCODE         AttributeType = 134
	AT_RESULT_IND        AttributeType = 135
)

// Attribute AKA packet attribute as described in RFC 4187 section 8.1
type Attribute struct {
	Type  AttributeType
	Value []byte
}

// AkaPacket the structure of an AKA packet as described in RFC 4187 section 8.1
type AkaPacket struct {
	Subtype    AkaSubtype
	Reserved   uint16
	Attributes []Attribute
}

// Bytes converts the AKA packet to its on-the-wire format
func (p AkaPacket) Bytes() []byte {
	result := []byte{byte(p.Subtype), byte((p.Reserved << 8) & 0xff), byte(p.Reserved & 0xff)}
	for _, attr := range p.Attributes {
		l := 2 + len(attr.Value)
		pad := (4 - (l & 3)) & 3
		result = append(result, byte(attr.Type), byte((l+pad)>>2))
		result = append(result, attr.Value...)
		if pad > 0 {
			result = append(result, make([]byte, pad)...)
		}
	}
	return result
}

// GetFirst returns the first attribute of the given type or nil if not exists
func (p AkaPacket) GetFirst(t AttributeType) *Attribute {
	for _, a := range p.Attributes {
		if a.Type == t {
			//nolint
			return &a
		}
	}
	return nil
}

// GetAll returns all attributes of the given type or nil if not exists
func (p AkaPacket) GetAll(t AttributeType) []Attribute {
	result := make([]Attribute, 0)
	for _, a := range p.Attributes {
		if a.Type == t {
			result = append(result, a)
		}
	}
	return result
}

// NewAkaPacket ...
func NewAkaPacket(data []byte) (*AkaPacket, error) {
	if len(data) < 3 {
		return nil, errors.New("packet len is too short")
	}

	// Create response
	result := &AkaPacket{
		Subtype:    AkaSubtype(data[0]),
		Reserved:   uint16(data[1])>>8 | uint16(data[2]),
		Attributes: []Attribute{},
	}

	// Parse Attributes
	p := 3
	for {
		// Break when we get to the end of the packet
		if p+1 >= len(data) {
			break
		}

		// Verify attribute length
		l := int(data[p+1]) * 4
		if p+l > len(data) {
			return nil, fmt.Errorf("attribute %d of length %d exceeds packet length (p=%d, len=%d)", data[p], l, p, len(data))
		}

		// TODO: verify packet further
		// (1) section 10.1: mutual exclusive attributes
		// (2) section 8.1: unknown attributes must error out and terminate

		// Append attribute to the packet
		attr := Attribute{
			Type:  AttributeType(data[p]),
			Value: data[p+2 : p+l],
		}
		result.Attributes = append(result.Attributes, attr)

		// Advance the pointer!
		p += l
	}

	return result, nil
}

// AppendMac adds AT_MAC attribute to the AKA packet
func (p *AkaPacket) AppendMac(eapCode eap.Code, identifier int, kaut []byte) error {
	// verify no AT_MAC attributes are there
	for _, attr := range p.Attributes {
		if attr.Type == AT_MAC {
			return errors.New("packet already contains AT_MAC attribute")
		}
	}

	// Append a new attribute
	p.Attributes = append(p.Attributes, Attribute{
		Type:  AT_MAC,
		Value: make([]byte, 18),
	})

	// Calculate AT_MAC
	// TODO: This code is a bit "hackey" - that is, it does not use any object model,
	// shared structures etc, and it uses "magic numbers" quite a bit.
	// However, those numbers & structures are well-established in the RFC set used
	// to implement this code, so it's good enough for POC-level code
	b := p.Bytes()
	totalLen := 5 + len(b)
	eapPacket := append(
		[]byte{
			byte(eapCode),
			byte(identifier),
			byte(((totalLen & 0xFF00) >> 8)),
			byte(totalLen & 0xFF),
			byte(eap.EAPTypeAKA),
		},
		b...,
	)
	mac := GenMac(eapPacket, kaut)

	// Write AT_MAC value back to packet
	updated := false
	for _, attr := range p.Attributes {
		if attr.Type == AT_MAC {
			copy(attr.Value[2:], mac)
			updated = true
			break
		}
	}
	if !updated {
		return errors.New("something bad happened. AT_MAC attribute was not found")
	}
	return nil
}
