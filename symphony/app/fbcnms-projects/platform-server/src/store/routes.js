/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import type {FBCNMSRequest} from '@fbcnms/auth/access';

const express = require('express');
const proxy = require('express-http-proxy');

const {STORE_HOST} = require('../config');

const router = express.Router();

const PROXY_OPTIONS = {
  proxyReqPathResolver: (req: FBCNMSRequest) =>
    req.originalUrl.replace(/^\/store/, ''),
  proxyReqOptDecorator: function(proxyReqOpts, srcReq: FBCNMSRequest) {
    proxyReqOpts.headers['x-auth-organization'] = srcReq.user.organization;
    return proxyReqOpts;
  },
};

router.use('/', proxy(STORE_HOST, PROXY_OPTIONS));

module.exports = router;
