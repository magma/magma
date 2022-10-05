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
	"magma/feg/cloud/go/protos"
)

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
	suite.EqualError(
		err, "rpc error: code = InvalidArgument desc = expected Visited PLMN to be 3 bytes, but got 0 bytes")
}

func (suite *EpsAuthTestSuite) TestUpdateLocation_Success() {
	ulr := &protos.UpdateLocationRequest{
		UserName:    "sub1",
		VisitedPlmn: []byte{0, 0, 0},
	}

	ula, err := suite.UpdateLocation(ulr)
	suite.NoError(err)
	suite.checkULA(
		ula, 7000, 5000, maxUlBitRateU32, maxDlBitRateU32, "apn", "172.16.254.1", "1.2.3.4")
}

func (suite *EpsAuthTestSuite) TestUpdateLocation_DefaultProfile() {
	ulr := &protos.UpdateLocationRequest{
		UserName:    "empty_sub",
		VisitedPlmn: []byte{0, 0, 0},
	}

	ula, err := suite.UpdateLocation(ulr)
	suite.NoError(err)
	suite.checkULA(
		ula, 1000, 2000, 1000, 2000, "magma.ipv4", "", "")
}

func (suite *EpsAuthTestSuite) TestUpdateLocation_UnknownSubscriber() {
	ulr := &protos.UpdateLocationRequest{
		UserName:    "sub_unknown",
		VisitedPlmn: []byte{0, 0, 0},
	}

	ula, err := suite.UpdateLocation(ulr)
	suite.EqualError(
		err,
		"rpc error: code = NotFound desc = error loading subscriber ent with assocs for network ID: test, SID: sub_unknown: Not found")
	suite.Equal(protos.ErrorCode_USER_UNKNOWN, ula.ErrorCode)
}

func (suite *EpsAuthTestSuite) checkULA(
	ula *protos.UpdateLocationAnswer,
	totalMaxUlBitRate, totalMaxDlBitRate,
	maxUlBitRate, maxDlBitRate uint32,
	apnName, apnResourceGwIp, staticUserIp string) {

	suite.Equal(protos.ErrorCode_SUCCESS, ula.GetErrorCode())
	suite.Equal(totalMaxDlBitRate, ula.GetTotalAmbr().GetMaxBandwidthDl())
	suite.Equal(totalMaxUlBitRate, ula.GetTotalAmbr().GetMaxBandwidthUl())
	suite.Equal(1, len(ula.Apn))

	apn := ula.Apn[0]
	suite.Equal(maxDlBitRate, apn.GetAmbr().GetMaxBandwidthDl())
	suite.Equal(maxUlBitRate, apn.GetAmbr().GetMaxBandwidthUl())
	suite.Equal(apnName, apn.GetServiceSelection())
	suite.Equal(protos.UpdateLocationAnswer_APNConfiguration_IPV4, apn.GetPdn())

	if len(apnResourceGwIp) > 0 {
		suite.Equal(apnResourceGwIp, apn.GetResource().GetGatewayIp())
	}
	if len(staticUserIp) > 0 {
		suip := "UNDEFINED"
		if len(apn.GetServedPartyIpAddress()) > 0 {
			suip = apn.GetServedPartyIpAddress()[0]
		}
		suite.Equal(staticUserIp, suip)
	}
	qos := apn.GetQosProfile()
	suite.Equal(int32(9), qos.GetClassId())
	suite.Equal(uint32(15), qos.GetPriorityLevel())
	suite.Equal(true, qos.GetPreemptionCapability())
	suite.Equal(false, qos.GetPreemptionVulnerability())
}
