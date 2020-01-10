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
import {AccessRoles} from '@fbcnms/auth/roles';

const {GRAPH_HOST} = require('../config');

const router = express.Router();

const proxyMiddleware = () => {
  return function(req: FBCNMSRequest, res, next) {
    let reqAsBuffer = false;
    let reqBodyEncoding = true;
    let parseReqBody = true;
    const contentType = req.headers['content-type'];
    if (contentType && contentType.indexOf('multipart') != -1) {
      reqAsBuffer = true;
      reqBodyEncoding = null;
      parseReqBody = false;
    }
    return proxy(GRAPH_HOST, {
      reqAsBuffer,
      reqBodyEncoding,
      parseReqBody,
      proxyReqPathResolver: req => req.originalUrl.replace(/^\/graph/, ''),
      proxyReqOptDecorator: async function(proxyReqOpts, srcReq) {
        const organization = await srcReq.organization();
        proxyReqOpts.headers['x-auth-organization'] = organization.name;
        proxyReqOpts.headers['x-auth-user-email'] = srcReq.user.email;
        proxyReqOpts.headers['x-auth-user-readonly'] =
          srcReq.user.role === AccessRoles.READ_ONLY_USER ? 'TRUE' : 'FALSE';

        return proxyReqOpts;
      },
    })(req, res, next);
  };
};

router.use('/', proxyMiddleware());

module.exports = router;
