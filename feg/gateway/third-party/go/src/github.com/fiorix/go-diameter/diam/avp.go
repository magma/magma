// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

// AVP is a Diameter attribute-value-pair.
type AVP struct {
	Code     uint32        // Code of this AVP
	Flags    uint8         // Flags of this AVP
	Length   int           // Length of this AVP's payload
	VendorID uint32        // VendorId of this AVP
	Data     datatype.Type // Data of this AVP (payload)
}

// NewAVP creates and initializes a new AVP.
func NewAVP(code uint32, flags uint8, vendor uint32, data datatype.Type) *AVP {
	a := &AVP{
		Code:     code,
		Flags:    flags,
		VendorID: vendor,
		Data:     data,
	}
	a.Length = a.headerLen() + a.Data.Len() // no padding length
	if vendor > 0 && flags&avp.Vbit != avp.Vbit {
		a.Flags |= avp.Vbit
	}
	return a
}

// DecodeAVP decodes the bytes of a Diameter AVP.
// It uses the given application id and dictionary for decoding the bytes.
func DecodeAVP(data []byte, application uint32, dictionary *dict.Parser) (*AVP, error) {
	avp := &AVP{}
	if err := avp.DecodeFromBytes(data, application, dictionary); err != nil {
		return avp, err
	}
	return avp, nil
}

// DecodeFromBytes decodes the bytes of a Diameter AVP.
// It uses the given application id and dictionary for decoding the bytes.
func (a *AVP) DecodeFromBytes(data []byte, application uint32, dictionary *dict.Parser) error {
	if len(data) < 8 {
		return fmt.Errorf("Not enough data to decode AVP header: %d bytes", len(data))
	}
	a.Code = binary.BigEndian.Uint32(data[0:4])
	a.Flags = data[4]
	a.Length = int(uint24to32(data[5:8]))
	if len(data) < a.Length {
		return fmt.Errorf("Not enough data to decode AVP: %d != %d",
			len(data), a.Length)
	}
	data = data[:a.Length] // this cuts padded bytes off
	if len(data) < 8 {
		return fmt.Errorf("Not enough data to decode AVP header: %d bytes", len(data))
	}

	var hdrLength int
	var payload []byte
	// Read VendorId when required.
	if a.Flags&avp.Vbit == avp.Vbit {
		a.VendorID = binary.BigEndian.Uint32(data[8:12])
		payload = data[12:]
		hdrLength = 12
	} else {
		payload = data[8:]
		hdrLength = 8
	}
	// Find this code in the dictionary.
	dictAVP, err := dictionary.FindAVPWithVendor(application, a.Code, a.VendorID)
	if err != nil && dictAVP == nil {
		return err
	}
	bodyLen := a.Length - hdrLength
	if n := len(payload); n < bodyLen {
		return fmt.Errorf(
			"Not enough data to decode AVP: %d != %d",
			hdrLength, n,
		)
	}
	a.Data, err = datatype.Decode(dictAVP.Data.Type, payload)
	if err != nil {
		return err
	}
	// Handle grouped AVPs.
	if a.Data.Type() == datatype.GroupedType {
		a.Data, err = DecodeGrouped(
			a.Data.(datatype.Grouped),
			application, dictionary,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// Serialize returns the byte sequence that represents this AVP.
// It requires at least the Code, Flags and Data fields set.
func (a *AVP) Serialize() ([]byte, error) {
	if a.Data == nil {
		return nil, errors.New("Failed to serialize AVP: Data is nil")
	}
	b := make([]byte, a.Len())
	err := a.SerializeTo(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// SerializeTo writes the byte sequence that represents this AVP to a byte array.
func (a *AVP) SerializeTo(b []byte) error {
	if a.Data == nil {
		return errors.New("Failed to serialize AVP: Data is nil")
	}
	binary.BigEndian.PutUint32(b[0:4], a.Code)
	b[4] = a.Flags
	hl := a.headerLen()
	copy(b[5:8], uint32to24(uint32(hl+a.Data.Len())))
	if a.Flags&avp.Vbit == avp.Vbit {
		binary.BigEndian.PutUint32(b[8:12], a.VendorID)
	}
	payload := a.Data.Serialize()
	copy(b[hl:], payload)
	// reset padding bytes
	b = b[hl+len(payload):]
	for i := 0; i < a.Data.Padding(); i++ {
		b[i] = 0
	}
	return nil
}

// Len returns the length of this AVP in bytes with padding.
func (a *AVP) Len() int {
	return a.headerLen() + a.Data.Len() + a.Data.Padding()
}

func (a *AVP) headerLen() int {
	if a.Flags&avp.Vbit == avp.Vbit {
		return 12
	}
	return 8
}

func (a *AVP) String() string {
	return fmt.Sprintf("{Code:%d,Flags:0x%x,Length:%d,VendorId:%d,Value:%s}",
		a.Code,
		a.Flags,
		a.Len(),
		a.VendorID,
		a.Data,
	)
}
