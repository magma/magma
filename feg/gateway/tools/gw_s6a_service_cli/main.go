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
	"os"
	"strconv"

	"magma/feg/cloud/go/protos"
	_ "magma/feg/gateway/registry"
	"magma/feg/gateway/services/s6a_proxy"
)

func main() {
	// setting up flags of the CLI
	helpPtr := flag.Bool("help", false, "[optional] Display this help message")
	cmdPtr := flag.String("rpcCall", "", "[required] The RPC call on the service. "+
		"{CLR|RSR}")

	// setting up helper message of the CLI
	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("	gw_s6a_service_cli [-h] " +
			"-rpcCall={CLR|RSR} <args for flags if any>")
		fmt.Println("Flags: ")
		fmt.Printf("	%s: %s\n", "rpcCall", flag.Lookup("rpcCall").Usage)
		fmt.Printf("	%s: %s\n", "help   ", flag.Lookup("help").Usage)
	}

	// parse command line inputs
	flag.Parse()

	// print usage if requested or if required arguments are not provided
	if *helpPtr || *cmdPtr == "" {
		flag.Usage()
		os.Exit(0)
	}

	// handle commands, make corresponding rpc calls, and print results
	err := handleCommands(*cmdPtr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func handleCommands(cmd string) error {

	switch cmd {
	case "CLR":
		if len(flag.Args()) != 2 {
			printSendCLRUsage()
			return fmt.Errorf("invalid args")
		}
		username := flag.Arg(0)
		clTypeStr := flag.Arg(1)
		clType, err := strconv.Atoi(clTypeStr)
		if err != nil {
			return fmt.Errorf("clType is not an integer")
		}
		if !isValidCLType(clType) {
			printSendCLRUsage()
			return fmt.Errorf("invalid cl type.")
		}
		return sendCLR(username, clType)
	case "RSR":
		req := &protos.ResetRequest{UserId: []string{}}
		req.UserId = append(req.UserId, flag.Args()...)
		rsa, err := s6a_proxy.GWS6AProxyReset(req)
		if err != nil {
			return fmt.Errorf("err sending RSR: %v", err)
		}
		fmt.Printf("Got rsa: %+v\n", rsa)
		return nil
	default:
		flag.Usage()
		return fmt.Errorf("command %s is not supported", cmd)
	}
}

func isValidCLType(clType int) bool {
	return clType >= 0 && clType <= 4
}

func printSendCLRUsage() {
	fmt.Println("Specify a username and a cltype" +
		" at the end of the command to set log_level")
	fmt.Println("Example: gw_s6a_service_cli -rpcCall=CLR user123 0")
}

func sendCLR(username string, clType int) error {
	cla, err := s6a_proxy.GWS6AProxyCancelLocation(
		&protos.CancelLocationRequest{
			UserName:         username,
			CancellationType: protos.CancelLocationRequest_CancellationType(clType),
		})
	if err != nil {
		return fmt.Errorf("err sending CLR: %v", err)
	}
	fmt.Printf("Got cla code: %v\n", cla)
	return nil
}
