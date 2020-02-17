/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
const express = require('express');
const proxy = require('http-proxy-middleware');
const {GRAPH_HOST} = require('../config');
import {accessRoleToString} from '@fbcnms/auth/roles';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

const router = express.Router();

router.use(
  '/',
  proxy({
    // hostname to the target server
    target: 'http://' + GRAPH_HOST,

    // enable websocket proxying
    ws: true,

    // needed for virtual hosted sites
    changeOrigin: true,

    // rewrite paths
    pathRewrite: (path: string): string => path.replace(/^\/graph/, ''),

    // subscribe to http-proxy's proxyReq event
    onProxyReq: (proxyReq, req: FBCNMSRequest): void => {
      proxyReq.setHeader('x-auth-organization', req.user.organization);
      proxyReq.setHeader('x-auth-user-email', req.user.email);
      proxyReq.setHeader('x-auth-user-role', accessRoleToString(req.user.role));
    },
  }),
);

module.exports = router;
