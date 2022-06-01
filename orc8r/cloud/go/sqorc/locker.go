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

package sqorc

import (
	"os"
	"strings"
)

// GetSqlLocker returns a lock statement generator for the configured SQL
// dialect as found in the SQL_DIALECT env var.
func GetSqlLocker() Locker {
	dialect := os.Getenv(SQLDialectEnv)
	// sqlite doesn't support locking
	switch strings.ToLower(dialect) {
	case PostgresDialect, MariaDialect:
		return SqlLocker{}
	default:
		return DummyLocker{}
	}
}

type Locker interface {
	WithLock() string
}

type SqlLocker struct{}

func (s SqlLocker) WithLock() string {
	return "FOR UPDATE"
}

type DummyLocker struct{}

func (d DummyLocker) WithLock() string {
	return ""
}
