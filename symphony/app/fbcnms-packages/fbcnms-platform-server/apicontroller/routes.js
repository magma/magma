/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

const express = require('express');
const proxy = require('express-http-proxy');
const HttpsProxyAgent = require('https-proxy-agent');
const url = require('url');
const {apiCredentials, API_HOST} = require('../config');
import auditLoggingDecorator from './auditLoggingDecorator';

import {intersection} from 'lodash';

const router: express.Router<FBCNMSRequest, ExpressResponse> = express.Router();

const PROXY_TIMEOUT_MS = 30000;

let agent = null;
if (process.env.HTTPS_PROXY) {
  const options = url.parse(process.env.HTTPS_PROXY);
  agent = new HttpsProxyAgent(options);
}
const PROXY_OPTIONS = {
  https: true,
  memoizeHost: false,
  timeout: PROXY_TIMEOUT_MS,
  proxyReqOptDecorator: (proxyReqOpts, _originalReq) => {
    return {
      ...proxyReqOpts,
      agent: agent,
      cert: apiCredentials().cert,
      key: apiCredentials().key,
      rejectUnauthorized: false,
    };
  },
  proxyReqPathResolver: req =>
    req.originalUrl.replace(/^\/nms\/apicontroller/, ''),
};

export async function networkIdFilter(req: FBCNMSRequest): Promise<boolean> {
  if (req.organization) {
    const organization = await req.organization();

    // If the request isn't an organization network, block
    // the request
    const isOrganizationAllowed = containsNetworkID(
      organization.networkIDs,
      req.params.networkID,
    );
    if (!isOrganizationAllowed) {
      return false;
    }
  }

  // super users on standalone deployments
  // have access to all proxied API requests
  // for the organization
  if (req.user.isSuperUser) {
    return true;
  }
  return containsNetworkID(req.user.networkIDs, req.params.networkID);
}

export async function networksResponseDecorator(
  _proxyRes: ExpressResponse,
  proxyResData: Buffer,
  userReq: FBCNMSRequest,
  _userRes: ExpressResponse,
) {
  let result = JSON.parse(proxyResData.toString('utf8'));
  if (userReq.organization) {
    const organization = await userReq.organization();
    result = intersection(result, organization.networkIDs);
  }
  if (!userReq.user.isSuperUser) {
    // the list of networks is further restricted to what the user
    // is allowed to see
    result = intersection(result, userReq.user.networkIDs);
  }
  return JSON.stringify(result);
}

const containsNetworkID = function(
  allowedNetworkIDs: string[],
  networkID: string,
): boolean {
  return (
    allowedNetworkIDs.indexOf(networkID) !== -1 ||
    // Remove secondary condition after T34404422 is addressed. Reason:
    //   Request needs to be lower cased otherwise calling
    //   MagmaAPIUrls.gateways() potentially returns missing devices.
    allowedNetworkIDs
      .map(id => id.toString().toLowerCase())
      .indexOf(networkID.toString().toLowerCase()) !== -1
  );
};

const proxyErrorHandler = (err, res, next) => {
  if (err.code === 'ENOTFOUND') {
    res.status(503).send('Cannot reach Orchestrator server');
  } else {
    next();
  }
};

router.use(
  /^\/magma\/v1\/networks$/,
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    userResDecorator: networksResponseDecorator,
    proxyErrorHandler,
  }),
);

router.use(
  '/magma/v1/networks/:networkID',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: networkIdFilter,
    userResDecorator: auditLoggingDecorator,
    proxyErrorHandler,
  }),
);

const networkTypeRegex = '(cwf|feg|lte|symphony|wifi)';
router.use(
  `/magma/v1/:networkType(${networkTypeRegex})/:networkID`,
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: networkIdFilter,
    userResDecorator: auditLoggingDecorator,
    proxyErrorHandler,
  }),
);

router.use(
  '/magma/channels/:channel',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: (req, _res) => req.method === 'GET',
  }),
);

router.use(
  '/magma/v1/channels/:channel',
  proxy(API_HOST, {
    ...PROXY_OPTIONS,
    filter: (req, _res) => req.method === 'GET',
  }),
);

router.use('', (req: FBCNMSRequest, res: ExpressResponse) => {
  res.status(404).send('Not Found');
});

export default router;
