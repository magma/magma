/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"flag"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/logger"
	"magma/orc8r/cloud/go/services/logger/exporters"
	"magma/orc8r/cloud/go/services/logger/nghttpxlogger"
	"magma/orc8r/cloud/go/services/logger/servicers"
	"magma/orc8r/lib/go/protos"

	"github.com/golang/glog"
)

const (
	SCRIBE_EXPORTER_EXPORT_INTERVAL = time.Second * 60
	SCRIBE_EXPORTER_QUEUE_LENGTH    = 100000
	NGHTTPX_LOG_FILE_PATH           = "/var/log/nghttpx.log"
)

var (
	tailNghttpx = flag.Bool("tailNghttpx", false, "Tail Nghttpx Logs and export")
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, logger.ServiceName)
	if err != nil || srv.Config == nil {
		glog.Fatalf("Error creating service: %s", err)
	}

	if *tailNghttpx {
		// run nghttpxlogger on its own goroutine
		parser, err := nghttpxlogger.NewNghttpParser()
		if err != nil {
			glog.Fatalf("Error creating Nghttp Parser: %v\n", err)
		}
		nghttpxLogger, err := nghttpxlogger.NewNghttpLogger(time.Minute, parser)
		if err != nil {
			glog.Fatalf("Error creating Nghttp Logger: %v\n", err)
		}
		glog.V(2).Infof("Running nghttpxlogger...\n")
		nghttpxLogger.Run(NGHTTPX_LOG_FILE_PATH)
	}

	scribeExportURL := srv.Config.GetRequiredStringParam("scribe_export_url")
	scribeAppID := srv.Config.GetRequiredStringParam("scribe_app_id")
	scribeAppSecret := srv.Config.GetRequiredStringParam("scribe_app_secret")

	// Initialize exporters
	scribeExporter := exporters.NewScribeExporter(
		scribeExportURL,
		scribeAppID,
		scribeAppSecret,
		SCRIBE_EXPORTER_QUEUE_LENGTH,
		SCRIBE_EXPORTER_EXPORT_INTERVAL,
	)
	logExporters := make(map[protos.LoggerDestination]exporters.Exporter)
	logExporters[protos.LoggerDestination_SCRIBE] = scribeExporter

	// Add servicers to the service
	loggingServ, err := servicers.NewLoggingService(logExporters)
	if err != nil {
		glog.Fatalf("LoggingService Initialization Error: %s", err)
	}
	// start exporting asynchronously
	scribeExporter.Start()

	protos.RegisterLoggingServiceServer(srv.GrpcServer, loggingServ)
	srv.GrpcServer.RegisterService(protos.GetLegacyLoggerDesc(), loggingServ)
	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
