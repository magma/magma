/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

if (!process.env.NODE_ENV) {
  process.env.BABEL_ENV = 'development';
  process.env.NODE_ENV = 'development';
} else {
  process.env.BABEL_ENV = process.env.NODE_ENV;
}

import app from '../server/app';
import logging from '@fbcnms/logging';
import {runMigrations} from './runMigrations';

const logger = logging.getLogger(module);
const port = parseInt(process.env.PORT || 80);

(async function main() {
  await runMigrations();
  app.listen(port, '', err => {
    if (err) {
      logger.error(err.toString());
    }
    if (process.env.NODE_ENV === 'development') {
      logger.info(`Development server started on port ${port}`);
    } else {
      logger.info(`Production server started on port ${port}`);
    }
  });
})().catch(error => {
  logger.error(error);
});
