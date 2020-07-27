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

package handlers

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/accessd"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/security/cert"
	"magma/orc8r/lib/go/security/key"

	"github.com/golang/protobuf/ptypes"
)

// Common to several add... handlers definitions & functions
var (
	networks, operators, gateways           Entities
	email, certFName, country, org, orgUnit string
	rsaBits                                 int
	validFor                                uint
)

// addInit is shared between add & add-admin
func addInit(f *flag.FlagSet) {
	f.StringVar(&email, "email", "",
		"Comma-separated email addresses for the certificate (optional)")

	f.UintVar(&validFor, "duration", 365,
		"Duration that the Operator's certificate is valid for (days)")

	f.IntVar(&rsaBits, "rsa-bits", 2048,
		"Size of RSA key to generate. Ignored if --ecdsa-curve is set")

	f.StringVar(&country, "C", "US", "Country (C)")
	f.StringVar(&org, "O", "", "Organization (O)")
	f.StringVar(&orgUnit, "OU", "", "Organizational Unit (OU)")
}

// addAcl Creates operator with oid, new certificate for this operator,
// associates the certificate with the operator (using Certifier) and creates
// accessd ACL for the operator using permissions provided by acl parameter
//
// addACL is used by add & add_admin handlers
func addACL(oid string, acl []*accessprotos.AccessControl_Entity) int {

	// Create Operator Identity for the oid
	operator := identity.NewOperator(oid)

	// Generate & sign a new Operator certificate
	priv, err := key.GenerateKey("", rsaBits)
	if err != nil {
		log.Fatalf("Failed to create key: %s", err)
	}
	outfile := strings.TrimSpace(certFName)
	if len(outfile) == 0 {
		outfile = oid + "_cert"
	}
	certFile := outfile + ".pem"
	keyFile := outfile + ".key.pem"
	// Create new certificate file
	certOut, err := os.Create(certFile)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %s", certFile, err)
	}
	// Save a new cert key before calling Certifier, it's easier to delete files
	// then remove Certificate and/or ACL from the services
	err = key.WriteKey(keyFile, priv)
	if err != nil {
		cleanupCertFiles(keyFile, certFile, certOut)
		log.Fatalf(
			"Failed to write private key into %s, Error: %s", keyFile, err)
	}
	// Create & sign new cert
	certDer, sn, err := createAndSignCert(operator, priv)
	if err != nil {
		cleanupCertFiles(keyFile, certFile, certOut)
		log.Fatalf(
			"Failed to write private key into %s, Error: %s", keyFile, err)
	}

	// Save new cert
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDer})
	certOut.Close()

	signedCert, err := x509.ParseCertificates(certDer)
	if err != nil || len(signedCert) == 0 {
		cleanupCertFiles(keyFile, certFile, nil)
		log.Fatalf("Certificate Parsing Error: %s", err)
	}

	if sn != cert.SerialToString(signedCert[0].SerialNumber) {
		log.Printf(
			"\nCertifier's & DER's Certificate Serial Number mismatch %s != %s\n",
			sn, cert.SerialToString(signedCert[0].SerialNumber))
	}

	err = accessd.SetOperator(operator, acl)
	if err != nil {
		certifier.RevokeCertificateSN(sn)
		cleanupCertFiles(keyFile, certFile, nil)
		log.Fatalf("Set Operator %s ACL Error: %s", operator.HashString(), err)
	}

	return 0
}

// createAndSignCert does it using values given by country, org, orgUnit and
// validFor command line flags or their defaults
func createAndSignCert(
	operator *protos.Identity,
	privKey interface{},
) (der []byte, sn string, err error) {

	cn := operator.ToCommonName()
	if cn == nil {
		return der, sn, fmt.Errorf("Invalid common name for %s", operator.HashString())
	}
	// Ask Certifier to sign the new CSR
	csrParams := x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{country},
			Organization:       []string{org},
			OrganizationalUnit: []string{orgUnit},
			CommonName:         *cn,
		},
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &csrParams, privKey)
	if err != nil {
		return der, sn, fmt.Errorf("CSR Create Error: %s", err)
	}
	csr := &protos.CSR{
		Id:        operator,
		ValidTime: ptypes.DurationProto(time.Hour * 24 * time.Duration(validFor)),
		CsrDer:    csrDER,
	}
	certMsg, err := certifier.SignCSR(csr)
	if err != nil {
		return der, sn, fmt.Errorf("Certifier Sign Error: %s", err)
	}
	return certMsg.GetCertDer(), certMsg.Sn.GetSn(), nil
}

// Cleanup generated files after an error
func cleanupCertFiles(keyFile, certFile string, certOut *os.File) {
	if certOut != nil {
		certOut.Close()
	}
	if len(certFile) > 0 {
		os.Remove(certFile)
	}
	if len(keyFile) > 0 {
		os.Remove(keyFile)
	}
}
