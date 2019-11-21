/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const morgan = require('morgan');
const winston = require('winston');
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

module.exports = {
  configure(options: Options) {
    Object.assign(globalOptions, options);
  },
  getHttpLogger: (callingModule: any) => {
    const logger = module.exports.getLogger(callingModule);
    return morgan('combined', {
      skip: (req, _) => req.baseUrl == '/healthz',
      stream: {
        write: message => {
          logger.info(message);
        },
      },
    });
  },
  getLogger: (callingModule: any): $winstonLogger<$winstonNpmLogLevels> => {
    return winston.createLogger({
      level: globalOptions.LOG_LEVEL,
      format: getLogFormat(callingModule),
      stderrLevels: ['error', 'warning'],
      transports: [new winston.transports.Console()],
    });
  },
  getValidLogLevel: (logLevel: ?string): $Keys<$winstonNpmLogLevels> => {
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
  },
};
