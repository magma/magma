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

import AuditLogEntryModel from './models/audit_log_entry';
import FeatureFlagModel from './models/featureflag';
import OrganizationModel from './models/organization';
import Sequelize from 'sequelize';
import UserModel from './models/user';
import sequelizeConfig from './sequelizeConfig';

const env = process.env.NODE_ENV || 'development';
const config = sequelizeConfig[env];

export const sequelize = new Sequelize(
  config.database || '',
  config.username,
  config.password,
  config,
);

const db = {
  AuditLogEntry: AuditLogEntryModel(sequelize, Sequelize),
  FeatureFlag: FeatureFlagModel(sequelize, Sequelize),
  Organization: OrganizationModel(sequelize, Sequelize),
  User: UserModel(sequelize, Sequelize),
};

Object.keys(db).forEach(
  modelName => db[modelName].associate != null && db[modelName].associate(db),
);

export const AuditLogEntry = db.AuditLogEntry;
export const Organization = db.Organization;
export const User = db.User;
export const FeatureFlag = db.FeatureFlag;

export function jsonArrayContains(column: string, value: string) {
  if (
    sequelize.getDialect() === 'mysql' ||
    sequelize.getDialect() === 'mariadb'
  ) {
    return Sequelize.fn('JSON_CONTAINS', Sequelize.col(column), `"${value}"`);
  } else {
    // sqlite
    const escapedColumn = sequelize
      .getQueryInterface()
      .quoteIdentifier(column, true);
    const innerQuery = Sequelize.literal(
      `(SELECT 1 FROM json_each(${escapedColumn})` +
        `WHERE json_each.value = ${sequelize.escape(value)})`,
    );
    return Sequelize.where(innerQuery, 'IS', Sequelize.literal('NOT NULL'));
  }
}
