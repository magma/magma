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
	"errors"
	"testing"

	"magma/orc8r/cloud/go/tools/migrations/m001_config_service/migration"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var mockResult = sqlmock.NewResult(1, 1)

func TestMigrateNetworkConfigs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	expectSelectExists(mock, "networks")
	mock.ExpectQuery("SELECT key, value FROM networks").
		WillReturnRows(
			sqlmock.NewRows([]string{"key", "value"}).
				AddRow("network1", getDefaultConfigFixture(t)).
				AddRow("network2", getConfigFixture(t, map[string][]byte{"cellular": []byte("cellular")})),
		)

	expectCreateTable(mock, "network1_configurations")
	network1Prepare := mock.ExpectPrepare("INSERT INTO network1_configurations")
	network1Prepare.ExpectExec().WithArgs("cellular_network", "network1", []byte("world"), []byte("world")).
		WillReturnResult(mockResult)
	network1Prepare.ExpectExec().WithArgs("magmad_network", "network1", []byte("hello"), []byte("hello")).
		WillReturnResult(mockResult)
	network1Prepare.WillBeClosed()

	expectCreateTable(mock, "network2_configurations")
	network2Prepare := mock.ExpectPrepare("INSERT INTO network2_configurations")
	network2Prepare.ExpectExec().WithArgs("cellular_network", "network2", []byte("cellular"), []byte("cellular")).
		WillReturnResult(mockResult)
	network2Prepare.WillBeClosed()

	tx := openMockDBTx(t, db)
	err = migration.MigrateNetworkConfigs(tx)
	assert.NoError(t, err)
}

func TestMigrateGatewayConfigs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT key FROM networks").
		WillReturnRows(sqlmock.NewRows([]string{"key"}).AddRow("network1").AddRow("network2"))

	expectSelectExists(mock, "network1_configs")
	mock.ExpectQuery("SELECT key, value FROM network1_configs").
		WillReturnRows(
			sqlmock.NewRows([]string{"key", "value"}).
				AddRow("gw1", getDefaultConfigFixture(t)).
				AddRow("gw2", getConfigFixture(t, map[string][]byte{"wifi": []byte("wifi")})),
		)
	expectCreateTable(mock, "network1_configurations")
	network1Prepare := mock.ExpectPrepare("INSERT INTO network1_configurations")
	network1Prepare.ExpectExec().WithArgs("cellular_gateway", "gw1", []byte("world"), []byte("world")).
		WillReturnResult(mockResult)
	network1Prepare.ExpectExec().WithArgs("magmad_gateway", "gw1", []byte("hello"), []byte("hello")).
		WillReturnResult(mockResult)
	network1Prepare.ExpectExec().WithArgs("wifi_gateway", "gw2", []byte("wifi"), []byte("wifi")).
		WillReturnResult(mockResult)
	network1Prepare.WillBeClosed()

	expectSelectExists(mock, "network2_configs")
	mock.ExpectQuery("SELECT key, value FROM network2_configs").
		WillReturnRows(
			sqlmock.NewRows([]string{"key", "value"}).
				AddRow("gw3", getConfigFixture(t, map[string][]byte{"magmad": []byte("magmad")})),
		)
	expectCreateTable(mock, "network2_configurations")
	network2Prepare := mock.ExpectPrepare("INSERT INTO network2_configurations")
	network2Prepare.ExpectExec().WithArgs("magmad_gateway", "gw3", []byte("magmad"), []byte("magmad")).
		WillReturnResult(mockResult)
	network2Prepare.WillBeClosed()

	tx := openMockDBTx(t, db)
	err = migration.MigrateGatewayConfigs(tx)
	assert.NoError(t, err)
}

func TestMigrateMeshConfigs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT key FROM networks").
		WillReturnRows(sqlmock.NewRows([]string{"key"}).AddRow("network1").AddRow("network2"))

	expectCreateTable(mock, "network1_configurations")
	network1Prepare := mock.ExpectPrepare("INSERT INTO network1_configurations")
	expectSelectExists(mock, "network1_mesh_config")
	mock.ExpectQuery("SELECT key, value FROM network1_mesh_config").
		WillReturnRows(
			sqlmock.NewRows([]string{"key", "value"}).
				AddRow("mesh1", []byte("hello")).
				AddRow("mesh2", []byte("world")),
		)
	network1Prepare.ExpectExec().WithArgs("mesh", "mesh1", []byte("hello"), []byte("hello")).
		WillReturnResult(mockResult)
	network1Prepare.ExpectExec().WithArgs("mesh", "mesh2", []byte("world"), []byte("world")).
		WillReturnResult(mockResult)
	network1Prepare.WillBeClosed()

	expectCreateTable(mock, "network2_configurations")
	network2Prepare := mock.ExpectPrepare("INSERT INTO network2_configurations")
	expectSelectExists(mock, "network2_mesh_config")
	mock.ExpectQuery("SELECT key, value FROM network2_mesh_config").
		WillReturnRows(
			sqlmock.NewRows([]string{"key", "value"}).AddRow("mesh3", []byte("mesh3")),
		)
	network2Prepare.ExpectExec().WithArgs("mesh", "mesh3", []byte("mesh3"), []byte("mesh3")).
		WillReturnResult(mockResult)
	network2Prepare.WillBeClosed()

	tx := openMockDBTx(t, db)
	err = migration.MigrateMeshConfigs(tx)
	assert.NoError(t, err)
}

// Same code flow as gateway config migration so we'll only run through the
// error cases for networks
func TestMigrateNetworkConfigs_Errors(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening stub DB conn: %s", err)
	}
	defer db.Close()

	// Read error
	mock.ExpectBegin()
	expectSelectExists(mock, "networks")
	mock.ExpectQuery("SELECT key, value FROM networks").WillReturnError(errors.New("mock select error"))

	tx := openMockDBTx(t, db)
	err = migration.MigrateNetworkConfigs(tx)
	assert.EqualError(t, err, "Error getting existing network configs: mock select error")
	assert.NoError(t, mock.ExpectationsWereMet())

	// Unrecognized config type
	mock.ExpectBegin()
	expectSelectExists(mock, "networks")
	mock.ExpectQuery("SELECT key, value FROM networks").
		WillReturnRows(
			sqlmock.NewRows([]string{"key", "value"}).
				AddRow("network1", getConfigFixture(t, map[string][]byte{"error": []byte("error")})),
		)
	expectCreateTable(mock, "network1_configurations")
	prepare := mock.ExpectPrepare("INSERT INTO network1_configurations")
	prepare.WillBeClosed()

	tx = openMockDBTx(t, db)
	err = migration.MigrateNetworkConfigs(tx)
	assert.EqualError(t, err, "Error migrating network config for network network1: Error migrating magmad config for network1 in network network1: No new config type defined for magmad key error")

	// Insert error
	mock.ExpectBegin()
	expectSelectExists(mock, "networks")
	mock.ExpectQuery("SELECT key, value FROM networks").
		WillReturnRows(
			sqlmock.NewRows([]string{"key", "value"}).
				AddRow("network1", getConfigFixture(t, map[string][]byte{"magmad": []byte("magmad")})),
		)
	expectCreateTable(mock, "network1_configurations")
	prepare = mock.ExpectPrepare("INSERT INTO network1_configurations")
	prepare.ExpectExec().WithArgs("magmad_network", "network1", []byte("magmad"), []byte("magmad")).
		WillReturnError(errors.New("Mock upsert error"))
	prepare.WillBeClosed()

	tx = openMockDBTx(t, db)
	err = migration.MigrateNetworkConfigs(tx)
	assert.EqualError(t, err, "Error migrating network config for network network1: Error migrating magmad config for network1 in network network1: Mock upsert error")
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

func expectCreateTable(mock sqlmock.Sqlmock, table string) {
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS " + table).WillReturnResult(mockResult)
}

func getDefaultConfigFixture(t *testing.T) []byte {
	return getConfigFixture(t, map[string][]byte{"magmad": []byte("hello"), "cellular": []byte("world")})
}

func getConfigFixture(t *testing.T, vals map[string][]byte) []byte {
	msg := &migration.Config{ConfigsByKey: vals}
	marshaled, err := marshalIntern(msg)
	assert.NoError(t, err)
	return marshaled
}

func marshalIntern(msg proto.Message) ([]byte, error) {
	var buff bytes.Buffer
	err := (&jsonpb.Marshaler{EmitDefaults: true, Indent: " "}).Marshal(
		&buff, msg)
	return buff.Bytes(), err
}
