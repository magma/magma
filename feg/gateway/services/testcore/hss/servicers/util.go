/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

// AllZero returns true if and only if the slice contains only zero bytes.
func AllZero(bytes []byte) bool {
	for _, b := range bytes {
		if b != 0 {
			return false
		}
	}
	return true
}

// BoolToInt converts true to a 1 and false to a 0.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
