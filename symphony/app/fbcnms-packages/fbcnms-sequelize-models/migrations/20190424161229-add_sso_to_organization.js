/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
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
