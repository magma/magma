/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

import ExpressApplication from 'express';
import proxy from './proxy/proxy';
import workflowRouter from './routes';
import {groupsForUser} from './proxy/graphqlGroups';

const app = ExpressApplication();

const OWNER_ROLE = 'OWNER';
const NETWORK_ADMIN_GROUP = 'network-admin';
const adminAccess = (role, groups) => {
  return role === OWNER_ROLE || groups.includes(NETWORK_ADMIN_GROUP);
};

const generalAccess = (role, groups) => {
  return true;
};
async function init() {
  const proxyTarget =
    process.env.PROXY_TARGET || 'http://conductor-server:8080';
  const schellarTarget = process.env.SCHELLAR_TARGET || 'http://schellar:3000';

  const proxyRouter = await proxy(proxyTarget, schellarTarget, adminAccess, groupsForUser);
  const rbacProxyRouter = await proxy(proxyTarget, schellarTarget, generalAccess, groupsForUser);

  app.use('/', workflowRouter);
  app.use('/proxy', proxyRouter);
  app.use("/rbac_proxy", rbacProxyRouter);
  app.listen(80);
}

init();
