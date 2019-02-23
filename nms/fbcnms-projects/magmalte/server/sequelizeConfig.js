/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

const {
  MYSQL_HOST,
  MYSQL_PORT,
  MYSQL_USER,
  MYSQL_PASS,
  MYSQL_DB,
} = require('./config');
const logger = require('@fbcnms/logging').getLogger(module);

module.exports = {
  development: {
    username: MYSQL_USER,
    password: MYSQL_PASS,
    database: MYSQL_DB,
    host: MYSQL_HOST,
    port: MYSQL_PORT,
    dialect: 'mysql',
    logging: msg => logger.debug(msg),
  },
  production: {
    username: MYSQL_USER,
    password: MYSQL_PASS,
    database: MYSQL_DB,
    host: MYSQL_HOST,
    port: MYSQL_PORT,
    dialect: 'mysql',
    logging: msg => logger.debug(msg),
  },
};
