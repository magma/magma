// Copyright 2013-2015 go-diameter authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package datatype

// pad4 returns n padded to 4 bytes.
func pad4(n int) int {
	return n + ((4 - n) & 3)
}
