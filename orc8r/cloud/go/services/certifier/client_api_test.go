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

package certifier_test

import (
	"crypto/x509"
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"magma/orc8r/cloud/go/services/certifier"
	"magma/orc8r/cloud/go/services/certifier/servicers"
	"magma/orc8r/cloud/go/services/certifier/test_init"
	"magma/orc8r/lib/go/protos"
	security_cert "magma/orc8r/lib/go/security/cert"
	certifier_test_utils "magma/orc8r/lib/go/security/csr"
)

func TestCertifier(t *testing.T) {
	test_init.StartTestService(t)

	// create and sign csr
	csrMsg, err := certifier_test_utils.CreateCSR(time.Hour*24*365, "cn", "cn")
	assert.NoError(t, err)
	certMsg, err := certifier.SignCSR(csrMsg)
	assert.NoError(t, err, "Failed to sign CSR")

	firstCertDer := certMsg.CertDer

	// get sn from cert
	cert, err := x509.ParseCertificates(certMsg.CertDer)
	assert.NoError(t, err, "Failed to parse cert")
	firstCertSN := cert[0].SerialNumber
	snMsg := &protos.Certificate_SN{
		Sn: security_cert.SerialToString(firstCertSN),
	}

	// test get identity
	certInfoMsg, err := certifier.GetIdentity(snMsg)
	assert.NoError(t, err, "Error getting identity")
	fmt.Printf("%+v\n", certInfoMsg)
	assert.True(t, proto.Equal(certInfoMsg.Id, csrMsg.Id))

	// test revoke cert
	err = certifier.RevokeCertificate(snMsg)
	assert.NoError(t, err, "Failed to revoke cert")
	_, err = certifier.GetIdentity(snMsg)
	assert.Error(t, err, "Error: no error getting revoked identity")

	// test collect garbage
	servicers.CollectGarbageAfter = time.Duration(0)

	csrMsg, err = certifier_test_utils.CreateCSR(time.Duration(0), "cn", "cn")
	assert.NoError(t, err)
	certMsg, err = certifier.SignCSR(csrMsg)
	assert.NoError(t, err, "Failed to sign CSR")
	cert, err = x509.ParseCertificates(certMsg.CertDer)
	assert.NoError(t, err, "Failed to parse cert")
	snMsg = &protos.Certificate_SN{
		Sn: security_cert.SerialToString(cert[0].SerialNumber),
	}

	err = certifier.CollectGarbage()
	assert.NoError(t, err, "Failed to collect garbage")
	_, err = certifier.GetIdentity(snMsg)
	assert.Equal(t, grpc.Code(err), codes.NotFound)

	oper := protos.NewOperatorIdentity("testOperator")
	assert.NoError(t,
		certifier.AddCertificate(oper, firstCertDer),
		"Failed to Add Existing Cert")

	certInfoMsg, err = certifier.GetCertificateIdentity(security_cert.SerialToString(firstCertSN))
	assert.NoError(t, err, "Error getting added cert identity")
	if err == nil {
		assert.Equal(t, oper.HashString(), certInfoMsg.Id.HashString())
	}

	sns, err := certifier.ListCertificates()
	assert.NoError(t, err, "Error Listing Certificates")
	assert.Equal(t, 1, len(sns))

	csrMsg, err = certifier_test_utils.CreateCSR(time.Hour*2, "cn1", "cn1")
	assert.NoError(t, err)
	_, err = certifier.SignCSR(csrMsg)
	assert.NoError(t, err, "Failed to sign CSR")

	sns, err = certifier.ListCertificates()
	assert.NoError(t, err, "Error Listing Certificates")
	assert.Equal(t, 2, len(sns))

	operSNs, err := certifier.FindCertificates(oper)
	assert.NoError(t, err, "Error Finding Operator Certificates")
	assert.Equal(t, 1, len(operSNs))
	if len(operSNs) > 0 {
		assert.Equal(t, security_cert.SerialToString(firstCertSN), operSNs[0])
	}
}
