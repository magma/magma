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

// package service implements the core of bootstrapper
package service

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"

	"magma/orc8r/lib/go/security/key"
)

func createCSRAndKey(hwId string) (privateKey interface{}, csr []byte, err error) {
	template := &x509.CertificateRequest{
		Subject: pkix.Name{
			Country:            []string{"US"},
			Organization:       []string{"MagmaGatewayClient"},
			OrganizationalUnit: []string{"go.gateway.magma.fb.com"},
			CommonName:         hwId,
		},
	}
	privateKey, err = key.GenerateKey(CertificateECKeyType, 0)
	if err != nil {
		return
	}
	csr, err = x509.CreateCertificateRequest(rand.Reader, template, privateKey)
	return
}
