/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// access_helper provides ToString() receiver for AccessControl_Permission mask
package protos

// ACCESS_CONTROL_ALL_PERMISSIONS is a bitmask for all existing permissions
// unfortunately, it cannot be const since it has to be 'built' by package's
// init to simplify future maintenance
var ACCESS_CONTROL_ALL_PERMISSIONS AccessControl_Permission

func init() {
	ACCESS_CONTROL_ALL_PERMISSIONS = AccessControl_NONE
	for _, val := range AccessControl_Permission_value {
		ACCESS_CONTROL_ALL_PERMISSIONS |= AccessControl_Permission(val)
	}
}

// ToString returns a string representation of AccessControl_Permission as a mask
// protoc generated String() receiver treats it as enum and does not represent
// the 'mask' use case
func (p AccessControl_Permission) ToString() string {
	res := ""
	for mask, name := range AccessControl_Permission_name {
		if int32(p)&mask != 0 {
			if len(res) == 0 {
				res = name
			} else {
				res += "|" + name
			}
		}
	}
	if len(res) == 0 {
		res = AccessControl_Permission_name[0]
	}
	return res
}
