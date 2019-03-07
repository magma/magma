/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import OrganizationModel from './models/organization';
import UserModel from './models/user';
import Sequelize from 'sequelize';

const env = process.env.NODE_ENV || 'development';
const config = require('./sequelizeConfig.js')[env];

export const sequelize = new Sequelize(
  config.database,
  config.username,
  config.password,
  config,
);

const db = {
  Organization: OrganizationModel(sequelize, Sequelize),
  User: UserModel(sequelize, Sequelize),
};

Object.keys(db).forEach(
  modelName => db[modelName].associate && db[modelName].associate(db),
);

export const Organization = db.Organization;
export const User = db.User;
