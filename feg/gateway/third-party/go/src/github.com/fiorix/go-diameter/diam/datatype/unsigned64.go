// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"encoding/binary"
	"fmt"
)

// Unsigned64 data type.
type Unsigned64 uint64

// DecodeUnsigned64 decodes an Unsigned64 data type from byte array.
func DecodeUnsigned64(b []byte) (Type, error) {
	if len(b) != 8 {
		return Unsigned64(0), nil
	}
	return Unsigned64(binary.BigEndian.Uint64(b)), nil
}

// Serialize implements the Type interface.
func (n Unsigned64) Serialize() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	return b
}

// Len implements the Type interface.
func (n Unsigned64) Len() int {
	return 8
}

// Padding implements the Type interface.
func (n Unsigned64) Padding() int {
	return 0
}

// Type implements the Type interface.
func (n Unsigned64) Type() TypeID {
	return Unsigned64Type
}

// String implements the Type interface.
func (n Unsigned64) String() string {
	return fmt.Sprintf("Unsigned64{%d}", n)
}
