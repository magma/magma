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

package sqorc

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"magma/orc8r/lib/go/definitions"

	"github.com/golang/glog"
	"github.com/stretchr/testify/require"
)

const (
	// hostEnvVar is the name of the environment variable holding the DB
	// server's hostname.
	hostEnvVar = "TEST_DATABASE_HOST"
	// postgresPortEnvVar is the name of the environment variable holding the
	// Postgres container's port number.
	postgresPortEnvVar = "TEST_DATABASE_PORT_POSTGRES"
	// mariaPortEnvVar is the name of the environment variable holding the
	// Maria container's port number.
	mariaPortEnvVar = "TEST_DATABASE_PORT_MARIA"

	// postgresTestHost is the hostname of the postgres_test container
	postgresTestHost = "postgres_test"
	// mariaTestHost is the hostname of the maria_test container
	mariaTestHost = "maria_test"

	// postgresDefaultPort is the default port exposed by the postgres container.
	postgresDefaultPort = "5432"
	// mariaDefaultPort is the default port exposed by the maria container.
	mariaDefaultPort = "3306"
)

// OpenCleanForTest is the same as OpenForTest, except it also drops then
// creates the underlying DB name before returning.
func OpenCleanForTest(t *testing.T, dbName, dbDriver string) *sql.DB {
	rootDB := OpenForTest(t, "", dbDriver)
	_, err := rootDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", dbName))
	if err != nil {
		t.Fatalf("Failed to drop test DB: %s", err)
	}
	_, err = rootDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("Failed to create test DB: %s", err)
	}
	require.NoError(t, rootDB.Close())

	db := OpenForTest(t, dbName, dbDriver)
	return db
}

// OpenForTest returns a new connection to a shared test DB.
// Does not guarantee the existence of the underlying DB name.
// The shared DB is part of the shared testing infrastructure, so care must
// be taken to avoid racing on the same DB name across testing code.
// Supported DB drivers include:
//	- postgres
//	- mysql
// Environment variables:
//	- SQL_DRIVER overrides the Go SQL driver
//	- TEST_DATABASE_HOST overrides the DB connection host
//	- TEST_DATABASE_PORT_POSTGRES overrides the port connected to for postgres driver
//	- TEST_DATABASE_PORT_MARIA overrides the port connected to for maria driver
func OpenForTest(t *testing.T, dbName, dbDriver string) *sql.DB {
	driver := definitions.GetEnvWithDefault("SQL_DRIVER", dbDriver)
	source := getSource(t, dbName, driver)
	setDialect(driver)

	db, err := Open(driver, source)
	if err != nil {
		t.Fatalf("Could not initialize %s DB connection: %v", driver, err)
	}
	return db
}

// setDialect sets an environment variable to inform all SQL builders.
func setDialect(driver string) {
	old, ok := os.LookupEnv("SQL_DIALECT")
	if ok {
		glog.Infof("Overwriting existing SQL_DIALECT %s", old)
	}
	switch driver {
	case PostgresDriver:
		_ = os.Setenv("SQL_DIALECT", PostgresDialect)
	case MariaDriver:
		_ = os.Setenv("SQL_DIALECT", MariaDialect)
	}
}

// getSource returns the driver-specific data source name.
func getSource(t *testing.T, dbName, driver string) string {
	host := getHost(t, driver)
	port := getPort(t, driver)

	var source string

	switch driver {
	case PostgresDriver:
		// Source: https://godoc.org/github.com/lib/pq#hdr-Connection_String_Parameters
		source = fmt.Sprintf("host=%s port=%s user=magma_test password=magma_test sslmode=disable", host, port)
		if dbName != "" {
			source = fmt.Sprintf("%s dbname=%s", source, dbName)
		}
	case MariaDriver:
		// Source: https://github.com/go-sql-driver/mysql#dsn-data-source-name
		// Format: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
		source = fmt.Sprintf("root:magma_test@(%s:%s)/%s", host, port, dbName)
	default:
		t.Fatalf("Unrecognized DB driver: %s", driver)
	}

	return source
}

func getHost(t *testing.T, driver string) string {
	var host string

	switch driver {
	case PostgresDriver:
		host = definitions.GetEnvWithDefault(hostEnvVar, postgresTestHost)
	case MariaDriver:
		host = definitions.GetEnvWithDefault(hostEnvVar, mariaTestHost)
	default:
		t.Fatalf("Unrecognized DB driver: %s", driver)
	}

	return host
}

func getPort(t *testing.T, driver string) string {
	var port string

	switch driver {
	case PostgresDriver:
		port = definitions.GetEnvWithDefault(postgresPortEnvVar, postgresDefaultPort)
	case MariaDriver:
		port = definitions.GetEnvWithDefault(mariaPortEnvVar, mariaDefaultPort)
	default:
		t.Fatalf("Unrecognized DB driver: %s", driver)
	}

	return port
}
