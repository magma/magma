/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import {TokenSet} from 'openid-client';
import {clientFromRequest} from './client';

import type {ExpressResponse, NextFunction} from 'express';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

// An OIDC middleware that will refresh an access token if it exists.
// If it's expired and can't be refreshed, the user will be logged out
export function oidcAuthMiddleware() {
  return async function access(
    req: FBCNMSRequest,
    res: ExpressResponse,
    next: NextFunction,
  ) {
    try {
      const passportTokenSet = req.session?.oidc?.tokenSet;
      if (!passportTokenSet) {
        next();
        return;
      }

      const tokenSet = new TokenSet(passportTokenSet);
      if (!tokenSet.expired()) {
        next();
        return;
      }

      const client = await clientFromRequest(req);
      const newToken = await client.refresh(tokenSet.refresh_token);
      req.session.oidc = {tokenSet: newToken};
      next();
    } catch (error) {
      if (error.name === 'OpenIdConnectError') {
        if (error.error === 'invalid_grant') {
          req.logout();
          delete req.session.oidc;
          res.redirect('/');
        }
      }
      throw error;
    }
  };
}
