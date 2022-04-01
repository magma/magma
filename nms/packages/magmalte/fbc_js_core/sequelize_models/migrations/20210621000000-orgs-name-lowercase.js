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
 * This migration changes all organization names to be lower case.
 * There are various ways in which Sequelize requires some customization
 * to integrate well with the idiosyncrasies of Postgres.
 *
 * One way to simplify this is to change string fields to be lower-case,
 * as in organization names, which is done here.
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

        // Check there won't be any duplicate organizations
        let query = '';
        let uniqueCount = 0;
        let totalCount = 0;
        try {
          query = `SELECT COUNT(DISTINCT LOWER(name)) FROM ${quote}Organizations${quote};`;
          let [results, _metadata] = await queryInterface.sequelize.query(
            query,
            {transaction},
          );
          uniqueCount = results[0].count;
          query = `SELECT COUNT(name) FROM ${quote}Organizations${quote};`;
          [results, _metadata] = await queryInterface.sequelize.query(query, {
            transaction,
          });
          totalCount = results[0].count;
        } catch (err) {
          console.error(
            `Failed to run query for migration: ${query}, error: ${err}`,
          );
        }
        if (uniqueCount < totalCount) {
          console.error(
            `There are ${totalCount} organizations and ${uniqueCount} unique organization names. ` +
              'Make sure there are no matching organization names before trying migration again.',
          );
          throw 'Failed to complete migration';
        }

        // Update organizations to lower-case
        try {
          for (const tableName: string of Object.keys(TableOrgNameColumn)) {
            const orgColName: string = TableOrgNameColumn[tableName];
            query = `UPDATE ${quote}${tableName}${quote} SET ${orgColName}=lower(${orgColName})`;
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

  down: (_queryInterface: QueryInterface, _types: DataTypes) => {
    return Promise.resolve();
  },
};
