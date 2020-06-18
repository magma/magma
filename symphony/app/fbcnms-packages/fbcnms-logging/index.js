/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import morgan from 'morgan';
import winston from 'winston';

import type {Morgan, StreamOptions} from 'morgan';

const {
  colorize,
  combine,
  json,
  label,
  printf,
  splat,
  timestamp,
} = winston.format;

function getLabel(callingModule: any) {
  const parts = callingModule.filename.split('/');
  return parts[parts.length - 2] + '/' + parts.pop();
}

const myFormat = printf(info => {
  return `${info.timestamp} [${info.label}] ${info.level}: ${info.message}`;
});

type Options = {
  LOG_FORMAT: 'json' | 'shell',
  LOG_LEVEL: $Keys<$winstonNpmLogLevels>,
};

const globalOptions: Options = {
  LOG_FORMAT: 'json',
  LOG_LEVEL: 'info',
};

function getLogFormat(callingModule) {
  switch (globalOptions.LOG_FORMAT) {
    case 'json':
      return combine(
        label({label: getLabel(callingModule)}),
        timestamp(),
        splat(),
        json(),
      );
    case 'shell':
      return combine(
        colorize(),
        label({label: getLabel(callingModule)}),
        timestamp(),
        splat(),
        myFormat,
      );
  }
}

export function configure(options: Options) {
  Object.assign(globalOptions, options);
}
export function getHttpLogger(callingModule: any): Morgan {
  const logger = getLogger(callingModule);
  const streamOptions: StreamOptions = {
    write: message => {
      logger.info(message);
    },
  };
  return morgan('combined', {
    skip: (req, _) => req.baseUrl == '/healthz',
    stream: streamOptions,
  });
}
export function getLogger(
  callingModule: any,
): $winstonLogger<$winstonNpmLogLevels> {
  return winston.createLogger({
    level: globalOptions.LOG_LEVEL,
    format: getLogFormat(callingModule),
    stderrLevels: ['error', 'warning'],
    transports: [new winston.transports.Console()],
  });
}
export function getValidLogLevel(
  logLevel: ?string,
): $Keys<$winstonNpmLogLevels> {
  switch (logLevel) {
    case 'error':
    case 'warn':
    case 'info':
    case 'verbose':
    case 'debug':
    case 'silly':
      return logLevel;
    case undefined:
    case null:
      return 'info';
    default:
      throw new Error('Invalid log level!');
  }
}

// export default enables unnamed es6 imports to work,
// export function above enables CommonJS (require) imports to work
export default {configure, getHttpLogger, getLogger, getValidLogLevel};
