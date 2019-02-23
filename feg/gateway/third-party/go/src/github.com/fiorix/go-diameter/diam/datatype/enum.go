// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// Enumerated data type.
type Enumerated Integer32

// DecodeEnumerated decodes an Enumerated data type from byte array.
func DecodeEnumerated(b []byte) (Type, error) {
	v, err := DecodeInteger32(b)
	if err != nil {
		return nil, err
	}
	return Enumerated(v.(Integer32)), nil
}

// Serialize implements the Type interface.
func (n Enumerated) Serialize() []byte {
	return Integer32(n).Serialize()
}

// Len implements the Type interface.
func (n Enumerated) Len() int {
	return 4
}

// Padding implements the Type interface.
func (n Enumerated) Padding() int {
	return 0
}

// Type implements the Type interface.
func (n Enumerated) Type() TypeID {
	return EnumeratedType
}

// String implements the Type interface.
func (n Enumerated) String() string {
	return fmt.Sprintf("Enumerated{%d}", n)
}
