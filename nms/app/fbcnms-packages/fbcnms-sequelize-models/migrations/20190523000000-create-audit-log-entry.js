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
    return queryInterface.createTable('AuditLogEntries', {
      id: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: types.INTEGER,
      },
      actingUserId: {
        type: types.INTEGER,
        allowNull: false,
      },
      organization: {
        type: types.STRING,
        allowNull: false,
      },
      mutationType: {
        type: types.STRING,
        allowNull: false,
      },
      objectId: {
        type: types.STRING,
        allowNull: false,
      },
      objectDisplayName: {
        type: types.STRING,
        allowNull: false,
      },
      objectType: {
        type: types.STRING,
        allowNull: false,
      },
      mutationData: {
        type: types.JSON,
        allowNull: false,
      },
      url: {
        type: types.STRING,
        allowNull: false,
      },
      ipAddress: {
        type: types.STRING,
        allowNull: false,
      },
      status: {
        type: types.STRING,
        allowNull: false,
      },
      statusCode: {
        type: types.STRING,
        allowNull: false,
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
    return queryInterface.dropTable('AuditLogEntries');
  },
};
