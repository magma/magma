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

import (
	"github.com/magma/milenage"
	"golang.org/x/net/context"

	fegprotos "magma/feg/cloud/go/protos"
)

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_NilRequest() {
	_, err := suite.AuthenticationInformation(nil)
	suite.EqualError(err, "rpc error: code = InvalidArgument desc = received a nil AuthenticationInformationRequest")
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_EmptyUserName() {
	air := &fegprotos.AuthenticationInformationRequest{
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	_, err := suite.AuthenticationInformation(air)
	suite.EqualError(err, "rpc error: code = InvalidArgument desc = user name was empty")
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_EmptyPlmm() {
	air := &fegprotos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		NumRequestedEutranVectors: 1,
	}

	_, err := suite.AuthenticationInformation(air)
	suite.EqualError(
		err, "rpc error: code = InvalidArgument desc = expected Visited PLMN to be 3 bytes, but got 0 bytes")
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_0RequestedVectors() {
	air := &fegprotos.AuthenticationInformationRequest{
		UserName:    "sub1",
		VisitedPlmn: []byte{0, 0, 0},
	}

	_, err := suite.AuthenticationInformation(air)
	suite.EqualError(err, "rpc error: code = InvalidArgument desc = 0 E-UTRAN vectors were requested")
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_UnknownGateway() {
	air := &fegprotos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	_, err := suite.Server.AuthenticationInformation(context.Background(), air)
	suite.EqualError(err, "rpc error: code = PermissionDenied desc = Missing Gateway Identity")
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_UnknownSubscriber() {
	air := &fegprotos.AuthenticationInformationRequest{
		UserName:                  "sub_unknown",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.EqualError(
		err,
		"rpc error: code = NotFound desc = error loading subscriber ent for network ID: test, SID: sub_unknown: Not found")
	suite.checkAIA(aia, fegprotos.ErrorCode_USER_UNKNOWN, 0)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_MissingAuthKey() {
	air := &fegprotos.AuthenticationInformationRequest{
		UserName:                  "missing_auth_key",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.EqualError(
		err,
		"rpc error: code = Unauthenticated desc = Authentication rejected: incorrect key size. Expected 16 bytes, but got 0 bytes")
	suite.checkAIA(aia, fegprotos.ErrorCode_AUTHORIZATION_REJECTED, 0)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_MissingSubscriberState() {
	air := &fegprotos.AuthenticationInformationRequest{
		UserName:                  "empty_sub",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.EqualError(
		err, "rpc error: code = Unauthenticated desc = Authentication rejected: Subscriber data missing LTE subscription")
	suite.checkAIA(aia, fegprotos.ErrorCode_AUTHORIZATION_REJECTED, 0)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_Success() {
	air := &fegprotos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 3,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.NoError(err)
	suite.checkAIA(aia, fegprotos.ErrorCode_SUCCESS, 3)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_InvalidResyncInfo() {
	resyncInfo := make([]byte, 10)
	resyncInfo[5] = 1
	air := &fegprotos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 3,
		ResyncInfo:                resyncInfo,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.EqualError(
		err,
		"rpc error: code = Unauthenticated desc = Authentication rejected: resync info incorrect length. expected 30 bytes, but got 10 bytes")
	suite.checkAIA(aia, fegprotos.ErrorCode_AUTHORIZATION_REJECTED, 0)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_InvalidResyncMacS() {
	resyncInfo := make([]byte, 30)
	macS := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	copy(resyncInfo[22:], macS)
	air := &fegprotos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 3,
		ResyncInfo:                resyncInfo,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.EqualError(
		err,
		"rpc error: code = Unauthenticated desc = Authentication rejected: Invalid resync authentication code")
	suite.checkAIA(aia, fegprotos.ErrorCode_AUTHORIZATION_REJECTED, 0)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_ResyncSuccess() {
	resyncInfo := make([]byte, 30)
	macS := []byte{47, 223, 5, 242, 77, 209, 76, 218}
	copy(resyncInfo[22:], macS)
	air := &fegprotos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 3,
		ResyncInfo:                resyncInfo,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.NoError(err)
	suite.checkAIA(aia, fegprotos.ErrorCode_SUCCESS, 3)
}

func (suite *EpsAuthTestSuite) checkAIA(
	aia *fegprotos.AuthenticationInformationAnswer, errorCode fegprotos.ErrorCode, numVectors int) {

	suite.Equal(errorCode, aia.ErrorCode)
	suite.Equal(numVectors, len(aia.EutranVectors))
	for _, vector := range aia.EutranVectors {
		suite.Equal(milenage.RandChallengeBytes, len(vector.Rand))
		suite.Equal(milenage.XresBytes, len(vector.Xres))
		suite.Equal(milenage.AutnBytes, len(vector.Autn))
		suite.Equal(milenage.KasmeBytes, len(vector.Kasme))
	}
}
