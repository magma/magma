/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package util

// MinInt returns the minimum of two ints.
func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
