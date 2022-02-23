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
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type pagination struct {
	limit   sql.NullInt64
	offset  sql.NullInt64
	orderBy string
	order   Order
}

type Order string

const (
	OrderAsc  Order = "ASC"
	OrderDesc Order = "DESC"
)

func addPagination(p *pagination, builder sq.SelectBuilder) sq.SelectBuilder {
	if p.limit.Valid {
		builder = builder.Limit(uint64(p.limit.Int64))
	}
	if p.offset.Valid {
		builder = builder.Offset(uint64(p.offset.Int64))
	}
	if p.orderBy != "" {
		builder = builder.OrderBy(fmt.Sprintf("%s %s", p.orderBy, p.order))
	}
	return builder
}
