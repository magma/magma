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
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"flag"
	"fmt"
	"log"
	"os"

	"magma/orc8r/lib/go/security/csr"
	"magma/orc8r/lib/go/security/key"
)

var (
	keyFile = flag.String("key", "default_csr.key.pem", "Name used for key file")
	rsaBits = flag.Int("rsa-bits", 2048,
		"Size of RSA key to generate. Ignored if --ecdsa-curve is set")
	ecdsaCurve = flag.String("ecdsa-curve", "",
		"ECDSA curve to use to generate a key. One of: P224, P256, P384, P521")

	country    = flag.String("C", "CL", "Country (C)")
	org        = flag.String("O", "MagmaClient", "Organization (O)")
	commonName = flag.String("CN", "", "Common Name (CN)")
	orgUnit    = flag.String("OU", "", "Organizational Unit (OU)")
)

const usageExamples string = `
Examples:

  Create a CSR using -key keyFile and write into csrFile

	$> %s csrFile

  If keyFile does not exist, it will create a new key using ecdsaCurve and rsaBits
  and write the key into keyFile.

`

func main() {
	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		cmd := os.Args[0]
		fmt.Printf(usageExamples, cmd)
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	csrFile := flag.Arg(0)
	var priv interface{}
	var err error

	if _, err = os.Stat(*keyFile); err == nil {
		priv, err = key.ReadKey(*keyFile)
		if err != nil {
			log.Fatalf("Failed to read private key from %s: %s", *keyFile, err)
		}
		fmt.Printf("Key read from: %s\n", *keyFile)
	} else {
		priv, err = key.GenerateKey(*ecdsaCurve, *rsaBits)
		if err != nil {
			log.Fatalf("Failed to generate private key: %s", err)
		}
		err = key.WriteKey(*keyFile, priv)
		if err != nil {
			log.Fatalf("Failed to write private key to %s: %s", *keyFile, err)
		}
		fmt.Printf("Key created and written into %s\n", *keyFile)
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{*country},
			Organization:       []string{*org},
			OrganizationalUnit: []string{*orgUnit},
			CommonName:         *commonName,
		},
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, priv)
	if err != nil {
		log.Fatalf("Failed to create certificate request: %s", err)
	}

	err = csr.WriteCSR(csrDER, csrFile)
	if err != nil {
		log.Fatalf("Failed to write certificate request: %s", err)
	}
	fmt.Printf("CSR written into %s\n", csrFile)
}
