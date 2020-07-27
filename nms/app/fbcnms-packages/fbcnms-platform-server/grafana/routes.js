/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
  syncDashboards,
  syncDatasource,
  syncGrafanaUser,
  syncTenants,
} from './handlers';

import type {Task} from './handlers';

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

const GRAFANA_PROTOCOL = 'http';
const GRAFANA_ADDRESS = process.env.USER_GRAFANA_ADDRESS ?? 'user-grafana:3000';
const GRAFANA_URL = `${GRAFANA_PROTOCOL}://${GRAFANA_ADDRESS}`;

const AUTH_PROXY_HEADER = 'X-WEBAUTH-USER';

const router: express.Router<FBCNMSRequest, ExpressResponse> = express.Router();

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
    // Create Dashboards
    const dbRes = await syncDashboards(grafanaAdminClient, req);
    tasksCompleted.push(...dbRes.completedTasks);
    if (dbRes.errorTask) {
      return await displayErrorMessage(res, tasksCompleted, dbRes.errorTask);
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
