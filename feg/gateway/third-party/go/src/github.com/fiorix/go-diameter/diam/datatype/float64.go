// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"encoding/binary"
	"fmt"
	"math"
)

// Float64 data type.
type Float64 float64

// DecodeFloat64 decodes a Float64 data type from byte array.
func DecodeFloat64(b []byte) (Type, error) {
	if len(b) != 8 {
		return Float64(0), nil
	}
	return Float64(math.Float64frombits(binary.BigEndian.Uint64(b))), nil
}

// Serialize implements the Type interface.
func (n Float64) Serialize() []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(float64(n)))
	return b
}

// Len implements the Type interface.
func (n Float64) Len() int {
	return 8
}

// Padding implements the Type interface.
func (n Float64) Padding() int {
	return 0
}

// Type implements the Type interface.
func (n Float64) Type() TypeID {
	return Float64Type
}

// String implements the Type interface.
func (n Float64) String() string {
	return fmt.Sprintf("Float64{%0.4f}", n)
}
