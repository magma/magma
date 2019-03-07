/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

const fs = require('fs');

require('dotenv').config();

const DEV_MODE = process.env.NODE_ENV !== 'production';
const LOG_LEVEL = process.env.LOG_LEVEL || 'info';
const LOG_FORMAT = DEV_MODE ? 'shell' : 'json';

const apiHostArg = process.env.API_HOST || 'staging';
let API_HOST;
switch (apiHostArg) {
  // shortcuts to make management a little easier
  case 'staging':
    API_HOST = 'api-staging.magma.test';
    break;
  case 'prod':
    API_HOST = 'api.magma.test';
    break;
  default:
    API_HOST = process.env.API_HOST;
}

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
