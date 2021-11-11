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

/*
	Package migrations houses the orc8r data migration executables.

	This doc.go describes how to write an orc8r data migration.

	Functionality

	Migrations are expected to conform to the following descriptions
	- Idempotent
		- Running the migration multiple times has same effect as running once
	- Self-contained
		- Migration itself doesn't rely on services or internal code
		- Edit SQL tables directly
		- Okay to use sqorc internal library as it's expected to remain stable
		- Okay to use internal code for validation step
	- Service-based validation
		- Provide option to call to running services and check their semantic
		  view of the updated tables

	Getting started

	- Migration gets its own, codebase-unique package name
		- MODULE/cloud/go/tools/migrations/m042_short_description
		- MODULE should be orc8r, lte, cwf, etc
		- Package name must be unique across the codebase
			- Use unique migration number
	- Emulate recent data migrations
		- Full migration should run in a single, sql.LevelSerializable tx
		- Minimize number of round-trip calls to DB
		- Provide abundant documentation on migration's purpose and mechanism

	During development

	Run the migration locally, for faster iteration, and against a recent copy
	of the prod DB, for visibility into the expected changes.
	- Ask teammate for copy of prod DB, then load it into your local Postgres
		- (host) $ docker cp ~/.magma/dbs/pgdump-prod-1595044229.sql orc8r_postgres_1:/var/lib/postgresql/data/
		- (orc8r_postgres_1) $ createdb -U magma_dev -T template0 pgdump-prod-1595044229
		- (orc8r_postgres_1) $ psql -U magma_dev pgdump-prod-1595044229 < /var/lib/postgresql/data/pgdump-prod-1595044229.sql
	- Run migration locally, without needing to rebuild all containers
		- Temporarily comment-out the ip.IsLoopback check in
		  unary/identity_decorator.go, then rebuild containers
		- Point editor/shell to prod DB via the DATABASE_SOURCE env variable
			- DATABASE_SOURCE='host=localhost dbname=prod-1595044229-july-17 user=magma_dev password=magma_dev sslmode=disable'
		- (optional) Update cloud docker-compose to expose the port of the
		  relevant controller service

	Manual verification

	Perform a final manual verification step against the prod DB, running the
	newly-built migration executable from within the controller container.
	- Point controller services to prod DB
		- Set DATABASE_SOURCE environment variable(s) in cloud docker-compose
		  mirroring above
		- Restart controller services
	- Run the migration from the controller container
		- Perform 2x to ensure idempotence
			- Run migration
			- Check relevant tables
			- Check relevant endpoints and logs
*/
package migrations
