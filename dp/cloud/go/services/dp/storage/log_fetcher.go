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
package storage

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/sqorc"
)

type LogFetcher interface {
	ListLogs(networkId string, filter *LogFilter, pagination *Pagination) ([]*DBLog, error)
}

type LogFilter struct {
	From           string
	To             string
	FccId          string
	SerialNumber   string
	Name           string
	ResponseCode   sql.NullInt64
	BeginTimestamp string
	EndTimestamp   string
}

func NewLogFetcher(db *sql.DB, builder sqorc.StatementBuilder) *logFetcher {
	return &logFetcher{
		db:      db,
		builder: builder,
	}
}

type logFetcher struct {
	db      *sql.DB
	builder sqorc.StatementBuilder
}

func (l *logFetcher) ListLogs(networkId string, filter *LogFilter, pagination *Pagination) ([]*DBLog, error) {
	res, err := sqorc.ExecInTx(l.db, nil, nil, func(tx *sql.Tx) (interface{}, error) {
		return buildPagination(db.NewQuery(), pagination).
			WithBuilder(l.builder.RunWith(tx)).
			From(&DBLog{}).
			Select(db.NewExcludeMask("id", "network_id")).
			Where(buildFilter(networkId, filter)).
			OrderBy("created_date", db.OrderDesc).
			List()
	})
	if err != nil {
		return nil, makeError(err)
	}
	models := res.([][]db.Model)
	logs := make([]*DBLog, len(models))
	for i, model := range models {
		logs[i] = model[0].(*DBLog)
	}
	return logs, nil
}

func buildFilter(networkId string, filter *LogFilter) sq.Sqlizer {
	res := sq.And{sq.Eq{"network_id": networkId}}
	if filter.From != "" {
		res = append(res, sq.Eq{"log_from": filter.From})
	}
	if filter.To != "" {
		res = append(res, sq.Eq{"log_to": filter.To})
	}
	if filter.FccId != "" {
		res = append(res, sq.Eq{"fcc_id": filter.FccId})
	}
	if filter.SerialNumber != "" {
		res = append(res, sq.Eq{"cbsd_serial_number": filter.SerialNumber})
	}
	if filter.Name != "" {
		res = append(res, sq.Eq{"log_name": filter.Name})
	}
	if filter.ResponseCode.Valid {
		res = append(res, sq.Eq{"response_code": filter.ResponseCode.Int64})
	}
	if filter.BeginTimestamp != "" {
		res = append(res, sq.GtOrEq{"created_date": filter.BeginTimestamp})
	}
	if filter.EndTimestamp != "" {
		res = append(res, sq.LtOrEq{"created_date": filter.EndTimestamp})
	}
	return res
}
