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

package tests

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/obsidian/access"
	"magma/orc8r/cloud/go/services/accessd"
	"magma/orc8r/cloud/go/services/accessd/protos"
	accessd_test_service "magma/orc8r/cloud/go/services/accessd/test_init"
	"magma/orc8r/cloud/go/services/certifier"
	certifier_test_service "magma/orc8r/cloud/go/services/certifier/test_init"
	security_cert "magma/orc8r/lib/go/security/cert"
	certifier_test_utils "magma/orc8r/lib/go/security/csr"
	"magma/orc8r/lib/go/util"
)

const (
	TEST_NETWORK_ID        = "N12345"
	TEST_OPERATOR_ID       = "bob"
	WRITE_TEST_NETWORK_ID  = "N6789"
	TEST_SUPER_OPERATOR_ID = "admin"
)

func testGet(t *testing.T, url string) {
	t.Logf("Testing URL: %s", url)
	s, _, err := util.SendHttpRequest("GET", url, "")
	assert.NoError(t, err)
	assert.Equal(t, s, 200)
}

func WaitForTestServer(t *testing.T, e *echo.Echo) net.Listener {
	assert.NotNil(t, e)
	if e != nil {
		for i := 10; i < 100; i++ {
			time.Sleep(time.Millisecond * time.Duration(i))
			if e.Listener != nil {
				break
			}
		}
		assert.NotNil(t, e.Listener)
		return e.Listener
	}
	return nil
}

func SendRequest(method, url, certSn string) (int, error) {
	var body io.Reader = nil
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(access.CLIENT_CERT_SN_KEY, certSn)

	var client = &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
	return response.StatusCode, err
}

// Starts access control services & adds a test admin operator with requested
// adminId
// Returns the Operator's certificate serial number
// NOTE: StartMockAccessControl is intended to be used
// by all REST specific unit tests
func StartMockAccessControl(t *testing.T, adminId string) string {

	// Start services
	certifier_test_service.StartTestService(t)
	accessd_test_service.StartTestService(t)

	// create and sign admin's csr
	csrMsg, err := certifier_test_utils.CreateCSR(
		time.Hour*4, adminId, adminId)
	assert.NoError(t, err)
	certMsg, err := certifier.SignCSR(context.Background(), csrMsg)
	assert.NoError(t, err, "Failed to sign Admin's CSR")
	// get sn from cert
	superCert, err := x509.ParseCertificates(certMsg.CertDer)
	assert.NoError(t, err, "Failed to parse Admin's cert")

	certSerialNum := security_cert.SerialToString(superCert[0].SerialNumber)

	t.Logf("Test Certificate SN: %s", certSerialNum)

	admin := identity.NewOperator(adminId)

	adminACL := []*protos.AccessControl_Entity{
		{Id: identity.NewOperatorWildcard(), Permissions: protos.ACCESS_CONTROL_ALL_PERMISSIONS},
		{Id: identity.NewNetworkWildcard(), Permissions: protos.ACCESS_CONTROL_ALL_PERMISSIONS},
		{Id: identity.NewGatewayWildcard(), Permissions: protos.ACCESS_CONTROL_ALL_PERMISSIONS},
	}

	// Verify that all wildcards are valid
	for idx, ent := range adminACL {
		if ent == nil || ent.Id == nil {
			t.Errorf("NIL Entity @ IDX: %d", idx)
		} else if ent.Id.ToCommonName() == nil {
			t.Errorf("NIL CN %s @ IDX: %d & ent: %+v (%q)",
				*ent.Id.ToCommonName(), idx, ent.Id, ent.Id)
		}
	}
	// Add ACL for supervisor
	assert.NoError(t, accessd.SetOperator(context.Background(), admin, adminACL))
	return certSerialNum
}

// MockAccessControl starts access control related services,
// creates a supervisor account as well as a test Operator
// It's intendend to be used by Access Control related unit tests.
func MockAccessControl(t *testing.T) (certSn string, superCertSn string) {

	// Start services & setup supervisor
	superCertSn = StartMockAccessControl(t, TEST_SUPER_OPERATOR_ID)

	// Get ACL for supervisor & verify, it's not corrupt
	super := identity.NewOperator(TEST_SUPER_OPERATOR_ID)
	testSuperAcl, err := accessd.GetOperatorACL(context.Background(), super)
	for hash, ent := range testSuperAcl {
		if ent == nil || ent.Id == nil {
			t.Errorf("NIL Entity @ KEY: %s", hash)
		} else if ent.Id.ToCommonName() == nil {
			t.Errorf("NIL CN @ KEY: %s & ent: %+v (%q)", hash, ent.Id, ent.Id)
		}
	}
	assert.NoError(t, err)

	// create and sign operator's csr
	csrMsg, err := certifier_test_utils.CreateCSR(
		time.Hour*12, TEST_OPERATOR_ID, TEST_OPERATOR_ID)
	assert.NoError(t, err)
	certMsg, err := certifier.SignCSR(context.Background(), csrMsg)
	assert.NoError(t, err, "Failed to sign CSR")
	// get sn from cert
	cert, err := x509.ParseCertificates(certMsg.CertDer)
	assert.NoError(t, err, "Failed to parse cert")
	certSn = security_cert.SerialToString(cert[0].SerialNumber)

	oper := identity.NewOperator(TEST_OPERATOR_ID)
	readNet := identity.NewNetwork(TEST_NETWORK_ID)
	writeNet := identity.NewNetwork(WRITE_TEST_NETWORK_ID)

	acl := []*protos.AccessControl_Entity{
		{Id: readNet, Permissions: protos.AccessControl_READ},
		{Id: writeNet, Permissions: protos.AccessControl_WRITE},
	}
	// Add ACL for operator
	assert.NoError(t, accessd.SetOperator(context.Background(), oper, acl))

	t.Logf(
		"Test Operator Cert SN: %s, Test Supervisor Cert SN: %s",
		certSn, superCertSn)
	return // return (certSn, superCertSn)
}

func SendRequestWithToken(method, url, username string, token string) (int, error) {
	var body io.Reader = nil
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, err
	}
	request.Header.Set("Content-Type", "application/json")
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + token))
	header := fmt.Sprintf("Basic %s", auth)
	request.Header.Set(echo.HeaderAuthorization, header)

	var client = &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()
	_, err = ioutil.ReadAll(response.Body)
	return response.StatusCode, err
}
