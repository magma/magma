/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export {default as organizationMiddleware} from './organizationMiddleware';
export {default as appMiddleware} from './appMiddleware';
export {default as csrfMiddleware} from './csrfMiddleware';
export {default as sessionMiddleware} from './sessionMiddleware';
export {default as webpackSmartMiddleware} from './webpackSmartMiddleware';

import type {OrganizationMiddlewareRequest} from './organizationMiddleware';

export type FBCNMSMiddleWareRequest = {
  csrfToken: () => string, // from csrf
  body: Object, // from bodyParser
  session: any,
} & OrganizationMiddlewareRequest;
