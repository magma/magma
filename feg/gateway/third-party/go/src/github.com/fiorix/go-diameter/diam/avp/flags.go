// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package avp

// AVP Flags. See section 4.1 of RFC 6733.
const (
	Pbit = 1 << 5 // The 'P' bit, reserved for future use.
	Mbit = 1 << 6 // The 'M' bit, known as the Mandatory bit.
	Vbit = 1 << 7 // The 'V' bit, known as the Vendor-Specific bit.
)
