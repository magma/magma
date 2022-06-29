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

import Umzug from 'umzug';
import {DataTypes} from 'sequelize';

const umzug = new Umzug({
  storage: 'sequelize',
  storageOptions: {
    sequelize,
  },
  // The logging function.
  // A function that gets executed everytime migrations start and have ended.
  logging: (msg: string) => logger.info(msg),
  // The name of the positive method in migrations.
  upName: 'up',
  // The name of the negative method in migrations.
  downName: 'down',
  migrations: {
    // The params that gets passed to the migrations.
    // Might be an array or a synchronous function which returns an array.
    params: [sequelize.getQueryInterface(), DataTypes],
    // The path to the migrations directory.
    path: sequelizerc['migrations-path'],
    // The pattern that determines whether or not a file is a migration.
    pattern: /^\d+[\w-]+\.[jt]s$/,
    // A function that receives and returns the to be executed function.
    // This can be used to modify the function.
    wrap(func) {
      return func;
    },
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
