/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package gateways

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/service/config"
	"magma/orc8r/cloud/go/services/materializer"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage/dynamo"
	storage_sql "magma/orc8r/cloud/go/services/materializer/gateways/storage/sql"
	"magma/orc8r/cloud/go/services/materializer/gateways/streaming"
)

// GetApplication returns the materializer Application for gateways
func GetApplication() materializer.Application {
	aggregator := streaming.NewStreamAggregator(
		streaming.NewStreamAggregatorConsumer,
		streaming.NewDecoder(),
		streaming.NewStreamAggregatorProducer,
	)

	store, err := getStorage()
	if err != nil {
		log.Fatalf("Error getting gateways materializer storage: %s", err)
	}
	recorder := streaming.NewStateRecorder(
		streaming.NewStateRecorderConsumer,
		streaming.NewDecoder(),
		store,
	)

	return materializer.Application{
		Name:       "gateways",
		Processors: []materializer.StreamProcessor{aggregator, recorder},
	}
}

// getStorage returns the appropriate storage impl for the gateways
// materializer application (the obsidian handler can use a different storage -
// see the exported GetStorage under the handlers package).
// This depends on the materializer service yml's field
// streamer_write_storage
func getStorage() (storage.GatewayViewStorage, error) {
	// Loading the config manually because the service is initialized separately
	// Hardcoding in orc8r to avoid cyclical import
	configMap, err := config.GetServiceConfig("orc8r", materializer.ServiceName)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving config map from materializer %v", err)
	}
	v := configMap.GetRequiredStringParam("obsidian_read_storage")
	if strings.ToLower(v) == "sql" {
		db, err := sql.Open(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
		if err != nil {
			return nil, fmt.Errorf("Could not initialize SQL connection: %s", err)
		}
		return storage_sql.NewSqlGatewayViewStorage(db), nil
	}

	// Default to dynamo for back-compat
	return dynamo.GetInitializedDynamoStorage()
}
