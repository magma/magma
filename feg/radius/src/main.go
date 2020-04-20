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
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"fbc/cwf/radius/config"
	"fbc/cwf/radius/loader"
	"fbc/cwf/radius/monitoring"
	"fbc/cwf/radius/server"

	_ "magma/orc8r/lib/go/initflag"

	"go.uber.org/zap"
)

const versionFile = "/app/VERSION"

func createLogger(encoding string) (*zap.Logger, error) {
	if encoding == "json" {
		return zap.NewProduction()
	}
	return zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
}

func main() {
	var configFilename, logEncoding string
	// Get configuration
	flag.StringVar(&configFilename, "config", "radius.config.json", "The configuration filename")
	flag.StringVar(&logEncoding, "log_fmt", "json", "Log encoding format, accepted values: 'json', 'console'")
	flag.Parse()

	// Get a simple stdout logger
	logger, err := createLogger(logEncoding)

	config, err := config.Read(configFilename)
	if err != nil {
		logger.Error("Failed to read configuration", zap.Error(err))
		return
	}

	// Initialize pprof debug interface
	if config.Debug != nil {
		if config.Debug.Enabled {
			logger.Info("Enabling Server Debugging", zap.Int("port", config.Debug.Port))
			go func() {
				err = http.ListenAndServe(fmt.Sprintf(":%d", config.Debug.Port), nil)
				if err != nil {
					logger.Fatal("Debug pprof endpint failed", zap.Error(err))
				}
			}()
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
	if config.Monitoring.Scuba != nil {
		logger = logger.With(zap.String("partner_shortname", config.Monitoring.Scuba.PartnerShortName))
	}
	logger = logger.With(zap.String("app_version", GetVersion()))
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

// GetVersion returns version of the app
func GetVersion() string {
	data, err := ioutil.ReadFile(versionFile)
	if err != nil {
		return "NA"
	}
	return strings.TrimSuffix(string(data), "\n")
}
