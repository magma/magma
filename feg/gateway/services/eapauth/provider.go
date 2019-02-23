/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package eapauth (EAP Authenticator) provides interface to supported & registered EAP Authenticator Providers
package eapauth

import (
	"fmt"

	"magma/feg/gateway/services/eapauth/protos"
)

// Provider is the Interface for Eap Provider
type Provider interface {
	// Stringer -> String() string with Provider Name/description
	fmt.Stringer
	// EAPType should return a valid EAP Type
	EAPType() uint8
	// Handle - handles EAP Resp message (protos.EapRequest)
	Handle(*protos.EapRequest) (*protos.EapResult, error)
}
