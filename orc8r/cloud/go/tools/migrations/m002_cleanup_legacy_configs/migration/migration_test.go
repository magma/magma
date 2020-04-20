/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package migration_test

import (
	"bytes"
	"database/sql"
	"testing"

	"magma/orc8r/cloud/go/tools/migrations/m002_cleanup_legacy_configs/migration"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var mockResult = sqlmock.NewResult(1, 1)

func TestMigrateNetworkConfigsToRecords(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	expectSelectExists(mock, "networks")
	// Cover 3 branches - 1 legacy config with magmad key, 1 legacy config without magmad key, 1 non-legacy record
	mock.ExpectQuery("SELECT key, value FROM networks").
		WillReturnRows(
			sqlmock.NewRows([]string{"key", "value"}).
				AddRow("network1", getDefaultConfigFixture(t)).
				AddRow("network2", getConfigFixture(t, map[string][]byte{"not_magmad": []byte("hello")})).
				AddRow("network3", getRecordFixture(t, "world")),
		)
	prepare := mock.ExpectPrepare("UPDATE networks")
	prepare.ExpectExec().WithArgs(getRecordFixture(t, "hello"), "network1").WillReturnResult(mockResult)
	prepare.ExpectExec().WithArgs(getRecordFixture(t, "network2"), "network2").WillReturnResult(mockResult)
	// No exec should be called for network3
	prepare.WillBeClosed()

	tx := openMockDBTx(t, db)
	err = migration.MigrateNetworkConfigsToRecords(tx)
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

func openMockDBTx(t *testing.T, db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error opening stub DB tx: %s", err)
	}
	return tx
}

func expectSelectExists(mock sqlmock.Sqlmock, table string) {
	mock.ExpectQuery("SELECT EXISTS").WithArgs(table).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
}

func getDefaultConfigFixture(t *testing.T) []byte {
	record := &migration.Record{Name: "hello"}
	val, err := marshalIntern(record)
	assert.NoError(t, err)
	return getConfigFixture(t, map[string][]byte{"magmad": val})
}

func getConfigFixture(t *testing.T, vals map[string][]byte) []byte {
	msg := &migration.Config{ConfigsByKey: vals}
	marshaled, err := marshalIntern(msg)
	assert.NoError(t, err)
	return marshaled
}

func getRecordFixture(t *testing.T, name string) []byte {
	record := &migration.Record{Name: name}
	val, err := marshalIntern(record)
	assert.NoError(t, err)
	return val
}

func marshalIntern(msg proto.Message) ([]byte, error) {
	var buff bytes.Buffer
	err := (&jsonpb.Marshaler{EmitDefaults: true, Indent: " "}).Marshal(
		&buff, msg)
	return buff.Bytes(), err
}
