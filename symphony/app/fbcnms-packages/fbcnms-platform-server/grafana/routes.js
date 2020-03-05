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

import Client from './GrafanaAPI';

import {HandleNewGrafanaUser, makeGrafanaUsername} from './handlers';

import type {FBCNMSRequest} from '@fbcnms/auth/access';
import type {GrafanaError} from './handlers';

const GRAFANA_PROTOCOL = 'http';
const GRAFANA_HOST = process.env.USER_GRAFANA_HOSTNAME || 'user-grafana';
const GRAFANA_PORT = process.env.USER_GRAFANA_PORT || '3000';
const GRAFANA_URL = `${GRAFANA_PROTOCOL}://${GRAFANA_HOST}:${GRAFANA_PORT}`;

const AUTH_PROXY_HEADER = 'X-WEBAUTH-USER';

const router = express.Router();

const grafanaAdminClient = Client(GRAFANA_URL, {
  [AUTH_PROXY_HEADER]: 'admin',
});

// Check that the NMS user and Org has been added to Grafana
const checkGrafanaUser = () => {
  return async function(req: FBCNMSRequest, res, next) {
    const userName = makeGrafanaUsername(req.user.id);
    const getUserResp = await grafanaAdminClient.getUser(userName);
    switch (getUserResp.status) {
      case 200:
        return next();
      case 404:
        const err: ?GrafanaError = await HandleNewGrafanaUser(
          grafanaAdminClient,
          req,
        );
        if (err) {
          const strData = JSON.stringify(err.response.data) || '';
          return res
            .status(err.response.status)
            .send(err.message + strData)
            .end();
        }
        return next();
      default:
        return res
          .status(getUserResp.status)
          .send(
            'Unexpected error getting user:' + JSON.stringify(getUserResp.data),
          )
          .end();
    }
  };
};

const proxyMiddleware = () => {
  return async function(req: FBCNMSRequest, res, next) {
    const userID = req.user.id;

    return proxy(GRAFANA_URL, {
      proxyReqOptDecorator: function(proxyReqOpts, _srcReq) {
        proxyReqOpts.headers[AUTH_PROXY_HEADER] = makeGrafanaUsername(userID);
        return proxyReqOpts;
      },
      proxyReqPathResolver: req => req.originalUrl.replace(/^\/grafana/, ''),
    })(req, res, next);
  };
};

// Only the root path should check for Grafana User
router.all('/', checkGrafanaUser());
// Use proxy on all paths
router.use('/', proxyMiddleware());

export default router;
