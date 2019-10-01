// Copyright 2013-2019 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"fmt"
	"net"
)

// IPv6 data type for Framed-IP-Address AVP.
type IPv6 net.IP

// DecodeIPv6 decodes an IPv4 data type from byte array.
func DecodeIPv6(b []byte) (Type, error) {
	if len(b) != net.IPv6len {
		return IPv6(make(net.IP, net.IPv6len)), nil
	}
	return IPv6(b), nil
}

// Serialize implements the Type interface.
func (ip IPv6) Serialize() []byte {
	if ip6 := net.IP(ip).To16(); ip6 != nil {
		return ip6
	}
	return ip
}

// Len implements the Type interface.
func (ip IPv6) Len() int {
	return net.IPv6len
}

// Padding implements the Type interface.
func (ip IPv6) Padding() int {
	return 0
}

// Type implements the Type interface.
func (ip IPv6) Type() TypeID {
	return IPv6Type
}

// String implements the Type interface.
func (ip IPv6) String() string {
	return fmt.Sprintf("IPv6{%s}", net.IP(ip))
}
