/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package dynamo

const configPrefix = "cfg_"

type gatewayState struct {
	NetworkID string
	GatewayID string
	Status    []byte `dynamodbav:",omitempty"`
	Record    []byte `dynamodbav:",omitempty"`
	Offset    int64
}
