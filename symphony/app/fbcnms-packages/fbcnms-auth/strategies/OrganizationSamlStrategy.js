/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import MultiSamlStrategy from 'passport-saml/multiSamlStrategy';
import {AccessRoles} from '../roles';
import {User} from '@fbcnms/sequelize-models';

import {getUserFromRequest} from '../util';
import {injectOrganizationParams} from '../organization';

const logger = require('@fbcnms/logging').getLogger(module);

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
            password: Math.random().toString(36),
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
