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
const proxy = require('express-http-proxy');
const HttpsProxyAgent = require('https-proxy-agent');
const url = require('url');
const {apiCredentials, API_HOST, NETWORK_FALLBACK} = require('../config');

import type {ExpressResponse} from 'express';
import type {NMSRequest} from '../../scripts/server';

const router = express.Router();

let agent = null;
if (process.env.HTTPS_PROXY) {
  const options = url.parse(process.env.HTTPS_PROXY);
  agent = new HttpsProxyAgent(options);
}
const PROXY_OPTIONS = {
  https: true,
  memoizeHost: false,
  proxyReqOptDecorator: (proxyReqOpts, _originalReq: NMSRequest) => {
    return {
      ...proxyReqOpts,
      agent: agent,
      cert: apiCredentials().cert,
      key: apiCredentials().key,
      rejectUnauthorized: false,
    };
  },
  proxyReqPathResolver: (req: NMSRequest) =>
    req.originalUrl.replace(/^\/nms\/apicontroller/, ''),
};

router.use(
  '/magma/networks/:networkID',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: (req: NMSRequest) => {
      // super users have access to all proxied API requests
      if (req.user.isSuperUser) {
        return true;
      }

      return (
        req.user.networkIDs.indexOf(req.params.networkID) !== -1 ||
        // Remove secondary condition after T34404422 is addressed. Reason:
        //   Request needs to be lower cased otherwise calling
        //   MagmaAPIUrls.gateways() potentially returns missing devices.
        req.user.networkIDs
          .map(id => id.toLowerCase())
          .indexOf(req.params.networkID.toLowerCase()) !== -1
      );
    },
  }),
);

router.use(
  '/magma/networks',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    userResDecorator: (
      proxyRes: ExpressResponse,
      proxyResData: Buffer,
      userReq: NMSRequest,
      userRes: ExpressResponse,
    ) => {
      let networkIds;
      if (
        (proxyRes.statusCode === 403 || proxyRes.statusCode === 401) &&
        NETWORK_FALLBACK.length > 0
      ) {
        // Temporary hack -- if you don't have a root magma cert,
        // it will return a 403.
        userRes.statusCode = 200;
        networkIds = NETWORK_FALLBACK;
      } else {
        networkIds = JSON.parse(proxyResData.toString('utf8'));
      }

      if (userReq.user.isSuperUser) {
        return JSON.stringify(networkIds);
      }

      // if a normal user is fetching the list of networks from the Magma
      // controller we return the intersection of the list from the controller
      // with the networks they're allowed to see
      const allNetworkIDs = new Set();
      networkIds.map(id => allNetworkIDs.add(id));

      const results = userReq.user.networkIDs.filter(id =>
        allNetworkIDs.has(id),
      );
      return JSON.stringify(results);
    },
  }),
);

router.use(
  '/magma/channels/:channel',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: (req: NMSRequest, _res: ExpressResponse) => req.method === 'GET',
  }),
);

router.use('', (req: NMSRequest, res: ExpressResponse) => {
  res.status(404).send('Not Found');
});

module.exports = router;
