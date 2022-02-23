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

package db

import (
	sq "github.com/Masterminds/squirrel"

	"magma/orc8r/cloud/go/sqorc"
)

func (q *Query) Insert() (int64, error) {
	var id int64
	fields := applyMask(q.arg.model.Fields(), q.arg.mask)
	err := q.builder.
		Insert(q.arg.model.GetMetadata().Table).
		SetMap(toValues(fields)).
		Suffix("RETURNING id").
		QueryRow().
		Scan(&id)
	return id, err
}

func (q *Query) Update() error {
	fields := applyMask(q.arg.model.Fields(), q.arg.mask)
	_, err := q.builder.
		Update(q.arg.model.GetMetadata().Table).
		SetMap(toValues(fields)).
		Where(q.arg.filter).
		Exec()
	return err
}

func (q *Query) Delete() error {
	_, err := q.builder.
		Delete(q.arg.model.GetMetadata().Table).
		Where(q.arg.filter).
		Exec()
	return err
}

func (q *Query) Fetch() ([]Model, error) {
	colsCollector := collectColumns(q)
	fieldsCollector := collectFields(q, colsCollector.columns)
	query := buildQuery(q.builder, q, q.arg.filter, colsCollector.getColumnNames())
	err := query.
		QueryRow().
		Scan(fieldsCollector.pointers...)
	return fieldsCollector.models, err
}

func (q *Query) List() ([][]Model, error) {
	colsCollector := collectColumns(q)
	query := buildQuery(q.builder, q, q.arg.filter, colsCollector.getColumnNames())
	query = addPagination(q.pagination, query)
	rows, err := query.Query()
	if err != nil {
		return nil, err
	}
	defer sqorc.CloseRowsLogOnError(rows, "List")
	var result [][]Model
	for rows.Next() {
		fieldsCollector := collectFields(q, colsCollector.columns)
		if err := rows.Scan(fieldsCollector.pointers...); err != nil {
			return nil, err
		}
		result = append(result, fieldsCollector.models)
	}
	return result, rows.Err()
}

func applyMask(fields FieldMap, mask FieldMask) FieldMap {
	m := make(FieldMap, len(fields))
	for k, v := range fields {
		if mask.ShouldInclude(k) {
			m[k] = v
		}
	}
	return m
}

func toValues(fields FieldMap) map[string]interface{} {
	m := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		m[k] = v.GetValue()
	}
	return m
}

func buildQuery(builder sq.StatementBuilderType, q *Query, filter sq.Sqlizer, columns []string) sq.SelectBuilder {
	return builder.
		Select(columns...).
		From(q.arg.model.GetMetadata().Table).
		JoinClause(&joinClause{query: q}).
		Where(filter)
}
