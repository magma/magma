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
	b "magma/dp/cloud/go/services/dp/builders"
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
	lastSeenTimestamp int64 = 1234567
	someSerialNumber        = "some_serial_number"
	someCbsdId              = "some_cbsd_id"
	registered              = "registered"
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
	protoPayloadBuilder := b.NewCbsdProtoPayloadBuilder()
	dbCbsdBuilder := b.NewDBCbsdBuilder()

	testCases := []struct {
		name     string
		input    *protos.CbsdData
		expected *storage.MutableCbsd
	}{{
		name:     "test create cbsd",
		input:    protoPayloadBuilder.Payload,
		expected: b.GetMutableDBCbsd(dbCbsdBuilder.Cbsd, registered),
	}, {
		name: "test create single step cbsd",
		input: protoPayloadBuilder.
			WithSingleStepEnabled().
			WithCbsdCategory("a").
			Payload,
		expected: b.GetMutableDBCbsd(
			dbCbsdBuilder.
				WithSingleStepEnabled(true).
				WithCbsdCategory("a").
				Cbsd,
			registered,
		),
	}}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			request := &protos.CreateCbsdRequest{
				NetworkId: networkId,
				Data:      tc.input,
			}
			_, err := s.manager.CreateCbsd(context.Background(), request)
			s.Require().NoError(err)

			s.Assert().Equal(networkId, s.store.networkId)
			s.Assert().Equal(tc.expected, s.store.mutableData)
		})
	}
}

func (s *CbsdManagerTestSuite) TestCreateWithDuplicateData() {
	s.store.err = merrors.ErrAlreadyExists

	request := &protos.CreateCbsdRequest{
		NetworkId: networkId,
		Data:      b.NewCbsdProtoPayloadBuilder().Payload,
	}
	_, err := s.manager.CreateCbsd(context.Background(), request)
	s.Require().Error(err)

	errStatus, _ := status.FromError(err)
	s.Assert().Equal(codes.AlreadyExists, errStatus.Code())
}

func (s *CbsdManagerTestSuite) TestUserUpdateCbsd() {
	// TODO adjust when User-triggered cbsd update is modified
	request := &protos.UpdateCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
		Data:      b.NewCbsdProtoPayloadBuilder().Payload,
	}
	_, err := s.manager.UserUpdateCbsd(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(cbsdId, s.store.id)
	s.Assert().Equal(b.GetMutableDBCbsd(b.NewDBCbsdBuilder().Cbsd, registered), s.store.mutableData)
}

func (s *CbsdManagerTestSuite) TestEnodebdUpdateCbsd() {
	testCases := []struct {
		name     string
		payload  *protos.CbsdData
		expected *storage.DBCbsd
	}{{
		name:     "update cbsd",
		payload:  b.NewCbsdProtoPayloadBuilder().WithEmptyInstallationParam().Payload,
		expected: b.NewDBCbsdBuilder().Cbsd,
	}, {
		name: "update cbsd with full installation param",
		payload: b.NewCbsdProtoPayloadBuilder().
			WithFullInstallationParam().Payload,
		expected: b.NewDBCbsdBuilder().
			WithFullInstallationParam().Cbsd,
	}, {
		name: "update cbsd with incomplete installation param",
		payload: b.NewCbsdProtoPayloadBuilder().
			WithIncompleteInstallationParam().Payload,
		expected: b.NewDBCbsdBuilder().
			WithIncompleteInstallationParam().Cbsd,
	}}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			request := &protos.EnodebdUpdateCbsdRequest{
				SerialNumber: tc.payload.SerialNumber,
				InstallationParam: &protos.InstallationParam{
					LatitudeDeg:      tc.payload.InstallationParam.LatitudeDeg,
					LongitudeDeg:     tc.payload.InstallationParam.LongitudeDeg,
					IndoorDeployment: tc.payload.InstallationParam.IndoorDeployment,
					HeightM:          tc.payload.InstallationParam.HeightM,
					HeightType:       tc.payload.InstallationParam.HeightType,
					AntennaGain:      tc.payload.InstallationParam.AntennaGain,
				},
				CbsdCategory: tc.payload.CbsdCategory,
			}
			s.store.data = tc.expected
			_, err := s.manager.EnodebdUpdateCbsd(context.Background(), request)
			s.Require().NoError(err)
			s.Assert().Equal(tc.expected, s.store.data)
		})
	}
}

func (s *CbsdManagerTestSuite) TestEnodebdUpdateWithError() {
	testCases := []struct {
		name                string
		storageError        error
		expectedErrorStatus codes.Code
	}{{
		name:                "update cbsd with duplicate data",
		storageError:        merrors.ErrAlreadyExists,
		expectedErrorStatus: codes.AlreadyExists,
	}, {
		name:                "update nonexistent cbsd",
		storageError:        merrors.ErrNotFound,
		expectedErrorStatus: codes.NotFound,
	}}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.store.err = tc.storageError
			s.store.data = b.NewDBCbsdBuilder().WithFullInstallationParam().Cbsd

			request := &protos.EnodebdUpdateCbsdRequest{
				SerialNumber:      someSerialNumber,
				InstallationParam: &protos.InstallationParam{},
				CbsdCategory:      "a",
			}
			_, err := s.manager.EnodebdUpdateCbsd(context.Background(), request)
			s.Require().Error(err)

			errStatus, _ := status.FromError(err)
			s.Assert().Equal(tc.expectedErrorStatus, errStatus.Code())
		})
	}
}

func (s *CbsdManagerTestSuite) TestUserUpdateWithError() {
	testCases := []struct {
		name                string
		storageError        error
		expectedErrorStatus codes.Code
	}{{
		name:                "update cbsd with duplicate data",
		storageError:        merrors.ErrAlreadyExists,
		expectedErrorStatus: codes.AlreadyExists,
	}, {
		name:                "update nonexistent cbsd",
		storageError:        merrors.ErrNotFound,
		expectedErrorStatus: codes.NotFound,
	}}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.store.err = tc.storageError

			request := &protos.UpdateCbsdRequest{
				NetworkId: networkId,
				Id:        cbsdId,
				Data:      b.NewCbsdProtoPayloadBuilder().Payload,
			}
			_, err := s.manager.UserUpdateCbsd(context.Background(), request)
			s.Require().Error(err)

			errStatus, _ := status.FromError(err)
			s.Assert().Equal(tc.expectedErrorStatus, errStatus.Code())
		})
	}
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
	testCases := []struct {
		name     string
		input    *storage.DetailedCbsd
		expected *protos.CbsdDetails
	}{{
		name: "fetch cbsd with default installation param",
		input: b.NewDetailedDBCbsdBuilder(
			b.NewDBCbsdBuilder()).
			WithDefaultTestData().Details,
		expected: b.NewDetailedProtoCbsdBuilder(
			b.NewCbsdProtoPayloadBuilder().
				WithEmptyInstallationParam()).
			WithGrant().Details,
	}, {
		name: "fetch cbsd with full installation param",
		input: b.NewDetailedDBCbsdBuilder(b.NewDBCbsdBuilder().
			WithFullInstallationParam()).
			WithDefaultTestData().Details,
		expected: b.NewDetailedProtoCbsdBuilder(
			b.NewCbsdProtoPayloadBuilder().
				WithFullInstallationParam()).
			WithGrant().Details,
	}, {
		name: "fetch cbsd with incomplete installation param",
		input: b.NewDetailedDBCbsdBuilder(b.NewDBCbsdBuilder().
			WithIncompleteInstallationParam()).
			WithDefaultTestData().Details,
		expected: b.NewDetailedProtoCbsdBuilder(
			b.NewCbsdProtoPayloadBuilder().
				WithIncompleteInstallationParam()).
			WithGrant().Details,
	}}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.store.details = tc.input

			request := &protos.FetchCbsdRequest{
				NetworkId: networkId,
				Id:        cbsdId,
			}
			actual, err := s.manager.FetchCbsd(context.Background(), request)
			s.Require().NoError(err)

			s.Assert().Equal(networkId, s.store.networkId)
			s.Assert().Equal(cbsdId, s.store.id)
			s.Assert().Equal(tc.expected, actual.Details)
		})
	}
}

func (s *CbsdManagerTestSuite) TestFetchNonActiveCbsd() {
	now := time.Unix(lastSeenTimestamp, 0).Add(interval)
	clock.SetAndFreezeClock(s.T(), now)
	cbsdBuilder := b.NewDBCbsdBuilder().
		WithId(cbsdId).
		WithCbsdId(someCbsdId).
		WithLastSeen(lastSeenTimestamp)
	s.store.details = b.NewDetailedDBCbsdBuilder(cbsdBuilder).
		WithDefaultTestData().
		Details

	request := &protos.FetchCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
	}
	actual, err := s.manager.FetchCbsd(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(cbsdId, s.store.id)
	expected := b.NewDetailedProtoCbsdBuilder(b.NewCbsdProtoPayloadBuilder().
		WithEmptyInstallationParam()).
		WithDefaultTestData().Details
	s.Assert().Equal(expected, actual.Details)
}

func (s *CbsdManagerTestSuite) TestFetchCbsdWithoutGrant() {
	builder := b.NewDBCbsdBuilder()
	s.store.details = &storage.DetailedCbsd{
		Cbsd:      builder.Cbsd,
		CbsdState: &storage.DBCbsdState{},
		DesiredState: &storage.DBCbsdState{
			Name: db.MakeString(registered),
		},
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
	expected := &protos.CbsdDetails{
		Data: b.NewCbsdProtoPayloadBuilder().
			WithEmptyInstallationParam().
			Payload,
	}
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
	dbCbsdBuilder := b.NewDBCbsdBuilder().
		WithId(cbsdId).
		WithCbsdId(someCbsdId).
		WithLastSeen(lastSeenTimestamp)
	s.store.list = b.GetDetailedDBCbsdList(b.NewDetailedDBCbsdBuilder(dbCbsdBuilder).
		WithDefaultTestData())

	request := &protos.ListCbsdRequest{
		NetworkId:  networkId,
		Pagination: &protos.Pagination{},
	}
	actual, err := s.manager.ListCbsds(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(&storage.Pagination{}, s.store.pagination)
	payloadBuilder := b.NewCbsdProtoPayloadBuilder().
		WithEmptyInstallationParam()
	detailsBuilder := b.NewDetailedProtoCbsdBuilder(payloadBuilder).
		WithDefaultTestData().
		Active()
	expected := b.GetDetailedProtoCbsdList(detailsBuilder)
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListCbsdWithPagination() {
	dbCbsdBuilder := b.NewDBCbsdBuilder().
		WithId(cbsdId).
		WithCbsdId(someCbsdId).
		WithLastSeen(lastSeenTimestamp)
	s.store.list = b.GetDetailedDBCbsdList(b.NewDetailedDBCbsdBuilder(dbCbsdBuilder).
		WithDefaultTestData())

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
	payloadBuilder := b.NewCbsdProtoPayloadBuilder().WithEmptyInstallationParam()
	detailsBuilder := b.NewDetailedProtoCbsdBuilder(payloadBuilder).
		WithDefaultTestData().
		Active()
	expected := b.GetDetailedProtoCbsdList(detailsBuilder)
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListCbsdWithFilter() {
	dbCbsdBuilder := b.NewDBCbsdBuilder().
		WithId(cbsdId).
		WithCbsdId(someCbsdId).
		WithLastSeen(lastSeenTimestamp)
	s.store.list = b.GetDetailedDBCbsdList(
		b.NewDetailedDBCbsdBuilder(dbCbsdBuilder).
			WithDefaultTestData(),
	)

	request := &protos.ListCbsdRequest{
		NetworkId:  networkId,
		Pagination: &protos.Pagination{},
		Filter: &protos.CbsdFilter{
			SerialNumber: someSerialNumber,
		},
	}
	actual, err := s.manager.ListCbsds(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	expectedFilter := &storage.CbsdFilter{
		SerialNumber: someSerialNumber,
	}
	s.Assert().Equal(expectedFilter, s.store.filter)
	payloadBuilder := b.NewCbsdProtoPayloadBuilder().
		WithEmptyInstallationParam()
	detailsBuilder := b.NewDetailedProtoCbsdBuilder(payloadBuilder).
		WithDefaultTestData().
		Active()
	expected := b.GetDetailedProtoCbsdList(detailsBuilder)
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestDeregisterCbsd() {
	request := &protos.DeregisterCbsdRequest{
		NetworkId: networkId,
		Id:        cbsdId,
	}
	_, err := s.manager.DeregisterCbsd(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(cbsdId, s.store.id)
}

type stubCbsdManager struct {
	networkId   string
	id          int64
	data        *storage.DBCbsd
	mutableData *storage.MutableCbsd
	details     *storage.DetailedCbsd
	list        *storage.DetailedCbsdList
	pagination  *storage.Pagination
	filter      *storage.CbsdFilter
	err         error
}

func (s *stubCbsdManager) CreateCbsd(networkId string, data *storage.MutableCbsd) error {
	s.networkId = networkId
	s.mutableData = data
	return s.err
}

func (s *stubCbsdManager) UpdateCbsd(networkId string, id int64, data *storage.MutableCbsd) error {
	s.networkId = networkId
	s.id = id
	s.mutableData = data
	return s.err
}

func (s *stubCbsdManager) EnodebdUpdateCbsd(data *storage.DBCbsd) error {
	s.data.CbsdCategory = data.CbsdCategory
	s.data.AntennaGain = data.AntennaGain
	s.data.LatitudeDeg = data.LatitudeDeg
	s.data.LongitudeDeg = data.LongitudeDeg
	s.data.HeightType = data.HeightType
	s.data.HeightM = data.HeightM
	s.data.IndoorDeployment = data.IndoorDeployment
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

func (s *stubCbsdManager) ListCbsd(networkId string, pagination *storage.Pagination, filter *storage.CbsdFilter) (*storage.DetailedCbsdList, error) {
	s.networkId = networkId
	s.pagination = pagination
	s.filter = filter
	return s.list, s.err
}

func (s *stubCbsdManager) DeregisterCbsd(networkId string, id int64) error {
	s.networkId = networkId
	s.id = id
	return s.err
}
