/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {Middleware} from 'express';

import express from 'express';

import csrf from 'csurf';

export default function csrfMiddleware(): Middleware {
  const router = express.Router();
  router.use(csrf());
  return router;
}
