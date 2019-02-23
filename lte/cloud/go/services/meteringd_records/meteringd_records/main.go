/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"log"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/meteringd_records"
	"magma/lte/cloud/go/services/meteringd_records/servicers"
	"magma/lte/cloud/go/services/meteringd_records/storage/dynamo"
	dynamo_common "magma/orc8r/cloud/go/dynamo"
	"magma/orc8r/cloud/go/service"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(lte.ModuleName, meteringd_records.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Init the Datastore
	sess, err := dynamo_common.GetAWSSession()
	if err != nil {
		log.Fatalf("Error creating AWS session: %s", err)
	}
	store := dynamo.NewDynamoDBMeteringRecordsStorage(dynamodb.New(sess), dynamo.NewEncoder(&dynamo.DefaultTimeProvider{}), dynamo.NewDecoder())
	// Should only be true for dev VM
	if dynamo_common.ShouldInitTables() {
		err = store.InitTables()
		if err != nil {
			log.Fatalf("Error initializing dynamoDB tables: %s", err)
		}
	}

	// Add servicers to the service
	servicer := servicers.NewMeteringdRecordsServer(store)
	protos.RegisterMeteringdRecordsControllerServer(srv.GrpcServer, servicer)
	srv.GrpcServer.RegisterService(protos.GetLegacyMeteringDesc(), servicer)

	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
