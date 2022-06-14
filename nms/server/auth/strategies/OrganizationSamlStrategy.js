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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AccessRoles} from '../../../shared/roles';
import {MultiSamlStrategy} from 'passport-saml'; // compatibility with breaking change in 3.1.0
import {User} from '../../../shared/sequelize_models';

import {getUserFromRequest} from '../util';
import {injectOrganizationParams} from '../organization';

import crypto from 'crypto';

const logger = require('../../../shared/logging').getLogger(module);

type Config = {
  urlPrefix: string,
};

export default function OrganizationSamlStrategy(config: Config) {
  return new MultiSamlStrategy(
    {
      path: `${config.urlPrefix}/login/saml/callback`,
      getSamlOptions: async (req, done) => {
        try {
          const host = req.get('host');
          const organization = await req.organization();
          const configuration = {
            callbackUrl: `https://${host}${
              config.urlPrefix
            }/login/saml/callback?to=${req.query.to || '/'}`,
            cert: organization.ssoCert,
            entryPoint: organization.ssoEntrypoint,
            issuer: organization.ssoIssuer,
          };
          return done(null, configuration);
        } catch (err) {
          return done(err);
        }
      },
      passReqToCallback: true,
    },
    async (req, profile, done) => {
      const email = profile.nameID;
      const organization = await req.organization();
      const ssoDefaultNetworkIDs = organization.ssoDefaultNetworkIDs;
      try {
        if (!email) {
          return done(null, false, {message: 'Failed to read user email'});
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
        done(null, user, {message: 'User logged in'});
      } catch (e) {
        logger.error('Error creating user', e);
        done(null, false, {message: 'Failed to login!'});
      }
    },
  );
}
