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

	"magma/dp/cloud/go/services/dp/storage/db"
)

type Pagination struct {
	Limit  sql.NullInt64
	Offset sql.NullInt64
}

func buildPagination(q *db.Query, pagination *Pagination) *db.Query {
	if pagination.Limit.Valid {
		q = q.Limit(pagination.Limit.Int64)
		if pagination.Offset.Valid {
			q = q.Offset(pagination.Offset.Int64)
		}
	}
	return q
}
