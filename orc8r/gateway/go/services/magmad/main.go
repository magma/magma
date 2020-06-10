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
	"os"
	"time"

	"github.com/golang/glog"

	"magma/gateway/config"
	"magma/gateway/services/bootstrapper/gateway_info"
	bootstrapper "magma/gateway/services/bootstrapper/service"
	configurator "magma/gateway/services/configurator/service"
	"magma/gateway/services/magmad/service"
	"magma/gateway/services/magmad/service_manager"
	"magma/gateway/services/magmad/status"
	sync_rpc "magma/gateway/services/sync_rpc/service"
)

const (
	BOOTSTRAP_RESTART_INTERVAL = time.Second * 120
)

const usageExamples string = `
Examples:

  1. Run magmad as a service:

    $> %s

    The command will run magmad service which will periodically 
    check and update the gateway certificats and cloud managed GW configs

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
			glog.Error(err)
			os.Exit(1)
		}
		fmt.Print(info)
		os.Exit(0)
	}

	eventChan := make(chan interface{}, 2)

	// Start event loop in a dedicated routine
	go mainEventLoop(eventChan)

	// Create bootstrapper
	b := bootstrapper.NewBootstrapper(eventChan)
	if err := b.Initialize(); err != nil {
		controlProxyConfigJson, _ := json.MarshalIndent(config.GetControlProxyConfigs(), "", "  ")
		magmadProxyConfigJson, _ := json.MarshalIndent(config.GetMagmadConfigs(), "", "  ")
		glog.Fatalf(
			"gateway '%s' bootstrap initialization error: %v, for configuration:\ncontrol_proxy: %s\nmagmad: %s",
			b.HardwareId, err, string(controlProxyConfigJson), string(magmadProxyConfigJson))
	}
	// Start bootstrapper
	glog.Info("Starting Bootstrapper")
	go func() {
		for {
			err := b.Start() // Start will only return on error
			if err != nil {
				glog.Error(err)
				time.Sleep(BOOTSTRAP_RESTART_INTERVAL)
				b.RefreshConfigs()
			} else {
				glog.Fatal("unexpected Bootstrapper state")
			}
		}
	}()

	// Start SyncRPC service if it's enabled
	if config.GetMagmadConfigs().EnableSyncRpc {
		glog.Info("Starting SynRPC service")
		syncRpcService := sync_rpc.NewClient(nil)
		go syncRpcService.Run()
	}

	// start service status collector & reporter
	go status.StartReporter()

	// Start configurator & block on main()
	cfg := configurator.NewConfigurator(eventChan)
	glog.Info("Starting Configurator")
	go func() {
		if err := cfg.Start(); err != nil {
			glog.Fatalf("configurator start error: %v", err)
		}
	}()
	if err := service.StartMagmadServer(); err != nil {
		glog.Fatalf("magmad start error: %v", err)
	}
}

func mainEventLoop(eventChan chan interface{}) {
	for i := range eventChan {
		switch e := i.(type) {
		case bootstrapper.BootstrapCompletion:
			if e.Result != nil {
				glog.Errorf("bootstrap failure: %v for Gateway ID: %s", e.Result, e.HardwareId)
			} else {
				glog.Infof("bootstrapped GW %s", e.HardwareId)
				if config.GetControlProxyConfigs().ProxyCloudConnection {
					// TODO: restart control proxy only
				} else {
					// Restart all magma services
					go func() {
						controller := service_manager.Get()
						for _, service := range config.GetMagmadConfigs().MagmaServices {
							controller.Restart(service)
						}
					}()
				}
			}
		case configurator.UpdateCompletion:
			glog.Verbose(len(e) > 0).Infof("mconfigs updated successfully for services: %v", e)
			// Restart all services with updated configs
			go func() {
				magmaServiceTable := map[string]struct{}{}
				for _, service := range config.GetMagmadConfigs().MagmaServices {
					magmaServiceTable[service] = struct{}{}
				}
				controller := service_manager.Get()
				for _, service := range e {
					// restart only if it's this GW's service
					if _, ok := magmaServiceTable[service]; ok {
						controller.Restart(service)
					}
				}
			}()
		default:
			glog.Errorf("unknown completion type: %T", e)
		}
	}
}
