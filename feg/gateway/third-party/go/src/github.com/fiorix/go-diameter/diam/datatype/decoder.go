// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// DecoderFunc is an adapter to decode a byte array to an AVP data type.
type DecoderFunc func([]byte) (Type, error)

// Decoder is a map of AVP data types indexed by TypeID.
var Decoder = map[TypeID]DecoderFunc{
	UnknownType:          DecodeUnknown,
	AddressType:          DecodeAddress,
	DiameterIdentityType: DecodeDiameterIdentity,
	DiameterURIType:      DecodeDiameterURI,
	EnumeratedType:       DecodeEnumerated,
	Float32Type:          DecodeFloat32,
	Float64Type:          DecodeFloat64,
	GroupedType:          DecodeGrouped,
	IPFilterRuleType:     DecodeIPFilterRule,
	IPv4Type:             DecodeIPv4,
	Integer32Type:        DecodeInteger32,
	Integer64Type:        DecodeInteger64,
	OctetStringType:      DecodeOctetString,
	TimeType:             DecodeTime,
	UTF8StringType:       DecodeUTF8String,
	Unsigned32Type:       DecodeUnsigned32,
	Unsigned64Type:       DecodeUnsigned64,
}

// Decode decodes a specific AVP data type from byte array to a DataType.
func Decode(Type TypeID, b []byte) (Type, error) {
	f, exists := Decoder[Type]
	if !exists {
		return nil, fmt.Errorf("Unknown data type: %d", Type)
	}
	return f(b)
}
