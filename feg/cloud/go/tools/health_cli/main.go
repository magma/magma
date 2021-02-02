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
	"time"

	"magma/feg/cloud/go/services/health"
)

func main() {
	// setting up flags of the CLI
	helpPtr := flag.Bool("help", false, "[optional] Display this help message")
	cmdPtr := flag.String("rpcCall", "", "[required] The RPC call on the health service. "+
		"{GetHealth, GetActiveGateway}")

	// setting up helper message of the CLI
	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("	health_cli [-h] " +
			"-rpcCall={GetHealth, GetActiveGateway} <args for flags if any>")
		fmt.Println("	Example usage: ./health_cli -rpcCall=GetHealth feg_test gw1")
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
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func handleCommands(cmd string) error {
	if cmd == "GetHealth" {
		if len(flag.Args()) != 2 {
			printGetHealthUsage()
			return fmt.Errorf("Please provide <networkID> <gatewayID>")
		}
		networkID := flag.Arg(0)
		gatewayID := flag.Arg(1)
		return getHealth(networkID, gatewayID)
	} else if cmd == "GetActiveGateway" {
		if len(flag.Args()) != 1 {
			printGetActiveGatewayUsage()
			return fmt.Errorf("Please provide <networkID>")
		}
		networkID := flag.Arg(0)
		return getActiveGateway(networkID)
	}
	return fmt.Errorf("Invalid rpcCall specified")
}

func getHealth(networkID string, gatewayID string) error {
	healthStats, err := health.GetHealth(networkID, gatewayID)
	if err != nil {
		return err
	}
	fmt.Println(" -------------------- HEALTH DATA --------------")
	fmt.Println(" -----------------------------------------------")
	fmt.Println("NetworkID: ", networkID)
	fmt.Println("GatewayID: ", gatewayID)
	fmt.Println("Health:    ", healthStats.Health.Health)
	fmt.Println(" -----------------------------------------------")
	fmt.Printf("Determined using the following data:\n\n")
	for name, status := range healthStats.ServiceStatus {
		fmt.Println("Service:", name)
		fmt.Println("\tAvailability:", status.ServiceState)
		fmt.Println("\tHealth:      ", status.ServiceHealthStatus.Health)
		fmt.Println("\tReason:      ", status.ServiceHealthStatus.HealthMessage)
	}
	fmt.Println("System")
	fmt.Println("\tCpu Util Pct:          ", healthStats.SystemStatus.CpuUtilPct)
	fmt.Println("\tMemory Available Bytes:", healthStats.SystemStatus.MemAvailableBytes)
	fmt.Println("\tMemory Total Bytes:    ", healthStats.SystemStatus.MemTotalBytes)
	fmt.Println("Time since last update (sec):", time.Now().Unix()-(int64(healthStats.Time)/1000))
	return nil
}

func getActiveGateway(networkID string) error {
	activeID, err := health.GetActiveGateway(networkID)
	if err != nil {
		return err
	}
	fmt.Println("NetworkID:        ", networkID)
	fmt.Println("Active GatewayID: ", activeID)
	return nil
}

func printGetHealthUsage() {
	fmt.Println("GetHealth rpcCall requires <networkID> and <gatewayID> args")
	fmt.Println("Example: ./health_cli -rpcCall=GetHealth example_network gw1")
}

func printGetActiveGatewayUsage() {
	fmt.Println("GetActiveGateway rpcCall requires <networkID> arg")
	fmt.Println("Example: ./health_cli -rpcCall=GetActiveGateway example_network")
}
