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

import crypto from 'crypto';
import logging from '../../../shared/logging';
import {AccessRoles} from '../../../shared/roles';
import {MultiSamlStrategy} from 'passport-saml'; // compatibility with breaking change in 3.1.0
import {Request} from 'express';
import {Strategy} from 'passport';
import {User} from '../../../shared/sequelize_models';
import {VerifyWithRequest} from 'passport-saml/lib/passport-saml/types';
import {getUserFromRequest} from '../util';
import {ignoreAsync} from '../../util/ignoreAsync';
import {injectOrganizationParams} from '../organization';

const logger = logging.getLogger(module);

type Config = {
  urlPrefix: string;
};

export default function OrganizationSamlStrategy(config: Config) {
  const verify: VerifyWithRequest = (req, profile, done) => {
    async function syncWrapper() {
      const email = profile!.nameID;
      const organization = await req.organization!();
      // ssoDefaultNetworkIDs is broken and only used in a migration revert.
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
      const ssoDefaultNetworkIDs: Array<string> =
        // @ts-ignore
        organization.ssoDefaultNetworkIDs;

      try {
        if (!email) {
          return done(null, undefined, {message: 'Failed to read user email'});
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
        done(null, (user as unknown) as Record<string, unknown>, {
          message: 'User logged in',
        });
      } catch (e) {
        logger.error('Error creating user', e);
        done(null, undefined, {message: 'Failed to login!'});
      }
    }

    void syncWrapper();
  };

  return (new MultiSamlStrategy(
    {
      path: `${config.urlPrefix}/login/saml/callback`,
      getSamlOptions: ignoreAsync(
        async (req: Request<never, any, never, {to?: string}>, done) => {
          try {
            const host = req.get('host')!;
            const organization = await req.organization!();
            const configuration = {
              callbackUrl: `https://${host}${
                config.urlPrefix
              }/login/saml/callback?to=${req.query.to || '/'}`,
              cert: organization.ssoCert,
              entryPoint: organization.ssoEntrypoint,
              issuer: organization.ssoIssuer,
            };
            return done(null, configuration);
          } catch (error) {
            return done(error as Error);
          }
        },
      ),
      passReqToCallback: true,
    },
    verify,
  ) as unknown) as Strategy;
}
