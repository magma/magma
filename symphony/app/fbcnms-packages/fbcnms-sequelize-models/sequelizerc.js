/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

require('@fbcnms/babel-register');

const path = require('path');

module.exports = {
  config: path.resolve(__dirname, 'sequelizeConfig.js'),
  'migrations-path': path.resolve(__dirname, 'migrations'),
  'models-path': path.resolve(__dirname, 'models'),
  'seeders-path': path.resolve(__dirname, 'seeders'),
};
