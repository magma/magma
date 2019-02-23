// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// OctetString data type.
type OctetString string

// DecodeOctetString decodes an OctetString from byte array.
func DecodeOctetString(b []byte) (Type, error) {
	return OctetString(b), nil
}

// Serialize implements the Type interface.
func (s OctetString) Serialize() []byte {
	return []byte(s)
}

// Len implements the Type interface.
func (s OctetString) Len() int {
	return len(s)
}

// Padding implements the Type interface.
func (s OctetString) Padding() int {
	l := len(s)
	return pad4(l) - l
}

// Type implements the Type interface.
func (s OctetString) Type() TypeID {
	return OctetStringType
}

// String implements the Type interface.
func (s OctetString) String() string {
	return fmt.Sprintf("OctetString{%#x},Padding:%d", string(s), s.Padding())
}
