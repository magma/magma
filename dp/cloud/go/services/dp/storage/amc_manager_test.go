package storage_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	b "magma/dp/cloud/go/services/dp/builders"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/dp/cloud/go/services/dp/storage/dbtest"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/sqorc"
)

const (
	requestPayload = "some payload"
	someGrantId    = "some_grant_id"
)

func TestAmcManager(t *testing.T) {
	suite.Run(t, &AmcManagerTestSuite{})
}

type AmcManagerTestSuite struct {
	suite.Suite
	database        *sql.DB
	amcManager      storage.AmcManager
	resourceManager dbtest.ResourceManager
	enumMaps        map[string]map[string]int64
}

func (s *AmcManagerTestSuite) SetupSuite() {
	builder := sqorc.GetSqlBuilder()
	errorChecker := sqorc.SQLiteErrorChecker{}
	locker := sqorc.GetSqlLocker()
	database, err := sqorc.Open("sqlite3", ":memory:")
	s.Require().NoError(err)
	s.database = database
	s.amcManager = storage.NewAmcManager(database, builder, errorChecker, locker)

	s.resourceManager = dbtest.NewResourceManager(s.T(), s.database, builder)
	err = s.resourceManager.CreateTables(
		&storage.DBCbsdState{},
		&storage.DBCbsd{},
		&storage.DBGrantState{},
		&storage.DBGrant{},
		&storage.DBRequest{},
		&storage.DBRequestType{},
	)
	s.Require().NoError(err)
	err = s.resourceManager.InsertResources(
		db.NewExcludeMask("id"),
		&storage.DBCbsdState{Name: db.MakeString(unregistered)},
		&storage.DBCbsdState{Name: db.MakeString(registered)},
		&storage.DBGrantState{Name: db.MakeString(idle)},
		&storage.DBGrantState{Name: db.MakeString(granted)},
		&storage.DBGrantState{Name: db.MakeString(authorized)},
		&storage.DBRequestType{Name: db.MakeString(grant)},
	)
	s.Require().NoError(err)
	s.enumMaps = map[string]map[string]int64{}
	for _, model := range []db.Model{
		&storage.DBCbsdState{},
		&storage.DBGrantState{},
		&storage.DBRequestType{},
	} {
		table := model.GetMetadata().Table
		s.enumMaps[table] = s.getNameIdMapping(model)
	}
}

func (s *AmcManagerTestSuite) TearDownTest() {
	clock.UnfreezeClock(s.T())
	err := s.resourceManager.DropResources(
		&storage.DBCbsd{},
		&storage.DBGrant{},
		&storage.DBRequest{},
	)
	s.Require().NoError(err)
}

func (s *AmcManagerTestSuite) TestCreateRequest() {
	request := storage.MutableRequest{
		Request: &storage.DBRequest{
			CbsdId:  db.MakeInt(1),
			Payload: requestPayload,
		},
		RequestType: &storage.DBRequestType{
			Name: db.MakeString(grant),
		},
	}

	_, err := storage.WithinTx(s.database, func(tx *sql.Tx) (interface{}, error) {
		return nil, s.amcManager.CreateRequest(tx, &request)
	})
	s.Require().NoError(err)

	err = s.resourceManager.InTransaction(func() {
		actual, err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&storage.DBRequest{}).
			Select(db.NewIncludeMask("type_id", "cbsd_id", "payload")).
			Fetch()
		s.Require().NoError(err)

		expected := []db.Model{&storage.DBRequest{
			CbsdId:  db.MakeInt(1),
			TypeId:  db.MakeInt(1),
			Payload: requestPayload,
		}}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *AmcManagerTestSuite) TestDeleteCbsd() {
	stateId := s.enumMaps[storage.CbsdStateTable][unregistered]
	cbsd1 := storage.DBCbsd{
		Id:                    db.MakeInt(1),
		NetworkId:             db.MakeString(someNetwork),
		StateId:               db.MakeInt(stateId),
		DesiredStateId:        db.MakeInt(stateId),
		PreferredBandwidthMHz: db.MakeInt(20),
	}
	cbsd2 := storage.DBCbsd{
		Id:                    db.MakeInt(2),
		NetworkId:             db.MakeString(someNetwork),
		StateId:               db.MakeInt(stateId),
		DesiredStateId:        db.MakeInt(stateId),
		PreferredBandwidthMHz: db.MakeInt(20),
	}
	err := s.resourceManager.InsertResources(db.NewExcludeMask(), &cbsd1, &cbsd2)
	s.Require().NoError(err)

	_, err = storage.WithinTx(s.database, func(tx *sql.Tx) (interface{}, error) {
		return nil, s.amcManager.DeleteCbsd(tx, &cbsd1)
	})
	s.Require().NoError(err)

	// only cbsd2 should exist
	err = s.resourceManager.InTransaction(func() {
		actual, err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&storage.DBCbsd{}).
			Select(db.NewIncludeMask("id")).
			Fetch()
		s.Require().NoError(err)

		expected := []db.Model{&storage.DBCbsd{Id: db.MakeInt(2)}}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)

	// delete on not existing cbsd should not return error
	_, err = storage.WithinTx(s.database, func(tx *sql.Tx) (interface{}, error) {
		return nil, s.amcManager.DeleteCbsd(tx, &cbsd1)
	})
	s.Require().NoError(err)
}

func (s *AmcManagerTestSuite) TestUpdateCbsd() {
	stateId := s.enumMaps[storage.CbsdStateTable][unregistered]
	cbsd := storage.DBCbsd{
		Id:                    db.MakeInt(1),
		NetworkId:             db.MakeString(someNetwork),
		StateId:               db.MakeInt(stateId),
		DesiredStateId:        db.MakeInt(stateId),
		PreferredBandwidthMHz: db.MakeInt(20),
	}
	err := s.resourceManager.InsertResources(db.NewExcludeMask(), &cbsd)
	s.Require().NoError(err)

	cbsdUpdate := storage.DBCbsd{
		Id:                    db.MakeInt(1),
		PreferredBandwidthMHz: db.MakeInt(30),
		MinPower:              db.MakeFloat(0),
		MaxPower:              db.MakeFloat(20),
	}
	mask := db.NewIncludeMask("preferred_bandwidth_mhz", "min_power", "max_power")
	_, err = storage.WithinTx(s.database, func(tx *sql.Tx) (interface{}, error) {
		return nil, s.amcManager.UpdateCbsd(tx, &cbsdUpdate, mask)
	})
	s.Require().NoError(err)

	err = s.resourceManager.InTransaction(func() {
		actual, err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&storage.DBCbsd{}).
			Select(db.NewIncludeMask("id", "preferred_bandwidth_mhz", "min_power", "max_power")).
			Fetch()
		s.Require().NoError(err)

		expected := []db.Model{&cbsdUpdate}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func TestWithinTx(t *testing.T) {
	testData := []struct {
		name            string
		prepareMockFunc func(sqlmock.Sqlmock)
		wrappedFunc     func(*sql.Tx) (any, error)
		resultCheckFunc func(any, error)
	}{{
		name: "test working insert",
		prepareMockFunc: func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO table").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			mock.ExpectClose()
		},
		wrappedFunc: func(tx *sql.Tx) (any, error) {
			res, _ := tx.Exec("INSERT INTO table (\"field\") VALUES (1);")
			lastId, err := res.LastInsertId()
			return lastId, err
		},
		resultCheckFunc: func(res any, err error) {
			assert.Equal(t, int64(1), res)
			assert.NoError(t, err)
		},
	}, {
		name: "test wrapped func error",
		prepareMockFunc: func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO table").WillReturnError(errors.New("exec error"))
			mock.ExpectRollback()
			mock.ExpectClose()
		},
		wrappedFunc: func(tx *sql.Tx) (any, error) {
			res, err := tx.Exec("INSERT INTO table (\"field\") VALUES (1);")
			return res, err
		},
		resultCheckFunc: func(res any, err error) {
			assert.Equal(t, nil, res)
			assert.Errorf(t, err, "exec error")
		},
	}, {
		name: "test commit error",
		prepareMockFunc: func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO table").WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			mock.ExpectClose()
		},
		wrappedFunc: func(tx *sql.Tx) (any, error) {
			res, _ := tx.Exec("INSERT INTO table (\"field\") VALUES (1);")
			lastId, err := res.LastInsertId()
			return lastId, err
		},
		resultCheckFunc: func(res any, err error) {
			assert.Equal(t, int64(1), res)
			assert.Errorf(t, err, "commit error")
		},
	}, {
		name: "test begin transaction error",
		prepareMockFunc: func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin().WillReturnError(errors.New("begin error"))
			mock.ExpectClose()
		},
		wrappedFunc: func(tx *sql.Tx) (any, error) {
			return nil, nil
		},
		resultCheckFunc: func(res any, err error) {
			assert.Equal(t, nil, res)
			assert.Errorf(t, err, "begin error")
		},
	}, {
		name: "test transaction rollback error",
		prepareMockFunc: func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectRollback().WillReturnError(errors.New("rollback error"))
			mock.ExpectClose()
		},
		wrappedFunc: func(tx *sql.Tx) (any, error) {
			return nil, errors.New("an error")
		},
		resultCheckFunc: func(res any, err error) {
			assert.Equal(t, nil, res)
			assert.Errorf(t, err, "rollback error")
		},
	}}

	for _, tc := range testData {
		database, mock, err := sqlmock.New()
		assert.NoError(t, err)

		t.Run(tc.name, func(t *testing.T) {
			tc.prepareMockFunc(mock)
			res, err := storage.WithinTx(database, tc.wrappedFunc)
			tc.resultCheckFunc(res, err)
		})

		err = database.Close()
		assert.NoError(t, err)
	}
}

func (s *AmcManagerTestSuite) TestGetState() {
	registeredId := s.enumMaps[storage.CbsdStateTable][registered]
	grantedId := s.enumMaps[storage.GrantStateTable][granted]
	authorizedId := s.enumMaps[storage.GrantStateTable][authorized]
	grantReqId := s.enumMaps[storage.RequestTypeTable]["grant"]
	preferences := []uint32{0b10101100, 0b00110, 0b0100000, 0b11010}
	availableFreqs := []uint32{0b10111100, 0b010110, 0b01001011, 0b11110}
	testCases := []struct {
		name     string
		input    []db.Model
		expected []*storage.DetailedCbsd
	}{{
		name: "test get basic state",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber).
				WithId(0).
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithAntennaGain(20).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(0).
						WithSerialNumber(someSerialNumber).
						WithCbsdId(someCbsdIdStr).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						WithAntennaGain(20).
						Cbsd,
					registered,
					registered).
				Details,
		},
	}, {
		name: "test get state with frequency preferences",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber+"1").
				WithId(1).
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithCbsdCategory("a").
				WithAntennaGain(20).
				WithDesiredStateId(registeredId).
				WithPreferences(15, []int64{3600, 3580, 3620}).
				WithAvailableFrequencies(preferences).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(1).
						WithSerialNumber(someSerialNumber+"1").
						WithCbsdId(someCbsdIdStr).
						WithCbsdCategory("a").
						WithAntennaGain(20).
						WithPreferences(15, []int64{3600, 3580, 3620}).
						WithAvailableFrequencies(preferences).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						Cbsd,
					registered,
					registered).
				Details,
		},
	}, {
		name: "test get state with grants",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithId(2).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber+"2").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithCbsdCategory("a").
				WithAntennaGain(20).
				WithDesiredStateId(registeredId).
				WithPreferences(15, []int64{3600, 3580, 3620}).
				WithAvailableFrequencies(availableFreqs).
				Cbsd,
			b.NewDBGrantBuilder().
				WithDefaultTestValues().
				WithCbsdId(2).
				WithStateId(grantedId).
				WithLastHeartbeatTime(time.Unix(111, 0).UTC()).
				Grant,
			b.NewDBGrantBuilder().
				WithDefaultTestValues().
				WithCbsdId(2).
				WithStateId(authorizedId).
				WithLastHeartbeatTime(time.Unix(112, 0).UTC()).
				Grant,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(2).
						WithSerialNumber(someSerialNumber+"2").
						WithCbsdId(someCbsdIdStr).
						WithCbsdCategory("a").
						WithAntennaGain(20).
						WithPreferences(15, []int64{3600, 3580, 3620}).
						WithAvailableFrequencies(availableFreqs).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						Cbsd,
					registered,
					registered).
				WithAmcGrant(granted, 3600, time.Unix(111, 0).UTC(), someGrantId, 1).
				WithAmcGrant(authorized, 3600, time.Unix(112, 0).UTC(), someGrantId, 1).
				Details,
		},
	}, {
		name: "test get state with channels",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithId(3).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber+"3").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithCbsdCategory("a").
				WithAntennaGain(20).
				WithDesiredStateId(registeredId).
				WithPreferences(15, []int64{3600, 3580, 3620}).
				WithAvailableFrequencies(availableFreqs).
				WithChannels([]storage.Channel{
					{
						LowFrequencyHz:  3590,
						HighFrequencyHz: 3610,
						MaxEirp:         15,
					},
					{
						LowFrequencyHz:  3600,
						HighFrequencyHz: 3620,
					},
				}).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(3).
						WithSerialNumber(someSerialNumber+"3").
						WithCbsdId(someCbsdIdStr).
						WithCbsdCategory("a").
						WithAntennaGain(20).
						WithPreferences(15, []int64{3600, 3580, 3620}).
						WithAvailableFrequencies(availableFreqs).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						WithChannels([]storage.Channel{
							{
								LowFrequencyHz:  3590,
								HighFrequencyHz: 3610,
								MaxEirp:         15,
							},
							{
								LowFrequencyHz:  3600,
								HighFrequencyHz: 3620,
							},
						}).
						Cbsd,
					registered,
					registered).
				Details,
		},
	}, {
		name: "test get state for cbsd marked for deletion",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithId(4).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "4").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithIsDeleted(true).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(4).
						WithSerialNumber(someSerialNumber+"4").
						WithCbsdId(someCbsdIdStr).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(true).
						Cbsd,
					registered,
					registered).
				Details,
		},
	}, {
		name: "test get state for cbsd marked for update",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithId(5).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "5").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithShouldDeregister(true).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(5).
						WithSerialNumber(someSerialNumber+"5").
						WithCbsdId(someCbsdIdStr).
						WithIndoorDeployment(false).
						WithShouldDeregister(true).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						Cbsd,
					registered,
					registered).
				Details,
		},
	}, {
		name: "test get state for cbsd marked for relinquish",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithId(6).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "6").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithShouldRelinquish(true).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(6).
						WithSerialNumber(someSerialNumber+"6").
						WithCbsdId(someCbsdIdStr).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(true).
						WithIsDeleted(false).
						Cbsd,
					registered,
					registered).
				Details,
		},
	}, {
		name: "test get state for cbsd with last seen",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithId(7).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "7").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithAntennaGain(20).
				WithLastSeen(1).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(7).
						WithSerialNumber(someSerialNumber+"7").
						WithCbsdId(someCbsdIdStr).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						WithAntennaGain(20).
						WithLastSeen(1).
						Cbsd,
					registered,
					registered).
				Details,
		},
	}, {
		name: "test get state for cbsd with pending requests",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithId(8).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "8").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithAntennaGain(20).
				Cbsd,
			b.NewRequestBuilder(1, 8, grantReqId, "{'some': 'payload'}").Request,
		},
		expected: []*storage.DetailedCbsd{},
	}, {
		name: "test get state with multiple cbsds",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "9").
				WithId(9).
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithAntennaGain(20).
				Cbsd,
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "10").
				WithId(10).
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithAntennaGain(20).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(9).
						WithSerialNumber(someSerialNumber+"9").
						WithCbsdId(someCbsdIdStr).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						WithAntennaGain(20).
						Cbsd,
					registered,
					registered).
				Details,
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(10).
						WithSerialNumber(someSerialNumber+"10").
						WithCbsdId(someCbsdIdStr).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						WithAntennaGain(20).
						Cbsd,
					registered,
					registered).
				Details,
		},
	}, {
		name: "test get state with single step enabled",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "11").
				WithId(11).
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithSingleStepEnabled(true).
				WithFullInstallationParam().
				WithCbsdCategory("a").
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(11).
						WithSerialNumber(someSerialNumber+"11").
						WithCbsdId(someCbsdIdStr).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						WithSingleStepEnabled(true).
						WithFullInstallationParam().
						WithCbsdCategory("a").
						Cbsd,
					registered,
					registered).
				Details,
		},
	}, {
		name: "test get state with single step enabled",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "12").
				WithId(12).
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithSingleStepEnabled(true).
				WithFullInstallationParam().
				WithCbsdCategory("a").
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithSerialNumber(someSerialNumber+"12").
						WithId(12).
						WithCbsdId(someCbsdIdStr).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						WithSingleStepEnabled(true).
						WithFullInstallationParam().
						WithCbsdCategory("a").
						Cbsd,
					registered,
					registered).
				Details,
		},
	}, {
		name: "test not get state without registration params",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				Empty().
				WithNetworkId(someNetwork).
				WithPreferences(15, []int64{3600, 3580, 3620}).
				WithStateId(registeredId).
				WithEirpCapabilities(1, 2, 2).
				WithDesiredStateId(registeredId).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{},
	}, {
		name: "test not get state without eirp capabilities",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				Empty().
				WithFccId(someFccId).
				WithUserId(someUserId).
				WithNetworkId(someNetwork).
				WithCbsdId(someCbsdIdStr).
				WithEirpCapabilities(1, 2, 2).
				WithSerialNumber(someSerialNumber+"14").
				WithPreferences(15, []int64{3600, 3580, 3620}).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{},
	}, {
		name: "test not get state with single step enabled and no installation params",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber+"14").
				WithFccId(someFccId).
				WithUserId(someUserId).
				WithNetworkId(someNetwork).
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithEirpCapabilities(1, 2, 2).
				WithDesiredStateId(registeredId).
				WithSingleStepEnabled(true).
				WithCbsdCategory("a").
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{},
	}, {
		name: "test not get state with single step enabled, category A and outdoor",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber+"15").
				WithFccId(someFccId).
				WithUserId(someUserId).
				WithNetworkId(someNetwork).
				WithCbsdId(someCbsdIdStr).
				WithEirpCapabilities(1, 2, 2).
				WithStateId(registeredId).
				WithDesiredStateId(registeredId).
				WithSingleStepEnabled(true).
				WithCbsdCategory("a").
				WithFullInstallationParam().
				WithIndoorDeployment(false).
				Cbsd,
		},
		expected: []*storage.DetailedCbsd{},
	}, {
		name: "test get state for multiple radios with multiple grants",
		input: []db.Model{
			b.NewDBCbsdBuilder().
				WithId(15).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "15").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithAntennaGain(20).
				WithDesiredStateId(registeredId).
				Cbsd,
			b.NewDBGrantBuilder().
				WithDefaultTestValues().
				WithStateId(grantedId).
				WithCbsdId(15).
				WithLastHeartbeatTime(time.Unix(113, 0).UTC()).
				Grant,
			b.NewDBGrantBuilder().
				WithDefaultTestValues().
				WithStateId(grantedId).
				WithCbsdId(15).
				WithLastHeartbeatTime(time.Unix(114, 0).UTC()).
				Grant,
			b.NewDBCbsdBuilder().
				WithId(16).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "16").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithAntennaGain(20).
				WithDesiredStateId(registeredId).
				Cbsd,
			b.NewDBGrantBuilder().
				WithDefaultTestValues().
				WithStateId(authorizedId).
				WithCbsdId(16).
				WithLastHeartbeatTime(time.Unix(115, 0).UTC()).
				Grant,
			b.NewDBGrantBuilder().
				WithDefaultTestValues().
				WithStateId(grantedId).
				WithCbsdId(16).
				WithLastHeartbeatTime(time.Unix(116, 0).UTC()).
				Grant,
			b.NewDBGrantBuilder().
				WithDefaultTestValues().
				WithStateId(grantedId).
				WithCbsdId(16).
				WithLastHeartbeatTime(time.Unix(117, 0).UTC()).
				Grant,
			b.NewDBGrantBuilder().
				WithDefaultTestValues().
				WithStateId(grantedId).
				WithCbsdId(16).
				WithLastHeartbeatTime(time.Unix(118, 0).UTC()).
				Grant,
			b.NewDBCbsdBuilder().
				WithId(17).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "17").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithAntennaGain(20).
				WithDesiredStateId(registeredId).
				Cbsd,
			b.NewDBCbsdBuilder().
				WithId(18).
				WithNetworkId(someNetwork).
				WithSerialNumber(someSerialNumber + "18").
				WithCbsdId(someCbsdIdStr).
				WithStateId(registeredId).
				WithAntennaGain(20).
				WithDesiredStateId(registeredId).
				Cbsd,
			b.NewDBGrantBuilder().
				WithDefaultTestValues().
				WithStateId(grantedId).
				WithCbsdId(18).
				WithLastHeartbeatTime(time.Unix(119, 0).UTC()).
				Grant,
			b.NewDBGrantBuilder().
				WithDefaultTestValues().
				WithStateId(authorizedId).
				WithCbsdId(18).
				WithLastHeartbeatTime(time.Unix(120, 0).UTC()).
				Grant,
		},
		expected: []*storage.DetailedCbsd{
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(15).
						WithSerialNumber(someSerialNumber+"15").
						WithCbsdId(someCbsdIdStr).
						WithAntennaGain(20).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						Cbsd,
					registered,
					registered).
				WithAmcGrant(granted, 3600, time.Unix(113, 0).UTC(), someGrantId, 1).
				WithAmcGrant(granted, 3600, time.Unix(114, 0).UTC(), someGrantId, 1).
				Details,
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(16).
						WithSerialNumber(someSerialNumber+"16").
						WithCbsdId(someCbsdIdStr).
						WithAntennaGain(20).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						Cbsd,
					registered,
					registered).
				WithAmcGrant(authorized, 3600, time.Unix(115, 0).UTC(), someGrantId, 1).
				WithAmcGrant(granted, 3600, time.Unix(116, 0).UTC(), someGrantId, 1).
				WithAmcGrant(granted, 3600, time.Unix(117, 0).UTC(), someGrantId, 1).
				WithAmcGrant(granted, 3600, time.Unix(118, 0).UTC(), someGrantId, 1).
				Details,
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(17).
						WithSerialNumber(someSerialNumber+"17").
						WithCbsdId(someCbsdIdStr).
						WithAntennaGain(20).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						Cbsd,
					registered,
					registered).
				Details,
			b.NewDetailedDBCbsdBuilder().
				WithCbsd(
					b.NewDBCbsdBuilder().
						WithId(18).
						WithSerialNumber(someSerialNumber+"18").
						WithCbsdId(someCbsdIdStr).
						WithAntennaGain(20).
						WithIndoorDeployment(false).
						WithShouldDeregister(false).
						WithShouldRelinquish(false).
						WithIsDeleted(false).
						Cbsd,
					registered,
					registered).
				WithAmcGrant(granted, 3600, time.Unix(119, 0).UTC(), someGrantId, 1).
				WithAmcGrant(authorized, 3600, time.Unix(120, 0).UTC(), someGrantId, 1).
				Details,
		},
	}}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.givenResourcesInserted(tc.input...)

			actual, err := storage.WithinTx(s.database, func(tx *sql.Tx) (interface{}, error) {
				actual, err := s.amcManager.GetState(tx)
				return actual, err
			})

			s.Require().NoError(err)
			s.Assert().Equal(tc.expected, actual)

			err = s.resourceManager.DropResources(
				&storage.DBCbsd{},
				&storage.DBGrant{},
			)
			s.Require().NoError(err)

		})
	}
}

func (s *AmcManagerTestSuite) givenResourcesInserted(models ...db.Model) {
	err := s.resourceManager.InsertResources(db.NewExcludeMask(), models...)
	s.Require().NoError(err)
}

func (s *AmcManagerTestSuite) getNameIdMapping(model db.Model) map[string]int64 {
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
