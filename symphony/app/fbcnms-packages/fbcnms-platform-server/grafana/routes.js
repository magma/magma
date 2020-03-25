/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import ReactDOM from 'react-dom/server';
import express from 'express';
import proxy from 'express-http-proxy';

import Client from './GrafanaAPI';

import GrafanaErrorMessage from './GrafanaErrorMessage';
import {
  makeGrafanaUsername,
  syncDatasource,
  syncGrafanaUser,
  syncTenants,
} from './handlers';

import type {Task} from './handlers';

import type {FBCNMSRequest} from '@fbcnms/auth/access';

const GRAFANA_PROTOCOL = 'http';
const GRAFANA_HOST = process.env.USER_GRAFANA_HOSTNAME ?? 'user-grafana';
const GRAFANA_PORT = process.env.USER_GRAFANA_PORT ?? '3000';
const GRAFANA_URL = `${GRAFANA_PROTOCOL}://${GRAFANA_HOST}:${GRAFANA_PORT}`;

const AUTH_PROXY_HEADER = 'X-WEBAUTH-USER';

const router = express.Router();

const grafanaAdminClient = Client(GRAFANA_URL, {
  [AUTH_PROXY_HEADER]: 'admin',
});

const syncGrafana = () => {
  return async function (req: FBCNMSRequest, res, next) {
    const tasksCompleted = [];
    // Sync User/Organization
    const userRes = await syncGrafanaUser(grafanaAdminClient, req);
    tasksCompleted.push(...userRes.completedTasks);
    if (userRes.errorTask) {
      return await displayErrorMessage(res, tasksCompleted, userRes.errorTask);
    }
    // Sync Datasource
    const dsRes = await syncDatasource(grafanaAdminClient, req);
    tasksCompleted.push(...dsRes.completedTasks);
    if (dsRes.errorTask) {
      return await displayErrorMessage(res, tasksCompleted, dsRes.errorTask);
    }
    // Sync Tenants
    const tenantsRes = await syncTenants();
    tasksCompleted.push(...tenantsRes.completedTasks);
    if (tenantsRes.errorTask) {
      return await displayErrorMessage(
        res,
        tasksCompleted,
        tenantsRes.errorTask,
      );
    }
    return next();
  };
};

async function displayErrorMessage(
  res: ExpressResponse,
  completedTasks: Array<Task>,
  errorTask: Task,
) {
  const healthResponse = await grafanaAdminClient.getHealth();
  const message = (
    <GrafanaErrorMessage
      completedTasks={completedTasks}
      errorTask={errorTask}
      grafanaHealth={healthResponse.data}
    />
  );
  res.status(errorTask.status).send(ReactDOM.renderToString(message)).end();
}

const proxyMiddleware = () => {
  return async function (req: FBCNMSRequest, res, next) {
    const userID = req.user.id;

    return proxy(GRAFANA_URL, {
      proxyReqOptDecorator: function (proxyReqOpts, _srcReq) {
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

// Only the root path should perform the sync operations
router.all('/', syncGrafana());
// Use proxy on all paths
router.use('/', proxyMiddleware());

export default router;
