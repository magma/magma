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

package storage_test

import (
	"fmt"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/suite"

	b "magma/dp/cloud/go/services/dp/builders"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/dp/cloud/go/services/dp/storage/dbtest"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/merrors"
)

const (
	registered          = "registered"
	unregistered        = "unregistered"
	idle                = "idle"
	someCbsdIdStr       = "some_cbsd_id"
	authorized          = "authorized"
	granted             = "granted"
	someNetwork         = "some_network"
	otherNetwork        = "other_network_id"
	someCbsdId          = 123
	otherCbsdId         = 456
	someFccId           = "some_fcc_id"
	someUserId          = "some_user_id"
	someSerialNumber    = "some_serial_number"
	anotherSerialNumber = "another_serial_number"
	nowTimestamp        = 12345678
	grant               = "grant"
)

func TestCbsdManager(t *testing.T) {
	suite.Run(t, &CbsdManagerTestSuite{})
}

type CbsdManagerTestSuite struct {
	suite.Suite
	cbsdManager     storage.CbsdManager
	resourceManager dbtest.ResourceManager
	enumMaps        map[string]map[string]int64
}

func (s *CbsdManagerTestSuite) SetupSuite() {
	builder := sqorc.GetSqlBuilder()
	errorChecker := sqorc.SQLiteErrorChecker{}
	locker := sqorc.GetSqlLocker()
	database, err := sqorc.Open("sqlite3", ":memory:")
	s.Require().NoError(err)
	s.cbsdManager = storage.NewCbsdManager(database, builder, errorChecker, locker)
	s.resourceManager = dbtest.NewResourceManager(s.T(), database, builder)
	err = s.resourceManager.CreateTables(
		&storage.DBCbsdState{},
		&storage.DBCbsd{},
		&storage.DBGrantState{},
		&storage.DBGrant{},
	)
	s.Require().NoError(err)
	err = s.resourceManager.InsertResources(
		db.NewExcludeMask("id"),
		&storage.DBCbsdState{Name: db.MakeString(unregistered)},
		&storage.DBCbsdState{Name: db.MakeString(registered)},
		&storage.DBGrantState{Name: db.MakeString("idle")},
		&storage.DBGrantState{Name: db.MakeString("granted")},
		&storage.DBGrantState{Name: db.MakeString(authorized)},
	)
	s.Require().NoError(err)
	s.enumMaps = map[string]map[string]int64{}
	for _, model := range []db.Model{
		&storage.DBCbsdState{},
		&storage.DBGrantState{},
	} {
		table := model.GetMetadata().Table
		s.enumMaps[table] = s.getNameIdMapping(model)
	}

}

func (s *CbsdManagerTestSuite) TearDownTest() {
	clock.UnfreezeClock(s.T())
	err := s.resourceManager.DropResources(
		&storage.DBCbsd{},
		&storage.DBGrant{},
	)
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) TestCreateCbsdWithDefaultValues() {
	err := s.cbsdManager.CreateCbsd(someNetwork, b.GetMutableDBCbsd(
		b.NewDBCbsdBuilder().Cbsd, registered))
	s.Require().NoError(err)
	s.thenCbsdIs(b.NewDBCbsdBuilder().WithIndoorDeployment(false).Cbsd)
}

func (s *CbsdManagerTestSuite) TestCreateCbsdWithCarrierAggregationFields() {
	err := s.cbsdManager.CreateCbsd(someNetwork, b.GetMutableDBCbsd(
		b.NewDBCbsdBuilder().
			WithCarrierAggregationEnabled(true).
			WithGrantRedundancy(true).
			WithMaxIbwMhx(140).Cbsd, registered))
	s.Require().NoError(err)
	s.thenCbsdIs(b.NewDBCbsdBuilder().
		WithIndoorDeployment(false).
		WithCarrierAggregationEnabled(true).
		WithGrantRedundancy(true).
		WithMaxIbwMhx(140).Cbsd)
}

func (s *CbsdManagerTestSuite) TestCreateSingleStepCbsd() {
	err := s.cbsdManager.CreateCbsd(someNetwork, b.GetMutableDBCbsd(
		b.NewDBCbsdBuilder().
			WithSingleStepEnabled(true).
			WithCbsdCategory("a").
			WithIndoorDeployment(true).Cbsd,
		registered))
	s.Require().NoError(err)
	s.thenCbsdIs(b.NewDBCbsdBuilder().
		WithSingleStepEnabled(true).
		WithCbsdCategory("a").
		WithIndoorDeployment(true).Cbsd)
}

func (s *CbsdManagerTestSuite) TestCreateCbsdWithExistingSerialNumber() {
	cbsd := b.NewDBCbsdBuilder().Cbsd
	err := s.cbsdManager.CreateCbsd(someNetwork, b.GetMutableDBCbsd(cbsd, registered))
	s.Require().NoError(err)
	err = s.cbsdManager.CreateCbsd(someNetwork, b.GetMutableDBCbsd(cbsd, registered))
	s.Assert().ErrorIs(err, merrors.ErrAlreadyExists)
}

func (s *CbsdManagerTestSuite) TestUpdateCbsdWithSerialNumberOfExistingCbsd() {
	stateId := s.enumMaps[storage.CbsdStateTable][registered]
	cbsd1 := b.NewDBCbsdBuilder().
		WithId(1).
		WithNetworkId(someNetwork).
		WithDesiredStateId(stateId).
		WithStateId(stateId).
		WithSerialNumber("some_serial_number_1").
		Cbsd
	cbsd2 := b.NewDBCbsdBuilder().
		WithId(2).
		WithNetworkId(someNetwork).
		WithDesiredStateId(stateId).
		WithStateId(stateId).
		WithSerialNumber("some_serial_number_2").
		Cbsd
	s.givenResourcesInserted(cbsd1, cbsd2)

	cbsd2.CbsdSerialNumber = cbsd1.CbsdSerialNumber
	m := b.GetMutableDBCbsd(cbsd2, registered)
	err := s.cbsdManager.UpdateCbsd(someNetwork, cbsd2.Id.Int64, m)
	s.Assert().ErrorIs(err, merrors.ErrAlreadyExists)
}

func (s *CbsdManagerTestSuite) TestUpdateCbsd() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	s.givenResourcesInserted(b.NewDBCbsdBuilder().
		WithId(someCbsdId).
		WithNetworkId(someNetwork).
		WithDesiredStateId(state).
		WithStateId(state).
		Cbsd)
	cbsdBuilder := b.NewDBCbsdBuilder()
	m := b.GetMutableDBCbsd(cbsdBuilder.
		WithUserId(fmt.Sprintf("%snew1", cbsdBuilder.Cbsd.UserId.String)).
		WithFccId(fmt.Sprintf("%snew2", cbsdBuilder.Cbsd.FccId.String)).
		WithSerialNumber(fmt.Sprintf("%snew3", cbsdBuilder.Cbsd.CbsdSerialNumber.String)).
		WithAntennaGain(1).
		WithEirpCapabilities(
			cbsdBuilder.Cbsd.MinPower.Float64+3,
			cbsdBuilder.Cbsd.MaxPower.Float64+2,
			cbsdBuilder.Cbsd.NumberOfPorts.Int64+4).
		WithSingleStepEnabled(true).
		WithIndoorDeployment(true).
		WithCarrierAggregationEnabled(true).
		WithMaxIbwMhx(140).
		WithGrantRedundancy(true).
		WithCbsdCategory("a").
		WithNetworkId(someNetwork).
		Cbsd,
		unregistered,
	)
	err := s.cbsdManager.UpdateCbsd(someNetwork, someCbsdId, m)
	s.Require().NoError(err)

	err = s.resourceManager.InTransaction(func() {
		actual, err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&storage.DBCbsd{}).
			Select(db.NewExcludeMask("id", "state_id",
				"cbsd_id", "is_deleted", "should_relinquish")).
			Where(sq.Eq{"id": someCbsdId}).
			Fetch()
		s.Require().NoError(err)
		expected := []db.Model{m.Cbsd}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) TestEnodebdUpdateCbsd() {
	now := time.Unix(nowTimestamp, 0)
	clock.SetAndFreezeClock(s.T(), now)
	registeredStateId := s.enumMaps[storage.CbsdStateTable][registered]
	authorizedStateId := s.enumMaps[storage.GrantStateTable][authorized]
	idleStateId := s.enumMaps[storage.GrantStateTable]["idle"]
	testCases := []struct {
		name          string
		inputCbsd     *storage.DBCbsd
		inputGrants   []*storage.DBGrant
		toUpdate      *storage.DBCbsd
		expected      *storage.DetailedCbsd
		expectedError string
	}{{
		name: "test enodebd update cbsd without a grant",
		inputCbsd: b.NewDBCbsdBuilder().
			WithId(0).
			WithNetworkId(someNetwork).
			WithSingleStepEnabled(true).
			WithSerialNumber(someSerialNumber).
			WithDesiredStateId(registeredStateId).
			WithStateId(registeredStateId).
			Cbsd,
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithSerialNumber(someSerialNumber).
			WithCbsdCategory("a").
			WithFullEnodebdAllowedInstallationParam().
			WithLastSeen(nowTimestamp).
			Cbsd,
		expected: b.NewDetailedDBCbsdBuilder().
			WithCbsd(
				b.NewDBCbsdBuilder().
					WithNetworkId(someNetwork).
					WithSingleStepEnabled(true).
					WithSerialNumber(someSerialNumber).
					WithDesiredStateId(registeredStateId).
					WithLastSeen(nowTimestamp).
					WithShouldDeregister(true).
					WithCbsdCategory("a").
					WithFullEnodebdAllowedInstallationParam().
					Cbsd,
				registered,
				registered).
			Details,
	}, {
		name: "test enodebd update cbsd with valid authorized grant",
		inputCbsd: b.NewDBCbsdBuilder().
			WithId(3).
			WithNetworkId(someNetwork).
			WithSingleStepEnabled(true).
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 3)).
			WithDesiredStateId(registeredStateId).
			WithStateId(registeredStateId).
			Cbsd,
		inputGrants: []*storage.DBGrant{
			b.NewDBGrantBuilder().
				WithId(3).
				WithCbsdId(3).
				WithGrantId("some_grant_id").
				WithFrequency(3600).
				WithMaxEirp(35).
				WithStateId(authorizedStateId).
				WithGrantExpireTime(clock.Now().UTC().Add(time.Hour)).
				WithTransmitExpireTime(clock.Now().UTC().Add(time.Hour)).
				Grant,
		},
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 3)).
			WithCbsdCategory("a").
			WithFullEnodebdAllowedInstallationParam().
			WithLastSeen(nowTimestamp).
			Cbsd,
		expected: b.NewDetailedDBCbsdBuilder().
			WithCbsd(
				b.NewDBCbsdBuilder().
					WithNetworkId(someNetwork).
					WithSingleStepEnabled(true).
					WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 3)).
					WithShouldDeregister(true).
					WithDesiredStateId(registeredStateId).
					WithCbsdCategory("a").
					WithFullEnodebdAllowedInstallationParam().
					WithLastSeen(nowTimestamp).
					Cbsd, registered, registered).
			WithGrant(authorized, 3600, clock.Now().UTC().Add(time.Hour), clock.Now().UTC().Add(time.Hour)).Details,
	}, {
		name: "test enodebd update cbsd having grant with expire times in the past",
		inputCbsd: b.NewDBCbsdBuilder().
			WithId(4).
			WithNetworkId(someNetwork).
			WithSingleStepEnabled(true).
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 4)).
			WithDesiredStateId(registeredStateId).
			WithStateId(registeredStateId).
			Cbsd,
		inputGrants: []*storage.DBGrant{
			b.NewDBGrantBuilder().
				WithId(4).
				WithCbsdId(4).
				WithStateId(authorizedStateId).
				WithGrantId("some_grant_id").
				WithFrequency(3610e6).
				WithMaxEirp(15).
				WithGrantExpireTime(clock.Now().UTC().Add(-time.Hour)).
				WithTransmitExpireTime(clock.Now().UTC().Add(-time.Hour)).
				Grant,
		},
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 4)).
			WithCbsdCategory("a").
			WithFullEnodebdAllowedInstallationParam().
			WithLastSeen(nowTimestamp).
			Cbsd,
		expected: &storage.DetailedCbsd{
			Cbsd: b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSingleStepEnabled(true).
				WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 4)).
				WithDesiredStateId(registeredStateId).
				WithShouldDeregister(true).
				WithCbsdCategory("a").
				WithFullEnodebdAllowedInstallationParam().
				WithLastSeen(nowTimestamp).
				Cbsd,
		},
	}, {
		name: "test enodebd update cbsd having grant with grant expire time in the future and transmit expire time in the past",
		inputCbsd: b.NewDBCbsdBuilder().
			WithId(5).
			WithNetworkId(someNetwork).
			WithSingleStepEnabled(true).
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 5)).
			WithDesiredStateId(registeredStateId).
			WithStateId(registeredStateId).
			Cbsd,
		inputGrants: []*storage.DBGrant{
			b.NewDBGrantBuilder().
				WithId(5).
				WithCbsdId(5).
				WithStateId(authorizedStateId).
				WithGrantId("some_grant_id").
				WithFrequency(3610e6).
				WithMaxEirp(15).
				WithGrantExpireTime(clock.Now().UTC().Add(time.Hour)).
				WithTransmitExpireTime(clock.Now().UTC().Add(-time.Hour)).
				Grant,
		},
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 5)).
			WithCbsdCategory("a").
			WithFullEnodebdAllowedInstallationParam().
			WithLastSeen(nowTimestamp).
			Cbsd,
		expected: b.NewDetailedDBCbsdBuilder().
			WithCbsd(
				b.NewDBCbsdBuilder().
					WithNetworkId(someNetwork).
					WithSingleStepEnabled(true).
					WithDesiredStateId(registeredStateId).
					WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 5)).
					WithShouldDeregister(true).
					WithCbsdCategory("a").
					WithFullEnodebdAllowedInstallationParam().
					WithLastSeen(nowTimestamp).
					Cbsd, registered, registered).Details,
	}, {
		name: "test enodebd update cbsd having grant with transmit expire time in the future and grant expire time in the past",
		inputCbsd: b.NewDBCbsdBuilder().
			WithId(6).
			WithNetworkId(someNetwork).
			WithSingleStepEnabled(true).
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 6)).
			WithDesiredStateId(registeredStateId).
			WithStateId(registeredStateId).
			Cbsd,
		inputGrants: []*storage.DBGrant{
			b.NewDBGrantBuilder().
				WithId(6).
				WithCbsdId(6).
				WithStateId(authorizedStateId).
				WithGrantId("some_grant_id").
				WithFrequency(3610e6).
				WithMaxEirp(15).
				WithGrantExpireTime(clock.Now().UTC().Add(-time.Hour)).
				WithTransmitExpireTime(clock.Now().UTC().Add(time.Hour)).
				Grant,
		},
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 6)).
			WithCbsdCategory("a").
			WithFullEnodebdAllowedInstallationParam().
			WithLastSeen(nowTimestamp).
			Cbsd,
		expected: b.NewDetailedDBCbsdBuilder().WithCbsd(
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSingleStepEnabled(true).
				WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 6)).
				WithDesiredStateId(registeredStateId).
				WithShouldDeregister(true).
				WithCbsdCategory("a").
				WithFullEnodebdAllowedInstallationParam().
				WithLastSeen(nowTimestamp).
				Cbsd, registered, registered).Details,
	}, {
		name: "test enodebd should not update if coordinates change is less than 10 m",
		inputCbsd: b.NewDBCbsdBuilder().
			WithId(7).
			WithNetworkId(someNetwork).
			WithSingleStepEnabled(true).
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 7)).
			WithDesiredStateId(registeredStateId).
			WithIndoorDeployment(true).
			WithLatitude(10).
			WithLongitude(100).
			WithStateId(registeredStateId).
			Cbsd,
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithCbsdCategory("a").
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 7)).
			WithIndoorDeployment(true).
			WithLatitude(10.00001).
			WithLongitude(100.00001).
			WithLastSeen(nowTimestamp).
			Cbsd,
		expected: b.NewDetailedDBCbsdBuilder().WithCbsd(
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSingleStepEnabled(true).
				WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 7)).
				WithDesiredStateId(registeredStateId).
				WithIndoorDeployment(true).
				WithShouldDeregister(false).
				WithLatitude(10).
				WithLongitude(100).
				WithLastSeen(nowTimestamp).
				Cbsd,
			registered,
			registered).Details,
	}, {
		name: "test enodebd only updates allowed fields",
		inputCbsd: b.NewDBCbsdBuilder().
			WithId(8).
			WithNetworkId(someNetwork).
			WithSingleStepEnabled(true).
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 8)).
			WithDesiredStateId(registeredStateId).
			WithStateId(registeredStateId).
			Cbsd,
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 8)).
			WithFccId(someFccId).
			WithUserId(someUserId).
			WithDesiredStateId(s.enumMaps[storage.CbsdStateTable][registered]).
			WithNetworkId(anotherSerialNumber).
			WithSingleStepEnabled(false).
			WithCbsdCategory("a").
			WithFullInstallationParam().
			WithLastSeen(nowTimestamp).
			Cbsd,
		expected: b.NewDetailedDBCbsdBuilder().
			WithCbsd(b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSingleStepEnabled(true).
				WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 8)).
				WithDesiredStateId(registeredStateId).
				WithFullEnodebdAllowedInstallationParam().
				WithCbsdCategory("a").
				WithLastSeen(nowTimestamp).
				WithShouldDeregister(true).
				Cbsd,
				registered,
				registered).
			Details,
	},
		{
			name: "test enodebd nulls out unfilled installation params",
			inputCbsd: b.NewDBCbsdBuilder().
				WithId(9).
				WithNetworkId(someNetwork).
				WithSingleStepEnabled(true).
				WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 9)).
				WithFullEnodebdAllowedInstallationParam().
				WithDesiredStateId(registeredStateId).
				WithStateId(registeredStateId).
				Cbsd,
			toUpdate: b.NewDBCbsdBuilder().
				Empty().
				WithNetworkId(someNetwork).
				WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 9)).
				WithCbsdCategory("a").
				WithIncompleteInstallationParam().
				WithLastSeen(nowTimestamp).
				Cbsd,
			expected: b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithNetworkId(someNetwork).
						WithSingleStepEnabled(true).
						WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 9)).
						WithCbsdCategory("a").
						WithDesiredStateId(registeredStateId).
						WithLastSeen(nowTimestamp).
						WithShouldDeregister(true).
						WithIncompleteInstallationParam().
						Cbsd, registered, registered).
				Details,
		},
		{
			name: "test enodebd update non-existent cbsd",
			inputCbsd: b.NewDBCbsdBuilder().
				WithId(10).
				WithNetworkId(someNetwork).
				WithSingleStepEnabled(true).
				WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 10)).
				WithFullEnodebdAllowedInstallationParam().
				WithDesiredStateId(registeredStateId).
				WithStateId(registeredStateId).
				Cbsd,
			toUpdate: b.NewDBCbsdBuilder().
				Empty().
				WithNetworkId(someNetwork).
				WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 101010101010)).
				WithCbsdCategory("a").
				WithFullEnodebdAllowedInstallationParam().
				WithLastSeen(nowTimestamp).
				Cbsd,
			expectedError: "Not found",
		},
		{
			name: "test enodebd update cbsd with idle grant",
			inputCbsd: b.NewDBCbsdBuilder().
				WithId(11).
				WithNetworkId(someNetwork).
				WithSingleStepEnabled(true).
				WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 11)).
				WithDesiredStateId(registeredStateId).
				WithStateId(registeredStateId).
				Cbsd,
			inputGrants: []*storage.DBGrant{b.NewDBGrantBuilder().
				WithId(11).
				WithCbsdId(11).
				WithGrantId("some_grant_id").
				WithFrequency(3610e6).
				WithMaxEirp(15).
				WithStateId(idleStateId).
				WithGrantExpireTime(clock.Now().UTC().Add(time.Hour)).
				WithTransmitExpireTime(clock.Now().UTC().Add(time.Hour)).
				Grant,
			},
			toUpdate: b.NewDBCbsdBuilder().
				Empty().
				WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 11)).
				WithCbsdCategory("a").
				WithFullEnodebdAllowedInstallationParam().
				WithLastSeen(nowTimestamp).
				Cbsd,
			expected: b.NewDetailedDBCbsdBuilder().
				WithCbsd(b.NewDBCbsdBuilder().
					WithNetworkId(someNetwork).
					WithSingleStepEnabled(true).
					WithSerialNumber(fmt.Sprintf(someSerialNumber+"%d", 11)).
					WithDesiredStateId(registeredStateId).
					WithShouldDeregister(true).
					WithCbsdCategory("a").
					WithFullEnodebdAllowedInstallationParam().
					WithLastSeen(nowTimestamp).
					Cbsd, registered, registered).
				Details,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.givenResourcesInserted(tc.inputCbsd)
			if tc.inputGrants != nil {
				for _, g := range tc.inputGrants {
					s.givenResourcesInserted(g)
				}
			}

			result, err := s.cbsdManager.EnodebdUpdateCbsd(tc.toUpdate)
			if tc.expectedError != "" {
				s.Assert().Errorf(err, tc.expectedError)
				return
			} else {
				s.Require().NoError(err)
			}

			s.Assert().Equal(tc.inputCbsd.CbsdSerialNumber, result.Cbsd.CbsdSerialNumber)
			s.Assert().Equal(tc.inputCbsd.NetworkId, result.Cbsd.NetworkId)
			s.Assert().Equal(tc.expected.Grants, result.Grants)

			err = s.resourceManager.InTransaction(func() {
				actual, err := db.NewQuery().
					WithBuilder(s.resourceManager.GetBuilder()).
					From(&storage.DBCbsd{}).
					Select(db.NewExcludeMask("id", "state_id",
						"cbsd_id", "is_deleted", "should_relinquish")).
					Where(sq.Eq{"cbsd_serial_number": tc.inputCbsd.CbsdSerialNumber}).
					Fetch()
				s.Require().NoError(err)

				s.Assert().Equal(tc.expected.Cbsd, actual[0].(*storage.DBCbsd))
			})
			s.Require().NoError(err)
		})
	}
}

func (s *CbsdManagerTestSuite) TestUpdateDeletedCbsd() {
	s.givenDeletedCbsd()

	err := s.cbsdManager.UpdateCbsd(someNetwork, someCbsdId, b.GetMutableDBCbsd(
		b.NewDBCbsdBuilder().Cbsd, registered))
	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

func (s *CbsdManagerTestSuite) TestUpdateNonExistentCbsd() {
	err := s.cbsdManager.UpdateCbsd(someNetwork, 0, b.GetMutableDBCbsd(
		b.NewDBCbsdBuilder().Cbsd, registered))
	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

func (s *CbsdManagerTestSuite) TestDeleteCbsd() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	s.givenResourcesInserted(b.NewDBCbsdBuilder().
		WithId(someCbsdId).
		WithNetworkId(someNetwork).
		WithDesiredStateId(state).
		WithStateId(state).
		Cbsd)

	err := s.cbsdManager.DeleteCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	err = s.resourceManager.InTransaction(func() {
		actual, err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&storage.DBCbsd{}).
			Select(db.NewIncludeMask("is_deleted")).
			Where(sq.Eq{"id": someCbsdId}).
			Fetch()
		s.Require().NoError(err)
		expected := []db.Model{
			&storage.DBCbsd{IsDeleted: db.MakeBool(true)},
		}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) TestDeletedDeletedCbsd() {
	s.givenDeletedCbsd()

	err := s.cbsdManager.DeleteCbsd(someNetwork, someCbsdId)
	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

func (s *CbsdManagerTestSuite) TestDeleteNonExistentCbsd() {
	err := s.cbsdManager.DeleteCbsd(someNetwork, 0)
	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

func (s *CbsdManagerTestSuite) TestFetchDeletedCbsd() {
	s.givenDeletedCbsd()

	_, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)
	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

func (s *CbsdManagerTestSuite) TestFetchCbsdFromDifferentNetwork() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	s.givenResourcesInserted(b.NewDBCbsdBuilder().
		WithNetworkId(otherNetwork).
		WithDesiredStateId(state).
		WithStateId(state).
		Cbsd)

	_, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)

	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

func (s *CbsdManagerTestSuite) TestFetchCbsdWithoutGrant() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	cbsd := b.NewDBCbsdBuilder().
		WithId(someCbsdId).
		WithNetworkId(someNetwork).
		WithDesiredStateId(state).
		WithStateId(state).
		WithCbsdId(someCbsdIdStr).
		Cbsd
	s.givenResourcesInserted(cbsd)

	actual, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	expected := b.NewDetailedDBCbsdBuilder().
		WithCbsd(
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithIndoorDeployment(false).
				WithId(someCbsdId).
				WithCbsdId(someCbsdIdStr).Cbsd,
			registered, registered).
		Details

	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestFetchCbsdWithGrant() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	grantState := s.enumMaps[storage.GrantStateTable][authorized]
	s.givenResourcesInserted(
		b.NewDBCbsdBuilder().
			WithNetworkId(someNetwork).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr).
			WithStateId(state).
			WithDesiredStateId(state).
			Cbsd,
		b.NewDBGrantBuilder().
			WithDefaultTestValues().
			WithStateId(grantState).
			WithCbsdId(someCbsdId).
			WithFrequency(3600).
			Grant,
	)

	actual, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	expected := b.NewDetailedDBCbsdBuilder().
		WithCbsd(b.NewDBCbsdBuilder().
			WithIndoorDeployment(false).
			WithNetworkId(someNetwork).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr).Cbsd,
			registered, registered).
		WithGrant(authorized, 3600, time.Unix(123, 0).UTC(), time.Unix(456, 0).UTC()).
		Details
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestFetchCbsdWithMultipleGrants() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	grantState := s.enumMaps[storage.GrantStateTable][authorized]
	s.givenResourcesInserted(
		b.NewDBCbsdBuilder().
			WithNetworkId(someNetwork).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr).
			WithStateId(state).
			WithDesiredStateId(state).
			Cbsd,
		b.NewDBGrantBuilder().
			WithDefaultTestValues().
			WithStateId(grantState).
			WithCbsdId(someCbsdId).
			WithFrequency(3600).
			Grant,
		b.NewDBGrantBuilder().
			WithDefaultTestValues().
			WithStateId(grantState).
			WithCbsdId(someCbsdId).
			WithFrequency(3650).
			Grant,
	)

	actual, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	expected := b.NewDetailedDBCbsdBuilder().
		WithCbsd(b.NewDBCbsdBuilder().
			WithIndoorDeployment(false).
			WithNetworkId(someNetwork).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr).Cbsd,
			registered, registered).
		WithGrant(authorized, 3600, time.Unix(123, 0).UTC(), time.Unix(456, 0).UTC()).
		WithGrant(authorized, 3650, time.Unix(123, 0).UTC(), time.Unix(456, 0).UTC()).
		Details
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListCbsdFromDifferentNetwork() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	s.givenResourcesInserted(b.NewDBCbsdBuilder().
		WithId(someCbsdId).
		WithNetworkId(otherNetwork).
		WithDesiredStateId(state).
		WithStateId(state).
		Cbsd)

	actual, err := s.cbsdManager.ListCbsd(someNetwork, &storage.Pagination{}, nil)
	s.Require().NoError(err)

	expected := &storage.DetailedCbsdList{
		Cbsds: []*storage.DetailedCbsd{},
		Count: 0,
	}
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListWithPagination() {
	const count = 4
	const limit = 2
	const offset = 1

	models := make([]db.Model, count)
	stateId := s.enumMaps[storage.CbsdStateTable][unregistered]
	for i := range models {
		cbsd := b.NewDBCbsdBuilder().
			WithId(int64(i + 1)).
			WithNetworkId(someNetwork).
			WithDesiredStateId(stateId).
			WithStateId(stateId).
			WithSerialNumber(fmt.Sprintf("some_serial_number%d", i+1)).
			Cbsd
		models[i] = cbsd
	}
	s.givenResourcesInserted(models...)

	pagination := &storage.Pagination{
		Limit:  db.MakeInt(limit),
		Offset: db.MakeInt(offset),
	}
	actual, err := s.cbsdManager.ListCbsd(someNetwork, pagination, nil)
	s.Require().NoError(err)

	expected := &storage.DetailedCbsdList{
		Count: count,
		Cbsds: make([]*storage.DetailedCbsd, limit),
	}

	for i := range expected.Cbsds {
		cbsd := b.NewDBCbsdBuilder().
			WithId(int64(i + 1 + offset)).
			WithIndoorDeployment(false).
			WithNetworkId(someNetwork).
			WithSerialNumber(fmt.Sprintf("some_serial_number%d", i+1+offset)).
			Cbsd
		expected.Cbsds[i] = b.NewDetailedDBCbsdBuilder().
			WithCbsd(cbsd, unregistered, unregistered).
			Details
	}
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListWithFilter() {
	const count = 1
	models := make([]db.Model, count)
	stateId := s.enumMaps[storage.CbsdStateTable][unregistered]
	for i := range models {
		cbsd := b.NewDBCbsdBuilder().
			WithId(int64(i + 1)).
			WithNetworkId(someNetwork).
			WithDesiredStateId(stateId).
			WithStateId(stateId).
			WithSerialNumber(fmt.Sprintf("some_serial_number%d", i+1)).
			Cbsd
		models[i] = cbsd
	}
	s.givenResourcesInserted(models...)

	pagination := &storage.Pagination{}
	filter := &storage.CbsdFilter{SerialNumber: "some_serial_number1"}
	actual, err := s.cbsdManager.ListCbsd(someNetwork, pagination, filter)
	s.Require().NoError(err)

	expected := &storage.DetailedCbsdList{
		Count: count,
		Cbsds: make([]*storage.DetailedCbsd, 1),
	}
	for i := range expected.Cbsds {
		cbsd := b.NewDBCbsdBuilder().
			WithId(int64(i + 1)).
			WithIndoorDeployment(false).
			WithNetworkId(someNetwork).
			WithSerialNumber(fmt.Sprintf("some_serial_number%d", i+1)).
			Cbsd
		expected.Cbsds[i] = b.NewDetailedDBCbsdBuilder().
			WithCbsd(cbsd, unregistered, unregistered).
			Details
	}
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListWithMultipleGrants() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	grantState := s.enumMaps[storage.GrantStateTable][authorized]
	s.givenResourcesInserted(
		b.NewDBCbsdBuilder().
			WithNetworkId(someNetwork).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr).
			WithStateId(state).
			WithDesiredStateId(state).
			Cbsd,
		b.NewDBCbsdBuilder().
			WithNetworkId(someNetwork).
			WithId(otherCbsdId).
			WithCbsdId(someCbsdIdStr).
			WithStateId(state).
			WithDesiredStateId(state).
			WithSerialNumber(anotherSerialNumber).
			Cbsd,
		b.NewDBGrantBuilder().
			WithDefaultTestValues().
			WithId(1).
			WithStateId(grantState).
			WithCbsdId(someCbsdId).
			WithFrequency(3530).
			Grant,
		b.NewDBGrantBuilder().
			WithDefaultTestValues().
			WithId(2).
			WithStateId(grantState).
			WithCbsdId(otherCbsdId).
			WithFrequency(3570).
			Grant,
		b.NewDBGrantBuilder().
			WithDefaultTestValues().
			WithId(3).
			WithStateId(grantState).
			WithCbsdId(otherCbsdId).
			WithFrequency(3610).
			Grant,
		b.NewDBGrantBuilder().
			WithDefaultTestValues().
			WithId(4).
			WithStateId(grantState).
			WithCbsdId(someCbsdId).
			WithFrequency(3650).
			Grant,
	)

	actual, err := s.cbsdManager.ListCbsd(someNetwork, &storage.Pagination{}, nil)
	s.Require().NoError(err)

	expected := &storage.DetailedCbsdList{
		Cbsds: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(b.NewDBCbsdBuilder().
					WithIndoorDeployment(false).
					WithId(someCbsdId).
					WithCbsdId(someCbsdIdStr).
					WithNetworkId(someNetwork).
					Cbsd,
					registered, registered).
				WithGrant(authorized, 3530, time.Unix(123, 0).UTC(), time.Unix(456, 0).UTC()).
				WithGrant(authorized, 3650, time.Unix(123, 0).UTC(), time.Unix(456, 0).UTC()).
				Details,
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(b.NewDBCbsdBuilder().
					WithIndoorDeployment(false).
					WithId(otherCbsdId).
					WithCbsdId(someCbsdIdStr).
					WithNetworkId(someNetwork).
					WithSerialNumber(anotherSerialNumber).
					Cbsd,
					registered, registered).
				WithGrant(authorized, 3570, time.Unix(123, 0).UTC(), time.Unix(456, 0).UTC()).
				WithGrant(authorized, 3610, time.Unix(123, 0).UTC(), time.Unix(456, 0).UTC()).
				Details,
		},
		Count: 2,
	}

	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListDeletedCbsd() {
	s.givenDeletedCbsd()

	actual, err := s.cbsdManager.ListCbsd(someNetwork, &storage.Pagination{}, nil)
	s.Require().NoError(err)

	expected := &storage.DetailedCbsdList{
		Cbsds: []*storage.DetailedCbsd{},
		Count: 0,
	}
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestDeregisterCbsd() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	s.givenResourcesInserted(b.NewDBCbsdBuilder().
		WithId(someCbsdId).
		WithNetworkId(someNetwork).
		WithDesiredStateId(state).
		WithStateId(state).
		Cbsd)

	err := s.cbsdManager.DeregisterCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	err = s.resourceManager.InTransaction(func() {
		actual, err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&storage.DBCbsd{}).
			Select(db.NewIncludeMask("should_deregister")).
			Where(sq.Eq{"id": someCbsdId}).
			Fetch()
		s.Require().NoError(err)
		cbsd := &storage.DBCbsd{ShouldDeregister: db.MakeBool(true)}
		expected := []db.Model{cbsd}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) TestDeregisterNonExistentCbsd() {
	err := s.cbsdManager.DeregisterCbsd(someNetwork, 0)
	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

func (s *CbsdManagerTestSuite) TestRelinquishCbsd() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	s.givenResourcesInserted(b.NewDBCbsdBuilder().
		WithId(someCbsdId).
		WithNetworkId(someNetwork).
		WithDesiredStateId(state).
		WithStateId(state).
		Cbsd)

	err := s.cbsdManager.RelinquishCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	err = s.resourceManager.InTransaction(func() {
		actual, err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&storage.DBCbsd{}).
			Select(db.NewIncludeMask("should_relinquish")).
			Where(sq.Eq{"id": someCbsdId}).
			Fetch()
		s.Require().NoError(err)
		cbsd := &storage.DBCbsd{ShouldRelinquish: db.MakeBool(true)}
		expected := []db.Model{cbsd}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) TestRelinquishNonExistentCbsd() {
	err := s.cbsdManager.RelinquishCbsd(someNetwork, 0)
	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

// TODO this modifies expected - needs fix!!!
func (s *CbsdManagerTestSuite) thenCbsdIs(expected *storage.DBCbsd) {
	err := s.resourceManager.InTransaction(func() {
		actual, err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&storage.DBCbsd{}).
			Select(db.NewExcludeMask("id", "state_id", "desired_state_id")).
			Join(db.NewQuery().
				From(&storage.DBCbsdState{}).
				As("t1").
				On(db.On(storage.CbsdTable, "state_id", "t1", "id")).
				Select(db.NewIncludeMask("name"))).
			Join(db.NewQuery().
				From(&storage.DBCbsdState{}).
				As("t2").
				On(db.On(storage.CbsdTable, "desired_state_id", "t2", "id")).
				Select(db.NewIncludeMask("name"))).
			Where(sq.Eq{"cbsd_serial_number": "some_serial_number"}).
			Fetch()
		s.Require().NoError(err)

		cbsd := expected
		cbsd.NetworkId = db.MakeString(someNetwork)
		cbsd.IsDeleted = db.MakeBool(false)
		cbsd.ShouldDeregister = db.MakeBool(false)
		cbsd.ShouldRelinquish = db.MakeBool(false)
		cbsd.Channels = []storage.Channel{}
		expected := []db.Model{
			cbsd,
			&storage.DBCbsdState{Name: db.MakeString(unregistered)},
			&storage.DBCbsdState{Name: db.MakeString(registered)},
		}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) givenResourcesInserted(models ...db.Model) {
	err := s.resourceManager.InsertResources(db.NewExcludeMask(), models...)
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) givenDeletedCbsd() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	cbsd := b.NewDBCbsdBuilder().
		WithId(someCbsdId).
		WithNetworkId(someNetwork).
		WithDesiredStateId(state).
		WithStateId(state).
		Cbsd
	cbsd.IsDeleted = db.MakeBool(true)
	s.givenResourcesInserted(cbsd)
}

func (s *CbsdManagerTestSuite) getNameIdMapping(model db.Model) map[string]int64 {
	var resources [][]db.Model
	err := s.resourceManager.InTransaction(func() {
		var err error
		resources, err = db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(model).
			Select(db.NewExcludeMask()).
			List()
		s.Require().NoError(err)
	})
	s.Require().NoError(err)
	m := make(map[string]int64, len(resources))
	for _, r := range resources {
		enum := r[0].(storage.EnumModel)
		m[enum.GetName()] = enum.GetId()
	}
	return m
}
