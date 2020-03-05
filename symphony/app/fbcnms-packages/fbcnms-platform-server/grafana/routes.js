/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import express from 'express';
import proxy from 'express-http-proxy';

const GRAFANA_PROTOCOL = 'http';
const GRAFANA_HOST = 'user-grafana';
const GRAFANA_PORT = 3000;
const GRAFANA_URL = `${GRAFANA_PROTOCOL}://${GRAFANA_HOST}:${GRAFANA_PORT}`;

const AUTH_PROXY_HEADER = 'X-WEBAUTH-USER';

const router = express.Router();

router.use(
  '*',
  proxy(GRAFANA_URL, {
    proxyReqOptDecorator: function(proxyReqOpts, _srcReq) {
      proxyReqOpts.headers[AUTH_PROXY_HEADER] = 'admin';
      return proxyReqOpts;
    },
    proxyReqPathResolver: req => req.originalUrl.replace(/^\/grafana/, ''),
  }),
);

export default router;
