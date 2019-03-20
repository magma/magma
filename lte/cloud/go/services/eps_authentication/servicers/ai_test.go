/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"magma/feg/cloud/go/protos"
	"magma/lte/cloud/go/services/eps_authentication/crypto"

	"golang.org/x/net/context"
)

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_NilRequest() {
	_, err := suite.AuthenticationInformation(nil)
	suite.EqualError(err, "rpc error: code = InvalidArgument desc = received a nil AuthenticationInformationRequest")
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_EmptyUserName() {
	air := &protos.AuthenticationInformationRequest{
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	_, err := suite.AuthenticationInformation(air)
	suite.EqualError(err, "rpc error: code = InvalidArgument desc = user name was empty")
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_EmptyPlmm() {
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		NumRequestedEutranVectors: 1,
	}

	_, err := suite.AuthenticationInformation(air)
	suite.EqualError(err, "rpc error: code = InvalidArgument desc = expected Visited PLMN to be 3 bytes, but got 0 bytes")
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_0RequestedVectors() {
	air := &protos.AuthenticationInformationRequest{
		UserName:    "sub1",
		VisitedPlmn: []byte{0, 0, 0},
	}

	_, err := suite.AuthenticationInformation(air)
	suite.EqualError(err, "rpc error: code = InvalidArgument desc = 0 E-UTRAN vectors were requested")
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_UnknownGateway() {
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	_, err := suite.Server.AuthenticationInformation(context.Background(), air)
	suite.EqualError(err, "rpc error: code = PermissionDenied desc = Missing Gateway Identity")
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_UnknownSubscriber() {
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "sub_unknown",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.EqualError(err, "rpc error: code = Aborted desc = Error fetching subscriber: IMSIsub_unknown, No record for query")
	suite.checkAIA(aia, protos.ErrorCode_USER_UNKNOWN, 0)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_MissingAuthKey() {
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "missing_auth_key",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.EqualError(err, "rpc error: code = Unauthenticated desc = Authentication rejected: incorrect key size. Expected 16 bytes, but got 0 bytes")
	suite.checkAIA(aia, protos.ErrorCode_AUTHORIZATION_REJECTED, 0)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_MissingSubscriberState() {
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "empty_sub",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 1,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.EqualError(err, "rpc error: code = Unavailable desc = Authentication data unavailable: subscriber state is nil")
	suite.checkAIA(aia, protos.ErrorCode_AUTHENTICATION_DATA_UNAVAILABLE, 0)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_Success() {
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 3,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.NoError(err)
	suite.checkAIA(aia, protos.ErrorCode_SUCCESS, 3)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_InvalidResyncInfo() {
	resyncInfo := make([]byte, 10)
	resyncInfo[5] = 1
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 3,
		ResyncInfo:                resyncInfo,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.EqualError(err, "rpc error: code = Unauthenticated desc = Authentication rejected: resync info incorrect length. expected 30 bytes, but got 10 bytes")
	suite.checkAIA(aia, protos.ErrorCode_AUTHORIZATION_REJECTED, 0)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_InvalidResyncMacS() {
	resyncInfo := make([]byte, 30)
	macS := []byte{1, 2, 3, 4, 5, 6, 7, 8}
	copy(resyncInfo[22:], macS)
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 3,
		ResyncInfo:                resyncInfo,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.EqualError(err, "rpc error: code = Unauthenticated desc = Authentication rejected: Invalid resync authentication code")
	suite.checkAIA(aia, protos.ErrorCode_AUTHORIZATION_REJECTED, 0)
}

func (suite *EpsAuthTestSuite) TestAuthenticationInformation_ResyncSuccess() {
	resyncInfo := make([]byte, 30)
	macS := []byte{47, 223, 5, 242, 77, 209, 76, 218}
	copy(resyncInfo[22:], macS)
	air := &protos.AuthenticationInformationRequest{
		UserName:                  "sub1",
		VisitedPlmn:               []byte{0, 0, 0},
		NumRequestedEutranVectors: 3,
		ResyncInfo:                resyncInfo,
	}

	aia, err := suite.AuthenticationInformation(air)
	suite.NoError(err)
	suite.checkAIA(aia, protos.ErrorCode_SUCCESS, 3)
}

func (suite *EpsAuthTestSuite) checkAIA(aia *protos.AuthenticationInformationAnswer, errorCode protos.ErrorCode, numVectors int) {
	suite.Equal(errorCode, aia.ErrorCode)
	suite.Equal(numVectors, len(aia.EutranVectors))
	for _, vector := range aia.EutranVectors {
		suite.Equal(crypto.RandChallengeBytes, len(vector.Rand))
		suite.Equal(crypto.XresBytes, len(vector.Xres))
		suite.Equal(crypto.AutnBytes, len(vector.Autn))
		suite.Equal(crypto.KasmeBytes, len(vector.Kasme))
	}
}
