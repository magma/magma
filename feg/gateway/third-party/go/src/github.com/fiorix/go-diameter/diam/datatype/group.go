// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// Grouped data type.
type Grouped []byte

// DecodeGrouped decodes a Grouped data type from byte array.
func DecodeGrouped(b []byte) (Type, error) {
	return Grouped(b), nil
}

// Serialize implements the Type interface.
func (g Grouped) Serialize() []byte {
	return g
}

// Len implements the Type interface.
func (g Grouped) Len() int {
	return len(g)
}

// Padding implements the Type interface.
func (g Grouped) Padding() int {
	return 0
}

// Type implements the Type interface.
func (g Grouped) Type() TypeID {
	return GroupedType
}

// String implements the Type interface.
func (g Grouped) String() string {
	return fmt.Sprint("Grouped{...}")
}
