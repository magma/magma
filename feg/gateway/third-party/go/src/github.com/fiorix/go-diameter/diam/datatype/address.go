// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

// Address data type.
type Address []byte

// DecodeAddress decodes an Address data type from byte array.
func DecodeAddress(b []byte) (Type, error) {
	if len(b) < 3 {
		return nil, fmt.Errorf("Not enough data to make an Address from byte[%d] = %+v", len(b), b)
	}
	if binary.BigEndian.Uint16(b[:2]) == 0 || binary.BigEndian.Uint16(b[:2]) == 65535 {
		return nil, errors.New("Invalid address type received")
	}
	switch binary.BigEndian.Uint16(b[:2]) {
	case 0x01:
		if len(b[2:]) != 4 {
			return nil, errors.New("Invalid length for IPv4")
		}
	case 0x02:
		if len(b[2:]) != 16 {
			return nil, errors.New("Invalid length for IPv6")
		}
	default:
		return Address(b), nil
	}
	return Address(b[2:]), nil
}

// Serialize implements the Type interface.
func (addr Address) Serialize() []byte {
	var b []byte
	if ip4 := net.IP(addr).To4(); ip4 != nil {
		b = make([]byte, 6)
		b[1] = 0x01
		copy(b[2:], ip4)
	} else if ip6 := net.IP(addr).To16(); ip6 != nil {
		b = make([]byte, 18)
		b[1] = 0x02
		copy(b[2:], addr)
	} else {
		b = make([]byte, len(addr))
		copy(b, addr)
	}
	return b
}

// Len implements the Type interface.
func (addr Address) Len() int {
	if ip4 := net.IP(addr).To4(); ip4 != nil {
		return len(ip4) + 2 // two bytes from the address family
	} else if ip6 := net.IP(addr).To16(); ip6 != nil {
		return len(ip6) + 2 // two bytes from the address family
	} else {
		return len(addr)
	}
}

// Padding implements the Type interface.
func (addr Address) Padding() int {
	var l int
	if ip4 := net.IP(addr).To4(); ip4 != nil {
		l = len(ip4) + 2 // two bytes from the address family
	} else if ip6 := net.IP(addr).To16(); ip6 != nil {
		l = len(ip6) + 2 // two bytes from the address family
	} else {
		l = len(addr)
	}
	return pad4(l) - l
}

// Type implements the Type interface.
func (addr Address) Type() TypeID {
	return AddressType
}

// String implements the Type interface.
func (addr Address) String() string {
	if ip4 := net.IP(addr).To4(); ip4 != nil {
		return fmt.Sprintf("Address{%s},Padding:%d", net.IP(addr), addr.Padding())
	}
	if ip6 := net.IP(addr).To16(); ip6 != nil {
		return fmt.Sprintf("Address{%s},Padding:%d", net.IP(addr), addr.Padding())
	}
	return fmt.Sprintf("Address{%#v}, Type{%#v} Padding:%d", addr[2:], addr[:2], addr.Padding())
}
