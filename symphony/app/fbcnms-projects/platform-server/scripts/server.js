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

import app from '../src/app';
import logging from '@fbcnms/logging';
import tracing from '@opencensus/nodejs';
import {PrometheusStatsExporter} from '@opencensus/exporter-prometheus';
import {globalStats} from '@opencensus/core';
import {registerAllViews} from '@opencensus/instrumentation-http';
import {runMigrations} from './runMigrations';

const {sequelize} = require('@fbcnms/sequelize-models');
const logger = logging.getLogger(module);
const port = parseInt(process.env.PORT || 80);

// Configure metrics
export const exporter = new PrometheusStatsExporter({
  startServer: true,
  logger,
});
globalStats.registerExporter(exporter);
registerAllViews(globalStats);

tracing.start({
  samplingRate: 1,
  plugins: {
    http: '@opencensus/instrumentation-http',
  },
  stats: globalStats,
});

(async function main() {
  for (;;) {
    try {
      await sequelize.authenticate();
      break;
    } catch (error) {
      logger.error('cannot connect to database', error);
      await new Promise(res => setTimeout(res, 2000));
    }
  }
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
