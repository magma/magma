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

module.exports = {
  /**
   * @param {QueryInterface} queryInterface
   * @param {DataTypes} Sequelize
   */
  up: (queryInterface, Sequelize) => {
    return queryInterface.sequelize.transaction(transaction =>
      Promise.all([
        queryInterface.addColumn(
          'Organizations',
          'ssoCert',
          {
            allowNull: false,
            defaultValue: '',
            type: Sequelize.TEXT,
          },
          {transaction},
        ),
        queryInterface.addColumn(
          'Organizations',
          'ssoEntrypoint',
          {
            allowNull: false,
            defaultValue: '',
            type: Sequelize.STRING,
          },
          {transaction},
        ),
        queryInterface.addColumn(
          'Organizations',
          'ssoIssuer',
          {
            allowNull: false,
            defaultValue: '',
            type: Sequelize.STRING,
          },
          {transaction},
        ),
      ]),
    );
  },

  /**
   * @param {QueryInterface} queryInterface
   */
  down: queryInterface => {
    return queryInterface.sequelize.transaction(transaction =>
      Promise.all([
        queryInterface.removeColumn('Organizations', 'ssoEntrypoint', {
          transaction,
        }),
        queryInterface.removeColumn('Organizations', 'ssoCert', {
          transaction,
        }),
        queryInterface.removeColumn('Organizations', 'ssoDefaultNetworkIDs', {
          transaction,
        }),
      ]),
    );
  },
};
