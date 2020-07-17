/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"os"
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/cloud/go/test_utils"
	"orc8r/fbinternal/cloud/go/services/fbinternal"
	"orc8r/fbinternal/cloud/go/services/fbinternal/servicers"
)

func StartTestService(t *testing.T) {
	exporterServicer := servicers.NewExporterServicer(
		os.Getenv("METRIC_EXPORT_URL"),
		os.Getenv("FACEBOOK_APP_ID"),
		os.Getenv("FACEBOOK_APP_SECRET"),
		os.Getenv("METRICS_PREFIX"),
		servicers.ODSMetricsQueueLength,
		servicers.ODSMetricsExportInterval,
	)
	StartTestServiceInternal(t, exporterServicer)
}

func StartTestServiceInternal(t *testing.T, exporter protos.MetricsExporterServer) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, fbinternal.ServiceName)
	protos.RegisterMetricsExporterServer(srv.GrpcServer, exporter)
	go srv.RunTest(lis)
}
