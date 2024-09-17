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
 */

import {DataTypes, Op} from 'sequelize';

module.exports = {
  /**
   * @param {{ context: QueryInterface}} params
   */
  up: ({context: queryInterface}) => {
    return queryInterface.sequelize.transaction(transaction =>
      Promise.all([
        queryInterface.addColumn(
          'Organizations',
          'ssoSelectedType',
          {
            allowNull: false,
            defaultValue: 'none',
            type: DataTypes.ENUM('none', 'oidc', 'saml'),
          },
          {transaction},
        ),
        queryInterface.addColumn(
          'Organizations',
          'ssoOidcClientID',
          {
            allowNull: false,
            defaultValue: '',
            type: DataTypes.STRING,
          },
          {transaction},
        ),
        queryInterface.addColumn(
          'Organizations',
          'ssoOidcClientSecret',
          {
            allowNull: false,
            defaultValue: '',
            type: DataTypes.STRING,
          },
          {transaction},
        ),
        queryInterface.addColumn(
          'Organizations',
          'ssoOidcConfigurationURL',
          {
            allowNull: false,
            defaultValue: '',
            type: DataTypes.STRING,
          },
          {transaction},
        ),
        queryInterface.bulkUpdate(
          'Organizations',
          {ssoSelectedType: 'none'},
          {},
          {transaction},
        ),
        queryInterface.bulkUpdate(
          'Organizations',
          {ssoSelectedType: 'saml'},
          {ssoEntrypoint: {[Op.ne]: ''}},
          {transaction},
        ),
      ]),
    );
  },

  /**
   * @param {{ context: QueryInterface}} params
   */
  down: ({context: queryInterface}) => {
    return queryInterface.sequelize.transaction(transaction =>
      Promise.all([
        queryInterface.removeColumn('Organizations', 'ssoSelectedType', {
          transaction,
        }),
        queryInterface.removeColumn('Organizations', 'ssoOidcClientID', {
          transaction,
        }),
        queryInterface.removeColumn('Organizations', 'ssoOidcClientSecret', {
          transaction,
        }),
        queryInterface.removeColumn(
          'Organizations',
          'ssoOidcConfigurationURL',
          {
            transaction,
          },
        ),
      ]),
    );
  },
};
