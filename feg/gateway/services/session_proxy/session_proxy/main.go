/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Central Session Controller is a service which proxies calls to the OCS and
// policydb to retrieve credit and traffic policy information and relay it to
// the gateway.
package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/policydb"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"
	"magma/feg/gateway/services/session_proxy/credit_control/gy"
	"magma/feg/gateway/services/session_proxy/servicers"
	lteprotos "magma/lte/cloud/go/protos"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/util"

	"github.com/golang/glog"
)

const (
	TestOCSIP = ":4444"
)

func init() {
	flag.Parse()
}

func main() {
	serviceBaseName := filepath.Base(os.Args[0])
	serviceBaseName = strings.TrimSuffix(serviceBaseName, filepath.Ext(serviceBaseName))
	if credit_control.SessionProxyServiceName != serviceBaseName {
		glog.Warningf(
			"Session Proxy Base Service name: %s does not match its managed configs key: %s",
			serviceBaseName, credit_control.SessionProxyServiceName)
	}
	// Create the service
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.SESSION_PROXY)
	if err != nil {
		glog.Fatalf("Error creating service: %s", err)
	}

	initMethod := gy.GetInitMethod()

	controllerCfg := &servicers.SessionControllerConfig{
		OCSConfig:        gy.GetOCSConfiguration(),
		PCRFConfig:       gx.GetPCRFConfiguration(),
		RequestTimeout:   3 * time.Second,
		InitMethod:       initMethod,
		UseGyForAuthOnly: util.IsTruthyEnv(gy.UseGyForAuthOnlyEnv),
	}
	cloudReg := registry.NewCloudRegistry()
	policyDBClient, err := policydb.NewRedisPolicyDBClient(cloudReg)
	if err != nil {
		glog.Fatalf("Error connecting to redis store: %s", err)
	}

	ocsDiamCfg := gy.GetOCSConfiguration()
	pcrfDiamCfg := gx.GetPCRFConfiguration()

	gyGlobalConfig := gy.GetGyGlobalConfig()

	var gxClnt *gx.GxClient
	var gyClnt *gy.GyClient

	gxClntCfg := gx.GetGxClientConfiguration()
	gyClntCfg := gy.GetGyClientConfiguration()

	if ocsDiamCfg.DiameterServerConnConfig == pcrfDiamCfg.DiameterServerConnConfig &&
		ocsDiamCfg != pcrfDiamCfg {

		glog.Infof("Using single Gy/Gx connection for server: %+v", ocsDiamCfg.DiameterServerConnConfig)

		var clientCfg = *gxClntCfg
		clientCfg.AuthAppID = gyClntCfg.AppID
		diamClient := diameter.NewClient(&clientCfg)
		diamClient.BeginConnection(ocsDiamCfg)

		gyClnt = gy.NewConnectedGyClient(
			diamClient,
			ocsDiamCfg,
			gy.GetGyReAuthHandler(cloudReg),
			cloudReg,
			gyGlobalConfig)
		gxClnt = gx.NewConnectedGxClient(
			diamClient,
			ocsDiamCfg,
			gx.GetGxReAuthHandler(cloudReg, policyDBClient), cloudReg)
	} else {
		glog.Infof("Using distinct Gy: %+v & Gx: %+v connection",
			ocsDiamCfg.DiameterServerConnConfig, pcrfDiamCfg.DiameterServerConnConfig)

		gyClnt = gy.NewGyClient(
			gy.GetGyClientConfiguration(),
			ocsDiamCfg,
			gy.GetGyReAuthHandler(cloudReg), cloudReg, gyGlobalConfig)
		gxClnt = gx.NewGxClient(
			gx.GetGxClientConfiguration(),
			pcrfDiamCfg,
			gx.GetGxReAuthHandler(cloudReg, policyDBClient), cloudReg)
	}
	// Add servicers to the service
	sessionManager := servicers.NewCentralSessionController(gyClnt, gxClnt, policyDBClient, controllerCfg)
	lteprotos.RegisterCentralSessionControllerServer(srv.GrpcServer, sessionManager)
	protos.RegisterServiceHealthServer(srv.GrpcServer, sessionManager)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}
