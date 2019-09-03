/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const fs = require('fs');
const nullthrows = require('@fbcnms/util/nullthrows').default;
const {getValidLogLevel} = require('@fbcnms/logging');

require('dotenv').config();

const DEV_MODE = process.env.NODE_ENV !== 'production';
const LOG_LEVEL = getValidLogLevel(process.env.LOG_LEVEL);
const LOG_FORMAT = DEV_MODE ? 'shell' : 'json';

const API_HOST = nullthrows(process.env.API_HOST);

const MAPBOX_ACCESS_TOKEN = process.env.MAPBOX_ACCESS_TOKEN || '';

const NETWORK_FALLBACK = process.env.NETWORK_FALLBACK
  ? process.env.NETWORK_FALLBACK.split(',')
  : [];

let _cachedApiCredentials = null;
function apiCredentials() {
  if (_cachedApiCredentials) {
    return _cachedApiCredentials;
  }

  const cert = process.env.API_CERT_FILENAME
    ? fs.readFileSync(process.env.API_CERT_FILENAME)
    : process.env.API_CERT;
  const key = process.env.API_PRIVATE_KEY_FILENAME
    ? fs.readFileSync(process.env.API_PRIVATE_KEY_FILENAME)
    : process.env.API_PRIVATE_KEY;

  return (_cachedApiCredentials = {
    cert,
    key,
  });
}

module.exports = {
  apiCredentials,
  API_HOST,
  DEV_MODE,
  LOG_FORMAT,
  LOG_LEVEL,
  MAPBOX_ACCESS_TOKEN,
  NETWORK_FALLBACK,
};
