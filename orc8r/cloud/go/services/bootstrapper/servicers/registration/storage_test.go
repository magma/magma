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
	"time"

	"github.com/go-openapi/errors"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/services/bootstrapper/servicers/registration"
	"magma/orc8r/cloud/go/test_utils"
	"magma/orc8r/lib/go/protos"
)

func TestBlobstoreStore(t *testing.T) {
	var (
		gatewayPreregisterInfo = &protos.GatewayDeviceInfo{
			NetworkId: "networkID1",
			LogicalId: "logicalID1",
		}

		tokenInfo1 = &protos.TokenInfo{
			GatewayDeviceInfo: gatewayPreregisterInfo,
			Nonce:             "someNonce1",
			Timeout: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
				Nanos:   int32(time.Now().Nanosecond()),
			},
		}

		tokenInfo2 = &protos.TokenInfo{
			GatewayDeviceInfo: gatewayPreregisterInfo,
			Nonce:             "someNonce2",
			Timeout: &timestamp.Timestamp{
				Seconds: time.Now().Unix(),
				Nanos:   int32(time.Now().Nanosecond()),
			}}
	)

	factory := test_utils.NewSQLBlobstore(t, bootstrapper.BlobstoreTableName)
	s := registration.NewBlobstoreStore(factory)

	// Asserts store works as expected when nonce is not saved in store

	tokenInfo, err := s.GetTokenInfoFromNonce(tokenInfo1.Nonce)
	assert.Error(t, errors.NotFound("Not found"), err)
	assert.Nil(t, tokenInfo)

	tokenInfo, err = s.GetTokenInfoFromLogicalID(gatewayPreregisterInfo.NetworkId, gatewayPreregisterInfo.LogicalId)
	assert.Error(t, errors.NotFound("Not found"), err)
	assert.Nil(t, tokenInfo)

	// Setup tokenInfo1 and test that store works as expected

	err = s.SetTokenInfo(tokenInfo1)
	assert.NoError(t, err)

	tokenInfo, err = s.GetTokenInfoFromNonce(tokenInfo1.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo1, tokenInfo)

	tokenInfo, err = s.GetTokenInfoFromLogicalID(gatewayPreregisterInfo.NetworkId, gatewayPreregisterInfo.LogicalId)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo1, tokenInfo)

	// Set new token info and see that store works as expected

	err = s.SetTokenInfo(tokenInfo2)
	assert.NoError(t, err)

	tokenInfo, err = s.GetTokenInfoFromNonce(tokenInfo1.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo1, tokenInfo)

	tokenInfo, err = s.GetTokenInfoFromNonce(tokenInfo2.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo2, tokenInfo)

	tokenInfo, err = s.GetTokenInfoFromLogicalID(gatewayPreregisterInfo.NetworkId, gatewayPreregisterInfo.LogicalId)
	assert.NoError(t, err)
	assert.Equal(t, tokenInfo2, tokenInfo)

	// Try setting with old nonce and see that store errors as expected

	err = s.SetTokenInfo(tokenInfo1)
	assert.Error(t, errors.NotFound("token is not unique"), err)
}
