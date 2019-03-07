/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DataTypes} from 'sequelize';

import Sequelize from 'sequelize';

export default (sequelize: Sequelize, types: DataTypes) => {
  const Organization = sequelize.define(
    'Organization',
    {
      name: types.STRING,
      tabs: {
        type: types.JSON,
        allowNull: false,
        defaultValue: [],
      },
      customDomains: {
        type: types.JSON,
        allowNull: false,
        defaultValue: [],
        get() {
          return this.getDataValue('customDomains') || [];
        },
      },
      networkIDs: {
        type: types.JSON,
        allowNull: false,
        defaultValue: [],
      },
    },
    {},
  );
  Organization.associate = function(_models) {
    // associations can be defined here
  };
  return Organization;
};
