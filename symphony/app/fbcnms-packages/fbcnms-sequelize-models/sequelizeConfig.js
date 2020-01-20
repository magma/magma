/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Options} from 'sequelize';

// TODO: Pull from shared config
const MYSQL_HOST = process.env.MYSQL_HOST || '127.0.0.1';
const MYSQL_PORT = parseInt(process.env.MYSQL_PORT || '3306');
const MYSQL_USER = process.env.MYSQL_USER || 'root';
const MYSQL_PASS = process.env.MYSQL_PASS || '';
const MYSQL_DB = process.env.MYSQL_DB || 'cxl';

const logger = require('@fbcnms/logging').getLogger(module);

const config: {[string]: Options} = {
  test: {
    username: '',
    password: '',
    database: 'db',
    dialect: 'sqlite',
    logging: false,
  },
  development: {
    username: MYSQL_USER,
    password: MYSQL_PASS,
    database: MYSQL_DB,
    host: MYSQL_HOST,
    port: MYSQL_PORT,
    dialect: 'mysql',
    logging: (msg: string) => logger.debug(msg),
  },
  production: {
    username: MYSQL_USER,
    password: MYSQL_PASS,
    database: MYSQL_DB,
    host: MYSQL_HOST,
    port: MYSQL_PORT,
    dialect: 'mysql',
    logging: (msg: string) => logger.debug(msg),
  },
};

export default config;
