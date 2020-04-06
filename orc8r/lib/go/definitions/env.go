/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package definitions

import "os"

func GetEnvWithDefault(variable string, defaultValue string) string {
	value := os.Getenv(variable)
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}
