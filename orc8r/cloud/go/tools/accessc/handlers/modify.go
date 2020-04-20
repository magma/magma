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

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/accessd"
	"magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/tools/commands"
)

// Add-existing command - Creates a new Operator and its ACL from specified
// Entities with required permissions and associates it with the given by -cert
// flag existing client Certificate
func init() {
	cmd := CommandRegistry.Add(
		"modify",
		"Modify an existing operator's Operator ACL",
		modify)
	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s %s [OPTIONS] <OperatorID>\n",
			os.Args[0], cmd.Name())
		f.PrintDefaults()
	}
	entHelp := "%s with required permissions in the form: <network Id|*>:" +
		"R|W|RW. Default (no networks, operators & gateways): '*:RW'"
	f.Var(&networks, "n", fmt.Sprintf(entHelp, "Networks"))
	f.Var(&operators, "o", fmt.Sprintf(entHelp, "Operators"))
	f.Var(&gateways, "g", fmt.Sprintf(entHelp, "Gateways"))
	f.BoolVar(&makeAdmin, "admin", false, "Make the Operator an administrator")
}

func modify(cmd *commands.Command, args []string) int {
	f := cmd.Flags()
	oid := strings.TrimSpace(f.Arg(0))
	if f.NArg() != 1 || len(oid) == 0 {
		f.Usage()
		log.Fatalf("A single Operator Id must be specified.")
	}
	// Find Operator Identity for the oid
	operator := identity.NewOperator(oid)
	cn := operator.ToCommonName()
	if cn == nil {
		log.Fatalf("Invalid common name for %s", oid)
	}
	opname := operator.HashString()
	aclMap, err := accessd.GetOperatorACL(operator)
	if err != nil {
		log.Fatalf("Operator %s does not exist. Error: %v", opname, err)
	}
	certSNs, err := certifier.FindCertificates(operator)
	if err != nil {
		log.Printf("Error %s getting certificates for %s\n", err, opname)
		certSNs = []string{}
	}
	fmt.Print("Updating Operator:\n")
	PrintACL(&protos.AccessControl_List{Operator: operator, Entities: aclMap}, certSNs)
	var acl []*protos.AccessControl_Entity
	if makeAdmin {
		acl = CreateAdminACL()
	} else {
		// Added a certificate for the operator, now add operator's ACL
		acl = BuildACLForEntities(networks, operators, gateways)
	}
	fmt.Print("\t\tTo ACL:\n")
	for _, ent := range acl {
		fmt.Printf(
			"\t\t  %s: %s (%d)\n",
			ent.GetId().HashString(),
			ent.Permissions.ToString(),
			ent.Permissions)
	}
	err = accessd.SetOperator(operator, acl)
	if err != nil {
		log.Fatalf("Set Operator %s ACL Error: %s", operator.HashString(), err)
	}
	return 0
}
