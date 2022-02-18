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
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/servicers"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

func TestLogFetcher(t *testing.T) {
	suite.Run(t, &LogFetcherTestSuite{})
}

type LogFetcherTestSuite struct {
	suite.Suite
	fetcher protos.LogFetcherServer
	store   *stubLogFetcher
}

func (s *LogFetcherTestSuite) SetupTest() {
	s.store = &stubLogFetcher{}
	s.fetcher = servicers.NewLogFetcher(s.store)
}

func (s *LogFetcherTestSuite) TestLogs() {
	s.store.logs = []*storage.DBLog{{
		From:         db.MakeString("some from"),
		To:           db.MakeString("some to"),
		Name:         db.MakeString("some name"),
		Message:      db.MakeString("some message"),
		SerialNumber: db.MakeString("some serial"),
		FccId:        db.MakeString("some fcc id"),
		ResponseCode: db.MakeInt(0),
		CreatedDate:  db.MakeTime(time.Unix(1e6, 123*int64(time.Millisecond))),
	}}
	request := &protos.ListLogsRequest{
		NetworkId:  networkId,
		Filter:     &protos.LogFilter{},
		Pagination: &protos.Pagination{},
	}
	logs, err := s.fetcher.ListLogs(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(&storage.LogFilter{}, s.store.filter)
	s.Assert().Equal(&storage.Pagination{}, s.store.pagination)
	expectedLogs := []*protos.Log{{
		From:           "some from",
		To:             "some to",
		Name:           "some name",
		Message:        "some message",
		SerialNumber:   "some serial",
		FccId:          "some fcc id",
		TimestampMilli: 1e9 + 123,
	}}
	s.Assert().Equal(expectedLogs, logs.Logs)
}

func (s *LogFetcherTestSuite) TestListLogsWithFilterAndPagination() {
	const format = "2006-01-02 15:04:05.999+00:00"
	request := &protos.ListLogsRequest{
		NetworkId: networkId,
		Filter: &protos.LogFilter{
			From:                "some from",
			To:                  "some to",
			Name:                "some name",
			SerialNumber:        "some serial",
			FccId:               "some fcc id",
			ResponseCode:        wrapperspb.Int64(0),
			BeginTimestampMilli: wrapperspb.Int64(1e9 + 123),
			EndTimestampMilli:   wrapperspb.Int64(2e9 + 456),
		},
		Pagination: &protos.Pagination{
			Limit:  wrapperspb.Int64(10),
			Offset: wrapperspb.Int64(20),
		},
	}
	_, err := s.fetcher.ListLogs(context.Background(), request)
	s.Require().NoError(err)

	s.Assert().Equal(networkId, s.store.networkId)
	s.Assert().Equal(&storage.LogFilter{
		From:           "some from",
		To:             "some to",
		FccId:          "some fcc id",
		SerialNumber:   "some serial",
		Name:           "some name",
		ResponseCode:   db.MakeInt(0),
		BeginTimestamp: time.Unix(1e6, 123*int64(time.Millisecond)).UTC().Format(format),
		EndTimestamp:   time.Unix(2e6, 456*int64(time.Millisecond)).UTC().Format(format),
	}, s.store.filter)
	s.Assert().Equal(&storage.Pagination{
		Limit:  db.MakeInt(10),
		Offset: db.MakeInt(20),
	}, s.store.pagination)
}

func (s *LogFetcherTestSuite) TestListLogsWithError() {
	s.store.err = errors.New("some error")
	request := &protos.ListLogsRequest{
		NetworkId:  networkId,
		Filter:     &protos.LogFilter{},
		Pagination: &protos.Pagination{},
	}
	_, err := s.fetcher.ListLogs(context.Background(), request)
	s.Assert().Error(err)
}

type stubLogFetcher struct {
	networkId  string
	filter     *storage.LogFilter
	pagination *storage.Pagination
	logs       []*storage.DBLog
	err        error
}

func (s *stubLogFetcher) ListLogs(networkId string, filter *storage.LogFilter, pagination *storage.Pagination) ([]*storage.DBLog, error) {
	s.networkId = networkId
	s.filter = filter
	s.pagination = pagination
	return s.logs, s.err
}
