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

	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/dp/cloud/go/services/dp/storage/dbtest"
	"magma/orc8r/cloud/go/sqorc"
	"magma/orc8r/lib/go/merrors"
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
	database, err := sqorc.Open("sqlite3", ":memory:")
	s.Require().NoError(err)
	s.cbsdManager = storage.NewCbsdManager(database, builder, errorChecker)
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
		&storage.DBCbsdState{Name: db.MakeString("unregistered")},
		&storage.DBCbsdState{Name: db.MakeString("registered")},
		&storage.DBGrantState{Name: db.MakeString("idle")},
		&storage.DBGrantState{Name: db.MakeString("granted")},
		&storage.DBGrantState{Name: db.MakeString("authorized")},
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

const (
	someNetwork  = "some_network"
	otherNetwork = "other_network_id"

	someCbsdId = 123
)

func (s *CbsdManagerTestSuite) TestCreateCbsd() {
	err := s.cbsdManager.CreateCbsd(someNetwork, getMutableCbsd())
	s.Require().NoError(err)

	err = s.resourceManager.InTransaction(func() {
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

		cbsd := getBaseCbsd()
		cbsd.NetworkId = db.MakeString(someNetwork)
		cbsd.GrantAttempts = db.MakeInt(0)
		cbsd.IsDeleted = db.MakeBool(false)
		cbsd.ShouldDeregister = db.MakeBool(false)
		expected := []db.Model{
			cbsd,
			&storage.DBCbsdState{Name: db.MakeString("unregistered")},
			&storage.DBCbsdState{Name: db.MakeString("registered")},
		}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) TestCreateCbsdWithExistingSerialNumber() {
	err := s.cbsdManager.CreateCbsd(someNetwork, getMutableCbsd())
	s.Require().NoError(err)
	err = s.cbsdManager.CreateCbsd(someNetwork, getMutableCbsd())
	s.Assert().ErrorIs(err, merrors.ErrAlreadyExists)
}

func (s *CbsdManagerTestSuite) TestUpdateCbsdWithSerialNumberOfExistingCbsd() {
	state := s.enumMaps[storage.CbsdStateTable]["registered"]
	cbsd1 := getCbsd(1, someNetwork, state)
	cbsd1.CbsdSerialNumber = db.MakeString("some_serial_number_1")
	cbsd2 := getCbsd(2, someNetwork, state)
	cbsd2.CbsdSerialNumber = db.MakeString("some_serial_number_2")
	s.givenResourcesInserted(cbsd1, cbsd2)

	cbsd2.CbsdSerialNumber = cbsd1.CbsdSerialNumber
	m := &storage.MutableCbsd{
		Cbsd:         cbsd2,
		DesiredState: &storage.DBCbsdState{Name: db.MakeString("registered")},
	}
	err := s.cbsdManager.UpdateCbsd(someNetwork, cbsd2.Id.Int64, m)
	s.Assert().ErrorIs(err, merrors.ErrAlreadyExists)
}

func (s *CbsdManagerTestSuite) givenResourcesInserted(models ...db.Model) {
	err := s.resourceManager.InsertResources(db.NewExcludeMask(), models...)
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) TestUpdateCbsd() {
	state := s.enumMaps[storage.CbsdStateTable]["registered"]
	s.givenResourcesInserted(getCbsd(someCbsdId, someNetwork, state))

	m := getMutableCbsd()
	m.Cbsd.UserId.String += "new1"
	m.Cbsd.FccId.String += "new2"
	m.Cbsd.CbsdSerialNumber.String += "new3"
	m.Cbsd.AntennaGain.Float64 += 1
	m.Cbsd.MaxPower.Float64 += 2
	m.Cbsd.MinPower.Float64 += 3
	m.Cbsd.NumberOfPorts.Int64 += 4
	m.DesiredState.Name = db.MakeString("unregistered")
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
		m.Cbsd.NetworkId = db.MakeString(someNetwork)
		m.Cbsd.ShouldDeregister = db.MakeBool(true)
		m.Cbsd.DesiredStateId = db.MakeInt(s.enumMaps[storage.CbsdStateTable]["unregistered"])
		expected := []db.Model{m.Cbsd}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *CbsdManagerTestSuite) TestUpdateDeletedCbsd() {
	s.givenDeletedCbsd()

	err := s.cbsdManager.UpdateCbsd(someNetwork, someCbsdId, getMutableCbsd())
	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

func (s *CbsdManagerTestSuite) TestUpdateNonExistentCbsd() {
	err := s.cbsdManager.UpdateCbsd(someNetwork, 0, getMutableCbsd())
	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

func (s *CbsdManagerTestSuite) TestDeleteCbsd() {
	state := s.enumMaps[storage.CbsdStateTable]["registered"]
	s.givenResourcesInserted(getCbsd(someCbsdId, someNetwork, state))

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
	state := s.enumMaps[storage.CbsdStateTable]["registered"]
	s.givenResourcesInserted(getCbsd(someCbsdId, otherNetwork, state))

	_, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)

	s.Assert().ErrorIs(err, merrors.ErrNotFound)
}

func (s *CbsdManagerTestSuite) TestFetchCbsdWithoutGrant() {
	state := s.enumMaps[storage.CbsdStateTable]["registered"]
	s.givenResourcesInserted(getCbsd(someCbsdId, someNetwork, state))

	actual, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	expected := &storage.DetailedCbsd{
		Cbsd:         getDetailedCbsd(someCbsdId),
		CbsdState:    &storage.DBCbsdState{Name: db.MakeString("registered")},
		DesiredState: &storage.DBCbsdState{Name: db.MakeString("registered")},
		Grant:        &storage.DBGrant{},
		GrantState:   &storage.DBGrantState{},
	}
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestFetchCbsdWithIdleGrant() {
	state := s.enumMaps[storage.CbsdStateTable]["registered"]
	grantState := s.enumMaps[storage.GrantStateTable]["idle"]
	s.givenResourcesInserted(
		getCbsd(someCbsdId, someNetwork, state),
		getGrant(1, grantState, someCbsdId),
	)

	actual, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	expected := &storage.DetailedCbsd{
		Cbsd:         getDetailedCbsd(someCbsdId),
		CbsdState:    &storage.DBCbsdState{Name: db.MakeString("registered")},
		DesiredState: &storage.DBCbsdState{Name: db.MakeString("registered")},
		Grant:        &storage.DBGrant{},
		GrantState:   &storage.DBGrantState{},
	}
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestFetchCbsdWithGrant() {
	state := s.enumMaps[storage.CbsdStateTable]["registered"]
	grantState := s.enumMaps[storage.GrantStateTable]["authorized"]
	s.givenResourcesInserted(
		getCbsd(someCbsdId, someNetwork, state),
		getGrant(1, grantState, someCbsdId),
	)

	actual, err := s.cbsdManager.FetchCbsd(someNetwork, someCbsdId)
	s.Require().NoError(err)

	expected := &storage.DetailedCbsd{
		Cbsd:         getDetailedCbsd(someCbsdId),
		CbsdState:    &storage.DBCbsdState{Name: db.MakeString("registered")},
		DesiredState: &storage.DBCbsdState{Name: db.MakeString("registered")},
		Grant:        getBaseGrant(),
		GrantState:   &storage.DBGrantState{Name: db.MakeString("authorized")},
	}
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListCbsdFromDifferentNetwork() {
	state := s.enumMaps[storage.CbsdStateTable]["registered"]
	s.givenResourcesInserted(getCbsd(someCbsdId, otherNetwork, state))

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
	models := make([]db.Model, count)
	stateId := s.enumMaps[storage.CbsdStateTable]["unregistered"]
	for i := range models {
		cbsd := getCbsd(int64(i+1), someNetwork, stateId)
		cbsd.CbsdSerialNumber = db.MakeString(fmt.Sprintf("some_serial_number%d", i+1))
		models[i] = cbsd
	}
	s.givenResourcesInserted(models...)

	const limit = 2
	const offset = 1
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
		cbsd := getDetailedCbsd(int64(i + 1 + offset))
		cbsd.CbsdSerialNumber = db.MakeString(fmt.Sprintf("some_serial_number%d", i+1+offset))
		expected.Cbsds[i] = &storage.DetailedCbsd{
			Cbsd:         cbsd,
			CbsdState:    &storage.DBCbsdState{Name: db.MakeString("unregistered")},
			DesiredState: &storage.DBCbsdState{Name: db.MakeString("unregistered")},
			Grant:        &storage.DBGrant{},
			GrantState:   &storage.DBGrantState{},
		}
	}
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListWithFilter() {
	const count = 1
	models := make([]db.Model, count)
	stateId := s.enumMaps[storage.CbsdStateTable]["unregistered"]
	for i := range models {
		cbsd := getCbsd(int64(i+1), someNetwork, stateId)
		cbsd.CbsdSerialNumber = db.MakeString(fmt.Sprintf("some_serial_number%d", i+1))
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
		cbsd := getDetailedCbsd(int64(i + 1))
		cbsd.CbsdSerialNumber = db.MakeString(fmt.Sprintf("some_serial_number%d", i+1))
		expected.Cbsds[i] = &storage.DetailedCbsd{
			Cbsd:         cbsd,
			CbsdState:    &storage.DBCbsdState{Name: db.MakeString("unregistered")},
			DesiredState: &storage.DBCbsdState{Name: db.MakeString("unregistered")},
			Grant:        &storage.DBGrant{},
			GrantState:   &storage.DBGrantState{},
		}
	}
	s.Assert().Equal(expected, actual)
}

func (s *CbsdManagerTestSuite) TestListNotIncludeIdleGrants() {
	state := s.enumMaps[storage.CbsdStateTable]["registered"]
	grantState := s.enumMaps[storage.GrantStateTable]["idle"]
	s.givenResourcesInserted(
		getCbsd(someCbsdId, someNetwork, state),
		getGrant(1, grantState, someCbsdId),
		getGrant(2, grantState, someCbsdId),
	)

	actual, err := s.cbsdManager.ListCbsd(someNetwork, &storage.Pagination{}, nil)
	s.Require().NoError(err)

	expected := &storage.DetailedCbsdList{
		Cbsds: []*storage.DetailedCbsd{{
			Cbsd:         getDetailedCbsd(someCbsdId),
			CbsdState:    &storage.DBCbsdState{Name: db.MakeString("registered")},
			DesiredState: &storage.DBCbsdState{Name: db.MakeString("registered")},
			Grant:        &storage.DBGrant{},
			GrantState:   &storage.DBGrantState{},
		}},
		Count: 1,
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

func (s *CbsdManagerTestSuite) givenDeletedCbsd() {
	state := s.enumMaps[storage.CbsdStateTable]["registered"]
	cbsd := getCbsd(someCbsdId, someNetwork, state)
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

func getBaseGrant() *storage.DBGrant {
	base := &storage.DBGrant{}
	base.GrantExpireTime = db.MakeTime(time.Unix(123, 0).UTC())
	base.TransmitExpireTime = db.MakeTime(time.Unix(456, 0).UTC())
	base.LowFrequency = db.MakeInt(3600 * 1e6)
	base.HighFrequency = db.MakeInt(3620 * 1e6)
	base.MaxEirp = db.MakeFloat(35)
	return base
}

func getGrant(id int64, stateId int64, cbsdId int64) *storage.DBGrant {
	base := getBaseGrant()
	base.Id = db.MakeInt(id)
	base.CbsdId = db.MakeInt(cbsdId)
	base.StateId = db.MakeInt(stateId)
	base.GrantId = db.MakeString("some_grant_id")
	return base
}

func getCbsd(id int64, networkId string, stateId int64) *storage.DBCbsd {
	base := getDetailedCbsd(id)
	base.NetworkId = db.MakeString(networkId)
	base.StateId = db.MakeInt(stateId)
	base.DesiredStateId = db.MakeInt(stateId)
	base.CbsdId = db.MakeString("some_cbsd_id")
	base.ShouldDeregister = db.MakeBool(false)
	base.IsDeleted = db.MakeBool(false)
	base.GrantAttempts = db.MakeInt(0)
	return base
}

func getDetailedCbsd(id int64) *storage.DBCbsd {
	base := getBaseCbsd()
	base.Id = db.MakeInt(id)
	base.CbsdId = db.MakeString("some_cbsd_id")
	return base
}

func getBaseCbsd() *storage.DBCbsd {
	base := &storage.DBCbsd{}
	base.UserId = db.MakeString("some_user_id")
	base.FccId = db.MakeString("some_fcc_id")
	base.CbsdSerialNumber = db.MakeString("some_serial_number")
	base.PreferredBandwidthMHz = db.MakeInt(20)
	base.PreferredFrequenciesMHz = db.MakeString("[3600]")
	base.MinPower = db.MakeFloat(10)
	base.MaxPower = db.MakeFloat(20)
	base.AntennaGain = db.MakeFloat(15)
	base.NumberOfPorts = db.MakeInt(2)
	return base
}

func getMutableCbsd() *storage.MutableCbsd {
	return &storage.MutableCbsd{
		Cbsd:         getBaseCbsd(),
		DesiredState: &storage.DBCbsdState{Name: db.MakeString("registered")},
	}
}
