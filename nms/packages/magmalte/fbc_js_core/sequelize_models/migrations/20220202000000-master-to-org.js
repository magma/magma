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
 * @flow strict-local
 * @format
 */
'use strict';

import type {DataTypes, QueryInterface, Transaction} from 'sequelize';

// Key: table name
// Value: column name which denotes the organization name
const TableOrgNameColumn: {[string]: string} = {
  AuditLogEntries: 'organization',
  FeatureFlags: 'organization',
  Organizations: 'name',
  Users: 'organization',
};

/**
 * This migration changes the name of the 'master' org to 'host'.
 *
 * This migration has two sets of nearly identical queries, one for
 * Postgres, and one for MySQL.
 */
module.exports = {
  up: (queryInterface: QueryInterface, _types: DataTypes) => {
    return queryInterface.sequelize.transaction(
      async (transaction: Transaction): Promise<void> => {
        // Postgres needs capitalized table names surrounded by quotations
        // This would cause an error in MySQL
        const dialect = queryInterface.sequelize.getDialect();
        let quote = '';
        switch (dialect) {
          case 'mysql':
          case 'mariadb':
            break;
          case 'postgres':
            quote = '"';
            break;
          default:
            console.error(
              `Unsupported DB dialect for migration: ${dialect}` +
                'Supported dialects are [mysql, mariadb, postgres]',
            );
        }

        // Update name of 'master' organization to 'host'
        let query = '';
        try {
          for (const tableName: string of Object.keys(TableOrgNameColumn)) {
            const orgColName: string = TableOrgNameColumn[tableName];
            query = `UPDATE ${quote}${tableName}${quote} SET ${orgColName}='host' WHERE ${orgColName}='master'`;
            await queryInterface.sequelize.query(query, {transaction});
          }
          return;
        } catch (err) {
          console.error(
            `Failed to run query for migration: ${query}, error: ${err}`,
          );
          throw 'Failed to complete migration';
        }
      },
    );
  },

  down: (queryInterface: QueryInterface, _types: DataTypes) => {
    return queryInterface.sequelize.transaction(
      async (transaction: Transaction): Promise<void> => {
        // Postgres needs capitalized table names surrounded by quotations
        // This would cause an error in MySQL
        const dialect = queryInterface.sequelize.getDialect();
        let quote = '';
        switch (dialect) {
          case 'mysql':
          case 'mariadb':
            break;
          case 'postgres':
            quote = '"';
            break;
          default:
            console.error(
              `Unsupported DB dialect for migration: ${dialect}` +
                'Supported dialects are [mysql, mariadb, postgres]',
            );
        }

        // Update name of 'host' organization to 'master'
        let query = '';
        try {
          for (const tableName: string of Object.keys(TableOrgNameColumn)) {
            const orgColName: string = TableOrgNameColumn[tableName];
            query = `UPDATE ${quote}${tableName}${quote} SET ${orgColName}='master' WHERE ${orgColName}='host'`;
            await queryInterface.sequelize.query(query, {transaction});
          }
          return;
        } catch (err) {
          console.error(
            `Failed to run query for migration: ${query}, error: ${err}`,
          );
          throw 'Failed to complete migration';
        }
      },
    );
  },
};
