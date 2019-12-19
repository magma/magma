/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

jest.mock('../sequelizeConfig', () => {
  process.env.NODE_ENV = 'test';
  return {
    [process.env.NODE_ENV]: {
      username: null,
      password: null,
      database: 'db',
      dialect: 'sqlite',
      logging: false,
    },
  };
});

beforeAll(async () => {
  const {sequelize} = jest.requireActual('../');
  // running sync instead of migrations because of weird foreign key issues
  await sequelize.sync({force: true});
});

const realModels = jest.requireActual('../');
module.exports = realModels;
