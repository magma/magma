/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */
'use strict';

import {AccessRoles} from '@fbcnms/auth/roles';

module.exports = (sequelize, DataTypes) => {
  const User = sequelize.define(
    'User',
    {
      email: DataTypes.STRING,
      password: DataTypes.STRING,
      role: DataTypes.INTEGER,
      networkIDs: {
        type: DataTypes.JSON,
        allowNull: false,
        defaultValue: [],
        get() {
          return this.getDataValue('networkIDs') || [];
        },
      },
    },
    {
      getterMethods: {
        isSuperUser() {
          return this.role === AccessRoles.SUPERUSER;
        },
      },
    },
  );
  User.associate = function(_models) {
    // associations can be defined here
  };
  return User;
};
