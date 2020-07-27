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

// package csr contains helper functions related to certificate signing requests
package csr

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
)

// WriteCSR write a DER encoded csr (e.g returned by x509.CreateCertificateRequest)
// into file specified by csrFile
func WriteCSR(csr []byte, csrFile string) error {
	csrOut, err := os.Create(csrFile)
	if err != nil {
		return fmt.Errorf("Failed to create CSR file %s: %s", csrFile, err)
	}

	err = pem.Encode(csrOut, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr})
	if err != nil {
		return fmt.Errorf("Failed to write CSR into file: %s", err)
	}
	return nil
}

// ReadCSR read and parse CSR from csrFile and return it as *x509.CertificateRequest
func ReadCSR(csrFile string) (*x509.CertificateRequest, error) {
	csrPEM, err := ioutil.ReadFile(csrFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to open CSR file %s: %s", csrFile, err)
	}

	csrBlock, _ := pem.Decode(csrPEM)
	if csrBlock == nil {
		return nil, fmt.Errorf("Failed to find CSR block in %s", csrFile)
	}

	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse certificate request: %s", err)
	}

	err = csr.CheckSignature()
	if err != nil {
		return nil, fmt.Errorf("Failed to check certificate request signature: %s", err)
	}
	return csr, nil
}
