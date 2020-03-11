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
	"magma/orc8r/cloud/go/services/logger"
	"magma/orc8r/cloud/go/services/logger/exporters"
	"magma/orc8r/cloud/go/services/logger/exporters/mocks"
	"magma/orc8r/cloud/go/services/logger/servicers"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

// this test init does not expose mockExporter, and caller does not do any handling
func StartTestService(t *testing.T) {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, logger.ServiceName)
	logExporters := make(map[protos.LoggerDestination]exporters.Exporter)
	logExporters[protos.LoggerDestination_SCRIBE] = mocks.NewMockExporter()
	loggingSrv, err := servicers.NewLoggingService(logExporters)
	if err != nil {
		t.Fatalf("Failed to create LoggingService")
	}
	protos.RegisterLoggingServiceServer(
		srv.GrpcServer,
		loggingSrv)
	go srv.RunTest(lis)
}

// this test init exposes mockExporter, but caller needs to define .on("Submit", <[]*LogEntry>).Return(<error>)
func StartTestServiceWithMockExporterExposed(t *testing.T) *mocks.ExposedMockExporter {
	srv, lis := test_utils.NewTestService(t, orc8r.ModuleName, logger.ServiceName)
	logExporters := make(map[protos.LoggerDestination]exporters.Exporter)
	exporter := mocks.NewExposedMockExporter()
	logExporters[protos.LoggerDestination_SCRIBE] = exporter
	loggingSrv, err := servicers.NewLoggingService(logExporters)
	if err != nil {
		t.Fatalf("Failed to create LoggingService")
	}
	protos.RegisterLoggingServiceServer(
		srv.GrpcServer,
		loggingSrv)
	go srv.RunTest(lis)
	return exporter
}
