// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import "fmt"

// QoSFilterRule data type.
type QoSFilterRule OctetString

// DecodeQoSFilterRule decodes an QoSFilterRule data type from byte array.
func DecodeQoSFilterRule(b []byte) (Type, error) {
	return QoSFilterRule(OctetString(b)), nil
}

// Serialize implements the Type interface.
func (s QoSFilterRule) Serialize() []byte {
	return OctetString(s).Serialize()
}

// Len implements the Type interface.
func (s QoSFilterRule) Len() int {
	return len(s)
}

// Padding implements the Type interface.
func (s QoSFilterRule) Padding() int {
	l := len(s)
	return pad4(l) - l
}

// Type implements the Type interface.
func (s QoSFilterRule) Type() TypeID {
	return QoSFilterRuleType
}

// String implements the Type interface.
func (s QoSFilterRule) String() string {
	return fmt.Sprintf("QoSFilterRule{%s},Padding:%d", string(s), s.Padding())
}
