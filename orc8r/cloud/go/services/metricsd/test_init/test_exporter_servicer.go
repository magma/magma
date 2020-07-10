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
	"magma/orc8r/cloud/go/services/metricsd/exporters"
	"magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/cloud/go/test_utils"
)

type exporterServicer struct {
	exporter exporters.Exporter
}

// StartNewTestExporter starts a new metrics exporter service which forwards
// calls to the passed exporter.
func StartNewTestExporter(t *testing.T, exporter exporters.Exporter) {
	labels := map[string]string{
		orc8r.MetricsExporterLabel: "true",
	}
	srv, lis := test_utils.NewTestOrchestratorService(t, orc8r.ModuleName, "MOCK_EXPORTER_SERVICE", labels, nil)
	servicer := &exporterServicer{exporter: exporter}
	protos.RegisterMetricsExporterServer(srv.GrpcServer, servicer)
	go srv.RunTest(lis)
}

func (e *exporterServicer) Submit(ctx context.Context, req *protos.SubmitMetricsRequest) (*protos.SubmitMetricsResponse, error) {
	err := e.exporter.Submit(exporters.MakeNativeMetrics(req.Metrics))
	return &protos.SubmitMetricsResponse{}, err
}
