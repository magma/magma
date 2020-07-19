/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package main

import (
	"os"

	"magma/fbinternal/cloud/go/fbinternal"
	fbinternal_service "magma/fbinternal/cloud/go/services/fbinternal"
	"magma/fbinternal/cloud/go/services/fbinternal/servicers"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/metricsd/protos"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(fbinternal.ModuleName, fbinternal_service.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating orc8r service for fbinternal: %s", err)
	}

	exporterServicer := servicers.NewExporterServicer(
		os.Getenv("METRIC_EXPORT_URL"),
		os.Getenv("FACEBOOK_APP_ID"),
		os.Getenv("FACEBOOK_APP_SECRET"),
		os.Getenv("METRICS_PREFIX"),
		servicers.ODSMetricsQueueLength,
		servicers.ODSMetricsExportInterval,
	)
	protos.RegisterMetricsExporterServer(srv.GrpcServer, exporterServicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
