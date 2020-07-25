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
//
// A helper tool to test/explore Go TLS related features & to create CA and
// client/server certificate chains for magma development and testing
// Crypto related implementation is derived from crypto/tls/generate_cert.go
package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"

	"magma/orc8r/lib/go/security/cert"
	"magma/orc8r/lib/go/security/key"
	"magma/orc8r/lib/go/util"
)

var (
	outName = flag.String("o", "",
		"Name used for certificate (N.pem) & key (N.key.pem) file names. "+
			"Defaults:'client_cert'/'server_cert' for client/server certificates")
	hosts = flag.String("hosts", "",
		"Comma-separated hostnames and IPs to generate a certificate for"+
			" (optional for client certificate)")
	sans = flag.String("sans", "",
		"Comma-separated SANs to generate a certificate for"+
			" (the names will be added to DNS SANs of the certificate)")
	email = flag.String("email", "",
		"Comma-separated email addresses for the certificate (optional)")
	validFor = flag.Duration("duration", 365,
		"Duration that certificate is valid for (days)")
	isCA     = flag.Bool("ca", false, "Whether this cert is a CA Certificate")
	isClient = flag.Bool("client", false, "Whether this cert is a Client Certificate")
	rsaBits  = flag.Int("rsa-bits", 2048,
		"Size of RSA key to generate. Ignored if --ecdsa-curve is set")
	ecdsaCurve = flag.String("ecdsa-curve", "",
		"ECDSA curve to use to generate a key. One of: P224, P256, P384, P521")

	caCertFile = flag.String("cac", "", "CA's Certificate file")
	caKeyFile  = flag.String("cak", "", "CA's Private Key file")

	country    = flag.String("C", "US", "Country (C)")
	org        = flag.String("O", "", "Organization (O)")
	commonName = flag.String("CN", "", "Common Name (CN)")
	orgUnit    = flag.String("OU", "", "Organizational Unit (OU)")

	printCert = flag.Bool("text", false,
		"Print out certificate information from given cert files")
	verifyCert = flag.Bool("verify", false,
		"Verify given certificates (optionally: with the CA cert)")
)

func printCertificate(certFile string) {
	certPEMBlock, err := ioutil.ReadFile(certFile)
	if err != nil {
		log.Fatalf("Cannot read certificate file '%s': %s\n", certFile, err)
	}
	for {
		var certDERBlock *pem.Block
		certDERBlock, certPEMBlock = pem.Decode(certPEMBlock)
		if certDERBlock == nil {
			break
		}
		fmt.Printf("Found '%s' DER Block\n", certDERBlock.Type)
		if certDERBlock.Type == "CERTIFICATE" {
			x509Cert, err := x509.ParseCertificate(certDERBlock.Bytes)
			if err != nil {
				log.Fatalf(
					"Cannot parse certificate from '%s': %s\n", certFile, err)
			}

			certJson, err := json.MarshalIndent(x509Cert, "", "    ")

			var strOut string
			if err != nil {
				strOut = "Certificate Marshaling Error: \n" + err.Error()
			} else {
				strOut = string(certJson)
			}
			fmt.Printf("\n%s\n", strOut)
		}
	}
}

const usageExamples string = `
Examples:

  1. Create a self-signed "root" certificate for cloud.magma domain:

    $> %s -ca -hosts cloud.magma.facebook.com

    The command will create server_cert.pem & server_cert.key.pem
    in the current directory

  2. Create a intermediate CA certificate for cloud2.magma domain signed by server_cert.pem:

    $> %s -ca -cac=server_cert.pem -cak=server_cert.key.pem -hosts cloud2.magma.facebook.com -o intermediate_ca_cert

    The command will create intermediate_ca_cert.pem & intermediate_ca_cert.key.pem
    in the current directory

  3. Create a client certificate signed by the server CA (server_cert from the
     example above):

    $> %s -client -cac=server_cert.pem -cak=server_cert.key.pem

    The command will create client_cert.pem & client_cert.key.pem
    in the current directory

  4. Print out JSON encoded representation of certificate from client_cert.pem:

    $> %s -text client_cert.pem

`

func main() {
	oldUsage := flag.Usage
	flag.Usage = func() {
		oldUsage()
		cmd := os.Args[0]
		fmt.Printf(usageExamples, cmd, cmd, cmd, cmd)
	}
	flag.Parse()

	if *verifyCert {
		verifyCertificate()
		return
	}

	if *printCert {
		for _, cf := range flag.Args() {
			fmt.Printf("Parsing '%s' certificate file\n", cf)
			printCertificate(cf)
		}
		return
	}

	if len(*outName) == 0 {
		if *isClient {
			*outName = "client_cert"
		} else {
			*outName = "server_cert"
		}
	}
	certFile := *outName + ".pem"
	keyFile := *outName + ".key.pem"

	if len(*hosts) == 0 && !(*isClient || *isCA) {
		log.Fatalf("Missing required -hosts parameter\n")
	}

	priv, err := key.GenerateKey(*ecdsaCurve, *rsaBits)
	if err != nil {
		log.Fatalf("Failed to generate private key: %s", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(*validFor * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %s", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{*country},
			Organization:       []string{*org},
			OrganizationalUnit: []string{*orgUnit},
			CommonName:         *commonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	fmt.Printf("\n\nCOMMON NAME: %s\n", *commonName)

	certHosts := strings.Split(*hosts, ",")
	for _, h := range certHosts {
		if len(h) == 0 {
			continue
		}
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
			fmt.Printf("Adding Cert IP SAN: %s\n", h)
		} else {
			template.DNSNames = append(template.DNSNames, h)
			fmt.Printf("Adding Cert DNS Name SAN: %s\n", h)
		}
	}

	emails := strings.Split(*email, ",")
	for _, e := range emails {
		if len(e) > 0 {
			template.EmailAddresses = append(template.EmailAddresses, e)
		}
	}

	certSans := strings.Split(*sans, ",")
	for _, san := range certSans {
		if len(san) == 0 {
			continue
		}
		template.DNSNames = append(template.DNSNames, san)
		fmt.Printf("Adding Cert DNS SAN: %s\n", san)
	}

	if *isClient {
		template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	} else {
		if *isCA {
			template.ExtKeyUsage = []x509.ExtKeyUsage{
				x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
		} else {
			template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		}
	}

	template.IsCA = *isCA
	if *isCA {
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	ski := make([]byte, 32)
	rand.Read(ski)
	template.SubjectKeyId = ski

	var caCert *x509.Certificate
	var caPrivKey interface{}
	if len(*caCertFile) > 0 && len(*caKeyFile) > 0 {
		caCert, caPrivKey, err = cert.LoadCertAndPrivKey(*caCertFile, *caKeyFile)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Creating Certificate for Subject: %s\n",
			util.FormatPkixSubject(&template.Subject))
		fmt.Printf(
			"Using CA Subject: %s\n", util.FormatPkixSubject(&caCert.Subject))
	} else {
		caCert = &template
		caPrivKey = priv
		fmt.Printf("Creating Self-signed certificate for Subject: %s\n",
			util.FormatPkixSubject(&template.Subject))
	}

	derBytes, err := x509.CreateCertificate(
		rand.Reader, &template, caCert, key.PublicKey(priv), caPrivKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}

	certOut, err := os.Create(certFile)
	if err != nil {
		log.Fatalf("Failed to open %s for writing: %s", certFile, err)
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	log.Printf("Written %s\n", certFile)

	err = key.WriteKey(keyFile, priv)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Written %s\n", keyFile)
}

func verifyCertificate() {
	var caPool *x509.CertPool = x509.NewCertPool()
	if len(*caCertFile) > 0 {
		// Read given CA PEM file
		certPEMBlock, err := ioutil.ReadFile(*caCertFile)
		if err != nil {
			log.Fatalf("Cannot read CA file '%s': %s\n", *caCertFile, err)
			return
		}
		// create ca pool from CA PEM file
		caPool = x509.NewCertPool()
		ok := caPool.AppendCertsFromPEM(certPEMBlock)
		if !ok {
			log.Printf("Cannot Load Certificates from: %s\n", *caCertFile)
		}
		caPoolLen := len(caPool.Subjects())

		if caPoolLen > 0 {
			fmt.Printf("Using CA Certificate Pool of size: %d\n", caPoolLen)
		} else {
			log.Printf("WARNING: Empty CA Certificate Pool\n")
		}
	} else {
		caPool = nil
		log.Printf("Using System Certificate Pool")
	}

	// Determine intended certificate usage
	extUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	certType := "SERVER"

	if *isClient {
		extUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
		certType = "CLIENT"
	} else if *isCA {
		extUsage = append(extUsage, x509.ExtKeyUsageClientAuth)
	}

	// Iterate over all given via command line cert files to verify
	for _, cf := range flag.Args() {
		fmt.Printf("Verifying %s certificate[s] from '%s'\n", certType, cf)
		certPEMBlock, err := ioutil.ReadFile(cf)
		if err != nil {
			log.Printf("Cannot read certificate file '%s': %s\n", cf, err)
			continue
		}

		var certs []*x509.Certificate
		var block *pem.Block
		for n := 0; len(certPEMBlock) > 0; n++ {
			block, certPEMBlock = pem.Decode(certPEMBlock)
			if block == nil {
				break
			}
			// find CERTIFICATE block
			if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
				continue
			}
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				continue
			}
			certs = append(certs, cert)
		}

		if len(certs) == 0 {
			log.Printf("No valid certificate found in '%s'\n", cf)
			printCertificate(cf)
			return
		}
		// verify parsed cert chain
		errs := verify(certs, extUsage, caPool)
		if len(errs) > 0 {
			log.Printf("Failed certificate '%s' verification:\n", cf)
			printCertificate(cf)
			for _, err = range errs {
				log.Printf("\tERROR: %s\n", err)
			}
			return
		}
		fmt.Printf("SUCCESS: Verified Certificate %s\n", cf)
	}
}

var vulnerableSigAlgos = map[x509.SignatureAlgorithm]string{
	x509.MD2WithRSA:    "MD2 with RSA",
	x509.MD5WithRSA:    "MD5 with RSA",
	x509.SHA1WithRSA:   "SHA1 with RSA",
	x509.DSAWithSHA1:   "DSA with SHA1",
	x509.ECDSAWithSHA1: "ECDSA with SHA1",
}

//
// verify - verefies given cert using caPool (or system CA pool if caPool is nil)
// for given usage flags
func verify(
	certs []*x509.Certificate,
	usage []x509.ExtKeyUsage,
	caPool *x509.CertPool) []error {

	result := []error{}
	if len(certs) == 0 {
		return result
	}

	opts := x509.VerifyOptions{
		Roots:         caPool,
		Intermediates: x509.NewCertPool(),
		KeyUsages:     usage,
	}

	for i, cert := range certs {
		if cert == nil {
			result = append(
				result,
				fmt.Errorf("Nil certificate at index: %d", i))
			if i == 0 {
				return result
			}
			continue
		}
		sigAlgoName, ok := vulnerableSigAlgos[cert.SignatureAlgorithm]
		if ok {
			result = append(
				result,
				fmt.Errorf("Depreciated signature algorithm '%s'", sigAlgoName))
		}
		switch cert.PublicKey.(type) {
		case *rsa.PublicKey, *ecdsa.PublicKey:
			break
		default:
			result = append(
				result,
				fmt.Errorf("Unsupported type of public key: %T  at index: %d",
					cert.PublicKey, i))
		}
		if i > 0 {
			opts.Intermediates.AddCert(cert)
		}
	}

	_, err := certs[0].Verify(opts)
	if err != nil {
		result = append(
			result,
			fmt.Errorf("Failed certificate verification: %s", err))
	}
	return result
}
