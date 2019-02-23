// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"fmt"
	"net"
)

// IPv4 data type for Framed-IP-Address AVP.
type IPv4 net.IP

// DecodeIPv4 decodes an IPv4 data type from byte array.
func DecodeIPv4(b []byte) (Type, error) {
	if len(b) != 4 {
		return IPv4{0, 0, 0, 0}, nil
	}
	return IPv4(b), nil
}

// Serialize implements the Type interface.
func (ip IPv4) Serialize() []byte {
	if ip4 := net.IP(ip).To4(); ip4 != nil {
		return ip4
	}
	return ip
}

// Len implements the Type interface.
func (ip IPv4) Len() int {
	return 4
}

// Padding implements the Type interface.
func (ip IPv4) Padding() int {
	return 0
}

// Type implements the Type interface.
func (ip IPv4) Type() TypeID {
	return IPv4Type
}

// String implements the Type interface.
func (ip IPv4) String() string {
	return fmt.Sprintf("IPv4{%s}", net.IP(ip))
}
