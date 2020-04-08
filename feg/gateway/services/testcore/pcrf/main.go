/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"flag"
	"log"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/testcore/pcrf/mock_pcrf"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

func init() {
	flag.Parse()
}

func main() {
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.MOCK_PCRF)
	if err != nil {
		log.Fatalf("Error creating mock PCRF service: %s", err)
	}

	// TODO: support multiple connections
	gxCliConf := gx.GetGxClientConfiguration()[0]
	gxServConf := gx.GetPCRFConfiguration()[0]

	pcrfServer := mock_pcrf.NewPCRFDiamServer(
		gxCliConf,
		&mock_pcrf.PCRFConfig{ServerConfig: gxServConf},
	)

	lis, err := pcrfServer.StartListener()
	if err != nil {
		log.Fatalf("Unable to start listener for mock PCRF: %s", err)
	}

	protos.RegisterMockPCRFServer(srv.GrpcServer, pcrfServer)

	go func() {
		glog.V(2).Infof("Starting mock PCRF server at %s", lis.Addr().String())
		glog.Errorf(pcrfServer.Start(lis).Error()) // blocks
	}()

	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running mock PCRF service: %s", err)
	}
}
