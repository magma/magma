/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
'use strict';

import Sequelize from 'sequelize';
import {AccessRoles} from '@fbcnms/auth/roles';
import {omit} from 'lodash';

import type {AssociateProp} from './AssociateTypes.flow';
import type {DataTypes, Model} from 'sequelize';

// This is the type required for creation
type UserRawInitType = {
  email: string,
  organization?: string,
  password: string,
  role: number,
  networkIDs?: Array<string>,
  tabs?: Array<string>,
};

// This is the type read back
export type UserRawType = {
  id: number,
  networkIDs: Array<string>,
  isSuperUser: boolean,
  isReadOnlyUser: boolean,
  role: number,
} & UserRawInitType;

export type UserModel = Model<UserRawType, UserRawInitType>;
export type StaticUserModel = Class<UserModel>;
export type UserType = UserModel & UserRawType;

export default (
  sequelize: Sequelize,
  types: DataTypes,
): StaticUserModel & AssociateProp => {
  const User = sequelize.define(
    'User',
    {
      email: types.STRING,
      organization: types.STRING,
      password: types.STRING,
      role: types.INTEGER,
      networkIDs: {
        type: types.JSON,
        allowNull: false,
        defaultValue: [],
        get() {
          return this.getDataValue('networkIDs') || [];
        },
      },
      tabs: {
        type: types.JSON,
        allowNull: true,
        get() {
          return this.getDataValue('tabs') || [];
        },
      },
    },
    {
      getterMethods: {
        isSuperUser() {
          return this.role === AccessRoles.SUPERUSER;
        },
        isReadOnlyUser() {
          return this.role === AccessRoles.READ_ONLY_USER;
        },
      },
    },
  );
  User.associate = function (_models) {
    // associations can be defined here
  };
  User.prototype.toJSON = function () {
    return omit(this.get(), 'password');
  };
  return User;
};
