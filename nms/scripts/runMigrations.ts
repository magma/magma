/**
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
import Logging from '../shared/logging';

const logger = Logging.getLogger(module);

import sequelizerc from '../shared/sequelize_models/sequelizerc';
import {sequelize} from '../shared/sequelize_models';

import {SequelizeStorage, Umzug} from 'umzug';

const umzug = new Umzug({
  storage: new SequelizeStorage({sequelize}),
  logger,

  // The context gets passed to the migrations.
  context: sequelize.getQueryInterface(),

  migrations: {
    glob: `${sequelizerc['migrations-path']}/*.{js,ts}`,
  },
});

export async function runMigrations() {
  const pendingMigrations = await umzug.pending();
  if (pendingMigrations) {
    await umzug.up();
  }
  // Sync defined models to the DB
  await sequelize.sync();
}

export async function rollbackMigrations() {
  const executedMigrations = await umzug.executed();
  if (executedMigrations) {
    await umzug.down();
  }

  // Sync defined models to the DB
  await sequelize.sync();
}

export default umzug;
