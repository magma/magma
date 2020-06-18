/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {DataTypes, QueryInterface} from 'sequelize';

module.exports = {
  up: (queryInterface: QueryInterface, types: DataTypes) => {
    return queryInterface.createTable('Organizations', {
      id: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: types.INTEGER,
      },
      customDomains: {
        allowNull: false,
        defaultValue: '[]',
        type: types.JSON,
      },
      name: {
        allowNull: false,
        type: types.STRING,
      },
      networkIDs: {
        allowNull: true,
        defaultValue: '[]',
        type: types.JSON,
      },
      tabs: {
        allowNull: false,
        defaultValue: '[]',
        type: types.JSON,
      },
      createdAt: {
        allowNull: false,
        type: types.DATE,
      },
      updatedAt: {
        allowNull: false,
        type: types.DATE,
      },
    });
  },

  down: (queryInterface: QueryInterface, _types: DataTypes) => {
    return queryInterface.dropTable('Organizations');
  },
};
