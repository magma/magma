/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

// BoolToInt converts true to a 1 and false to a 0.
func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
