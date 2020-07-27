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
	cmdInfo := &cobra.Command{
		Use:   "info <service> [--gateway-service (--hwid=<hardware-id> | --network=<network-id> --gateway=<gateway-id>)]",
		Short: "Get service info",
		Args:  validateInfoArgs,
		Run:   infoCmd,
	}

	rootCmd.AddCommand(cmdInfo)
}

func validateInfoArgs(cmd *cobra.Command, args []string) error {
	if err := validateGlobalFlags(); err != nil {
		return err
	}
	if err := setHwIdFlag(); err != nil {
		return err
	}
	if len(args) != 1 {
		return errors.New("requires 1 arg")
	}
	if !isGatewayServiceQuery && !isValidService(args[0], services) {
		return fmt.Errorf("service %s is invalid, needs to match one of %v", args[0], services)
	}
	if isGatewayServiceQuery && !isValidGwService(gateway_registry.GwServiceType(args[0]), gwServices) {
		return fmt.Errorf("service %s is invalid, needs to match one of %v", args[0], gwServices)
	}
	return nil
}

func infoCmd(cmd *cobra.Command, args []string) {
	err := getInfo(args[0])
	if err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}

func getInfo(service string) error {
	serviceInfo, err := getInfoOrGwInfo(service)
	if err != nil {
		return fmt.Errorf("Failed to GetServiceInfo for %s: %s", service, err)
	}
	fmt.Printf("%v\n", serviceInfo)
	return nil
}

func getInfoOrGwInfo(service string) (*protos.ServiceInfo, error) {
	if isGatewayServiceQuery {
		return service303.GWService303GetServiceInfo(gateway_registry.GwServiceType(service), hardwareID)
	} else {
		return client.Service303GetServiceInfo(service)
	}
}
