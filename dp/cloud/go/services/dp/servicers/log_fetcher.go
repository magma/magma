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
	"context"
	"time"

	"magma/dp/cloud/go/protos"
	"magma/dp/cloud/go/services/dp/storage"
	"magma/dp/cloud/go/services/dp/storage/db"
)

type logFetcher struct {
	protos.UnimplementedLogFetcherServer
	store storage.LogFetcher
}

func NewLogFetcher(store storage.LogFetcher) protos.LogFetcherServer {
	return &logFetcher{
		store: store,
	}
}

func (l *logFetcher) ListLogs(_ context.Context, request *protos.ListLogsRequest) (*protos.ListLogsResponse, error) {
	filter := dbFilter(request.Filter)
	pagination := dbPagination(request.Pagination)
	dbLogs, err := l.store.ListLogs(request.NetworkId, filter, pagination)
	if err != nil {
		return nil, makeErr(err, "list logs")
	}
	logs := make([]*protos.Log, len(dbLogs))
	for i, log := range dbLogs {
		logs[i] = logFromDatabase(log)
	}
	return &protos.ListLogsResponse{Logs: logs}, nil
}

func dbFilter(filter *protos.LogFilter) *storage.LogFilter {
	f := &storage.LogFilter{
		From:         filter.From,
		To:           filter.To,
		FccId:        filter.FccId,
		SerialNumber: filter.SerialNumber,
		Name:         filter.Name,
	}
	if filter.ResponseCode != nil {
		f.ResponseCode = db.MakeInt(filter.ResponseCode.Value)
	}
	if filter.BeginTimestampMilli != nil {
		f.BeginTimestamp = formatTimestamp(filter.BeginTimestampMilli.Value)
	}
	if filter.EndTimestampMilli != nil {
		f.EndTimestamp = formatTimestamp(filter.EndTimestampMilli.Value)
	}
	return f
}

func formatTimestamp(timestamp int64) string {
	const format = "2006-01-02 15:04:05.999+00:00"
	const milli = 1000
	sec, msec := timestamp/milli, timestamp%milli
	return time.Unix(sec, msec*int64(time.Millisecond)).UTC().Format(format)
}

func logFromDatabase(data *storage.DBLog) *protos.Log {
	return &protos.Log{
		From:           data.From.String,
		To:             data.To.String,
		Name:           data.Name.String,
		Message:        data.Message.String,
		SerialNumber:   data.SerialNumber.String,
		FccId:          data.FccId.String,
		TimestampMilli: data.CreatedDate.Time.UnixNano() / int64(time.Millisecond),
	}
}
