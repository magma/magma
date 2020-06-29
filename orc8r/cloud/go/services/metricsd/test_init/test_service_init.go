/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package test_init

import (
	"testing"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/metricsd"
	"magma/orc8r/cloud/go/services/metricsd/protos"
	"magma/orc8r/cloud/go/services/metricsd/servicers"
	"magma/orc8r/cloud/go/test_utils"
)

func StartTestService(t *testing.T) {
	exporterServicer := servicers.NewPushExporterServicer([]string{})
	StartTestServiceInternal(t, exporterServicer)
}

func StartTestServiceInternal(t *testing.T, exporter protos.MetricsExporterServer) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, metricsd.ServiceName)
	protos.RegisterMetricsExporterServer(srv.GrpcServer, exporter)
	go srv.RunTest(lis)
}
