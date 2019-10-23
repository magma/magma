/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/loader"
	"fbc/cwf/radius/monitoring"
	"fbc/cwf/radius/server"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sort"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// Get a simple stdout logger
	logger, err := zap.NewProduction()

	// Get configuration
	var configFilename string
	flag.StringVar(&configFilename, "config", "radius.config.json", "The configuration filename")
	flag.Parse()
	config, err := config.Read(configFilename)
	if err != nil {
		logger.Error("Failed to read configuration", zap.Error(err))
		return
	}

	// Initialize pprof debug interface
	if config.Debug != nil {
		if config.Debug.Enabled {
			logger.Info("Enabling Server Debugging", zap.Int("port", config.Debug.Port))
			go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Debug.Port), nil))
		} else {
			logger.Info("Server Debugging interface is disabled")
		}
	}

	// Initialize monitoring
	logger, err = monitoring.Init(config.Monitoring, logger)
	if err != nil {
		fmt.Println("Failed initializing monitoring", zap.Error(err))
		return
	}

	logger = logger.With(zap.String("host", getHostIdentifier()))

	loader := loader.NewStaticLoader(logger)

	// Create server
	radiusServer, err := server.New(config.Server, logger, loader)
	if err != nil {
		logger.Error("Failed creating server", zap.Error(err))
		return
	}

	// Capture CTRL+C
	sigtermChannel := make(chan os.Signal, 1)
	signal.Notify(sigtermChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigtermChannel
		logger.Info("Received SIGTERM, existing")
		radiusServer.Stop()
		logger.Sync()
	}()

	// Start the server
	radiusServer.Start()
}

func getHostIdentifier() string {
	hostname, err := os.Hostname()
	if err == nil {
		return hostname
	}

	// Get the MAC address with the lowest lexicographical index
	// This is some sort of stable host identifier...
	interfaces, err := net.Interfaces()
	if err == nil && len(interfaces) > 0 {
		var macs []string
		for _, ifa := range interfaces {
			mac := ifa.HardwareAddr.String()
			if mac != "" {
				macs = append(macs, mac)
			}
		}
		sort.Strings(macs)
		return macs[0]
	}

	// Just a random, unstable identifer
	return fmt.Sprintf("random:%d", rand.Intn(9999999))
}
