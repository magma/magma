/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

const logger = require('@fbcnms/logging').getLogger(module);
const {sequelize} = require('../server/models');

const path = require('path');
const {DataTypes} = require('sequelize');
const Umzug = require('umzug');

const umzug = new Umzug({
  storage: 'sequelize',
  storageOptions: {
    sequelize,
  },
  // The logging function.
  // A function that gets executed everytime migrations start and have ended.
  logging: msg => logger.info(msg),
  // The name of the positive method in migrations.
  upName: 'up',
  // The name of the negative method in migrations.
  downName: 'down',
  migrations: {
    // The params that gets passed to the migrations.
    // Might be an array or a synchronous function which returns an array.
    params: [sequelize.getQueryInterface(), DataTypes],
    // The path to the migrations directory.
    path: path.join(__dirname, '..', 'server/migrations'),
    // The pattern that determines whether or not a file is a migration.
    pattern: /^\d+[\w-]+\.js$/,
    // A function that receives and returns the to be executed function.
    // This can be used to modify the function.
    wrap(func) {
      return func;
    },
  },
});

export async function runMigrations() {
  const pendingMigrations = await umzug.pending();
  if (pendingMigrations) {
    await umzug.up();
  }
  // Sync defined models to the DB
  await sequelize.sync();
}

export async function rollbackMigrations() {
  const executedMigrations = await umzug.executed();
  if (executedMigrations) {
    await umzug.down();
  }

  // Sync defined models to the DB
  await sequelize.sync();
}

export default umzug;
