// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package diam

// uint24to32 converts b from []byte in network byte order to uint32.
func uint24to32(b []byte) uint32 {
	if len(b) != 3 {
		return 0
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// uint32to24 converts b from uint32 to []byte in network byte order.
func uint32to24(n uint32) []byte {
	return []byte{uint8(n >> 16), uint8(n >> 8), uint8(n)}
}
