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
import workflowRouter from './routes';

const app = ExpressApplication();

app.use('/', workflowRouter);

app.listen(80);
