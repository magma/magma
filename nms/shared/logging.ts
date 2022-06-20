/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow
 * @format
 */

import morgan from 'morgan';
import winston from 'winston';

import type {Request} from 'express';
import type {StreamOptions} from 'morgan';

const {
  colorize,
  combine,
  json,
  label,
  printf,
  splat,
  timestamp,
} = winston.format;

function getLabel(callingModule: NodeModule) {
  const parts = callingModule.filename.split('/');
  return `${parts[parts.length - 2]}/${parts.pop() || ''}`;
}

const myFormat = printf(info => {
  // eslint-disable-next-line @typescript-eslint/restrict-template-expressions
  return `${info.timestamp} [${info.label}] ${info.level}: ${info.message}`;
});

type Options = {
  LOG_FORMAT: 'json' | 'shell';
  LOG_LEVEL: string;
};

const globalOptions: Options = {
  LOG_FORMAT: 'json',
  LOG_LEVEL: 'info',
};

function getLogFormat(callingModule: NodeModule) {
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
export function getHttpLogger(callingModule: NodeModule) {
  const logger = getLogger(callingModule);
  const streamOptions: StreamOptions = {
    write: message => {
      logger.info(message);
    },
  };
  return morgan<Request>('combined', {
    skip: req => req.baseUrl == '/healthz',
    stream: streamOptions,
  });
}
export function getLogger(callingModule: NodeModule) {
  return winston.createLogger({
    level: globalOptions.LOG_LEVEL,
    format: getLogFormat(callingModule),
    transports: [new winston.transports.Console()],
  });
}
export function getValidLogLevel(logLevel: string | null | undefined) {
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
