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
 *
 * We are using JSDoc type annotations because renaming this file will cause
 * the migration to be re-executed.
 *
 * NEW MIGRATIONS SHOULD BE WRITTEN IN TYPESCRIPT!
 *
 * @typedef { import("sequelize").QueryInterface } QueryInterface
 * @typedef { import("sequelize").DataTypes } DataTypes
 */

const SequelizeTables = [
  'AuditLogEntries',
  'FeatureFlags',
  'Organizations',
  'Users',
];

const SequenceColumn = 'id';

/**
 * For Postgres, fixes the id sequence for tables that have been migrated.
 * When migrating table data to Postgres, it is necessary to reset the id
 * sequence, otherwise inserts will fail due to unique ID constraint.
 * This resetting of ID sequences was not accounted for in this package's
 * db data migration function, so this migration can account for that.
 */
module.exports = {
  /**
   * @param {QueryInterface} queryInterface
   */
  up: queryInterface => {
    return queryInterface.sequelize.transaction(async transaction => {
      try {
        for (const table of SequelizeTables) {
          // Get current highest value from the table

          // prettier-ignore
          const [
            [{max}],
          ] = /** @type {[[{max:number}],unknown]} */ (await queryInterface.sequelize.query(
            `SELECT MAX("${SequenceColumn}") AS max FROM "${table}";`,
            {transaction},
          ));

          // Set the autoincrement current value to highest value + 1
          await queryInterface.sequelize.query(
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
    });
  },

  down: () => {
    return Promise.resolve();
  },
};
