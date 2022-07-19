/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package servicers

import "magma/feg/cloud/go/protos"

func (suite *EpsAuthTestSuite) TestPurgeUE_UnknownSubscriber() {
	purge := &protos.PurgeUERequest{UserName: "sub_unknown"}
	answer, err := suite.PurgeUE(purge)
	suite.EqualError(
		err,
		"rpc error: code = NotFound desc = error loading subscriber ent for network ID: test, SID: sub_unknown: Not found")
	suite.Equal(protos.ErrorCode_USER_UNKNOWN, answer.ErrorCode)
}

func (suite *EpsAuthTestSuite) TestPurgeUE_Success() {
	purge := &protos.PurgeUERequest{UserName: "sub1"}
	answer, err := suite.PurgeUE(purge)
	suite.NoError(err)
	suite.Equal(protos.ErrorCode_SUCCESS, answer.ErrorCode)
}
