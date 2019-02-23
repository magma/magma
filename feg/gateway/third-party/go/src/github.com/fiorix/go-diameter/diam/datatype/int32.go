// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"encoding/binary"
	"fmt"
)

// Integer32 data type.
type Integer32 int32

// DecodeInteger32 decodes an Integer32 data type from byte array.
func DecodeInteger32(b []byte) (Type, error) {
	if len(b) != 4 {
		return Integer32(0), nil
	}
	return Integer32(binary.BigEndian.Uint32(b)), nil
}

// Serialize implements the Type interface.
func (n Integer32) Serialize() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

// Len implements the Type interface.
func (n Integer32) Len() int {
	return 4
}

// Padding implements the Type interface.
func (n Integer32) Padding() int {
	return 0
}

// Type implements the Type interface.
func (n Integer32) Type() TypeID {
	return Integer32Type
}

// String implements the Type interface.
func (n Integer32) String() string {
	return fmt.Sprintf("Integer32{%d}", n)
}
