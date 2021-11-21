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

package registration_test

import (
	"testing"

	"github.com/go-openapi/errors"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers/registration"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/protos"
)

var (
	unusedNonce = "unusedNonce"

	gatewayPreregisterInfo = protos.GatewayPreregisterInfo{
		NetworkId: "networkID1",
		LogicalId: "logicalID1",
	}

	tokenInfo1 = protos.TokenInfo{
		GatewayPreregisterInfo: &gatewayPreregisterInfo,
		Nonce:                  "someNonce",
		Timeout:                nil,
	}

	tokenInfo2 = protos.TokenInfo{
		GatewayPreregisterInfo: &gatewayPreregisterInfo,
		Nonce:                  "someNonce2",
		Timeout:                nil,
	}
)

func TestBlobstoreStore(t *testing.T) {
	db, err := sqorc.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	factory := blobstore.NewSQLStoreFactory(bootstrapper.DBTableName, db, sqorc.GetSqlBuilder())
	assert.NoError(t, factory.InitializeFactory())
	s := registration.NewBlobstoreStore(factory)

	// Asserts store works as expected when nonce is not saved in store

	isUnique, err := s.IsNonceUnique(tokenInfo1.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, true, isUnique)

	tokenInfo, err := s.GetTokenInfoFromNonce(tokenInfo1.Nonce)
	assert.Error(t, errors.NotFound("Not found"), err)
	assert.Nil(t, tokenInfo)

	tokenInfo, err = s.GetTokenInfoFromLogicalID(gatewayPreregisterInfo.NetworkId, gatewayPreregisterInfo.LogicalId)
	assert.Error(t, errors.NotFound("Not found"), err)
	assert.Nil(t, tokenInfo)

	// Setup tokenInfo1 and test that store works as expected

	err = s.SetTokenInfo("", tokenInfo1)
	assert.NoError(t, err)

	// try SetTokenInfo with an oldNonce that isn't in the store
	err = s.SetTokenInfo(unusedNonce, tokenInfo1)
	assert.NoError(t, err)

	tokenInfo, err = s.GetTokenInfoFromNonce(tokenInfo1.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo1, *tokenInfo)

	tokenInfo, err = s.GetTokenInfoFromLogicalID(gatewayPreregisterInfo.NetworkId, gatewayPreregisterInfo.LogicalId)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo1, *tokenInfo)

	isUnique, err = s.IsNonceUnique(tokenInfo1.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, false, isUnique)

	// Replace old nonce with new nonce value and see that store works as expected

	err = s.SetTokenInfo(tokenInfo1.Nonce, tokenInfo2)
	assert.NoError(t, err)

	tokenInfo, err = s.GetTokenInfoFromNonce(tokenInfo1.Nonce)
	assert.Error(t, errors.NotFound("Not found"), err)
	assert.Nil(t, tokenInfo)
	isUnique, err = s.IsNonceUnique(tokenInfo1.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, true, isUnique)

	tokenInfo, err = s.GetTokenInfoFromNonce(tokenInfo2.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo2, *tokenInfo)
	isUnique, err = s.IsNonceUnique(tokenInfo2.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, false, isUnique)

	tokenInfo, err = s.GetTokenInfoFromLogicalID(gatewayPreregisterInfo.NetworkId, gatewayPreregisterInfo.LogicalId)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo2, *tokenInfo)
}
