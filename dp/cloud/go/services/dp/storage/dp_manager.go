/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package storage

import (
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"magma/orc8r/cloud/go/sqorc"
)

type dpManager struct {
	db           *sql.DB
	builder      sqorc.StatementBuilder
	cache        *enumCache
	errorChecker sqorc.ErrorChecker
	locker       sqorc.Locker
}

type queryRunner struct {
	builder sq.StatementBuilderType
	cache   *enumCache
	locker  sqorc.Locker
}

func (m *dpManager) getQueryRunner(tx sq.BaseRunner) *queryRunner {
	return &queryRunner{
		builder: m.builder.RunWith(tx),
		cache:   m.cache,
		locker:  m.locker,
	}
}
