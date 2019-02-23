// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// DiameterIdentity data type.
type DiameterIdentity OctetString

// DecodeDiameterIdentity decodes a DiameterIdentity from byte array.
func DecodeDiameterIdentity(b []byte) (Type, error) {
	return DiameterIdentity(b), nil
}

// Serialize implements the Type interface.
func (s DiameterIdentity) Serialize() []byte {
	return []byte(s)
}

// Len implements the Type interface.
func (s DiameterIdentity) Len() int {
	return len(s)
}

// Padding implements the Type interface.
func (s DiameterIdentity) Padding() int {
	l := len(s)
	return pad4(l) - l
}

// Type implements the Type interface.
func (s DiameterIdentity) Type() TypeID {
	return DiameterIdentityType
}

// String implements the Type interface.
func (s DiameterIdentity) String() string {
	return fmt.Sprintf("DiameterIdentity{%s},Padding:%d", string(s), s.Padding())
}
