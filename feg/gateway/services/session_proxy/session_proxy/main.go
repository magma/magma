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
	"fmt"
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
	"magma/orc8r/lib/go/service"
	"magma/orc8r/lib/go/util"

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

	// Create configs for each server and start diam connections
	controllerParms, policyDBClient, err := generateClientsConfsAndDiameterConnection()
	if err != nil {
		glog.Fatal(err)
		return
	}

	// Add servicers to the service
	sessionManagerAndHealthServer, err := servicers.
		NewCentralSessionControllerDefaultMultiplexWithHealth(controllerParms, policyDBClient)
	lteprotos.RegisterCentralSessionControllerServer(srv.GrpcServer, sessionManagerAndHealthServer)
	protos.RegisterServiceHealthServer(srv.GrpcServer, sessionManagerAndHealthServer)

	// Run the service
	err = srv.Run()
	if err != nil {
		glog.Fatalf("Error running service: %s", err)
	}
}

// TODO: move this to servicers and add testing
// generateClientsConfsAndDiameterConnection reads configurations for all GXs and GYs connections configured
// at gateway.mconfig and creates a slice containing all the requiered parameters to start CentralSessionControllers
func generateClientsConfsAndDiameterConnection() (
	[]*servicers.ControllerParam, *policydb.RedisPolicyDBClient, error) {
	cloudReg := registry.Get()
	policyDBClient, err := policydb.NewRedisPolicyDBClient(cloudReg)
	if err != nil {
		return nil, nil, fmt.Errorf("Error connecting to redis store: %s", err)
	}

	// ---- Read configus from gateway.mconfig  ----
	glog.Info("------ Reading Gx and Gy configuration ------")
	// Global config, init Method and policyDb (static routes) are shared by all the controllers
	gyGlobalConf := gy.GetGyGlobalConfig()
	gxGlobalConf := gx.GetGxGlobalConfig()
	initMethod := gy.GetInitMethod()

	// Each controller will take one entry of PCRF, OCS, and gx/gy clients confs
	gxCliConfs := gx.GetGxClientConfiguration()
	gyCLiConfs := gy.GetGyClientConfiguration()
	OCSConfs := gy.GetOCSConfiguration()
	PCRFConfs := gx.GetPCRFConfiguration()

	// this is a new copy needed to fill in the controllerParms
	OCSConfsCopy := gy.GetOCSConfiguration()
	PCRFConfsCopy := gx.GetPCRFConfiguration()

	// Exit if the number of GX and GY configurations are different
	if len(OCSConfs) != len(PCRFConfs) {
		return nil, nil, fmt.Errorf(
			"Number of Gx and Gy servers configured must be equal Gx:%d Gx:%d",
			len(OCSConfs), len(PCRFConfs))
	}
	glog.Info("------ Done reading configuration ------")

	// ---- Create diammeter connections and build parameters for CentralSessionControllersn ----
	glog.Info("------ Create diameter connexions ------")
	totalLen := len(OCSConfs)
	controllerParms := make([]*servicers.ControllerParam, 0, totalLen)
	for i := 0; i < totalLen; i++ {
		controlParam := &servicers.ControllerParam{}
		// Fill in general parameters for controler i
		controlParam.Config = &servicers.SessionControllerConfig{
			OCSConfig:        OCSConfs[i],
			PCRFConfig:       PCRFConfs[i],
			RequestTimeout:   3 * time.Second,
			InitMethod:       initMethod,
			UseGyForAuthOnly: util.IsTruthyEnv(gy.UseGyForAuthOnlyEnv),
		}
		// Fill in gx and gy config for controller i
		if OCSConfsCopy[i].DiameterServerConnConfig == PCRFConfsCopy[i].DiameterServerConnConfig &&
			OCSConfsCopy[i] != PCRFConfsCopy[i] {
			glog.Infof("Using single Gy/Gx connection for server: %+v",
				OCSConfsCopy[i].DiameterServerConnConfig)
			var clientCfg = *gxCliConfs[i]
			clientCfg.AuthAppID = gyCLiConfs[i].AppID
			diamClient := diameter.NewClient(&clientCfg)
			diamClient.BeginConnection(OCSConfsCopy[i])
			controlParam.CreditClient = gy.NewConnectedGyClient(
				diamClient,
				OCSConfsCopy[i],
				gy.GetGyReAuthHandler(cloudReg),
				cloudReg,
				gyGlobalConf)
			controlParam.PolicyClient = gx.NewConnectedGxClient(
				diamClient,
				OCSConfsCopy[i],
				gx.GetGxReAuthHandler(cloudReg, policyDBClient),
				cloudReg,
				gxGlobalConf)
		} else {
			glog.Infof("Using distinct Gy: %+v & Gx: %+v connection",
				OCSConfsCopy[i].DiameterServerConnConfig, PCRFConfsCopy[i].DiameterServerConnConfig)
			controlParam.CreditClient = gy.NewGyClient(
				gy.GetGyClientConfiguration()[i],
				OCSConfsCopy[i],
				gy.GetGyReAuthHandler(cloudReg),
				cloudReg,
				gyGlobalConf)
			controlParam.PolicyClient = gx.NewGxClient(
				gx.GetGxClientConfiguration()[i],
				PCRFConfsCopy[i],
				gx.GetGxReAuthHandler(cloudReg, policyDBClient),
				cloudReg,
				gxGlobalConf)
		}
		controllerParms = append(controllerParms, controlParam)
	}
	glog.Infof("------ Done creating %d diameter connexions ------", totalLen)
	return controllerParms, policyDBClient, nil
}
