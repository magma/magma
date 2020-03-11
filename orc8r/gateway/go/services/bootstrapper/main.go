/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"magma/gateway/config"
	"magma/gateway/services/bootstrapper/gateway_info"
	"magma/gateway/services/bootstrapper/service"
)

const usageExamples string = `
Examples:

  1. Run Bootstrapper as a service:

    $> %s

    The command will run Bootstrapper service which will periodically 
    check the gateway certificats and update them when needed

  2. Show the gateway information needed for the gateway registration and exit:

    $> %s -show
    OR
    $> %s -s

    The command will print the gateway hardware ID and challenge key and exit

`

var showGwInfo = flag.Bool("show", false, "Print out gateway information needed for GW registration")

func main() {
	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		cmd := os.Args[0]
		fmt.Printf(usageExamples, cmd, cmd, cmd)
	}
	flag.BoolVar(showGwInfo, "s", *showGwInfo, "Print out gateway information needed for GW registration (shortcut)")
	flag.Parse()

	if *showGwInfo {
		info, err := gateway_info.GetFormatted()
		if err != nil {
			log.Print(err)
			os.Exit(1)
		}
		fmt.Print(info)
		os.Exit(0)
	}

	b := service.NewBootstrapper(nil)
	if err := b.Initialize(); err != nil {
		controlProxyConfigJson, _ := json.MarshalIndent(config.GetControlProxyConfigs(), "", "  ")
		magmadProxyConfigJson, _ := json.MarshalIndent(config.GetMagmadConfigs(), "", "  ")
		log.Fatalf(
			"gateway '%s' bootstrap initialization error: %v, for configuration:\ncontrol_proxy: %s\nmagmad: %s",
			b.HardwareId, err, string(controlProxyConfigJson), string(magmadProxyConfigJson))
	}
	// Main bootstrapper loop
	log.Print("Starting Bootstrapper")
	for {
		err := b.Start() // Start will only return on error
		if err != nil {
			log.Print(err)
			time.Sleep(service.BOOTSTRAP_RETRY_INTERVAL)
			b.RefreshConfigs()
		} else {
			log.Fatal("unexpected Bootstrapper state")
		}
	}
}
