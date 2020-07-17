/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package ods

// Datapoint is used to Marshal JSON encoding for ODS data submission
type Datapoint struct {
	Entity string   `json:"entity"`
	Key    string   `json:"key"`
	Value  string   `json:"value"`
	Time   int      `json:"time"`
	Tags   []string `json:"tags"`
}
