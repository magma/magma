/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ExpressRequest, Middleware} from 'express';

import csrf from 'csurf';

export default function csrfMiddleware(): Middleware {
  return csrf({
    cookie: true,
    value: (req: ExpressRequest) => req.cookies.csrfToken,
  });
}
