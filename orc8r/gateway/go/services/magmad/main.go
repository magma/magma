/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime/debug"
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
	"magma/orc8r/lib/go/profile"
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

var (
	showGwInfo      = flag.Bool("show", false, "Print out gateway information needed for GW registration")
	gcPercent       = flag.Int("gc_percent", 20, "GC Percent")
	freeMemInterval = flag.Duration("memory_purge_interval", time.Hour*6, "Force GC & unused memory purge interval")
)

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

	if _, isset := os.LookupEnv("GOGC"); !isset {
		debug.SetGCPercent(*gcPercent)
	}
	eventChan := make(chan interface{}, 2)

	// Start event loop in a dedicated routine
	go mainEventLoop(eventChan, time.Tick(*freeMemInterval))

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

func mainEventLoop(eventChan <-chan interface{}, freeMemChan <-chan time.Time) {
	for {
		select {
		case evnt := <-eventChan:
			switch e := evnt.(type) {
			case bootstrapper.BootstrapCompletion:
				if e.Result != nil {
					glog.Errorf("bootstrap failure: %v for Gateway ID: %s", e.Result, e.HardwareId)
				} else {
					glog.Infof("bootstrapped GW %s", e.HardwareId)
					controller := service_manager.Get()
					if config.GetControlProxyConfigs().ProxyCloudConnection {
						glog.Info("restarting control_proxy")
						if err := controller.Restart("control_proxy"); err != nil {
							glog.Errorf("failed to restart control_proxy: %v", err)
						}
					} else {
						// Restart all magma services
						go func() {
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
			} // switch
		case _, ok := <-freeMemChan:
			if ok {
				glog.Info("purging unused memory")
				debug.FreeOSMemory()
				if glog.V(2) {
					profile.LogMemStats()
				}
				// write out heap profile if built with -tags with_profiler, noop otherwise
				// to use:
				//    go tool pprof -http=127.0.0.1:9999 <path/to/magmad> <profiles_dir/memory_MMDD_HH.mm.SS.pprof>
				profile.MemWrite()
			}
		} // select
	}
}
