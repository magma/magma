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
import {B3Format} from '@opencensus/propagation-b3';

import {JaegerTraceExporter} from '@opencensus/exporter-jaeger';
import {PrometheusStatsExporter} from '@opencensus/exporter-prometheus';
import {globalStats} from '@opencensus/core';
import {registerAllViews} from '@opencensus/instrumentation-http';
import {runMigrations} from './runMigrations';

const {sequelize} = require('@fbcnms/sequelize-models');
const logger = logging.getLogger(module);
const port = parseInt(process.env.PORT || 80);

// Configure metrics
const prometheusExporter = new PrometheusStatsExporter({
  startServer: true,
  logger,
});
globalStats.registerExporter(prometheusExporter);
registerAllViews(globalStats);

let jaegerExporter = null;
if (process.env.TELEMETRY_TRACE_EXPORTER == 'jaeger') {
  if (
    !process.env.JAEGER_AGENT_ENDPOINT &&
    !process.env.JAEGER_COLLECTOR_ENDPOINT
  ) {
    throw new Error(
      'When using TELEMETRY_TRACE_EXPORTER = "jaeger", you ' +
        'must set either JAEGER_AGENT_ENDPOINT or JAEGER_COLLECTOR_ENDPOINT',
    );
  }
  // Configure opencensus for jaeger
  const agentEndpoint = process.env.JAEGER_AGENT_ENDPOINT;
  const [agentHost, agentPort] = agentEndpoint
    ? agentEndpoint.split(':')
    : ['', ''];
  const jaegerOptions = {
    serviceName: process.env.TELEMETRY_TRACE_SERVICE || 'front',
    host: agentHost,
    port: agentPort,
    tags: [{key: 'opencensus-exporter-jaeger', value: '0.0.22'}],
    bufferTimeout: 10000, // time in milliseconds
    logger,
  };
  jaegerExporter = new JaegerTraceExporter(jaegerOptions);

  if (process.env.JAEGER_COLLECTOR_ENDPOINT) {
    // This is a hack for when using JAEGER_COLLECTOR_ENDPOINT, because
    // otherwise it uses UDPSender.  This will be updated as soon as the
    // opencensus library upgrades jaeger-client
    const HTTPSender = require('jaeger-client/dist/src/reporters/http_sender')
      .default;
    jaegerExporter.sender = new HTTPSender({
      ...jaegerOptions,
      endpoint: process.env.JAEGER_COLLECTOR_ENDPOINT,
    });
    jaegerExporter.sender.setProcess(jaegerExporter.process);
  }
}

tracing.start({
  samplingRate: 1,
  plugins: {
    http: '@opencensus/instrumentation-http',
  },
  stats: globalStats,
  exporter: jaegerExporter,
  propagation: jaegerExporter ? new B3Format() : null,
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
