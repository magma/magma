/*
 Copyright (c) Facebook, Inc. and its affiliates.
 All rights reserved.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.
*/

package main

import (
	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/policydb"
	"magma/lte/cloud/go/services/policydb/servicers"
	"magma/orc8r/cloud/go/service"

	"github.com/golang/glog"
)

func main() {
	// Create the service
	srv, err := service.NewOrchestratorService(lte.ModuleName, policydb.ServiceName)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}

	assignmentServicer := servicers.NewPolicyAssignmentServer()
	protos.RegisterPolicyAssignmentControllerServer(srv.GrpcServer, assignmentServicer)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
