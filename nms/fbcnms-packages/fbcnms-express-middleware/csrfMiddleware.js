/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  ExpressRequest,
  ExpressResponse,
  Middleware,
  NextFunction,
} from 'express';
import type {FBCNMSMiddleWareRequest} from './index';

import express from 'express';

import csrf from 'csurf';

export default function csrfMiddleware(): Middleware {
  const router = express.Router();
  router.use(
    csrf({
      cookie: true,
      value: (req: ExpressRequest) => req.cookies.csrfToken,
    }),
  );
  router.use(
    (
      req: FBCNMSMiddleWareRequest,
      res: ExpressResponse,
      next: NextFunction,
    ) => {
      res.cookie('csrfToken', req.csrfToken ? req.csrfToken() : '', {
        sameSite: true,
        httpOnly: true,
      });
      next();
    },
  );
  return router;
}
