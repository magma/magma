package storage

import (
	"database/sql"
	"errors"
	"testing"

	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/dp/cloud/go/services/dp/storage/dbtest"
	"magma/orc8r/cloud/go/sqorc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestAmcManager(t *testing.T) {
	suite.Run(t, &AmcManagerTestSuite{})
}

type AmcManagerTestSuite struct {
	suite.Suite
	amcManager      AmcManager
	database        *sql.DB
	resourceManager dbtest.ResourceManager
	enumMaps        map[string]map[string]int64
}

func (s *AmcManagerTestSuite) SetupSuite() {
	database, err := sqorc.Open("sqlite3", ":memory:")
	s.Require().NoError(err)
	s.database = database

	builder := sqorc.GetSqlBuilder()
	errorChecker := sqorc.SQLiteErrorChecker{}
	locker := sqorc.GetSqlLocker()
	s.amcManager = NewAmcManager(database, builder, errorChecker, locker)

	s.resourceManager = dbtest.NewResourceManager(s.T(), database, builder)
	err = s.resourceManager.CreateTables(
		&DBRequest{},
		&DBRequestType{},
	)
	s.Require().NoError(err)

	err = s.resourceManager.InsertResources(
		db.NewExcludeMask("id"),
		&DBRequestType{Name: db.MakeString("request type")},
	)
	s.Require().NoError(err)
	s.enumMaps = map[string]map[string]int64{}
	for _, model := range []db.Model{
		&DBRequestType{},
	} {
		table := model.GetMetadata().Table
		s.enumMaps[table] = s.getNameIdMapping(model)
	}
}

func (s *AmcManagerTestSuite) TestCreateRequest() {
	request := MutableRequest{
		Request: &DBRequest{
			CbsdId:  db.MakeInt(1),
			Payload: "some payload",
		},
		DesiredTypeId: &DBRequestType{
			Name: db.MakeString("request type"),
		},
	}

	_, err := WithinTx(s.database, func(tx *sql.Tx) (interface{}, error) {
		return nil, s.amcManager.CreateRequest(tx, &request)
	})
	s.Require().NoError(err)

	err = s.resourceManager.InTransaction(func() {
		actual, err := db.NewQuery().
			WithBuilder(s.resourceManager.GetBuilder()).
			From(&DBRequest{}).
			Select(db.NewIncludeMask("type_id", "cbsd_id", "payload")).
			Fetch()
		s.Require().NoError(err)

		expected := []db.Model{&DBRequest{
			CbsdId:  db.MakeInt(1),
			TypeId:  db.MakeInt(1),
			Payload: "some payload",
		}}
		s.Assert().Equal(expected, actual)
	})
	s.Require().NoError(err)
}

func (s *AmcManagerTestSuite) TearDownTest() {
	err := s.resourceManager.DropResources(
		&DBRequest{},
	)
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
		enum := r[0].(EnumModel)
		m[enum.GetName()] = enum.GetId()
	}
	return m
}

func TestWithinTx(t *testing.T) {
	type testCase struct {
		prepareMockFunc func(sqlmock.Sqlmock)
		wrappedFunc     func(*sql.Tx) (any, error)
		resultCheckFunc func(any, error)
	}
	testCases := []testCase{
		{ // test working insert
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
		},

		{ // test wrapped func error
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
		},

		{ // test commit error
			prepareMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO table").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
				mock.ExpectClose()
			},
			wrappedFunc: func(tx *sql.Tx) (any, error) {
				return tx.Exec("INSERT INTO table (\"field\") VALUES (1);")
			},
			resultCheckFunc: func(res any, err error) {
				assert.Equal(t, nil, res)
				assert.Errorf(t, err, "commit error")
			},
		},

		{ // test begin transaction error
			prepareMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
				mock.ExpectRollback()
				mock.ExpectClose()
			},
			wrappedFunc: func(tx *sql.Tx) (any, error) {
				return nil, nil
			},
			resultCheckFunc: func(res any, err error) {
				assert.Equal(t, nil, res)
				assert.Errorf(t, err, "begin error")
			},
		},

		{ // test transaction rollback error
			prepareMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("begin error"))
				mock.ExpectRollback().WillReturnError(errors.New("rollback error"))
				mock.ExpectClose()
			},
			wrappedFunc: func(tx *sql.Tx) (any, error) {
				return nil, nil
			},
			resultCheckFunc: func(res any, err error) {
				assert.Equal(t, nil, res)
				assert.Errorf(t, err, "rollback error")
			},
		},

		{ // test wrapped func panic
			prepareMockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback()
				mock.ExpectClose()
			},
			wrappedFunc: func(tx *sql.Tx) (any, error) {
				panic("I am panicking")
			},
			resultCheckFunc: func(res any, err error) {},
		},
	}

	database, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer func() {
		err = database.Close()
		assert.NoError(t, err)
	}()

	for _, tc := range testCases {
		tc.prepareMockFunc(mock)
		res, err := WithinTx(database, tc.wrappedFunc)
		tc.resultCheckFunc(res, err)
	}
}
