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

const app = ExpressApplication();

async function init() {
  const proxyTarget =
    process.env.PROXY_TARGET || 'http://conductor-server:8080';
  const proxyRouter = await proxy(proxyTarget);

  app.use('/', workflowRouter);
  app.use('/proxy', proxyRouter);
  app.listen(80);
}

init();
