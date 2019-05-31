/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package aaa provides Carrier WiFi related services
//
//go:generate protoc -I protos --go_out=plugins=grpc,paths=source_relative:protos protos/context.proto protos/eap.proto protos/accounting.proto
//
package aaa

import (
	"fmt"
	"math/rand"
	"time"
)

// CreateSessionId creates & returns unique session ID string
func CreateSessionId() string {
	return fmt.Sprintf("%X-%X", time.Now().UnixNano()>>16, rand.Uint32())
}
