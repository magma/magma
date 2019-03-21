/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import "magma/feg/cloud/go/protos"

func (suite *EpsAuthTestSuite) TestPurgeUE_UnknownSubscriber() {
	purge := &protos.PurgeUERequest{UserName: "sub_unknown"}
	answer, err := suite.PurgeUE(purge)
	suite.EqualError(err, "rpc error: code = NotFound desc = Error fetching subscriber: IMSIsub_unknown, No record for query")
	suite.Equal(protos.ErrorCode_USER_UNKNOWN, answer.ErrorCode)
}

func (suite *EpsAuthTestSuite) TestPurgeUE_Success() {
	purge := &protos.PurgeUERequest{UserName: "sub1"}
	answer, err := suite.PurgeUE(purge)
	suite.NoError(err)
	suite.Equal(protos.ErrorCode_SUCCESS, answer.ErrorCode)
}
