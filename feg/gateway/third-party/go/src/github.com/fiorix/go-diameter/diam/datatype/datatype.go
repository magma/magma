// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

// Type is an interface to support Diameter AVP data types.
type Type interface {
	Serialize() []byte
	Len() int
	Padding() int
	Type() TypeID
	String() string
}

// TypeID is the identifier of an AVP data type.
type TypeID int

// List of available AVP data types.
const (
	UnknownType TypeID = iota
	AddressType
	DiameterIdentityType
	DiameterURIType
	EnumeratedType
	Float32Type
	Float64Type
	GroupedType
	IPFilterRuleType
	IPv4Type
	Integer32Type
	Integer64Type
	OctetStringType
	QoSFilterRuleType
	TimeType
	UTF8StringType
	Unsigned32Type
	Unsigned64Type
)

// Available is a map of data types available, indexed by name.
var Available = map[string]TypeID{
	"Address":          AddressType,
	"DiameterIdentity": DiameterIdentityType,
	"DiameterURI":      DiameterURIType,
	"Enumerated":       EnumeratedType,
	"Float32":          Float32Type,
	"Float64":          Float64Type,
	"Grouped":          GroupedType,
	"IPFilterRule":     IPFilterRuleType,
	"IPv4":             IPv4Type,
	"Integer32":        Integer32Type,
	"Integer64":        Integer64Type,
	"OctetString":      OctetStringType,
	"QoSFilterRule":    QoSFilterRuleType,
	"Time":             TimeType,
	"UTF8String":       UTF8StringType,
	"Unsigned32":       Unsigned32Type,
	"Unsigned64":       Unsigned64Type,
}
