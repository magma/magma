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

// Command Line Tool to create & manage Operators, ACLs and Certificates
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"magma/orc8r/cloud/go/tools/accessc/handlers"
	"magma/orc8r/lib/go/registry"
)

func main() {
	flag.Parse()
	registry.MustPopulateServices()

	// Init help for all commands
	flag.Usage = func() {
		cmd := os.Args[0]
		fmt.Printf(
			"\nUsage: \033[1m%s [GENERAL OPTIONS] command [COMMAND OPTIONS]\033[0m\n\n",
			filepath.Base(cmd))
		flag.PrintDefaults()
		fmt.Println("\nCommands:")
		handlers.CommandRegistry.Usage()
	}

	cmd, args := handlers.CommandRegistry.GetCommand()
	if cmd == nil {
		cmdName := strings.ToLower(flag.Arg(0))
		if cmdName != "" && cmdName != "help" && cmdName != "h" {
			fmt.Println("\nInvalid Command: ", cmdName)
		}
		flag.Usage()
		os.Exit(1)
	}
	cmd.Flags().Parse(args)
	os.Exit(cmd.Handle(args))
}
