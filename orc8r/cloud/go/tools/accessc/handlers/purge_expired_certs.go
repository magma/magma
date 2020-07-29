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

// Package handlers implements individual accessc commands as well as common
// across multiple commands functionality
package handlers

import (
	"fmt"
	"log"
	"os"

	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/tools/commands"
)

// List-certs command - prints out all registered certificates & associated with
// them Identities
func init() {
	cmd := CommandRegistry.Add(
		"collect-garbage",
		"Remove all expired certificates",
		collectGarbage)
	cmd.Flags().Usage = func() {
		fmt.Printf("\tUsage: %s %s\n", os.Args[0], cmd.Name())
	}
}

func collectGarbage(cmd *commands.Command, args []string) int {
	err := certifier.CollectGarbage()
	if err != nil {
		log.Fatalf("Garbage Collection Error: %s", err)
	}
	return 0
}
