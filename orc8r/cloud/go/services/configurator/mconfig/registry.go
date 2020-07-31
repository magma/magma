/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package mconfig

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/lib/go/registry"
)

// File registry.go provides an mconfig builder registry by forwarding calls to
// the service registry.

// GetBuilders returns all registered mconfig builders.
func GetBuilders() ([]Builder, error) {
	services := registry.FindServices(orc8r.MconfigBuilderLabel)

	var builders []Builder
	for _, s := range services {
		builders = append(builders, NewRemoteBuilder(s))
	}

	return builders, nil
}

// DeregisterAllForTest deregisters all previously-registered mconfig builders.
// This should only be called by test code.
func DeregisterAllForTest(t *testing.T) {
	if t == nil {
		panic("for tests only")
	}
	registry.RemoveServicesWithLabel(orc8r.MconfigBuilderLabel)
}
