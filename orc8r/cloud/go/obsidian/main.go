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

	"magma/orc8r/cloud/go/datastore"
	"magma/orc8r/cloud/go/obsidian/config"
	"magma/orc8r/cloud/go/obsidian/server"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/service"
)

func main() {
	flag.IntVar(&config.Port, "port", -1, "HTTP (REST) Server Port")
	flag.IntVar(&config.Port, "p", -1, "HTTP (REST) Server Port (shorthand)")
	flag.StringVar(&config.MagmadDBDriver, "md_db_driver", datastore.SQL_DRIVER, "Magmad DB Driver")
	flag.StringVar(&config.MagmadDBSource, "md_db_source", datastore.DATABASE_SOURCE, "Magmad DB Source")

	// HTTPS settings
	flag.BoolVar(&config.TLS, "tls", false, "HTTPS only access")
	flag.StringVar(
		&config.ServerCertPemPath, "cert",
		datastore.GetEnvWithDefault("REST_CERT", config.DefaultServerCert),
		"Server's certificate PEM file",
	)
	flag.StringVar(
		&config.ServerKeyPemPath, "cert_key",
		datastore.GetEnvWithDefault("REST_CERT_KEY", config.DefaultServerCertKey),
		"Server's certificate private key PEM file",
	)
	flag.StringVar(
		&config.ClientCAPoolPath, "client_ca",
		datastore.GetEnvWithDefault("REST_CLIENT_CERT", config.DefaultClientCAs),
		"Client certificate CA pool PEM file",
	)
	flag.BoolVar(
		&config.AllowAnyClientCert, "client_cert_any", false,
		"Accept Any Client Certificate (Do not verify with given client CAs)",
	)
	flag.StringVar(
		&config.StaticFolder, "static_folder", config.DefaultStaticFolder,
		"Folder containing the static files served",
	)

	srv, err := service.NewOrchestratorService(orc8r.ModuleName, config.ServiceName)
	if err != nil {
		log.Fatalf("Error creating service: %s", err)
	}

	if config.Port == -1 {
		config.Port = config.DefaultPort
		if config.TLS {
			config.Port = config.DefaultHttpsPort
		}
	}

	go srv.Run()
	server.Start()
}
