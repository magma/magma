/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// This starts the mconfig stream processor and a wrapper grpc service to
// report metrics
package main

import (
	"context"
	"log"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/config/streaming"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	"magma/orc8r/cloud/go/services/mconfig_streamer"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, mconfig_streamer.ServiceName)
	if err != nil {
		log.Fatalf("Error creating mconfig streamer service: %s", err)
	}
	db, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err)
	}

	// Start the service (just a service303 servicer for metrics)
	serviceError := make(chan error)
	go func() {
		err := srv.Run()
		serviceError <- err
	}()

	// start the stream processor
	processorError := make(chan error)
	store := storage.NewDatastoreMconfigStorage(db)
	streamProcessor := streaming.NewStreamProcessor(store, streaming.NewDecoder(), streaming.ConsumerFactoryImpl)
	go func() {
		err := streamProcessor.Run()
		processorError <- err
	}()

	for {
		select {
		case err = <-serviceError:
			streamProcessor.Stop()
			log.Fatalf("Error while running mconfig streamer service: %s", err)
		case err = <-processorError:
			srv.StopService(context.Background(), &protos.Void{})
			log.Fatalf("Error while running mconfig stream processor: %s", err)
		}
	}
}
