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
 *
 * @flow
 * @format
 */

const logger = require('@fbcnms/logging').getLogger(module);

const {importFromDatabase} = require('@fbcnms/sequelize-models');
import {User} from '@fbcnms/sequelize-models';
import type {Options} from 'sequelize';

export async function runDataMigration() {
  const sourceDbHost = process.env.MYSQL_SRC_HOST;
  if (typeof sourceDbHost === 'undefined') {
    return;
  }
  const allUsers = await User.findAll();
  if (allUsers.length > 0) {
    logger.info('Users found in NMS DB. Skipping migration from source DB');
    return;
  }

  const dbOptions: Options = {
    username: process.env.MYSQL_SRC_USER || 'root',
    password: process.env.MYSQL_SRC_PASS || '',
    database: process.env.MYSQL_SRC_DB || 'nms',
    host: process.env.MYSQL_SRC_HOST,
    port: parseInt(process.env.MYSQL_SRC_PORT || '3306'),
    dialect: process.env.MYSQL_SRC_DIALECT || 'mysql',
    logging: (msg: string) => logger.debug(msg),
  };
  try {
    await importFromDatabase(dbOptions);
    logger.info('Completed data migration to current NMS DB');
  } catch (error) {
    logger.error(
      `Unable to connect to source database with specified options for migration: ${error}`,
    );
  }
}
