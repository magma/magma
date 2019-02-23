// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// IPFilterRule data type.
type IPFilterRule OctetString

// DecodeIPFilterRule decodes an IPFilterRule data type from byte array.
func DecodeIPFilterRule(b []byte) (Type, error) {
	return IPFilterRule(OctetString(b)), nil
}

// Serialize implements the Type interface.
func (s IPFilterRule) Serialize() []byte {
	return OctetString(s).Serialize()
}

// Len implements the Type interface.
func (s IPFilterRule) Len() int {
	return len(s)
}

// Padding implements the Type interface.
func (s IPFilterRule) Padding() int {
	l := len(s)
	return pad4(l) - l
}

// Type implements the Type interface.
func (s IPFilterRule) Type() TypeID {
	return IPFilterRuleType
}

// String implements the Type interface.
func (s IPFilterRule) String() string {
	return fmt.Sprintf("IPFilterRule{%s},Padding:%d", string(s), s.Padding())
}
