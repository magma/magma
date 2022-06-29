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

import DynamicStrategy from './DynamicStrategy';
import crypto from 'crypto';
import logging from '../../../shared/logging';
import {AccessRoles} from '../../../shared/roles';
import {FBCNMSMiddleWareRequest} from '../../middleware';
import {
  Strategy as OidcStrategy,
  OpenidUserInfoClaims,
  TokenSet,
} from 'openid-client';
import {User} from '../../../shared/sequelize_models';
import {UserModel} from '../../../shared/sequelize_models/models/user';
import {clientFromRequest} from '../oidc/client';
import {getUserFromRequest} from '../util';
import {injectOrganizationParams} from '../organization';

type Config = {urlPrefix: string};

const logger = logging.getLogger(module);

export default function OrganizationOIDCStrategy(config: Config) {
  const verify = async (
    req: FBCNMSMiddleWareRequest,
    tokenSet: TokenSet,
    userInfo: OpenidUserInfoClaims,
    done: (error: Error | void, user?: UserModel) => void,
  ) => {
    const email = userInfo.email;
    const organization = await req.organization!();
    // ssoDefaultNetworkIDs is broken and only used in a migration revert.
    // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
    const ssoDefaultNetworkIDs: Array<string> =
      // @ts-ignore
      organization.ssoDefaultNetworkIDs;

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

      req.session.oidc = {tokenSet};
      done(undefined, user);
    } catch (e) {
      logger.error('Error creating user', e);
      done(new Error('Failed to login!'));
    }
  };

  return new DynamicStrategy(
    async req => (await req.organization!()).id.toString(),
    async req => {
      const client = await clientFromRequest(req);
      const redirectTo = Array.isArray(req.query.to)
        ? '/'
        : (req.query.to as string);
      const host = req.get('host')!;

      return new OidcStrategy<UserModel | null | undefined>(
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
