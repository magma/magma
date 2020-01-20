/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const {LOG_FORMAT, LOG_LEVEL} = require('@fbcnms/platform-server/config');

// This must be done before any module imports to configure
// logging correctly
const logging = require('@fbcnms/logging');
logging.configure({
  LOG_FORMAT,
  LOG_LEVEL,
});

const {runMigrations} = require('./runMigrations');

runMigrations()
  .then(_ => console.log('Ran migrations successfully'))
  .catch(_ => console.error('Failed to run migrations'));
