/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"magma/gateway/services/bootstrapper/service"
	"magma/orc8r/lib/go/security/key"
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

	b := service.NewBootstrapper()
	if err := b.Initialize(); err != nil {
		log.Fatalf("Init error: %v", err)
	}
	if *showGwInfo {
		fmt.Printf("\nHardware ID:\n------------\n%s\n", b.HardwareId)
		ck, err := key.ReadKey(b.ChallengeKeyFile)
		if err != nil {
			log.Printf("Failed to load Challenge Key from %s: %v", b.ChallengeKeyFile, err)
			os.Exit(1)
			return
		}
		marshaledPubKey, err := x509.MarshalPKIXPublicKey(key.PublicKey(ck))
		if err != nil {
			log.Printf("Failed to marshal Public Challenge Key from %s: %v", b.ChallengeKeyFile, err)
			os.Exit(2)
			return
		}
		fmt.Printf("\nChallenge Key:\n--------------\n%s\n", base64.StdEncoding.EncodeToString(marshaledPubKey))
		os.Exit(0)
	}

	// Main bootstrapper loop
	configJson, _ := json.MarshalIndent(b, "", "  ")
	if err := b.Initialize(); err != nil {
		log.Fatalf("Bootstrapper Initialization error: %v, for configuration: %s", err, string(configJson))
	}
	log.Printf("Starting Bootstrapper with configuration: %s\n", string(configJson))
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
