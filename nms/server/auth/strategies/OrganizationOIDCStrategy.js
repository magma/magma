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
 * @flow
 * @format
 */

import DynamicStrategy from './DynamicStrategy';
// $FlowFixMe migrated to typescript
import {AccessRoles} from '../../../shared/roles';
import {Strategy as OidcStrategy} from 'openid-client';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {User} from '../../../shared/sequelize_models';

import {TokenSet} from 'openid-client';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {clientFromRequest} from '../oidc/client';
// $FlowFixMe migrated to typescript
import {getUserFromRequest} from '../util';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {injectOrganizationParams} from '../organization';

import type {OpenidUserInfoClaims} from 'openid-client';

import crypto from 'crypto';
// $FlowFixMe migrated to typescript
const logger = require('../../../shared/logging.ts').getLogger(module);

type Config = {
  urlPrefix: string,
};

export default function OrganizationOIDCStrategy(config: Config) {
  const verify = async (
    req,
    tokenSet: TokenSet,
    userInfo: OpenidUserInfoClaims,
    done: (error: Error | void, user?: User) => void,
  ) => {
    const email = userInfo.email;
    const organization = await req.organization();
    const ssoDefaultNetworkIDs = organization.ssoDefaultNetworkIDs;
    try {
      if (!email) {
        return done(new Error('Failed to read user email'));
      }
      let user = await getUserFromRequest(req, email);
      if (!user) {
        const createArgs = await injectOrganizationParams(req, {
          email: email.toLowerCase(),
          password: crypto.randomBytes(16).toString('hex'),
          // Hardcoded role for now, should be configurable
          role: AccessRoles.SUPERUSER,
          ssoDefaultNetworkIDs,
        });
        user = await User.create(createArgs);
      }
      req.session.oidc = {
        tokenSet,
      };
      done(undefined, user);
    } catch (e) {
      logger.error('Error creating user', e);
      done(new Error('Failed to login!'));
    }
  };

  return new DynamicStrategy(
    async req => (await req.organization()).id.toString(),
    async req => {
      const client = await clientFromRequest(req);
      const redirectTo = Array.isArray(req.query.to) ? '/' : req.query.to;
      // $FlowFixMe: req.get exists and is typed, this is a bug in flow
      const host = req.get('host');

      return new OidcStrategy<?User>(
        {
          client,
          path: `${config.urlPrefix}/login/oidc/callback`,
          passReqToCallback: true,
          params: {
            redirect_uri:
              `https://${host}${config.urlPrefix}/login/oidc/callback?to=` +
              encodeURIComponent(redirectTo),
          },
        },
        verify,
      );
    },
  );
}
