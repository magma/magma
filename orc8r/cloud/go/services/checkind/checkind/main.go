/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

/*
Checkind is a dedicated Magma Cloud service which maintains Gateways' runtime
state. The service is intended to be independent, lightweight
and easily scalable. It'll rely on it's own storage which is not required to be
long term persistent or 100% reliable
*/

package main

import (
	"log"
	"time"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/checkind"
	"magma/orc8r/cloud/go/services/checkind/metrics"
	"magma/orc8r/cloud/go/services/checkind/servicers"
	"magma/orc8r/cloud/go/services/checkind/store"
)

const (
	// how often to report checkin status
	GATEWAY_CHECKIN_STATUS_REPORT_INTERVAL = time.Second * 60
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, checkind.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Init the Datastore
	ds, err :=
		datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
	if err != nil {
		log.Fatalf("Failed to initialize datastore: %s", err)
	}

	checkinStore, err := store.NewCheckinStore(ds)
	if err != nil {
		log.Fatalf("Failed to initialize checkin store: %s", err)
	}

	// Add servicers to the service
	checkindServer, err := servicers.NewCheckindServer(checkinStore)
	if err != nil {
		log.Fatalf("Checkin Servicer Initialization Error: %s", err)
	}
	protos.RegisterCheckindServer(srv.GrpcServer, checkindServer)
	srv.GrpcServer.RegisterService(protos.GetLegacyCheckinDesc(), checkindServer)

	// create a gatewayStatusReporter to monitor and periodically log metrics
	// on if the gateway checked in recently on all gateways across all networks
	gwStatusReporter, err := metrics.NewGatewayStatusReporter(checkinStore)
	if err != nil {
		log.Fatalf("GatewayStatusReporter Initialization Error: %s\n", err)
	}
	go gwStatusReporter.ReportCheckinStatus(GATEWAY_CHECKIN_STATUS_REPORT_INTERVAL)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
