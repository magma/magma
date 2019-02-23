/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"context"
	"log"
	"time"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/materializer"
	materializer_registry "magma/orc8r/cloud/go/services/materializer/registry"

	"github.com/golang/glog"
)

func main() {
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, materializer.ServiceName)
	if err != nil || srv.Config == nil {
		log.Fatalf("Error creating materializer service: %s", err)
	}

	errors := make(chan error)
	timeout := 10 * time.Second

	go func() {
		err := srv.Run()
		if err != nil {
			select {
			case errors <- err:
			case <-time.After(timeout):
				glog.Warningf("Service write to error channel timed out: %s", err)
			}
		}
	}()

	go func() {
		err := materializer_registry.RunAll()
		if err != nil {
			select {
			case errors <- err:
			case <-time.After(timeout):
				glog.Warningf("Materializer registry write to error channel timed out: %s", err)
			}
		}
	}()

	err = <-errors
	srv.StopService(context.Background(), &protos.Void{})
	materializer_registry.StopAll()
	log.Fatalf("Error in materializer: %s", err)
}
