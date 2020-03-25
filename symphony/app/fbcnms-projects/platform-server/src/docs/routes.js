/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {FBCNMSRequest} from '@fbcnms/auth/access';

const express = require('express');
const proxy = require('express-http-proxy');

const {DOCS_HOST} = require('../config');

const router = express.Router();

const PROXY_OPTIONS = {
  proxyReqPathResolver: (req: FBCNMSRequest) =>
    req.originalUrl.replace(/^\/docs/, ''),
  proxyReqOptDecorator: function (proxyReqOpts, srcReq: FBCNMSRequest) {
    proxyReqOpts.headers['x-auth-organization'] = srcReq.user.organization;
    return proxyReqOpts;
  },
};

router.use('/', proxy(DOCS_HOST, PROXY_OPTIONS));

module.exports = router;
