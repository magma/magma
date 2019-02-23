// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"encoding/binary"
	"fmt"
)

// Integer64 data type.
type Integer64 int64

// DecodeInteger64 decodes an Integer64 data type from byte array.
func DecodeInteger64(b []byte) (Type, error) {
	if len(b) != 8 {
		return Integer64(0), nil
	}
	return Integer64(binary.BigEndian.Uint64(b)), nil
}

// Serialize implements the Type interface.
func (n Integer64) Serialize() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	return b
}

// Len implements the Type interface.
func (n Integer64) Len() int {
	return 8
}

// Padding implements the Type interface.
func (n Integer64) Padding() int {
	return 0
}

// Type implements the Type interface.
func (n Integer64) Type() TypeID {
	return Integer64Type
}

// String implements the Type interface.
func (n Integer64) String() string {
	return fmt.Sprintf("Integer64{%d}", n)
}
