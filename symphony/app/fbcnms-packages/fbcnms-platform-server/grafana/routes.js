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

import {
  HandleNewDatasource,
  HandleNewGrafanaUser,
  HandleSyncOrganizations,
  ORC8R_DATASOURCE_NAME,
  makeGrafanaUsername,
} from './handlers';

import type {FBCNMSRequest} from '@fbcnms/auth/access';
import type {GetDatasourcesResponse} from './GrafanaAPIType';
import type {GrafanaResponse} from './GrafanaAPI';

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
        const err = await HandleNewGrafanaUser(grafanaAdminClient, req);
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

const checkOrchestratorDatasource = () => {
  return async function(req: FBCNMSRequest, res, next) {
    // Get Grafana OrgID from this organization's Name
    const nmsOrg = await req.organization();
    const orgResp = await grafanaAdminClient.getOrg(nmsOrg.name);
    const orgIDForUser = orgResp.data.id;

    // Check if this organization has a datasource for Orchestrator
    const getDSResp: GrafanaResponse<GetDatasourcesResponse> = await grafanaAdminClient.getDatasources(
      orgIDForUser,
    );
    for (const ds of getDSResp.data) {
      if (
        ds.orgId == orgIDForUser &&
        ds.name.startsWith(ORC8R_DATASOURCE_NAME)
      ) {
        return next();
      }
    }

    // If not, create orchestrator datasource
    const err = await HandleNewDatasource(grafanaAdminClient, req);
    if (err) {
      const strData = JSON.stringify(err.response.data) || '';
      return res
        .status(err.response.status)
        .send(err.message + strData)
        .end();
    }
    return next();
  };
};

// Ensure that organizations in the NMS are in sync with
// tenants in Orchestrator
const syncOrganizations = () => {
  return async function(req: FBCNMSRequest, res, next) {
    const err = await HandleSyncOrganizations();
    if (err) {
      const strData = JSON.stringify(err.response.data) || '';
      return res
        .status(err.response.status)
        .send(err.message + strData)
        .end();
    }
    next();
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
      userResDecorator: (proxyRes, proxyResData, userReq, userRes) => {
        userRes.set('X-Frame-Options', 'allow');
        return proxyResData;
      },
    })(req, res, next);
  };
};

// Only the root path should check for Grafana User
router.all('/', checkGrafanaUser());
router.all('/', syncOrganizations());
router.all('/', checkOrchestratorDatasource());
// Use proxy on all paths
router.use('/', proxyMiddleware());

export default router;
