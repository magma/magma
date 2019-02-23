// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// DiameterURI data type.
type DiameterURI OctetString

// DecodeDiameterURI decodes a DiameterURI from byte array.
func DecodeDiameterURI(b []byte) (Type, error) {
	return DiameterURI(OctetString(b)), nil
}

// Serialize implements the Type interface.
func (s DiameterURI) Serialize() []byte {
	return OctetString(s).Serialize()
}

// Len implements the Type interface.
func (s DiameterURI) Len() int {
	return len(s)
}

// Padding implements the Type interface.
func (s DiameterURI) Padding() int {
	l := len(s)
	return pad4(l) - l
}

// Type implements the Type interface.
func (s DiameterURI) Type() TypeID {
	return DiameterURIType
}

// String implements the Type interface.
func (s DiameterURI) String() string {
	return fmt.Sprintf("DiameterURI{%s},Padding:%d", string(s), s.Padding())
}
