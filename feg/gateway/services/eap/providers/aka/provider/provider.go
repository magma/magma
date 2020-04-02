/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package aka implements EAP-AKA provider
package provider

import (
	"sync"

	"magma/feg/gateway/services/eap/providers"
	"magma/feg/gateway/services/eap/providers/aka"
	"magma/feg/gateway/services/eap/providers/aka/servicers"
)

// AKA Provider Implementation
type providerImpl struct {
	sync.RWMutex
	*servicers.EapAkaSrv
}

func New() providers.Method {
	return &providerImpl{}
}

// String returns EAP AKA Provider name/info
func (*providerImpl) String() string {
	return "<Magma EAP-AKA Method Provider>"
}

// EAPType returns EAP AKA Type - 23
func (*providerImpl) EAPType() uint8 {
	return aka.TYPE
}
