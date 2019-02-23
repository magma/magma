/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package migrations

import "os"

func GetEnvWithDefault(variable string, defaultValue string) string {
	val, set := os.LookupEnv(variable)
	if set {
		return val
	}
	return defaultValue
}
