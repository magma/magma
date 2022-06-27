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
 */

import {
  DataTypes,
  Model,
  Options,
  QueryTypes,
  Sequelize,
  Transaction,
  WhereOptions,
} from 'sequelize';

import createAuditLogEntryModel from './models/audit_log_entry';
import createFeatureFlagModel from './models/featureflag';
import createOrganizationModel from './models/organization';
import createUserModel, {UserModelStatic} from './models/user';
import sequelizeConfig from './sequelizeConfig';

const env = process.env.NODE_ENV || 'development';
const config = sequelizeConfig[env];

export const sequelize: Sequelize = new Sequelize(
  config.database || '',
  config.username || '',
  config.password,
  config,
);

const SequelizeTables = [
  'AuditLogEntries',
  'FeatureFlags',
  'Organizations',
  'Users',
];
const SequenceColumn = 'id';

const db = createNmsDb(sequelize);

function createNmsDb(sequelize: Sequelize) {
  const db = {
    AuditLogEntry: createAuditLogEntryModel(sequelize),
    FeatureFlag: createFeatureFlagModel(sequelize),
    Organization: createOrganizationModel(sequelize),
    User: createUserModel(sequelize),
  };

  (Object.keys(db) as Array<keyof typeof db>).forEach(
    modelName => db[modelName].associate != null && db[modelName].associate(db),
  );
  return db;
}

async function isMigrationNeeded(userModel: UserModelStatic): Promise<boolean> {
  try {
    const allUsers = await userModel.findAll();

    if (allUsers.length > 0) {
      console.warn('Users found in target DB. Migration may already have run');
      return false;
    }
    return true;
  } catch (e) {
    const message = e instanceof Error ? e.message : 'unknown error';
    console.error(
      `Unable to run migration. Connection error to specified database: \n` +
        `------------------------\n` +
        `${message} \n` +
        `------------------------\n`,
    );
    process.exit(1);
  }
  return false;
}

export const AuditLogEntry = db.AuditLogEntry;
export const Organization = db.Organization;
export const User = db.User;
export const FeatureFlag = db.FeatureFlag;

export function jsonArrayContains(column: string, value: string): WhereOptions {
  if (
    sequelize.getDialect() === 'mysql' ||
    sequelize.getDialect() === 'mariadb'
  ) {
    // TODO[ts-migration]
    return (Sequelize.fn(
      'JSON_CONTAINS',
      Sequelize.col(column),
      `"${value}"`,
    ) as unknown) as WhereOptions;
  } else if (sequelize.getDialect() === 'postgres') {
    const escapedColumn = sequelize
      .getQueryInterface()
      .quoteIdentifier(column, true);
    const escapedValue = sequelize
      .getQueryInterface()
      .quoteIdentifier(value, true);
    return Sequelize.literal(`${escapedColumn}::jsonb @> '${escapedValue}'`);
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

export async function importFromDatabase(sourceConfig: Options) {
  const sourceSequelize = new Sequelize(
    sourceConfig.database || '',
    sourceConfig.username || '',
    sourceConfig.password,
    sourceConfig,
  );
  const sourceDb = createNmsDb(sourceSequelize);

  await sourceDb.AuditLogEntry.sync();
  await sourceDb.FeatureFlag.sync();
  await sourceDb.Organization.sync();
  await sourceDb.User.sync();

  const willRunMigration = await isMigrationNeeded(User);
  if (!willRunMigration) {
    console.log('Skipping DB migration');
    return;
  }

  await migrateMeta(sourceSequelize, sequelize);

  const auditLogEntries = await sourceDb.AuditLogEntry.findAll();
  await AuditLogEntry.bulkCreate(getDataValues(auditLogEntries));

  const featureFlags = await sourceDb.FeatureFlag.findAll();
  await FeatureFlag.bulkCreate(getDataValues(featureFlags));

  const organizations = await sourceDb.Organization.findAll();
  await Organization.bulkCreate(getDataValues(organizations));

  const users = await sourceDb.User.findAll();
  await User.bulkCreate(getDataValues(users));
}

export async function exportToDatabase(targetConfig: Options) {
  const targetSequelize = new Sequelize(
    targetConfig.database || '',
    targetConfig.username || '',
    targetConfig.password,
    targetConfig,
  );
  const targetDb = createNmsDb(targetSequelize);

  await targetDb.AuditLogEntry.sync();
  await targetDb.FeatureFlag.sync();
  await targetDb.Organization.sync();
  await targetDb.User.sync();

  const willRunMigration = await isMigrationNeeded(targetDb.User);
  if (!willRunMigration) {
    console.log('Skipping DB migration');
    return;
  }

  await migrateMeta(sequelize, targetSequelize);

  const auditLogEntries = await AuditLogEntry.findAll();
  await targetDb.AuditLogEntry.bulkCreate(getDataValues(auditLogEntries));

  const featureFlags = await FeatureFlag.findAll();
  await targetDb.FeatureFlag.bulkCreate(getDataValues(featureFlags));

  const organizations = await Organization.findAll();
  await targetDb.Organization.bulkCreate(getDataValues(organizations));

  const users = await User.findAll();
  await targetDb.User.bulkCreate(getDataValues(users));

  await resetPgIdSeq(targetSequelize);
}

/**
 * Reset the Postgres ID sequence.
 * When migrating table data to Postgres, it is necessary to reset the id
 * sequence, otherwise inserts will fail due to unique ID constraint.
 */
async function resetPgIdSeq(sequelize: Sequelize) {
  await sequelize.transaction(
    async (transaction: Transaction): Promise<void> => {
      try {
        for (const table of SequelizeTables) {
          // Get current highest value from the table
          const [
            [{max}],
          ] = (await sequelize.query(
            `SELECT MAX("${SequenceColumn}") AS max FROM "${table}";`,
            {transaction},
          )) as [[{max: number}], unknown];

          // Set the autoincrement current value to highest value + 1
          await sequelize.query(
            `ALTER SEQUENCE "${table}_id_seq" RESTART WITH ${max + 1};`,
            {transaction},
          );
        }
      } catch (exception) {
        // This likely means we're just not running in Postgres.
        // Do nothing.
        console.error(
          'Had an issue resetting ID sequences. ',
          'If you are running Postgres, you may need to reset ID sequences manually. ',
          'Otherwise, ignore this error',
          'Exception: ',
          exception,
        );
      }
    },
  );
}

async function migrateMeta(source: Sequelize, target: Sequelize) {
  // Read in the current SequelizeMeta data
  const rows = await sequelize.query('SELECT * FROM `SequelizeMeta`', {
    type: QueryTypes.SELECT,
  });

  // Write SequelizeMeta data
  const targetInterface = target.getQueryInterface();
  await targetInterface.createTable('SequelizeMeta', {
    name: {
      type: DataTypes.STRING,
      allowNull: false,
      unique: true,
      primaryKey: true,
    },
  });
  await targetInterface.bulkInsert('SequelizeMeta', rows);
}

function getDataValues(sequelizeModels: Array<Model>): Array<object> {
  // @ts-ignore
  return sequelizeModels.map(model => model.dataValues); // eslint-disable-line @typescript-eslint/no-unsafe-return
}
