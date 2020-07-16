/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {Middleware} from 'express';

export default function asyncHandler(fn: Middleware): Middleware {
  return (req, res, next) => Promise.resolve(fn(req, res, next)).catch(next);
}
