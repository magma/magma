/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package testlib

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"golang.org/x/crypto/pkcs12"
	"io/ioutil"
)

var certsDir = flag.String("certs-dir", "", "Location of the certs")

func GetTLSConfig() (*tls.Config, error) {
	rootCAFilename := fmt.Sprintf("%s/rootCA.pem", *certsDir)
	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}

	pemBytes, err := ioutil.ReadFile(rootCAFilename)
	if err != nil {
		return nil, err
	}

	ok := rootCAs.AppendCertsFromPEM(pemBytes)
	if !ok {
		return nil, fmt.Errorf("failed appending certs from pem")
	}
	// Read client cert
	adminCertPFXFilename := fmt.Sprintf("%s/admin_operator.pfx", *certsDir)
	pfxBytes, err := ioutil.ReadFile(adminCertPFXFilename)
	if err != nil {
		return nil, err
	}

	key, cert, err := pkcs12.Decode(pfxBytes, "")
	if err != nil {
		return nil, err
	}

	// Get REST API client
	tlsConfig := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{{Certificate: [][]byte{cert.Raw}, PrivateKey: key}},
	}
	tlsConfig.BuildNameToCertificate()
	return tlsConfig, nil
}
