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
    return queryInterface.addColumn('Organizations', 'csvCharset', {
      allowNull: true,
      defaultValue: '',
      type: types.STRING,
    });
  },

  down: (queryInterface: QueryInterface, _: DataTypes) => {
    return queryInterface.removeColumn('Organizations', 'csvCharset');
  },
};
