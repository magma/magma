// Copyright 2013-2018 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// Unknown data type.
type Unknown []byte

// DecodeUnknown decodes an Unknown from byte array.
func DecodeUnknown(b []byte) (Type, error) {
	return Unknown(b), nil
}

// Serialize implements the Type interface.
func (u Unknown) Serialize() []byte {
	return []byte(u)
}

// Len implements the Type interface.
func (u Unknown) Len() int {
	return len(u)
}

// Padding implements the Type interface.
func (u Unknown) Padding() int {
	l := len(u)
	return pad4(l) - l
}

// Type implements the Type interface.
func (u Unknown) Type() TypeID {
	return UnknownType
}

// String implements the Type interface.
func (u Unknown) String() string {
	return fmt.Sprintf("Unknown{%#x},Padding:%d", string(u), u.Padding())
}
