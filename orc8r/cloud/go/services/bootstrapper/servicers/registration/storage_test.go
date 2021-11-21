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
	gatewayPreregisterInfo1 = protos.GatewayPreregisterInfo{
		NetworkId: "networkID1",
		LogicalId: "logicalID1",
	}
	tokenInfo1 = protos.TokenInfo{
		GatewayPreregisterInfo: &gatewayPreregisterInfo1,
		Nonce:                  "someNonce",
		Timeout:                nil,
	}

	tokenInfo2 = protos.TokenInfo{
		GatewayPreregisterInfo: &gatewayPreregisterInfo1,
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

	// assert returns right value when nonce is in store
	isUnique, err := s.IsNonceUnique(tokenInfo1.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, true, isUnique)

	tokenInfo, err := s.GetTokenInfoFromNonce(tokenInfo1.Nonce)
	assert.Error(t, errors.NotFound("Not found"), err)
	assert.Nil(t, tokenInfo)

	tokenInfo, err = s.GetTokenInfoFromLogicalID(gatewayPreregisterInfo1.NetworkId, gatewayPreregisterInfo1.LogicalId)
	assert.Error(t, errors.NotFound("Not found"), err)
	assert.Nil(t, tokenInfo)

	// set tokenInfo1 and test
	err = s.SetTokenInfo("", tokenInfo1)
	assert.NoError(t, err)

	tokenInfo, err = s.GetTokenInfoFromNonce(tokenInfo1.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo1, *tokenInfo)

	tokenInfo, err = s.GetTokenInfoFromLogicalID(gatewayPreregisterInfo1.NetworkId, gatewayPreregisterInfo1.LogicalId)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo1, *tokenInfo)

	isUnique, err = s.IsNonceUnique(tokenInfo1.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, false, isUnique)

	// replace old nonce with new nonce value and see that old nonce was removed
	err = s.SetTokenInfo(tokenInfo1.Nonce, tokenInfo2)
	assert.NoError(t, err)

	tokenInfo, err = s.GetTokenInfoFromNonce(tokenInfo1.Nonce)
	assert.Error(t, errors.NotFound("Not found"), err)
	assert.Nil(t, tokenInfo)

	tokenInfo, err = s.GetTokenInfoFromNonce(tokenInfo2.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo2, *tokenInfo)

	tokenInfo, err = s.GetTokenInfoFromLogicalID(gatewayPreregisterInfo1.NetworkId, gatewayPreregisterInfo1.LogicalId)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo2, *tokenInfo)

	isUnique, err = s.IsNonceUnique(tokenInfo1.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, true, isUnique)

	isUnique, err = s.IsNonceUnique(tokenInfo2.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, false, isUnique)
}
