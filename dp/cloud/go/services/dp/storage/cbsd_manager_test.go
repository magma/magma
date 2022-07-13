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

	sq "github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/suite"

	b "magma/dp/cloud/go/services/dp/builders"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/dp/cloud/go/services/dp/storage/dbtest"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/merrors"
)

const (
	registered            = "registered"
	unregistered          = "unregistered"
	someCbsdIdStr         = "some_cbsd_id"
	authorized            = "authorized"
	someNetwork           = "some_network"
	otherNetwork          = "other_network_id"
	someCbsdId            = 123
	someFccId             = "some_fcc_id"
	someUserId            = "some_user_id"
	someSerialNumber      = "some_serial_number"
	someOtherSerialNumber = "some_other_serial_number"
	anotherSerialNumber   = "another_serial_number"
	differentSerialNumber = "different_serial_number"
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
	s.verifyCbsdCreation(b.NewDBCbsdBuilder().WithIndoorDeployment(false).Cbsd)
}

func (s *CbsdManagerTestSuite) TestCreateSingleStepCbsd() {
	err := s.cbsdManager.CreateCbsd(someNetwork, b.GetMutableDBCbsd(
		b.NewDBCbsdBuilder().
			WithSingleStepEnabled(true).
			WithCbsdCategory("a").
			WithIndoorDeployment(true).Cbsd,
		registered))
	s.Require().NoError(err)
	s.verifyCbsdCreation(b.NewDBCbsdBuilder().
		WithSingleStepEnabled(true).
		WithCbsdCategory("a").
		WithIndoorDeployment(true).Cbsd)
}

func (s *CbsdManagerTestSuite) verifyCbsdCreation(expected *storage.DBCbsd) {
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
		cbsd.GrantAttempts = db.MakeInt(0)
		cbsd.IsDeleted = db.MakeBool(false)
		cbsd.ShouldDeregister = db.MakeBool(false)
		expected := []db.Model{
			cbsd,
			&storage.DBCbsdState{Name: db.MakeString(unregistered)},
			&storage.DBCbsdState{Name: db.MakeString(registered)},
		}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
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
		WithMaxPower(cbsdBuilder.Cbsd.MaxPower.Float64+2).
		WithMinPower(cbsdBuilder.Cbsd.MinPower.Float64+3).
		WithNumberOfPorts(cbsdBuilder.Cbsd.NumberOfPorts.Int64+4).
		WithSingleStepEnabled(true).
		WithIndoorDeployment(true).
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
				"cbsd_id", "grant_attempts", "is_deleted")).
			Where(sq.Eq{"id": someCbsdId}).
			Fetch()
		s.Require().NoError(err)
		expected := []db.Model{m.Cbsd}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) TestEnodebdUpdateCbsd() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	testCases := []struct {
		name     string
		input    *storage.DBCbsd
		toUpdate *storage.DBCbsd
		expected *storage.DBCbsd
	}{{
		name: "test enodebd update",
		input: b.NewDBCbsdBuilder().
			WithId(1).
			WithNetworkId(someNetwork).
			WithSerialNumber(someSerialNumber).
			WithDesiredStateId(state).
			WithStateId(state).
			Cbsd,
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithSerialNumber(someSerialNumber).
			WithCbsdCategory("a").
			WithFullInstallationParam().
			Cbsd,
		expected: b.NewDBCbsdBuilder().
			WithNetworkId(someNetwork).
			WithSerialNumber(someSerialNumber).
			WithDesiredStateId(state).
			WithShouldDeregister(true).
			WithCbsdCategory("a").
			WithFullInstallationParam().
			Cbsd,
	}, {
		name: "test enodebd should not update if coordinates change is less than 10 m",
		input: b.NewDBCbsdBuilder().
			WithId(2).
			WithNetworkId(someNetwork).
			WithSerialNumber(someOtherSerialNumber).
			WithDesiredStateId(state).
			WithIndoorDeployment(true).
			WithLatitude(10).
			WithLongitude(100).
			WithStateId(state).
			Cbsd,
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithSerialNumber(someOtherSerialNumber).
			WithCbsdCategory("a").
			WithIndoorDeployment(true).
			WithLatitude(10.00001).
			WithLongitude(100.00001).
			Cbsd,
		expected: b.NewDBCbsdBuilder().
			WithNetworkId(someNetwork).
			WithSerialNumber(someOtherSerialNumber).
			WithDesiredStateId(state).
			WithIndoorDeployment(true).
			WithShouldDeregister(false).
			WithLatitude(10).
			WithLongitude(100).
			Cbsd,
	}, {
		name: "test enodebd only updates allowed fields",
		input: b.NewDBCbsdBuilder().
			WithId(3).
			WithNetworkId(someNetwork).
			WithSerialNumber(anotherSerialNumber).
			WithDesiredStateId(state).
			WithStateId(state).
			Cbsd,
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithSerialNumber(anotherSerialNumber).
			WithFccId(someFccId).
			WithUserId(someUserId).
			WithDesiredStateId(s.enumMaps[storage.CbsdStateTable][registered]).
			WithNetworkId(anotherSerialNumber).
			WithSingleStepEnabled(false).
			WithCbsdCategory("a").
			WithFullInstallationParam().
			Cbsd,
		expected: b.NewDBCbsdBuilder().
			WithNetworkId(someNetwork).
			WithSerialNumber(anotherSerialNumber).
			WithDesiredStateId(state).
			WithFullInstallationParam().
			WithCbsdCategory("a").
			WithShouldDeregister(true).
			Cbsd,
	}, {
		name: "test enodebd nulls out unfilled installation params",
		input: b.NewDBCbsdBuilder().
			WithId(4).
			WithNetworkId(someNetwork).
			WithSerialNumber(differentSerialNumber).
			WithFullInstallationParam().
			WithDesiredStateId(state).
			WithStateId(state).
			Cbsd,
		toUpdate: b.NewDBCbsdBuilder().
			Empty().
			WithNetworkId(someNetwork).
			WithSerialNumber(differentSerialNumber).
			WithCbsdCategory("a").
			WithIncompleteInstallationParam().
			Cbsd,
		expected: b.NewDBCbsdBuilder().
			WithNetworkId(someNetwork).
			WithSerialNumber(differentSerialNumber).
			WithCbsdCategory("a").
			WithDesiredStateId(state).
			WithShouldDeregister(true).
			WithIncompleteInstallationParam().
			Cbsd,
	}}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.givenResourcesInserted(tc.input)

			cbsd, err := s.cbsdManager.EnodebdUpdateCbsd(tc.toUpdate)
			s.Require().NoError(err)
			s.Assert().Equal(tc.input.CbsdSerialNumber, cbsd.CbsdSerialNumber)
			s.Assert().Equal(tc.input.NetworkId, cbsd.NetworkId)

			err = s.resourceManager.InTransaction(func() {
				actual, err := db.NewQuery().
					WithBuilder(s.resourceManager.GetBuilder()).
					From(&storage.DBCbsd{}).
					Select(db.NewExcludeMask("id", "state_id",
						"cbsd_id", "grant_attempts", "is_deleted")).
					Where(sq.Eq{"cbsd_serial_number": tc.input.CbsdSerialNumber}).
					Fetch()
				s.Require().NoError(err)
				expected := []db.Model{tc.expected}
				s.Assert().Equal(expected, actual)
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

	expected := b.NewDetailedDBCbsdBuilder(
		b.NewDBCbsdBuilder().
			WithIndoorDeployment(false).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr)).
		WithEmptyGrant().
		WithEmptyGrantState().
		WithCbsdState(registered).
		WithDesiredState(registered).
		Details

	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestFetchCbsdWithIdleGrant() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	grantState := s.enumMaps[storage.GrantStateTable]["idle"]

	s.givenResourcesInserted(
		b.NewDBCbsdBuilder().
			WithNetworkId(someNetwork).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr).
			WithStateId(state).
			WithDesiredStateId(state).
			Cbsd,
		b.NewDBGrantBuilder().
			WithStateId(grantState).
			WithCbsdId(someCbsdId).
			Grant,
	)

	actual, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	expected := b.NewDetailedDBCbsdBuilder(
		b.NewDBCbsdBuilder().
			WithIndoorDeployment(false).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr)).
		WithCbsdState(registered).
		WithDesiredState(registered).
		WithEmptyGrant().
		WithEmptyGrantState().
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
			WithStateId(grantState).
			WithCbsdId(someCbsdId).
			Grant,
	)

	actual, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	expected := b.NewDetailedDBCbsdBuilder(
		b.NewDBCbsdBuilder().
			WithIndoorDeployment(false).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr)).
		WithCbsdState(registered).
		WithDesiredState(registered).
		WithGrant().
		WithGrantState(authorized).
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
		cbsdBuilder := b.NewDBCbsdBuilder().
			WithId(int64(i + 1 + offset)).
			WithIndoorDeployment(false).
			WithSerialNumber(fmt.Sprintf("some_serial_number%d", i+1+offset))
		expected.Cbsds[i] = b.NewDetailedDBCbsdBuilder(cbsdBuilder).
			WithCbsdState(unregistered).
			WithDesiredState(unregistered).
			WithEmptyGrant().
			WithEmptyGrantState().
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
		cbsdBuilder := b.NewDBCbsdBuilder().
			WithId(int64(i + 1)).
			WithIndoorDeployment(false).
			WithSerialNumber(fmt.Sprintf("some_serial_number%d", i+1))
		expected.Cbsds[i] = b.NewDetailedDBCbsdBuilder(cbsdBuilder).
			WithCbsdState(unregistered).
			WithDesiredState(unregistered).
			WithEmptyGrant().
			WithEmptyGrantState().
			Details
	}
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListNotIncludeIdleGrants() {
	state := s.enumMaps[storage.CbsdStateTable][registered]
	grantState := s.enumMaps[storage.GrantStateTable]["idle"]
	s.givenResourcesInserted(
		b.NewDBCbsdBuilder().
			WithNetworkId(someNetwork).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr).
			WithStateId(state).
			WithDesiredStateId(state).
			Cbsd,
		b.NewDBGrantBuilder().
			WithId(1).
			WithStateId(grantState).
			WithCbsdId(someCbsdId).
			Grant,
		b.NewDBGrantBuilder().
			WithId(2).
			WithStateId(grantState).
			WithCbsdId(someCbsdId).
			Grant,
	)

	actual, err := s.cbsdManager.ListCbsd(someNetwork, &storage.Pagination{}, nil)
	s.Require().NoError(err)

	builder := b.NewDetailedDBCbsdBuilder(
		b.NewDBCbsdBuilder().
			WithIndoorDeployment(false).
			WithId(someCbsdId).
			WithCbsdId(someCbsdIdStr)).
		WithEmptyGrant().
		WithEmptyGrantState().
		WithCbsdState(registered).
		WithDesiredState(registered)

	expected := b.GetDetailedDBCbsdList(builder)

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
