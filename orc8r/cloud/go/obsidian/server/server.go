/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Server's main package, run with obsidian -h to see all available options
package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/access"
	"magma/orc8r/cloud/go/obsidian/reverse_proxy"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Start() {
	e := echo.New()

	obsidian.AttachAll(e)
	// metrics middleware is used before all other middlewares
	e.Use(CollectStats)
	e.Use(middleware.Recover())

	// Serve static pages for the API docs
	e.Static(obsidian.StaticURLPrefix, obsidian.StaticFolder+"/apidocs")
	e.Static(obsidian.StaticURLPrefix+"/swagger-ui/dist", obsidian.StaticFolder+"/swagger-ui/dist")

	portStr := fmt.Sprintf(":%d", obsidian.Port)
	log.Printf("Starting %s on %s", obsidian.Product, portStr)

	var err error
	if obsidian.TLS {
		var caCerts []byte
		caCerts, err = ioutil.ReadFile(obsidian.ClientCAPoolPath)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(caCerts)
		if ok {
			log.Printf("Loaded %d Client CA Certificate[s] from '%s'", len(caCertPool.Subjects()), obsidian.ClientCAPoolPath)
		} else {
			log.Printf(
				"ERROR: No Certificates found in '%s'", obsidian.ClientCAPoolPath)
		}
		// Possible clientCertVerification values:
		// 	NoClientCert
		// 	RequestClientCert
		// 	RequireAnyClientCert
		// 	VerifyClientCertIfGiven
		// 	RequireAndVerifyClientCert
		clientCertVerification := tls.RequireAndVerifyClientCert
		if obsidian.AllowAnyClientCert {
			clientCertVerification = tls.RequireAnyClientCert
		}
		s := e.TLSServer
		s.TLSConfig = &tls.Config{
			Certificates: make([]tls.Certificate, 1),
			ClientCAs:    caCertPool,
			ClientAuth:   clientCertVerification,
			// Limit versions & Ciphers to our preferred list
			MinVersion: tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256, // 4 HTTP2 support
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				//tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				//tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		s.TLSConfig.Certificates[0], err = tls.LoadX509KeyPair(obsidian.ServerCertPemPath, obsidian.ServerKeyPemPath)
		if err != nil {
			log.Fatalf(
				"ERROR loading server certificate ('%s') and/or key ('%s'): %s",
				obsidian.ServerCertPemPath, obsidian.ServerKeyPemPath, err,
			)
		}
		s.TLSConfig.BuildNameToCertificate()
		s.Addr = portStr
		if !e.DisableHTTP2 {
			s.TLSConfig.NextProtos = append(s.TLSConfig.NextProtos, "h2")
		}
	} else {
		e.Use(access.Middleware)
	}

	e.Use(reverse_proxy.ReverseProxy)
	if obsidian.TLS {
		err = e.StartServer(e.TLSServer)
	} else {
		err = e.Start(portStr)
	}
	if err != nil {
		log.Println(err)
	}
}
