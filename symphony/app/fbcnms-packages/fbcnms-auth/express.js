/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import bcrypt from 'bcryptjs';
import express from 'express';
import expressOnboarding from './expressOnboarding';
import logging from '@fbcnms/logging';
import passport from 'passport';
import staticDist from 'fbcnms-webpack-config/staticDist';
import {AccessRoles} from './roles';
import {User} from '@fbcnms/sequelize-models';
import {access} from './access';
import {
  addQueryParamsToUrl,
  getPropsToUpdate,
  validateAndHashPassword,
} from './util';
import {injectOrganizationParams} from './organization';
import {isEmpty} from 'lodash';

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from './access';

const logger = logging.getLogger(module);

type Options = {|
  loginSuccessUrl: string,
  loginFailureUrl: string,
  onboardingUrl?: string,
|};

export function unprotectedUserRoutes() {
  const router = express.Router();
  router.post(
    '/login/saml/callback',
    passport.authenticate('saml', {
      failureRedirect: '/user/login?failure=true',
    }),
    (req, res: ExpressResponse) => {
      const redirectTo = ensureRelativeUrl(req.query.to) || '/';
      res.redirect(redirectTo);
    },
  );
  return router;
}

function userMiddleware(options: Options): express.Router {
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

  router.get('/login', async (req: FBCNMSRequest, res) => {
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

    let isSSO = false;
    try {
      if (req.organization) {
        const org = await req.organization();
        isSSO = !!org.ssoEntrypoint;
      }
    } catch (e) {
      logger.error('Error getting organization', e);
    }

    res.render('login', {
      staticDist,
      configJson: JSON.stringify({
        appData: {
          csrfToken: req.csrfToken(),
          isSSO,
        },
      }),
    });
  });

  router.get(
    '/login/saml',
    passport.authenticate('saml', {
      failureRedirect: options.loginFailureUrl,
      authnRequestBinding: 'HTTP-Redirect',
    }),
  );

  router.get('/logout', (req: FBCNMSRequest, res) => {
    if (req.isAuthenticated()) {
      req.logout();
    }
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

        const allowedProps = [
          'email',
          'networkIDs',
          'password',
          'role',
          'tabs',
        ];
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
            userProperties.password = Math.random().toString(36);
          }
        }
        const user = await User.create(userProperties);

        res.status(201).send({user});
      } catch (error) {
        res.status(400).send({error: error.toString()});
      }
    },
  );

  router.put(
    '/async/:id',
    access(AccessRoles.SUPERUSER),
    async (req: FBCNMSRequest, res) => {
      try {
        const {id} = req.params;

        const where = await injectOrganizationParams(req, {id});
        const user = await User.findOne({where});

        // Check if user exists
        if (!user) {
          throw new Error('User does not exist!');
        }

        // Create object to pass into update()
        const allowedProps = ['networkIDs', 'password', 'role', 'tabs'];

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
        res.status(200).send({user});
      } catch (error) {
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
        const where = await injectOrganizationParams(req, {id});
        await User.destroy({where});
        res.status(200).send();
      } catch (error) {
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
      res.status(200).send();
    } catch (error) {
      res.status(400).send({error: error.toString()});
    }
  });

  return router;
}

function ensureRelativeUrl(url: ?string): ?string {
  if (url && (url.indexOf('/') !== 0 || url.indexOf('//') === 0)) {
    return null;
  }
  return url;
}

export default userMiddleware;
