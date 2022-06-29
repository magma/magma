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

import bcrypt from 'bcryptjs';
import express from 'express';
import expressOnboarding from './expressOnboarding';
// $FlowFixMe migrated to typescript
import logging from '../../shared/logging';
import passport from 'passport';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import staticDist from '../../config/staticDist';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AccessRoles} from '../../shared/roles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {AuditLogEntry, User} from '../../shared/sequelize_models';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {access} from './access';
import {
  addQueryParamsToUrl,
  getPropsToUpdate,
  validateAndHashPassword,
  // $FlowFixMe[cannot-resolve-module] for TypeScript migration
} from './util';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {injectOrganizationParams} from './organization';
import {isEmpty} from 'lodash';
// $FlowFixMe migrated to typescript
import type {EmbeddedData} from '../../shared/types/embeddedData';
import type {ExpressRequest, ExpressResponse} from 'express';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {FBCNMSRequest} from './access';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {UserType} from '../../shared/sequelize_models/models/user';

import crypto from 'crypto';

const logger = logging.getLogger(module);
const PASSWORD_FOR_LOGGING = '<SECRET>';

type Options = {|
  loginSuccessUrl: string,
  loginFailureUrl: string,
  loginView?: string,
  onboardingUrl?: string,
|};

function accessRoleToString(role: number): string {
  if (role === AccessRoles.SUPERUSER) {
    return 'OWNER';
  }
  return 'USER';
}

export function unprotectedUserRoutes() {
  const router: express.Router<
    ExpressRequest,
    ExpressResponse,
  > = express.Router();
  router.post(
    '/login/saml/callback',
    passport.authenticate('saml', {
      failureRedirect: '/user/login?failure=true',
    }),
    (req, res: ExpressResponse) => {
      /* $FlowFixMe req.query is user-controlled input, properties and values
       in this object are untrusted and should be validated before trusting */
      const redirectTo = ensureRelativeUrl(req.query.to) || '/';
      res.redirect(redirectTo);
    },
  );
  return router;
}

function userMiddleware(options: Options): express.Router<FBCNMSRequest, *> {
  const router = express.Router();

  // Login / Logout Routes
  router.post('/login', (req: FBCNMSRequest, res, next) => {
    const redirectTo = ensureRelativeUrl(req.body.to);

    const loginSuccessUrl = redirectTo || options.loginSuccessUrl;
    const loginFailureUrl = redirectTo
      ? addQueryParamsToUrl(options.loginFailureUrl, {to: redirectTo})
      : options.loginFailureUrl;

    passport.authenticate('local', (err, user, _info) => {
      if (!user || err) {
        logger.error('Failed login: ' + err);
        return res.redirect(loginFailureUrl);
      }
      req.logIn(user, err => {
        if (err) {
          next(err);
          return;
        }
        res.redirect(loginSuccessUrl);
      });
    })(req, res, next);
  });

  if (options.onboardingUrl) {
    router.use(expressOnboarding());
  }

  router.get('/login', async (req: FBCNMSRequest, res, next) => {
    try {
      if (req.isAuthenticated()) {
        const to = req.query.to;
        const next = ensureRelativeUrl(Array.isArray(to) ? null : to);
        res.redirect(next || '/');
        return;
      }

      if (options.onboardingUrl && !(await User.findOne())) {
        res.redirect(options.onboardingUrl);
        return;
      }

      let ssoSelectedType = 'none';
      try {
        if (req.organization) {
          const org = await req.organization();
          ssoSelectedType = org.ssoSelectedType || 'none';
        }
      } catch (e) {
        logger.error('Error getting organization', e);
      }

      const appData: EmbeddedData = {
        csrfToken: req.csrfToken(),
        ssoEnabled: ssoSelectedType !== 'none',
        ssoSelectedType,
        csvCharset: null,
        enabledFeatures: [],
        user: {
          tenant: '',
          email: '',
          isSuperUser: false,
          isReadOnlyUser: false,
        },
      };

      res.render(options.loginView || 'login', {
        staticDist,
        configJson: JSON.stringify({
          appData,
        }),
      });
    } catch (e) {
      next(e);
    }
  });

  router.get(
    '/login/saml',
    passport.authenticate('saml', {
      failureRedirect: options.loginFailureUrl,
      authnRequestBinding: 'HTTP-Redirect',
    }),
  );

  router.get(
    '/login/oidc',
    passport.authenticate('oidc', {
      failureRedirect: options.loginFailureUrl,
    }),
  );
  router.get('/login/oidc/callback', (req: FBCNMSRequest, res, next) => {
    const to = req.query.to;
    const loginSuccessUrl =
      ensureRelativeUrl(Array.isArray(to) ? null : to) || '/';
    passport.authenticate('oidc', (err, user, _info) => {
      if (!user || err) {
        logger.error('Error logging in with oidc: ' + err);
        return res.redirect(options.loginFailureUrl);
      }
      req.logIn(user, err => {
        if (err) {
          next(err);
          return;
        }
        res.redirect(loginSuccessUrl);
      });
    })(req, res, next);
  });

  router.get('/logout', (req: FBCNMSRequest, res) => {
    if (req.isAuthenticated()) {
      req.logout();
    }
    delete req.session.oidc;
    res.redirect('/');
  });

  // User Details
  router.get(
    '/list',
    access(AccessRoles.USER),
    async (req: FBCNMSRequest, res) => {
      try {
        let users;
        if (req.organization) {
          const organization = await req.organization();
          users = await User.findAll({
            where: {
              organization: organization.name,
            },
          });
        } else {
          users = await User.findAll();
        }
        users = users.map(user => {
          return {
            id: user.id,
            email: user.email,
          };
        });
        res.status(200).send({users});
      } catch (error) {
        res.status(400).send({error: error.toString()});
      }
    },
  );

  // Current User Details
  router.get(
    '/me',
    passport.authenticate(['basic_local', 'session'], {session: false}),
    access(AccessRoles.USER),
    (req: FBCNMSRequest, res: ExpressResponse): void => {
      res.status(200).send({
        organization: req.user.organization,
        email: req.user.email,
        role: accessRoleToString(req.user.role),
      });
    },
  );

  // User Routes
  router.get(
    '/async/',
    access(AccessRoles.SUPERUSER),
    async (req: FBCNMSRequest, res) => {
      try {
        let users;
        if (req.organization) {
          const organization = await req.organization();
          users = await User.findAll({
            where: {
              organization: organization.name,
            },
          });
        } else {
          users = await User.findAll();
        }
        res.status(200).send({users});
      } catch (error) {
        res.status(400).send({error: error.toString()});
      }
    },
  );

  router.post(
    '/async/',
    access(AccessRoles.SUPERUSER),
    async (req: FBCNMSRequest, res) => {
      try {
        const body = req.body;
        if (!body.email) {
          throw new Error('Email not included!');
        }

        const allowedProps = ['email', 'networkIDs', 'password', 'role'];
        let userProperties = await getPropsToUpdate(
          allowedProps,
          body,
          params => injectOrganizationParams(req, params),
        );
        userProperties = await injectOrganizationParams(req, userProperties);

        // this happens when the user is being added to an organization that
        // uses SSO for login, give it a random password
        if (req.organization && userProperties.password === undefined) {
          const organization = await req.organization();
          if (organization.ssoEntrypoint) {
            userProperties.password = crypto.randomBytes(16).toString('hex');
          }
        }
        const user = await User.create(userProperties);
        await logUserChange(
          req,
          req.user,
          'CREATE',
          {...userProperties},
          'SUCCESS',
        );
        res.status(201).send({user});
      } catch (error) {
        res.status(400).send({error: error.toString()});
        await logUserChange(req, req.user, 'CREATE', req.body, 'FAILURE');
      }
    },
  );

  router.put(
    '/async/:id',
    access(AccessRoles.SUPERUSER),
    async (req: FBCNMSRequest, res) => {
      try {
        const {id} = req.params;
        const user = await User.findOne({where: {id}});

        // Check if user exists
        if (!user) {
          throw new Error('User does not exist!');
        }

        // Create object to pass into update()
        const allowedProps = ['networkIDs', 'password', 'role'];

        const userProperties = await getPropsToUpdate(
          allowedProps,
          req.body,
          params => injectOrganizationParams(req, params),
        );
        if (isEmpty(userProperties)) {
          throw new Error('No valid properties to edit!');
        }

        // Update user's password
        await user.update(userProperties);
        await logUserChange(req, req.user, 'UPDATE', req.body, 'SUCCESS');
        res.status(200).send({user});
      } catch (error) {
        await logUserChange(req, req.user, 'UPDATE', req.body, 'FAILURE');
        res.status(400).send({error: error.toString()});
      }
    },
  );

  router.delete(
    '/async/:id/',
    access(AccessRoles.SUPERUSER),
    async (req: FBCNMSRequest, res) => {
      const {id} = req.params;

      try {
        await User.destroy({where: {id}});
        await logUserChange(req, req.user, 'DELETE', {}, 'SUCCESS');
        res.status(200).send();
      } catch (error) {
        await logUserChange(req, req.user, 'DELETE', {}, 'FAILURE');
        res.status(400).send({error: error.toString()});
      }
    },
  );

  router.post('/change_password', async (req: FBCNMSRequest, res) => {
    try {
      const {currentPassword, newPassword} = req.body;
      const verified = await bcrypt.compare(currentPassword, req.user.password);
      if (!verified) {
        throw new Error('Incorrect password');
      }

      const hashedPassword = await validateAndHashPassword(newPassword);
      await req.user.update({password: hashedPassword});
      await logUserChange(req, req.user, 'UPDATE', {password: ''}, 'SUCCESS');
      res.status(200).send();
    } catch (error) {
      await logUserChange(req, req.user, 'UPDATE', {password: ''}, 'FAILURE');
      res.status(400).send({error: error.toString()});
    }
  });

  router.put(
    '/set/:email',
    access(AccessRoles.SUPERUSER),
    async (req: FBCNMSRequest, res) => {
      try {
        const {email} = req.params;

        const where = await injectOrganizationParams(req, {email});
        const user = await User.findOne({where});

        // Check if user exists
        if (!user) {
          throw new Error('User does not exist!');
        }

        // Create object to pass into update()
        const allowedProps = ['password', 'role'];

        const userProperties = await getPropsToUpdate(
          allowedProps,
          req.body,
          params => injectOrganizationParams(req, params),
        );
        if (isEmpty(userProperties)) {
          throw new Error('No valid properties to edit!');
        }

        await user.update(userProperties);
        await logUserChange(req, req.user, 'UPDATE', req.body, 'SUCCESS');
        res.status(200).send({user});
      } catch (error) {
        await logUserChange(req, req.user, 'UPDATE', req.body, 'FAILURE');
        res.status(400).send({error: error.toString()});
      }
    },
  );

  return router;
}

function ensureRelativeUrl(url: ?string): ?string {
  if (url && (url.indexOf('/') !== 0 || url.indexOf('//') === 0)) {
    return null;
  }
  return url;
}

async function logUserChange(
  req: FBCNMSRequest,
  target: UserType,
  mutationType: 'CREATE' | 'UPDATE' | 'DELETE',
  data: {[string]: mixed},
  status: 'SUCCESS' | 'FAILURE',
) {
  let org;
  if (req.organization) {
    org = await req.organization();
  }

  const mutationData = {...data};
  if (data.password != null) {
    mutationData.password = PASSWORD_FOR_LOGGING;
  }

  const auditLog = {
    actingUserId: req.user.id,
    organization: org?.name || '<NO_ORGANIZATION>',
    mutationType,
    objectId: `${target.id}`,
    objectType: 'USER',
    objectDisplayName: target.email,
    mutationData,
    url: req.originalUrl,
    ipAddress: req.ip,
    status,
    statusCode: 'N/A',
  };
  await AuditLogEntry.create(auditLog);
}

export default userMiddleware;
