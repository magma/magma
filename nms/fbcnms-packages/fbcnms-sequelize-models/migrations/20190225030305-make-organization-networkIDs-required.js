/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DataTypes, QueryInterface} from 'sequelize';

module.exports = {
  up: (queryInterface: QueryInterface, types: DataTypes) => {
    return queryInterface.changeColumn('Organizations', 'networkIDs', {
      allowNull: false,
      defaultValue: '[]',
      type: types.JSON,
    });
  },

  down: (queryInterface: QueryInterface, types: DataTypes) => {
    return queryInterface.changeColumn('Organizations', 'networkIDs', {
      allowNull: true,
      defaultValue: '[]',
      type: types.JSON,
    });
  },
};
