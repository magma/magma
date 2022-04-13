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

package servicers_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/servicers"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/lib/go/merrors"
)

func TestCbsdManager(t *testing.T) {
	suite.Run(t, &CbsdManagerTestSuite{})
}

type CbsdManagerTestSuite struct {
	suite.Suite
	manager protos.CbsdManagementServer
	store   *stubCbsdManager
}

const (
	networkId               = "some_network"
	cbsdId            int64 = 123
	interval                = time.Hour
	lastSeenTimestamp       = 1234567
)

func (s *CbsdManagerTestSuite) SetupTest() {
	s.store = &stubCbsdManager{}
	s.manager = servicers.NewCbsdManager(s.store, interval)

	now := time.Unix(lastSeenTimestamp, 0).Add(interval - time.Second)
	clock.SetAndFreezeClock(s.T(), now)
}

func (s *CbsdManagerTestSuite) TearDownTest() {
	clock.UnfreezeClock(s.T())
}

func (s *CbsdManagerTestSuite) TestCreateCbsd() {
	request := &protos.CreateCbsdRequest{
		NetworkId: networkId,
		Data:      getProtoCbsd(),
	}
	_, err := s.manager.CreateCbsd(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(getDBCbsd(), s.store.data)
}

func (s *CbsdManagerTestSuite) TestCreateWithDuplicateData() {
	s.store.err = merrors.ErrAlreadyExists

	request := &protos.CreateCbsdRequest{
		NetworkId: networkId,
		Data:      getProtoCbsd(),
	}
	_, err := s.manager.CreateCbsd(context.Background(), request)
	s.Require().Error(err)

	errStatus, _ := status.FromError(err)
	s.Assert().Equal(codes.AlreadyExists, errStatus.Code())
}

func (s *CbsdManagerTestSuite) TestUpdateCbsd() {
	request := &protos.UpdateCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
		Data:      getProtoCbsd(),
	}
	_, err := s.manager.UpdateCbsd(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(cbsdId, s.store.id)
	s.Assert().Equal(getDBCbsd(), s.store.data)
}

func (s *CbsdManagerTestSuite) TestUpdateNonexistentCbsd() {
	s.store.err = merrors.ErrNotFound

	request := &protos.UpdateCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
		Data:      getProtoCbsd(),
	}
	_, err := s.manager.UpdateCbsd(context.Background(), request)
	s.Require().Error(err)

	errStatus, _ := status.FromError(err)
	s.Assert().Equal(codes.NotFound, errStatus.Code())
}

func (s *CbsdManagerTestSuite) TestUpdateWithDuplicateData() {
	s.store.err = merrors.ErrAlreadyExists

	request := &protos.UpdateCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
		Data:      getProtoCbsd(),
	}
	_, err := s.manager.UpdateCbsd(context.Background(), request)
	s.Require().Error(err)

	errStatus, _ := status.FromError(err)
	s.Assert().Equal(codes.AlreadyExists, errStatus.Code())
}

func (s *CbsdManagerTestSuite) TestDeleteCbsd() {
	request := &protos.DeleteCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
	}
	_, err := s.manager.DeleteCbsd(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(cbsdId, s.store.id)
}

func (s *CbsdManagerTestSuite) TestDeleteNonExistentCbsd() {
	s.store.err = merrors.ErrNotFound

	request := &protos.DeleteCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
	}
	_, err := s.manager.DeleteCbsd(context.Background(), request)
	s.Require().Error(err)

	errStatus, _ := status.FromError(err)
	s.Assert().Equal(codes.NotFound, errStatus.Code())
}

func (s *CbsdManagerTestSuite) TestFetchCbsd() {
	s.store.details = getDetailedCbsd()

	request := &protos.FetchCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
	}
	actual, err := s.manager.FetchCbsd(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(cbsdId, s.store.id)
	s.Assert().Equal(getProtoDetailedCbsd(), actual.Details)
}

func (s *CbsdManagerTestSuite) TestFetchNonActiveCbsd() {
	now := time.Unix(lastSeenTimestamp, 0).Add(interval)
	clock.SetAndFreezeClock(s.T(), now)
	s.store.details = getDetailedCbsd()

	request := &protos.FetchCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
	}
	actual, err := s.manager.FetchCbsd(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(cbsdId, s.store.id)
	expected := getProtoDetailedCbsd()
	expected.IsActive = false
	s.Assert().Equal(expected, actual.Details)
}

func (s *CbsdManagerTestSuite) TestFetchCbsdWithoutGrant() {
	s.store.details = &storage.DetailedCbsd{
		Cbsd:       getDBCbsd(),
		CbsdState:  &storage.DBCbsdState{},
		Grant:      &storage.DBGrant{},
		GrantState: &storage.DBGrantState{},
	}

	request := &protos.FetchCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
	}
	actual, err := s.manager.FetchCbsd(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(cbsdId, s.store.id)
	expected := &protos.CbsdDetails{Data: getProtoCbsd()}
	s.Assert().Equal(expected, actual.Details)
}

func (s *CbsdManagerTestSuite) TestFetchNonexistentCbsd() {
	s.store.err = merrors.ErrNotFound

	request := &protos.FetchCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
	}
	_, err := s.manager.FetchCbsd(context.Background(), request)
	s.Require().Error(err)

	errStatus, _ := status.FromError(err)
	s.Assert().Equal(codes.NotFound, errStatus.Code())
}

func (s *CbsdManagerTestSuite) TestListCbsd() {
	s.store.list = getDetailedCbsdList()

	request := &protos.ListCbsdRequest{
		NetworkId:  networkId,
		Pagination: &protos.Pagination{},
	}
	actual, err := s.manager.ListCbsds(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(&storage.Pagination{}, s.store.pagination)
	expected := getProtoDetailedCbsdList()
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListCbsdWithPagination() {
	s.store.list = getDetailedCbsdList()

	request := &protos.ListCbsdRequest{
		NetworkId: networkId,
		Pagination: &protos.Pagination{
			Limit:  wrapperspb.Int64(10),
			Offset: wrapperspb.Int64(20),
		},
	}
	actual, err := s.manager.ListCbsds(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	expectedPagination := &storage.Pagination{
		Limit:  db.MakeInt(10),
		Offset: db.MakeInt(20),
	}
	s.Assert().Equal(expectedPagination, s.store.pagination)
	expected := getProtoDetailedCbsdList()
	s.Assert().Equal(expected, actual)
}

func getProtoCbsd() *protos.CbsdData {
	return &protos.CbsdData{
		UserId:       "some_user_id",
		FccId:        "some_fcc_id",
		SerialNumber: "some_serial_number",
		Capabilities: &protos.Capabilities{
			MinPower:         10,
			MaxPower:         20,
			NumberOfAntennas: 2,
			AntennaGain:      15,
		},
		Preferences: &protos.FrequencyPreferences{
			BandwidthMhz:   20,
			FrequenciesMhz: []int64{3600},
		},
	}
}

func getProtoDetailedCbsd() *protos.CbsdDetails {
	return &protos.CbsdDetails{
		Id:       cbsdId,
		Data:     getProtoCbsd(),
		CbsdId:   "some_cbsd_id",
		State:    "registered",
		IsActive: true,
		Grant: &protos.GrantDetails{
			BandwidthMhz:            20,
			FrequencyMhz:            3610,
			MaxEirp:                 35,
			State:                   "authorized",
			TransmitExpireTimestamp: 1e9,
			GrantExpireTimestamp:    2e9,
		},
	}
}

func getProtoDetailedCbsdList() *protos.ListCbsdResponse {
	return &protos.ListCbsdResponse{
		Details:    []*protos.CbsdDetails{getProtoDetailedCbsd()},
		TotalCount: 1,
	}
}

func getDBCbsd() *storage.DBCbsd {
	return &storage.DBCbsd{
		UserId:                  db.MakeString("some_user_id"),
		FccId:                   db.MakeString("some_fcc_id"),
		CbsdSerialNumber:        db.MakeString("some_serial_number"),
		PreferredBandwidthMHz:   db.MakeInt(20),
		PreferredFrequenciesMHz: db.MakeString("[3600]"),
		MinPower:                db.MakeFloat(10),
		MaxPower:                db.MakeFloat(20),
		AntennaGain:             db.MakeFloat(15),
		NumberOfPorts:           db.MakeInt(2),
	}
}

func getDetailedCbsd() *storage.DetailedCbsd {
	cbsd := getDBCbsd()
	cbsd.Id = db.MakeInt(cbsdId)
	cbsd.CbsdId = db.MakeString("some_cbsd_id")
	cbsd.LastSeen = db.MakeTime(time.Unix(lastSeenTimestamp, 0).UTC())
	return &storage.DetailedCbsd{
		Cbsd: cbsd,
		CbsdState: &storage.DBCbsdState{
			Name: db.MakeString("registered"),
		},
		Grant: &storage.DBGrant{
			GrantExpireTime:    db.MakeTime(time.Unix(2e9, 0).UTC()),
			TransmitExpireTime: db.MakeTime(time.Unix(1e9, 0).UTC()),
			LowFrequency:       db.MakeInt(3600 * 1e6),
			HighFrequency:      db.MakeInt(3620 * 1e6),
			MaxEirp:            db.MakeFloat(35),
		},
		GrantState: &storage.DBGrantState{
			Name: db.MakeString("authorized"),
		},
	}
}

func getDetailedCbsdList() *storage.DetailedCbsdList {
	return &storage.DetailedCbsdList{
		Cbsds: []*storage.DetailedCbsd{getDetailedCbsd()},
		Count: 1,
	}
}

type stubCbsdManager struct {
	networkId  string
	id         int64
	data       *storage.DBCbsd
	details    *storage.DetailedCbsd
	list       *storage.DetailedCbsdList
	pagination *storage.Pagination
	err        error
}

func (s *stubCbsdManager) CreateCbsd(networkId string, data *storage.DBCbsd) error {
	s.networkId = networkId
	s.data = data
	return s.err
}

func (s *stubCbsdManager) UpdateCbsd(networkId string, id int64, data *storage.DBCbsd) error {
	s.networkId = networkId
	s.id = id
	s.data = data
	return s.err
}

func (s *stubCbsdManager) DeleteCbsd(networkId string, id int64) error {
	s.networkId = networkId
	s.id = id
	return s.err
}

func (s *stubCbsdManager) FetchCbsd(networkId string, id int64) (*storage.DetailedCbsd, error) {
	s.networkId = networkId
	s.id = id
	return s.details, s.err
}

func (s *stubCbsdManager) ListCbsd(networkId string, pagination *storage.Pagination) (*storage.DetailedCbsdList, error) {
	s.networkId = networkId
	s.pagination = pagination
	return s.list, s.err
}
