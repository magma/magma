/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package main

import (
	"fbc/cwf/radius/config"
	"fbc/cwf/radius/counters"
	"fbc/cwf/radius/loader"
	"fbc/cwf/radius/server"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// First there was a logger...
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}

	// Get configuration
	var configFilename string
	flag.StringVar(&configFilename, "config", "radius.config.json", "The configuration filename")
	flag.Parse()
	config, err := config.Read(configFilename)
	if err != nil {
		logger.Error("Failed to read configuration", zap.Error(err))
		return
	}

	// Initialize counters
	counters.Init(config.Counters, logger)

	// Prepare dependencies
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
	}()

	// Start the server
	radiusServer.Start()
}
