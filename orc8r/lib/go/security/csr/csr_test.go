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

package csr_test

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"os"
	"testing"

	"magma/orc8r/lib/go/security/csr"
	"magma/orc8r/lib/go/security/key"

	"github.com/stretchr/testify/assert"
)

func getDefaultCertificateRequestTemplate() (template *x509.CertificateRequest) {
	template = &x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{"MagmaClient"},
			OrganizationalUnit: []string{"magma_client.foo-bar.com"},
			CommonName:         "",
		},
	}
	return
}

func TestWriteAndReadCSR(t *testing.T) {
	csrFile := "test_csr.pem"

	priv, err := key.GenerateKey("", 2048)
	assert.NoError(t, err, "Failed to generate private key")

	template := getDefaultCertificateRequestTemplate()
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, template, priv)
	assert.NoError(t, err, "Failed to create certificate request")

	err = csr.WriteCSR(csrDER, csrFile)
	assert.NoError(t, err, "Failed to write csr to file")

	trueCSR, err := x509.ParseCertificateRequest(csrDER)
	assert.NoError(t, err, "Failed to parse certificate request")

	retrievedCSR, err := csr.ReadCSR(csrFile)
	assert.NoError(t, err, "Failed to read CSR from file")

	if !bytes.Equal(trueCSR.Raw, retrievedCSR.Raw) {
		t.Fatalf("CSR read from file does not match original CSR")
	}
	os.Remove(csrFile)
}
