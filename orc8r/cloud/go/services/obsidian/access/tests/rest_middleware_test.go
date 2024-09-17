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
	"fmt"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	certprotos "magma/orc8r/cloud/go/services/certifier/protos"
	certifier_test_service "magma/orc8r/cloud/go/services/certifier/test_init"
	"magma/orc8r/cloud/go/services/certifier/test_utils"
	"magma/orc8r/cloud/go/services/obsidian/access"
	"magma/orc8r/cloud/go/services/tenants"
	tenantsh "magma/orc8r/cloud/go/services/tenants/obsidian/handlers"
	tenant_protos "magma/orc8r/cloud/go/services/tenants/protos"
	tenants_test_init "magma/orc8r/cloud/go/services/tenants/test_init"
)

func TestMiddlewareWithoutCertifier(t *testing.T) {
	e := startTestMiddlewareServer(t)
	e.Use(access.CertificateMiddleware) // inject obsidian access control middleware
	listener := WaitForTestServer(t, e)

	if listener == nil {
		return // WaitForTestServer should have 'logged' error already
	}

	urlPrefix := "http://" + listener.Addr().String()

	// Test if we set httpCode to be 503 when certifier is down
	s, err := SendRequest(
		"GET", // READ
		urlPrefix+RegisterNetworkV1+"/"+TEST_NETWORK_ID,
		"test cert string",
	)
	assert.NoError(t, err)
	assert.Equal(t, 503, s)
}

func TestAuthMiddleware(t *testing.T) {
	// Set up auth middleware by creating root user, non-admin user bob, and their respective policies
	certifier_test_service.StartTestService(t)
	tenants_test_init.StartTestService(t)

	store := test_utils.GetCertifierBlobstore(t)

	rootToken := test_utils.CreateTestUser(t, store, test_utils.TestRootUsername, test_utils.TestPassword, []*certprotos.Policy{
		{
			Effect:   certprotos.Effect_ALLOW,
			Action:   certprotos.Action_WRITE,
			Resource: &certprotos.Policy_Path{Path: &certprotos.PathResource{Path: "**"}},
		},
	})

	userToken := test_utils.CreateTestUser(t, store, test_utils.TestUsername, test_utils.TestPassword, []*certprotos.Policy{
		{
			Effect:   certprotos.Effect_ALLOW,
			Action:   certprotos.Action_READ,
			Resource: &certprotos.Policy_Path{Path: &certprotos.PathResource{Path: "**"}},
		},
		{
			Effect:   certprotos.Effect_DENY,
			Action:   certprotos.Action_WRITE,
			Resource: &certprotos.Policy_Network{Network: &certprotos.NetworkResource{Networks: []string{test_utils.WriteTestNetworkId}}},
		},
	})

	queryUserToken := test_utils.CreateTestUser(t, store, test_utils.TestQueryUsername, test_utils.TestPassword, []*certprotos.Policy{
		{
			Effect:   certprotos.Effect_ALLOW,
			Action:   certprotos.Action_READ,
			Resource: &certprotos.Policy_Network{Network: &certprotos.NetworkResource{Networks: []string{"foo"}}},
		},
		{
			Effect:   certprotos.Effect_DENY,
			Action:   certprotos.Action_WRITE,
			Resource: &certprotos.Policy_Network{Network: &certprotos.NetworkResource{Networks: []string{test_utils.WriteTestNetworkId}}},
		},
	})

	tenants.CreateTenant(context.Background(), test_utils.TestTenantId, &tenant_protos.Tenant{
		Name:     fmt.Sprint(test_utils.TestTenantId),
		Networks: []string{test_utils.TestTenantNetworkId},
	})
	tenantUserToken := test_utils.CreateTestUser(t, store, test_utils.TestTenantUsername, test_utils.TestPassword, []*certprotos.Policy{
		{
			Effect:   certprotos.Effect_ALLOW,
			Action:   certprotos.Action_WRITE,
			Resource: &certprotos.Policy_Tenant{Tenant: &certprotos.TenantResource{Tenants: []int64{test_utils.TestTenantId}}},
		},
	})

	e := startTestMiddlewareServer(t)
	e.Use(access.TokenMiddleware)
	listener := WaitForTestServer(t, e)
	if listener == nil {
		return
	}

	urlPrefix := fmt.Sprintf("http://%s", listener.Addr().String())
	tests := []struct {
		method   string
		url      string
		user     string
		token    string
		expected int
	}{
		// Test admin user
		{"GET", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), test_utils.TestRootUsername, rootToken, http.StatusOK},
		{"POST", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), test_utils.TestRootUsername, rootToken, http.StatusOK},
		{"GET", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, TEST_NETWORK_ID), test_utils.TestRootUsername, rootToken, http.StatusOK},
		{"PUT", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, TEST_NETWORK_ID), test_utils.TestRootUsername, rootToken, http.StatusOK},
		{"GET", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, WRITE_TEST_NETWORK_ID), test_utils.TestRootUsername, rootToken, http.StatusOK},
		{"PUT", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, WRITE_TEST_NETWORK_ID), test_utils.TestRootUsername, rootToken, http.StatusOK},
		{"GET", fmt.Sprintf("%s%s/%d", urlPrefix, TenantRootPathV1, test_utils.TestTenantId), test_utils.TestRootUsername, rootToken, http.StatusOK},
		{"POST", fmt.Sprintf("%s%s/%d", urlPrefix, TenantRootPathV1, test_utils.TestTenantId), test_utils.TestRootUsername, rootToken, http.StatusOK},
		{"GET", fmt.Sprintf("%s%s", urlPrefix, "/malformed/url"), test_utils.TestRootUsername, rootToken, http.StatusOK},
		{"PUT", fmt.Sprintf("%s%s", urlPrefix, "/malformed/url"), test_utils.TestRootUsername, rootToken, http.StatusOK},

		// Test non-admin user
		// User has read access to all URI endpoints
		{"GET", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), test_utils.TestUsername, userToken, http.StatusOK},
		{"POST", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), test_utils.TestUsername, userToken, http.StatusForbidden},
		// User has read access to all networks and is denied read/write access to WRITE_TEST_NETWORK_ID
		{"GET", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, TEST_NETWORK_ID), test_utils.TestUsername, userToken, http.StatusOK},
		{"PUT", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, TEST_NETWORK_ID), test_utils.TestUsername, userToken, http.StatusForbidden},
		{"GET", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, WRITE_TEST_NETWORK_ID), test_utils.TestUsername, userToken, http.StatusForbidden},
		{"PUT", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, WRITE_TEST_NETWORK_ID), test_utils.TestUsername, userToken, http.StatusForbidden},

		// Test query parameters, ensuring network_id is pulled only from path params
		// User is allowed to read random network foo
		{"GET", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterLteNetworkV1, "foo"), test_utils.TestQueryUsername, queryUserToken, http.StatusOK},
		// User is not allowed to read network foo when it is included in query param
		{"GET", fmt.Sprintf("%s%s?network_id=%s", urlPrefix, RegisterLteNetworkV1, "foo"), test_utils.TestQueryUsername, queryUserToken, http.StatusForbidden},
		// User is denied write access to WRITE_TEST_NETWORK_ID when it is included in query param
		{"PUT", fmt.Sprintf("%s%s/%s?network_id=%s", urlPrefix, RegisterLteNetworkV1, "foo", WRITE_TEST_NETWORK_ID), test_utils.TestQueryUsername, queryUserToken, http.StatusForbidden},

		// User does not have write access to any tenants
		{"GET", fmt.Sprintf("%s%s/%d", urlPrefix, TenantRootPathV1, test_utils.TestTenantId), test_utils.TestUsername, userToken, http.StatusOK},
		{"POST", fmt.Sprintf("%s%s/%d", urlPrefix, TenantRootPathV1, test_utils.TestTenantId), test_utils.TestUsername, userToken, http.StatusForbidden},

		// Test non-admin user who only has tenant-based access
		// User can access endpoints that manage tenant 0
		{"GET", fmt.Sprintf("%s%s/%d", urlPrefix, TenantRootPathV1, test_utils.TestTenantId), test_utils.TestTenantUsername, tenantUserToken, http.StatusOK},
		{"POST", fmt.Sprintf("%s%s/%d", urlPrefix, TenantRootPathV1, test_utils.TestTenantId), test_utils.TestTenantUsername, tenantUserToken, http.StatusOK},
		// User can access tenant 0's networks
		{"GET", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, test_utils.TestTenantNetworkId), test_utils.TestTenantUsername, tenantUserToken, http.StatusOK},
		{"PUT", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, test_utils.TestTenantNetworkId), test_utils.TestTenantUsername, tenantUserToken, http.StatusOK},
		// User cannot access any other resources
		{"GET", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), test_utils.TestTenantUsername, tenantUserToken, http.StatusForbidden},
		{"POST", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), test_utils.TestTenantUsername, tenantUserToken, http.StatusForbidden},
		{"GET", fmt.Sprintf("%s%s/%d", urlPrefix, RegisterNetworkV1, test_utils.TestDenyTenantId), test_utils.TestTenantUsername, tenantUserToken, http.StatusForbidden},
		{"PUT", fmt.Sprintf("%s%s/%d", urlPrefix, RegisterNetworkV1, test_utils.TestDenyTenantId), test_utils.TestTenantUsername, tenantUserToken, http.StatusForbidden},
		{"GET", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, test_utils.TestDenyTenantNetworkId), test_utils.TestTenantUsername, tenantUserToken, http.StatusForbidden},
		{"PUT", fmt.Sprintf("%s%s/%s", urlPrefix, RegisterNetworkV1, test_utils.TestDenyTenantNetworkId), test_utils.TestTenantUsername, tenantUserToken, http.StatusForbidden},
	}
	for _, tt := range tests {
		s, err := SendRequestWithToken(tt.method, tt.url, tt.user, tt.token)
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, s)
	}
}

func TestMiddleware(t *testing.T) {
	operCertSn, superCertSn := MockAccessControl(t)

	e := startTestMiddlewareServer(t)
	e.Use(access.CertificateMiddleware) // inject obsidian access control middleware
	listener := WaitForTestServer(t, e)

	if listener == nil {
		return // WaitForTestServer should have 'logged' error already
	}

	urlPrefix := fmt.Sprintf("http://%s", listener.Addr().String())
	tests := []struct {
		method   string
		url      string
		certSn   string
		expected int
	}{
		// Test regular operator wildcard failures
		{"GET", fmt.Sprintf("%s%s%s%s", urlPrefix, RegisterNetworkV1, "/", TEST_NETWORK_ID), operCertSn, 200},
		{"PUT", fmt.Sprintf("%s%s%s%s", urlPrefix, RegisterNetworkV1, "/", TEST_NETWORK_ID), operCertSn, 403},
		{"GET", fmt.Sprintf("%s%s%s%s", urlPrefix, RegisterNetworkV1, "/", WRITE_TEST_NETWORK_ID), operCertSn, 403},
		{"PUT", fmt.Sprintf("%s%s%s%s", urlPrefix, RegisterNetworkV1, "/", WRITE_TEST_NETWORK_ID), operCertSn, 200},
		{"GET", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), operCertSn, 403},
		{"POST", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), operCertSn, 403},
		{"GET", fmt.Sprintf("%s%s", urlPrefix, tenantsh.TenantInfoURL), operCertSn, 403},
		{"POST", fmt.Sprintf("%s%s", urlPrefix, tenantsh.TenantInfoURL), operCertSn, 403},
		// Test Supervisor Permissions
		{"GET", fmt.Sprintf("%s%s%s%s", urlPrefix, RegisterNetworkV1, "/", WRITE_TEST_NETWORK_ID), superCertSn, 200},
		{"PUT", fmt.Sprintf("%s%s%s%s", urlPrefix, RegisterNetworkV1, "/", TEST_NETWORK_ID), superCertSn, 200},
		{"GET", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), superCertSn, 200},
		{"POST", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), superCertSn, 200},
		{"GET", fmt.Sprintf("%s%s", urlPrefix, RegisterNetworkV1), superCertSn, 200},
		{"GET", fmt.Sprintf("%s%s", urlPrefix, "/malformed/url"), superCertSn, 200},
		{"PUT", fmt.Sprintf("%s%s", urlPrefix, "/malformed/url"), superCertSn, 200},
		{"GET", fmt.Sprintf("%s%s", urlPrefix, tenantsh.TenantInfoURL), superCertSn, 200},
		{"POST", fmt.Sprintf("%s%s", urlPrefix, tenantsh.TenantInfoURL), superCertSn, 200},
	}
	for _, tt := range tests {
		s, err := SendRequest(tt.method, tt.url, tt.certSn)
		assert.NoError(t, err)
		assert.Equal(t, tt.expected, s)
	}
}

func startTestMiddlewareServer(t *testing.T) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	assert.NotNil(t, e)

	dummyHandlerFunc := func(c echo.Context) error {
		return c.String(http.StatusOK, "All good!")
	}

	// Endpoint requiring Network Wildcard READ Access Permissions
	e.GET(RegisterNetworkV1, dummyHandlerFunc)

	// Endpoint requiring Network Wildcard WRITE Access Permissions
	e.POST(RegisterNetworkV1, dummyHandlerFunc)

	// Endpoint requiring Network Wildcard READ Access Permissions
	e.GET(RegisterLteNetworkV1, dummyHandlerFunc)

	// Endpoint requiring Network Wildcard WRITE Access Permissions
	e.POST(RegisterLteNetworkV1, dummyHandlerFunc)

	// Endpoint requiring Network Wildcard READ Access Permissions
	e.GET(ManageLteNetworkV1, dummyHandlerFunc)

	// Endpoint requiring Network Wildcard WRITE Access Permissions
	e.POST(ManageLteNetworkV1, dummyHandlerFunc)

	// Endpoint requiring a specific Network READ Entity Access Permissions
	e.GET(ManageNetworkV1, dummyHandlerFunc)

	// Endpoint requiring a specific Network WRITE Entity Access Permissions
	e.PUT(ManageNetworkV1, dummyHandlerFunc)

	// Endpoint requiring supervisor permissions
	e.GET("/malformed/url", dummyHandlerFunc)

	// Endpoint requiring Write supervisor permissions
	e.PUT("/malformed/url", dummyHandlerFunc)

	// Tenants Endpoint requiring Network Wildcard WRITE access permissions
	e.POST(tenantsh.TenantInfoURL, dummyHandlerFunc)

	// Tenants Endpoint requiring Network Wildcard READ access permissions
	e.GET(tenantsh.TenantInfoURL, dummyHandlerFunc)

	go func(t *testing.T) {
		assert.NoError(t, e.Start(""))
	}(t)

	return e
}
