/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"context"
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	"magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	"magma/orc8r/cloud/go/test_utils"
)

type builderServicer struct {
	builder mconfig.Builder
}

// StartNewTestBuilder starts a new mconfig builder service which forwards
// calls to the passed builder.
func StartNewTestBuilder(t *testing.T, builder mconfig.Builder) {
	labels := map[string]string{
		orc8r.MconfigBuilderLabel: "true",
	}
	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "test_mconfig_builder_service", labels, nil)
	servicer := &builderServicer{builder: builder}
	protos.RegisterBuilderServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
}

func (b builderServicer) Build(ctx context.Context, request *protos.BuildRequest) (*protos.BuildResponse, error) {
	configs, err := b.builder.Build(request.Network, request.Graph, request.GatewayId)
	if err != nil {
		return nil, err
	}
	res := &protos.BuildResponse{ConfigsByKey: configs}
	return res, nil
}
