/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

// Package test provides common definitions and function for eap related tests
package mock_sessiond

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"

	"magma/lte/cloud/go/protos"
	orc_protos "magma/orc8r/lib/go/protos"

	"magma/feg/gateway/registry"
	"magma/orc8r/cloud/go/test_utils"
)

// MockSessionManager test sessiond  implementation
type MockSessionManager struct {
	returnErrors bool
}

func NewRunningSessionManager(t *testing.T) *MockSessionManager {
	srv, lis := test_utils.NewTestService(t, registry.ModuleName, registry.SESSION_MANAGER)
	service := &MockSessionManager{
		returnErrors: false,
	}
	protos.RegisterLocalSessionManagerServer(srv.GrpcServer, service)
	go srv.RunTest(lis)
	return service
}

func (c *MockSessionManager) ReportRuleStats(ctx context.Context, in *protos.RuleRecordTable) (*orc_protos.Void, error) {
	out := new(orc_protos.Void)
	if c.returnErrors {
		return out, fmt.Errorf("CreateSession returnErrors enabled")
	}
	err := fmt.Errorf("ReportRuleStats not implemented on test sessionManager")
	return out, err
}

func (c *MockSessionManager) CreateSession(ctx context.Context, in *protos.LocalCreateSessionRequest) (*protos.LocalCreateSessionResponse, error) {
	if c.returnErrors {
		return nil, fmt.Errorf("CreateSession returnErrors enabled")
	}

	out := &protos.LocalCreateSessionResponse{
		SessionId: fmt.Sprintf("%s-12345678", in.Sid.Id),
	}
	return out, nil
}

func (c *MockSessionManager) EndSession(ctx context.Context, in *protos.LocalEndSessionRequest) (*protos.LocalEndSessionResponse, error) {
	if c.returnErrors {
		return nil, fmt.Errorf("CreateSession returnErrors enabled")
	}
	return &protos.LocalEndSessionResponse{}, nil
}

func (c *MockSessionManager) ReturnErrors(enable bool) {
	c.returnErrors = enable
}
