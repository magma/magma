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

/**
 * This migration removes 'tabs' from organizations and users.
 */
module.exports = {
  up: async (queryInterface: QueryInterface, _types: DataTypes) => {
    return queryInterface.sequelize.transaction(
      async (transaction: Transaction): Promise<void> => {
        await queryInterface.removeColumn('Users', 'tabs', {transaction});
        await queryInterface.removeColumn('Organizations', 'tabs', {
          transaction,
        });
      },
    );
  },

  down: async (queryInterface: QueryInterface, types: DataTypes) => {
    return queryInterface.sequelize.transaction(
      async (transaction: Transaction): Promise<void> => {
        await queryInterface.addColumn(
          'Users',
          'tabs',
          {
            allowNull: true,
            defaultValue: ['nms'],
            type: types.JSON,
          },
          {transaction},
        );

        await queryInterface.addColumn(
          'Organizations',
          'tabs',
          {
            allowNull: true,
            defaultValue: ['nms'],
            type: types.JSON,
          },
          {transaction},
        );
      },
    );
  },
};
