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
	"net/http"
	"testing"

	"magma/orc8r/cloud/go/obsidian/access"
	"magma/orc8r/lib/go/protos"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const RegisterLteNetworkV1 = "/magma/v1/lte"
const ManageLteNetworkV1 = RegisterLteNetworkV1 + "/:network_id"
const RegisterNetworkV1 = "/magma/v1/networks"
const ManageNetworkV1 = RegisterNetworkV1 + "/:network_id"

func TestIdentityFinder(t *testing.T) {
	e := startTestIdentityServer(t)
	listener := WaitForTestServer(t, e)

	if listener != nil {
		urlPrefix := "http://" + listener.Addr().String()

		// Test V1 network entity
		testGet(t, urlPrefix+RegisterNetworkV1+"/"+TEST_NETWORK_ID)
		// Test V1 LTE network entity
		testGet(t, urlPrefix+RegisterLteNetworkV1+"/"+TEST_NETWORK_ID)
		// Test operator entity
		testGet(t, urlPrefix+"/magma/operators/"+TEST_OPERATOR_ID)
		// Test supervisor wildcards (non magma URL)
		testGet(t, urlPrefix+"/malformed/url")
		// Test supervisor wildcards (magma URL)
		testGet(t, urlPrefix+"/magma/malformed/url")
	}
}

// testSupervisorWildcards verifies that the ctx 'resolves' to supervisor ents
// list and the list itself includes all known wildcards
func testSupervisorWildcards(t *testing.T, c echo.Context) {
	assert.NotNil(t, c)
	ents := access.FindRequestedIdentities(c)
	assert.Len(t, ents, len(protos.Identity_Wildcard_Type_value))

	testMap := map[int32]string{} // copy of Wildcards map
	for key, nm := range protos.Identity_Wildcard_Type_name {
		testMap[key] = nm
	}
	t.Logf("All Wildcard Types: %v", protos.Identity_Wildcard_Type_value)
	for i, ent := range ents {
		wildcard := ent.GetWildcard()
		assert.NotNil(t, wildcard, "Invalid Entity at %d position", i)
		delete(testMap, int32(wildcard.Type))
	}
	assert.Len(t,
		testMap,
		0,
		"Supervisor Wildcards are missing the following types: %q", testMap)
}

func startTestIdentityServer(t *testing.T) *echo.Echo {
	e := echo.New()

	assert.NotNil(t, e)

	// V1 Endpoint requiring a specific Network Entity Access Permissions
	e.GET(ManageNetworkV1, func(c echo.Context) error {
		assert.NotNil(t, c)
		ents := access.FindRequestedIdentities(c)
		assert.Len(t, ents, 1)
		networkIdentity, ok := ents[0].Value.(*protos.Identity_Network)
		assert.True(t, ok)
		assert.Equal(t, networkIdentity.Network, TEST_NETWORK_ID)
		return c.String(http.StatusOK, "All good!")
	})

	// V1 LTE Endpoint requiring a specific Network Entity Access Permissions
	e.GET(ManageLteNetworkV1, func(c echo.Context) error {
		assert.NotNil(t, c)
		ents := access.FindRequestedIdentities(c)
		assert.Len(t, ents, 1)
		networkIdentity, ok := ents[0].Value.(*protos.Identity_Network)
		assert.True(t, ok)
		assert.Equal(t, networkIdentity.Network, TEST_NETWORK_ID)
		return c.String(http.StatusOK, "All good!")
	})

	// Endpoint requiring specific Network Entity Access Permissions
	e.GET("magma/operators/:operator_id", func(c echo.Context) error {
		assert.NotNil(t, c)
		ents := access.FindRequestedIdentities(c)
		assert.Len(t, ents, 1)
		operatorIdentity, ok := ents[0].Value.(*protos.Identity_Operator)
		assert.True(t, ok)
		assert.Equal(t, operatorIdentity.Operator, TEST_OPERATOR_ID)
		return c.String(http.StatusOK, "All good!")
	})

	// Endpoint requiring supervisor permissions
	e.GET("/malformed/url", func(c echo.Context) error {
		testSupervisorWildcards(t, c)
		return c.String(http.StatusOK, "All good!")
	})

	// Endpoint requiring supervisor permissions
	e.GET("/magma/malformed/url", func(c echo.Context) error {
		testSupervisorWildcards(t, c)
		return c.String(http.StatusOK, "All good!")
	})

	go func(t *testing.T) {
		assert.NoError(t, e.Start(""))
	}(t)

	return e
}
