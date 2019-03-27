/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import "magma/feg/cloud/go/protos"

func (suite *EpsAuthTestSuite) TestUpdateLocation_NilRequest() {
	_, err := suite.UpdateLocation(nil)
	suite.EqualError(err, "rpc error: code = InvalidArgument desc = received a nil UpdateLocationRequest")
}

func (suite *EpsAuthTestSuite) TestUpdateLocation_EmptyUserName() {
	ulr := &protos.UpdateLocationRequest{
		VisitedPlmn: []byte{0, 0, 0},
	}

	_, err := suite.UpdateLocation(ulr)
	suite.EqualError(err, "rpc error: code = InvalidArgument desc = user name was empty")
}

func (suite *EpsAuthTestSuite) TestUpdateLocation_EmptyPlmm() {
	ulr := &protos.UpdateLocationRequest{
		UserName: "sub1",
	}

	_, err := suite.UpdateLocation(ulr)
	suite.EqualError(err, "rpc error: code = InvalidArgument desc = expected Visited PLMN to be 3 bytes, but got 0 bytes")
}

func (suite *EpsAuthTestSuite) TestUpdateLocation_Success() {
	ulr := &protos.UpdateLocationRequest{
		UserName:    "sub1",
		VisitedPlmn: []byte{0, 0, 0},
	}

	ula, err := suite.UpdateLocation(ulr)
	suite.NoError(err)
	suite.checkULA(ula, 7000, 5000)
}

func (suite *EpsAuthTestSuite) TestUpdateLocation_DefaultProfile() {
	ulr := &protos.UpdateLocationRequest{
		UserName:    "empty_sub",
		VisitedPlmn: []byte{0, 0, 0},
	}

	ula, err := suite.UpdateLocation(ulr)
	suite.NoError(err)
	suite.checkULA(ula, 1000, 2000)
}

func (suite *EpsAuthTestSuite) TestUpdateLocation_UnknownSubscriber() {
	ulr := &protos.UpdateLocationRequest{
		UserName:    "sub_unknown",
		VisitedPlmn: []byte{0, 0, 0},
	}

	ula, err := suite.UpdateLocation(ulr)
	suite.EqualError(err, "rpc error: code = NotFound desc = Error fetching subscriber: IMSIsub_unknown, No record for query")
	suite.Equal(protos.ErrorCode_USER_UNKNOWN, ula.ErrorCode)
}

func (suite *EpsAuthTestSuite) checkULA(ula *protos.UpdateLocationAnswer, maxUlBitRate, maxDlBitRate uint32) {
	suite.Equal(protos.ErrorCode_SUCCESS, ula.GetErrorCode())
	suite.Equal(maxDlBitRate, ula.GetTotalAmbr().GetMaxBandwidthDl())
	suite.Equal(maxUlBitRate, ula.GetTotalAmbr().GetMaxBandwidthUl())
	suite.Equal(1, len(ula.Apn))

	apn := ula.Apn[0]
	suite.Equal(maxDlBitRate, apn.GetAmbr().GetMaxBandwidthDl())
	suite.Equal(maxUlBitRate, apn.GetAmbr().GetMaxBandwidthUl())
	suite.Equal("oai.ipv4", apn.GetServiceSelection())
	suite.Equal(protos.UpdateLocationAnswer_APNConfiguration_IPV4, apn.GetPdn())

	qos := apn.GetQosProfile()
	suite.Equal(int32(9), qos.GetClassId())
	suite.Equal(uint32(15), qos.GetPriorityLevel())
	suite.Equal(true, qos.GetPreemptionCapability())
	suite.Equal(false, qos.GetPreemptionVulnerability())
}
