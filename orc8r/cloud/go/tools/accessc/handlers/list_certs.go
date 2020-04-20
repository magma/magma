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
	"time"

	"github.com/golang/protobuf/ptypes"

	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/tools/commands"
)

// List-certs command - prints out all registered certificates & associated with
// them Identities
func init() {
	cmd := CommandRegistry.Add(
		"list-certs",
		"List all registered certificates & associated and their Identities",
		listCerts)
	cmd.Flags().Usage = func() {
		fmt.Printf("\tUsage: %s %s\n", os.Args[0], cmd.Name())
	}
}

func listCerts(cmd *commands.Command, args []string) int {
	certs, err := certifier.ListCertificates()
	if err != nil {
		log.Fatalf("List Certificates Error: %s", err)
	}
	fmt.Println()
	for _, csn := range certs {
		fmt.Printf("Serial Number: %s", csn)
		info, err := certifier.GetCertificateIdentity(csn)
		if err != nil || info == nil {
			log.Printf("\nError %s gettting certificate %s info\n", err, csn)
			continue
		}
		before, _ := ptypes.Timestamp(info.NotBefore)
		after, _ := ptypes.Timestamp(info.NotAfter)
		fmt.Printf(
			"; Identity: %s; Not Before: %s; Not After: %s\n",
			info.Id.HashString(),
			before.In(time.Local),
			after.In(time.Local))
	}
	fmt.Println()
	return 0
}
