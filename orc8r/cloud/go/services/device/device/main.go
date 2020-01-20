/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/device"
	"magma/orc8r/cloud/go/services/device/protos"
	"magma/orc8r/cloud/go/services/device/servicers"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, device.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating device service %s", err)
	}
	db, err := sqorc.Open(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE)
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
	}
	store := blobstore.NewEntStorage(device.DBTableName, db, sqorc.GetSqlBuilder())
	err = store.InitializeFactory()
	if err != nil {
		glog.Fatalf("Failed to initialize device database: %s", err)
	}
	// Add servicers to the service
	deviceServicer, err := servicers.NewDeviceServicer(store)
	if err != nil {
		glog.Fatalf("Failed to instantiate the device servicer: %v", deviceServicer)
	}
	protos.RegisterDeviceServer(srv.GrpcServer, deviceServicer)

	err = srv.Run()
	if err != nil {
		glog.Fatalf("Failed to start device service: %v", err)
	}
}
