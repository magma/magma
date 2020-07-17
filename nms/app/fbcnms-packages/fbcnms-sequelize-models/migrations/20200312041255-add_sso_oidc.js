/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {DataTypes, QueryInterface, Transaction} from 'sequelize';

// $FlowFixMe: Op is in sequelize 5, but only sequelize 4 flow types exist
import {Op} from 'sequelize';

module.exports = {
  up: (queryInterface: QueryInterface, Sequelize: DataTypes) => {
    return queryInterface.sequelize.transaction(
      (transaction: Transaction): Promise<void[]> =>
        Promise.all([
          queryInterface.addColumn(
            'Organizations',
            'ssoSelectedType',
            {
              allowNull: false,
              defaultValue: 'none',
              type: Sequelize.ENUM('none', 'oidc', 'saml'),
            },
            {transaction},
          ),
          queryInterface.addColumn(
            'Organizations',
            'ssoOidcClientID',
            {
              allowNull: false,
              defaultValue: '',
              type: Sequelize.STRING,
            },
            {transaction},
          ),
          queryInterface.addColumn(
            'Organizations',
            'ssoOidcClientSecret',
            {
              allowNull: false,
              defaultValue: '',
              type: Sequelize.STRING,
            },
            {transaction},
          ),
          queryInterface.addColumn(
            'Organizations',
            'ssoOidcConfigurationURL',
            {
              allowNull: false,
              defaultValue: '',
              type: Sequelize.STRING,
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

  down: (queryInterface: QueryInterface, _Sequelize: DataTypes) => {
    return queryInterface.sequelize.transaction(
      (transaction: Transaction): Promise<void[]> =>
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
