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
	"time"

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
	"magma/orc8r/cloud/go/services/certifier"
	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/services/certifier/servicers"
	"magma/orc8r/cloud/go/services/certifier/storage"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/security/cert"

	"github.com/golang/glog"
	"golang.org/x/net/context"
)

var (
	bootstrapCACertFile = flag.String("cac", "server_cert.pem", "Signer CA's Certificate file")
	bootstrapCAKeyFile  = flag.String("cak", "server_cert.key.pem", "Signer CA's Private Key file")

	vpnCertFile = flag.String("vpnc", "vpn_ca.crt", "VPN CA's Certificate file")
	vpnKeyFile  = flag.String("vpnk", "vpn_ca.key", "VPN CA's Private Key file")

	gcHours = flag.Int64("gc-hours", 12, "Garbage Collection time interval (in hours)")
)

func main() {
	// Create the service, flag will be parsed inside this function
	srv, err := service.NewOrchestratorService(orc8r.ModuleName, certifier.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	// Init the datastore
	store, err := datastore.NewSqlDb(datastore.SQL_DRIVER, datastore.DATABASE_SOURCE, sqorc.GetSqlBuilder())
	if err != nil {
		log.Fatalf("Failed to initialize datastore: %s", err)
	}
	caMap := map[protos.CertType]*servicers.CAInfo{}

	// Add servicers to the service
	bootstrapCert, bootstrapPrivKey, err := cert.LoadCertAndPrivKey(*bootstrapCACertFile, *bootstrapCAKeyFile)
	if err != nil {
		log.Printf("ERROR: Failed to load bootstrap CA cert and key: %v", err)
	} else {
		caMap[protos.CertType_DEFAULT] = &servicers.CAInfo{Cert: bootstrapCert, PrivKey: bootstrapPrivKey}
	}
	vpnCert, vpnPrivKey, vpnErr := cert.LoadCertAndPrivKey(*vpnCertFile, *vpnKeyFile)
	if vpnErr != nil {
		fmtstr := "ERROR: Failed to load VPN cert and key: %v"
		if err != nil {
			log.Fatalf(fmtstr, vpnErr)
		} else {
			log.Printf(fmtstr, vpnErr)
		}
	} else {
		caMap[protos.CertType_VPN] = &servicers.CAInfo{Cert: vpnCert, PrivKey: vpnPrivKey}
	}
	certStore := storage.NewCertifierDatastore(store)
	servicer, err := servicers.NewCertifierServer(certStore, caMap)
	if err != nil {
		log.Fatalf("Failed to create certifier server: %s", err)
	}
	certprotos.RegisterCertifierServer(srv.GrpcServer, servicer)

	// Start Garbage Collector Ticker
	gc := time.Tick(time.Hour * time.Duration(*gcHours))
	go func() {
		for now := range gc {
			log.Printf("%v - Removing Stale Certificates", now)
			_, err := servicer.CollectGarbage(context.Background(), &protos.Void{})
			if err != nil {
				glog.Errorf("error collecting garbage for certifier: %s", err)
			}
		}
	}()

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running service: %s", err)
	}
}
