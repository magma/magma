/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package protos

import "github.com/golang/protobuf/ptypes/wrappers"

func GetStringWrapper(v *string) *wrappers.StringValue {
	if v == nil {
		return nil
	}
	return &wrappers.StringValue{Value: *v}
}
