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
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/testcore/ocs/mock_ocs"
	"magma/orc8r/lib/go/service"

	"github.com/golang/glog"
)

const (
	MaxUsageBytes = 2048
	MaxUsageTime  = 1000 // in second
	ValidityTime  = 60   // in second
)

func init() {
	flag.Parse()
}

func main() {
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.MOCK_OCS)
	if err != nil {
		log.Fatalf("Error creating mock OCS service: %s", err)
	}

	// TODO: support multiple connections
	gyCliConf := gy.GetGyClientConfiguration()[0]
	gyServConf := gy.GetOCSConfiguration()[0]

	diamServer := mock_ocs.NewOCSDiamServer(
		gyCliConf,
		&mock_ocs.OCSConfig{
			ServerConfig:   gyServConf,
			MaxUsageOctets: &protos.Octets{TotalOctets: MaxUsageBytes},
			MaxUsageTime:   MaxUsageTime,
			ValidityTime:   ValidityTime,
			GyInitMethod:   gy.PerSessionInit,
		},
	)

	lis, err := diamServer.StartListener()
	if err != nil {
		log.Fatalf("Unable to start listener for mock OCS: %s", err)
	}

	protos.RegisterMockOCSServer(srv.GrpcServer, diamServer)

	go func() {
		glog.V(2).Infof("Starting mock OCS server at %s", lis.Addr().String())
		glog.Errorf(diamServer.Start(lis).Error()) // blocks
	}()

	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running mock OCS service: %s", err)
	}
}
