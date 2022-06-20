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

import fs from 'fs';
import {getLogger} from '../../shared/logging';
import type {Dialect, Options} from 'sequelize';

// TODO: Pull from shared config
const MYSQL_HOST = process.env.MYSQL_HOST || '127.0.0.1';
const MYSQL_PORT = parseInt(process.env.MYSQL_PORT || '3306');
const MYSQL_USER = process.env.MYSQL_USER || 'root';
const MYSQL_PASS = process.env.MYSQL_PASS || '';
const MYSQL_DB = process.env.MYSQL_DB || 'cxl';
const MYSQL_DIALECT: Dialect =
  (process.env.MYSQL_DIALECT as Dialect) || 'mysql';

const logger = getLogger(module);
let ssl_required = false;
let CAcert: string | Buffer | undefined = process.env.CA_FILE;
let Ckey: string | Buffer | undefined = process.env.KEY_FILE;
let Ccert: string | Buffer | undefined = process.env.CERT_FILE;
let dialectOptions = {};

if (process.env.CA_FILE) {
  try {
    CAcert = fs.readFileSync(process.env.CA_FILE);
    ssl_required = true;
  } catch (e) {
    console.warn('cannot read ca cert file', e);
  }
}

if (process.env.KEY_FILE) {
  try {
    Ckey = fs.readFileSync(process.env.KEY_FILE);
    ssl_required = true;
  } catch (e) {
    console.warn('cannot read key file', e);
  }
}

if (process.env.CERT_FILE) {
  try {
    Ccert = fs.readFileSync(process.env.CERT_FILE);
    ssl_required = true;
  } catch (e) {
    console.warn('cannot read cert file', e);
  }
}

if (ssl_required) {
  dialectOptions = {
    ssl: {
      ca: CAcert,
      key: Ckey,
      cert: Ccert,
    },
  };
}

const config: {[config: string]: Options} = {
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
    dialect: MYSQL_DIALECT,
    ssl: ssl_required,
    dialectOptions,
    logging: (msg: string) => logger.debug(msg),
  },
  production: {
    username: MYSQL_USER,
    password: MYSQL_PASS,
    database: MYSQL_DB,
    host: MYSQL_HOST,
    port: MYSQL_PORT,
    dialect: MYSQL_DIALECT,
    ssl: ssl_required,
    dialectOptions,
    logging: (msg: string) => logger.debug(msg),
  },
};

export default config;
