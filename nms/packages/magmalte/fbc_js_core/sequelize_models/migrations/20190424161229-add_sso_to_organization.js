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

import type {DataTypes, QueryInterface, Transaction} from 'sequelize';

module.exports = {
  up: (queryInterface: QueryInterface, Sequelize: DataTypes) => {
    return queryInterface.sequelize.transaction(
      (transaction: Transaction): Promise<void[]> =>
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

  down: (queryInterface: QueryInterface, _Sequelize: DataTypes) => {
    return queryInterface.sequelize.transaction(
      (transaction: Transaction): Promise<void[]> =>
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
