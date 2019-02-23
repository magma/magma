// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// UTF8String data type.
type UTF8String OctetString

// DecodeUTF8String decodes an UTF8String data type from byte array.
func DecodeUTF8String(b []byte) (Type, error) {
	return UTF8String(OctetString(b)), nil
}

// Serialize implements the Type interface.
func (s UTF8String) Serialize() []byte {
	return OctetString(s).Serialize()
}

// Len implements the Type interface.
func (s UTF8String) Len() int {
	return len(s)
}

// Padding implements the Type interface.
func (s UTF8String) Padding() int {
	l := len(s)
	return pad4(l) - l
}

// Type implements the Type interface.
func (s UTF8String) Type() TypeID {
	return UTF8StringType
}

// String implements the Type interface.
func (s UTF8String) String() string {
	return fmt.Sprintf("UTF8String{%s},Padding:%d", string(s), s.Padding())
}
