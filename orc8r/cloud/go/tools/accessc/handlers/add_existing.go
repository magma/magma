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
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/accessd"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/tools/commands"
	"magma/orc8r/lib/go/security/cert"
)

// Add-existing command - Creates a new Operator and its ACL from specified
// Entities with required permissions and associates it with the given by -cert
// flag existing client Certificate
func init() {
	cmd := CommandRegistry.Add(
		"add-existing",
		"Add a new Operator and associate it with an existing Certificate",
		addExisting)
	f := cmd.Flags()
	f.Usage = func() {
		fmt.Fprintf(os.Stderr, // std Usage() & PrintDefaults() use Stderr
			"\tUsage: %s %s [OPTIONS] -cert <client cert PEM> <OperatorID>\n",
			os.Args[0], cmd.Name())
		f.PrintDefaults()
	}

	f.StringVar(&certFName, "cert", "",
		"The existing certificate PEM file to use for the new Operator")

	entHelp := "%s with required permissions in the form: <network Id|*>:" +
		"R|W|RW. Default (no networks, operators & gateways): '*:RW'"
	f.Var(&networks, "n", fmt.Sprintf(entHelp, "Networks"))
	f.Var(&operators, "o", fmt.Sprintf(entHelp, "Operators"))
	f.Var(&gateways, "g", fmt.Sprintf(entHelp, "Gateways"))
	f.BoolVar(&makeAdmin, "admin", false, "Make the new Operator an administrator")
}

func addExisting(cmd *commands.Command, args []string) int {
	f := cmd.Flags()
	oid := strings.TrimSpace(f.Arg(0))
	if f.NArg() != 1 || len(oid) == 0 {
		f.Usage()
		log.Fatalf("A single Operator Id must be specified.")
	}
	// Create Operator Identity for the oid
	operator := identity.NewOperator(oid)
	cn := operator.ToCommonName()
	if cn == nil {
		log.Fatalf("Invalid common name for %s", oid)
	}
	certFName := strings.TrimSpace(certFName)
	if len(certFName) == 0 {
		f.Usage()
		log.Fatalf("Certificate file name (-cert ...) must be provided.")
	}
	certPEMBlock, err := ioutil.ReadFile(certFName)
	if err != nil {
		log.Fatalf("Cannot read certificate file '%s': %s", certFName, err)
	}
	// Find first valid certificate
	for certDERBlock, _ := pem.Decode(certPEMBlock); certDERBlock != nil; certDERBlock, _ = pem.Decode(certPEMBlock) {

		if certDERBlock.Type == "CERTIFICATE" {
			x509Cert, err := x509.ParseCertificate(certDERBlock.Bytes)
			if err != nil {
				log.Fatalf(
					"Cannot parse certificate from '%s': %s\n", certFName, err)
			}
			err = certifier.AddCertificate(operator, certDERBlock.Bytes)
			if err != nil {
				log.Fatalf(
					"Error '%s' Adding Certificate From %s", err, certFName)
			}

			var acl []*accessprotos.AccessControl_Entity
			if makeAdmin {
				acl = CreateAdminACL()
			} else {
				// Added a certificate for the operator, now add operator's ACL
				acl = BuildACLForEntities(networks, operators, gateways)
			}

			err = accessd.SetOperator(operator, acl)
			if err != nil {
				certifier.RevokeCertificateSN(cert.SerialToString(x509Cert.SerialNumber))
				log.Fatalf("Set Operator %s ACL Error: %s", operator.HashString(), err)
			}
			return 0
		}
	}
	log.Printf("ERROR: No Certificates Found in '%s'", certFName)
	return 2
}
