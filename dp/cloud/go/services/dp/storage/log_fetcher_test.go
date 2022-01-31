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
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/dp/cloud/go/services/dp/storage/dbtest"
	"magma/orc8r/cloud/go/sqorc"
)

func TestLogFetcher(t *testing.T) {
	suite.Run(t, &LogFetcherTestSuite{})
}

type LogFetcherTestSuite struct {
	suite.Suite
	logFetcher      storage.LogFetcher
	logs            []*storage.DBLog
	resourceManager dbtest.ResourceManager
}

const (
	someSource       = "some_source"
	otherSource      = "other_source"
	someName         = "some_name"
	otherName        = "other_name"
	someMessage      = "some_message"
	someSerial       = "some_serial"
	otherSerial      = "other_serial"
	someFcc          = "some_fcc"
	otherFcc         = "other_fcc"
	someResponseCode = 0
)

func (s *LogFetcherTestSuite) SetupSuite() {
	builder := sqorc.GetSqlBuilder()
	database, err := sqorc.Open("sqlite3", ":memory:")
	s.Require().NoError(err)
	s.logFetcher = storage.NewLogFetcher(database, builder)
	s.resourceManager = dbtest.NewResourceManager(s.T(), database, builder)

	err = s.resourceManager.CreateTables(&storage.DBLog{})
	s.Require().NoError(err)

	s.logs = []*storage.DBLog{{
		From:         db.MakeString(someSource),
		To:           db.MakeString(otherSource),
		Name:         db.MakeString(someName),
		Message:      db.MakeString(someMessage),
		SerialNumber: db.MakeString(someSerial),
		FccId:        db.MakeString(someFcc),
		CreatedDate:  db.MakeTime(getTime(0)),
	}, {
		From:         db.MakeString(someSource),
		To:           db.MakeString(otherSource),
		Name:         db.MakeString(someName),
		Message:      db.MakeString(someMessage),
		SerialNumber: db.MakeString(otherSerial),
		FccId:        db.MakeString(otherFcc),
		CreatedDate:  db.MakeTime(getTime(1)),
	}, {
		From:         db.MakeString(otherSource),
		To:           db.MakeString(someSource),
		Name:         db.MakeString(otherName),
		Message:      db.MakeString(someMessage),
		SerialNumber: db.MakeString(someSerial),
		FccId:        db.MakeString(someFcc),
		ResponseCode: db.MakeInt(someResponseCode),
		CreatedDate:  db.MakeTime(getTime(2)),
	}, {
		From:         db.MakeString(otherSource),
		To:           db.MakeString(someSource),
		Name:         db.MakeString(otherName),
		Message:      db.MakeString(someMessage),
		SerialNumber: db.MakeString(otherSerial),
		FccId:        db.MakeString(otherFcc),
		ResponseCode: db.MakeInt(someResponseCode),
		CreatedDate:  db.MakeTime(getTime(3)),
	}}
	models := make([]db.Model, len(s.logs))
	for i, log := range s.logs {
		models[i] = withNetworkId(someNetwork, log)
	}
	err = s.resourceManager.InsertResources(db.NewExcludeMask("id"), models...)
	s.Require().NoError(err)
}

func (s *LogFetcherTestSuite) TestFetch() {
	const format = "2006-01-02 15:04:05.999+00:00"
	testCases := []struct {
		name       string
		filter     storage.LogFilter
		pagination storage.Pagination
		expected   []int
	}{
		{
			name:     "Should list all logs sorted by date newest first",
			expected: []int{3, 2, 1, 0},
		},
		{
			name: "Should list all logs with limit",
			pagination: storage.Pagination{
				Limit: db.MakeInt(2),
			},
			expected: []int{3, 2},
		},
		{
			name: "Should list all logs with limit and offset",
			pagination: storage.Pagination{
				Limit:  db.MakeInt(2),
				Offset: db.MakeInt(1),
			},
			expected: []int{2, 1},
		},
		{
			name:     "Should filter by from",
			filter:   storage.LogFilter{From: otherSource},
			expected: []int{3, 2},
		},
		{
			name:     "Should filter by to",
			filter:   storage.LogFilter{To: otherSource},
			expected: []int{1, 0},
		},
		{
			name:     "Should filter by fcc id",
			filter:   storage.LogFilter{FccId: someFcc},
			expected: []int{2, 0},
		},
		{
			name:     "Should filter by serial number",
			filter:   storage.LogFilter{SerialNumber: otherSerial},
			expected: []int{3, 1},
		},
		{
			name:     "Should filter by name",
			filter:   storage.LogFilter{Name: someName},
			expected: []int{1, 0},
		},
		{
			name:     "Should filter by response code",
			filter:   storage.LogFilter{ResponseCode: db.MakeInt(someResponseCode)},
			expected: []int{3, 2},
		},
		{
			name:     "Should filter by begin date",
			filter:   storage.LogFilter{BeginTimestamp: getTime(1).Format(format)},
			expected: []int{3, 2, 1},
		},
		{
			name:     "Should filter by end date",
			filter:   storage.LogFilter{EndTimestamp: getTime(2).Format(format)},
			expected: []int{2, 1, 0},
		},
		{
			name: "Should combine all filters",
			filter: storage.LogFilter{
				FccId:          someFcc,
				SerialNumber:   someSerial,
				Name:           someName,
				BeginTimestamp: getTime(0).Format(format),
				EndTimestamp:   getTime(3).Format(format),
			},
			expected: []int{0},
		},
		{
			name: "Should return empty list when no matching logs",
			filter: storage.LogFilter{
				BeginTimestamp: getTime(2).Format(format),
				EndTimestamp:   getTime(1).Format(format),
			},
			expected: nil,
		},
	}
	for _, tt := range testCases {
		s.Run(tt.name, func() {
			actual, err := s.logFetcher.ListLogs(someNetwork, &tt.filter, &tt.pagination)
			s.Require().NoError(err)
			expected := s.getExpectedLogs(tt.expected)
			s.Assert().Equal(expected, actual)
		})
	}
}

func getTime(index int64) time.Time {
	return time.Unix(index*100, index*int64(time.Millisecond)).UTC()
}

func withNetworkId(networkId string, model *storage.DBLog) *storage.DBLog {
	newModel := *model
	newModel.NetworkId = db.MakeString(networkId)
	return &newModel
}

func (s *LogFetcherTestSuite) getExpectedLogs(indices []int) []*storage.DBLog {
	logs := make([]*storage.DBLog, 0, len(indices))
	for _, i := range indices {
		logs = append(logs, s.logs[i])
	}
	return logs
}
