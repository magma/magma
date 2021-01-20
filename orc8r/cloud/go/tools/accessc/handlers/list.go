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

	"magma/orc8r/cloud/go/services/accessd"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/tools/commands"
)

// List command - prints out all registered Operators and their attributes
func init() {
	cmd := CommandRegistry.Add(
		"list",
		"List all Operators, their Certificate Serial Numbers and ACLs",
		list)
	cmd.Flags().Usage = func() {
		fmt.Fprintf(os.Stderr, "\tUsage: %s %s\n", os.Args[0], cmd.Name())
	}
}

func list(cmd *commands.Command, args []string) int {
	operators, err := accessd.ListOperators()
	if err != nil {
		log.Fatalf("List Operators Error: %s", err)
	}
	acls, err := accessd.GetOperatorsACLs(operators)
	if err != nil {
		log.Fatalf("Get Operators Error: %s", err)
	}
	// Retrieve all certificate records
	// It can be memory intensive, but Cert Info record is relatively small ~ 90 bytes on average
	certMap, err := certifier.GetAll()
	if err != nil {
		log.Fatalf("Get Certificates Error: %s", err)
	}
	// operator hash to list of cert SNs map
	// this map should be small, it will only include system's Operator entities
	certSnsByOperMap := map[string][]string{}
	// Create an empty list entry for each existing operator
	for _, operator := range operators {
		certSnsByOperMap[operator.HashString()] = []string{}
	}
	// Reverse mapping - add every found Cert SN to its operator's list
	for sn, cinfo := range certMap {
		op := cinfo.GetId()
		hash := op.HashString()
		if snlist, ok := certSnsByOperMap[hash]; ok {
			certSnsByOperMap[hash] = append(snlist, sn)
		}
	}
	fmt.Println("Operators:")
	for idx, acl := range acls {
		if acl == nil {
			log.Printf("ERROR: Nil ACL at %d index", idx)
			continue
		}
		// Get Cert SNs for this operator
		opname := acl.GetOperator().HashString()
		certSNs, ok := certSnsByOperMap[opname]
		if !ok {
			certSNs = []string{}
			log.Printf("Error Finding certificates for %s", opname)
		}
		PrintACL(acl, certSNs)
	}
	return 0
}
