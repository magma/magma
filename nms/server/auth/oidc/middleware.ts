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
 */

import {TokenSet} from 'openid-client';
import {clientFromRequest} from './client';

import type {NextFunction, Request, Response} from 'express';

export function oidcAccessToken(req: Request) {
  return req.session?.oidc?.tokenSet?.access_token;
}

// An OIDC middleware that will refresh an access token if it exists.
// If it's expired and can't be refreshed, the user will be logged out
export function oidcAuthMiddleware() {
  return async function access(
    req: Request,
    res: Response,
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
      req.session!.oidc = {tokenSet: newToken};
      next();
    } catch (error) {
      if ((error as Error).name === 'OpenIdConnectError') {
        if ((error as {error: string}).error === 'invalid_grant') {
          req.logout();
          delete req.session!.oidc;
          res.redirect('/');
        }
      }
      throw error;
    }
  };
}
