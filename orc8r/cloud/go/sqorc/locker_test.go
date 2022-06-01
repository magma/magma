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

package sqorc_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/sqorc"
)

func TestGetLocker(t *testing.T) {
	prev := os.Getenv(sqorc.SQLDialectEnv)
	defer os.Setenv(sqorc.SQLDialectEnv, prev)

	testData := []struct {
		name     string
		dialect  string
		expected sqorc.Locker
	}{{
		name:     "test dummy locker for no dialect",
		dialect:  "",
		expected: sqorc.DummyLocker{},
	}, {
		name:     "test dummy locker for sqlite dialect",
		dialect:  sqorc.SQLiteDialect,
		expected: sqorc.DummyLocker{},
	}, {
		name:     "test sql locker for maria dialect",
		dialect:  sqorc.MariaDialect,
		expected: sqorc.SqlLocker{},
	}, {
		name:     "test sql locker for postgres dialect",
		dialect:  sqorc.PostgresDialect,
		expected: sqorc.SqlLocker{},
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv(sqorc.SQLDialectEnv, tt.dialect)
			actual := sqorc.GetSqlLocker()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestLockerWithLock(t *testing.T) {
	testData := []struct {
		name     string
		locker   sqorc.Locker
		expected string
	}{{
		name:     "test dummy locker generate nothing",
		locker:   sqorc.DummyLocker{},
		expected: "",
	}, {
		name:     "test sql locker generate sql lock",
		locker:   sqorc.SqlLocker{},
		expected: "FOR UPDATE",
	}}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.locker.WithLock()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
