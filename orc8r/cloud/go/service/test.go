/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package service

import (
	"testing"

	"google.golang.org/grpc"

	"magma/orc8r/cloud/go/service/middleware/unary"
	platform_srvc "magma/orc8r/lib/go/service"
)

// NewTestOrchestratorService returns a new GRPC orchestrator service without
// loading Orchestrator plugins from disk. This should only be used in test
// contexts, where plugins are registered manually.
func NewTestOrchestratorService(t *testing.T, moduleName string, serviceType string) (*platform_srvc.Service, error) {
	if t == nil {
		panic("Nice try, but *testing.T must be non-nil. NewTestOrchestratorService can only be used in a test context.")
	}
	return platform_srvc.NewServiceWithOptions(moduleName, serviceType, grpc.UnaryInterceptor(unary.MiddlewareHandler))
}
