/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */

import {TokenSet} from 'openid-client';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {clientFromRequest} from './client';

import type {ExpressRequest, ExpressResponse, NextFunction} from 'express';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {FBCNMSRequest} from '../access';

type OIDCTokenSet = {
  access_token: string,
};

type OIDCRequest = ExpressRequest & {
  session: {
    oidc?: {
      tokenSet: OIDCTokenSet,
    },
  },
};

export function oidcAccessToken(req: OIDCRequest) {
  return req.session?.oidc?.tokenSet?.access_token;
}

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
