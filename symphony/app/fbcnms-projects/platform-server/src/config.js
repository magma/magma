/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

require('dotenv').config();

const {getValidLogLevel} = require('@fbcnms/logging');

const DEV_MODE = process.env.NODE_ENV !== 'production';
const LOG_FORMAT = DEV_MODE ? 'shell' : 'json';
const LOG_LEVEL = getValidLogLevel(process.env.LOG_LEVEL);

const MAPBOX_ACCESS_TOKEN = process.env.MAPBOX_ACCESS_TOKEN || '';

const GRAPH_HOST = process.env.GRAPH_HOST || 'graph';
const STORE_HOST = process.env.STORE_HOST || 'store';
const DOCS_HOST = process.env.DOCS_HOST || 'docs';
const ID_HOST = process.env.ID_HOST || 'id';

module.exports = {
  DEV_MODE,
  GRAPH_HOST,
  LOG_FORMAT,
  LOG_LEVEL,
  MAPBOX_ACCESS_TOKEN,
  STORE_HOST,
  DOCS_HOST,
  ID_HOST,
};
