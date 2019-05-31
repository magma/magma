/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
// Package providers encapsulates supported EAP Authenticator Providers
//
//go:generate protoc -I ../protos -I . --go_out=plugins=grpc,paths=source_relative:. protos/eap_provider.proto
//
package providers

import (
	"fmt"

	"magma/feg/gateway/services/aaa/protos"
)

// Method is the Interface for Eap Provider
type Method interface {
	// Stringer -> String() string with Provider Name/description
	fmt.Stringer
	// EAPType should return a valid EAP Type
	EAPType() uint8
	// Handle - handles EAP Resp message (protos.EapRequest)
	Handle(*protos.Eap) (*protos.Eap, error)
}
