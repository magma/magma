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
	"strings"

	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/tools/commands"
)

// Revoke command - prints out all registered Operators and their attributes
func init() {
	cmd := CommandRegistry.Add(
		"revoke",
		"Revoke Given Certificate & remove it from Certifier",
		revoke)

	cmd.Flags().Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s %s <Certificate Serial Number>\n",
			os.Args[0], cmd.Name())
	}
}

func revoke(cmd *commands.Command, args []string) int {
	f := cmd.Flags()
	csn := strings.TrimSpace(f.Arg(0))
	if f.NArg() != 1 || len(csn) == 0 {
		f.Usage()
		log.Fatalf("A single Certificate Serial Number must be specified.")
	}
	fmt.Printf("Revoking Certificate Serial Number: %s\n", csn)

	err := certifier.RevokeCertificateSN(csn)
	if err != nil {
		log.Fatalf("Error %s revoking certificate %s", csn, err)
	}
	return 0
}
