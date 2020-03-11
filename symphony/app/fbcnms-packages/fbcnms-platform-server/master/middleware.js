/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ExpressResponse, NextFunction} from 'express';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

import asyncHandler from '@fbcnms/util/asyncHandler';

export const masterOrgMiddleware = asyncHandler(
  async (req: FBCNMSRequest, res: ExpressResponse, next: NextFunction) => {
    if (req.organization) {
      const organization = await req.organization();
      if (organization.isMasterOrg) {
        next();
        return;
      }
    }

    return res.redirect(req.access?.loginUrl || '/user/login');
  },
);
