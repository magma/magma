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
	"errors"
	"fmt"
	"os"

	"magma/orc8r/cloud/go/services/dispatcher/gateway_registry"
	"magma/orc8r/cloud/go/services/dispatcher/gw_client_apis/service303"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/service/client"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

func init() {
	cmdLogLevel := &cobra.Command{
		Use:   "log_level {DEBUG,INFO,WARNING,ERROR,FATAL} <service> [--gateway-service (--hwid=<hardware-id> | --network=<network-id> --gateway=<gateway-id>)]",
		Short: "Set log level",
		Args:  validateLogLevelArgs,
		Run:   logLevelCmd,
	}

	rootCmd.AddCommand(cmdLogLevel)
}

func validateLogLevelArgs(cmd *cobra.Command, args []string) error {
	if err := validateGlobalFlags(); err != nil {
		return err
	}
	if err := setHwIdFlag(); err != nil {
		return err
	}
	if len(args) != 2 {
		return errors.New("requires 2 args")
	}
	if !isValidLogLevel(args[0]) {
		return fmt.Errorf("log_level provided is not valid")
	}
	if !isGatewayServiceQuery && !isValidService(args[1], services) {
		return fmt.Errorf("service %s is invalid, needs to match one of %v", args[1], services)
	}
	if isGatewayServiceQuery && !isValidGwService(gateway_registry.GwServiceType(args[1]), gwServices) {
		return fmt.Errorf("service %s is invalid, needs to match one of %v", args[1], gwServices)
	}
	return nil
}

func logLevelCmd(cmd *cobra.Command, args []string) {
	err := setLogLevel(args[1], args[0])
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}

func setLogLevel(service string, logLevel string) error {
	err := setLogLevelOrGwLogLevel(service, logLevel)
	if err != nil {
		return fmt.Errorf("Failed to SetLogLevel for %s: %s", service, err)
	}
	return nil
}

func setLogLevelOrGwLogLevel(service string, logLevel string) error {
	if isGatewayServiceQuery {
		return service303.GWService303SetLogLevel(gateway_registry.GwServiceType(service), hardwareID, &protos.LogLevelMessage{Level: protos.LogLevel(protos.LogLevel_value[logLevel])})
	} else {
		return client.Service303SetLogLevel(service, &protos.LogLevelMessage{Level: protos.LogLevel(protos.LogLevel_value[logLevel])})
	}
}

func isValidLogLevel(level string) bool {
	return level == "DEBUG" || level == "INFO" || level == "WARNING" || level == "ERROR" || level == "FATAL"
}
