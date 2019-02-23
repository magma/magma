// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"encoding/binary"
	"fmt"
)

// Unsigned32 data type.
type Unsigned32 uint32

// DecodeUnsigned32 decodes an Unsigned32 data type from byte array.
func DecodeUnsigned32(b []byte) (Type, error) {
	if len(b) != 4 {
		return Unsigned32(0), nil
	}
	return Unsigned32(binary.BigEndian.Uint32(b)), nil
}

// Serialize implements the Type interface.
func (n Unsigned32) Serialize() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

// Len implements the Type interface.
func (n Unsigned32) Len() int {
	return 4
}

// Padding implements the Type interface.
func (n Unsigned32) Padding() int {
	return 0
}

// Type implements the Type interface.
func (n Unsigned32) Type() TypeID {
	return Unsigned32Type
}

// String implements the Type interface.
func (n Unsigned32) String() string {
	return fmt.Sprintf("Unsigned32{%d}", n)
}
