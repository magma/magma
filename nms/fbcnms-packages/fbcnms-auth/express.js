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
import {injectOrganizationParams} from './organization';
import {isEmpty} from 'lodash-es';
import express from 'express';
import logging from '@fbcnms/logging';
import passport from 'passport';
import staticDist from 'fbcnms-webpack-config/staticDist';
import {access} from './access';
import {AccessRoles} from './roles';
import {addQueryParamsToUrl} from './util';
import EmailValidator from 'email-validator';
import {User} from '@fbcnms/sequelize-models';

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from './access';
import type {UserRawType} from '@fbcnms/sequelize-models/models/user';

const SALT_GEN_ROUNDS = 10;
const MIN_PASSWORD_LENGTH = 10;

const logger = logging.getLogger(module);

type Options = {|
  loginSuccessUrl: string,
  loginFailureUrl: string,
|};

const FIELD_MAP = {
  email: 'email',
  networkIDs: 'networkIDs',
  organization: 'organization',
  password: 'password',
  superUser: 'role',
};

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

export async function getPropsToUpdate(
  allowedProps: $Keys<typeof FIELD_MAP>[],
  body: {[string]: mixed},
  organizationInjector: ({[string]: any}) => Promise<{
    [string]: any,
    organization?: string,
  }>,
): Promise<$Shape<UserRawType>> {
  allowedProps = allowedProps.filter(prop =>
    User.rawAttributes.hasOwnProperty(FIELD_MAP[prop]),
  );
  const userProperties = {};
  for (const prop of allowedProps) {
    if (body.hasOwnProperty(prop)) {
      switch (prop) {
        case 'email':
          const emailUnsafe = body[prop];
          if (
            typeof emailUnsafe !== 'string' ||
            !EmailValidator.validate(body.email)
          ) {
            throw new Error('Please enter a valid email');
          }
          const email = emailUnsafe.toLowerCase();

          // Check if user exists
          const where = await organizationInjector({email});
          if (await User.findOne({where})) {
            throw new Error(`${email} already exists`);
          }
          userProperties[prop] = email;
          break;
        case 'password':
          userProperties[prop] = await validateAndHashPassword(body[prop]);
          break;
        case 'superUser':
          userProperties.role =
            body[prop] == true ? AccessRoles.SUPERUSER : AccessRoles.USER;
          break;
        case 'networkIDs':
          const networkIDsunsafe = body[prop];
          if (Array.isArray(networkIDsunsafe)) {
            const networkIDs: Array<string> = networkIDsunsafe.map(it => {
              if (typeof it !== 'string') {
                throw new Error('Please enter valid network IDs');
              }
              return it;
            });
            userProperties[prop] = networkIDs;
            break;
          }
          throw new Error('Please enter valid network IDs');
        case 'organization':
          if (typeof body[prop] !== 'string') {
            throw new Error('Invalid Organization!');
          }
          userProperties[prop] = body[prop];
          break;
        default:
          userProperties[prop] = body[prop];
          break;
      }
    }
  }
  return userProperties;
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

  router.get('/login', async (req: FBCNMSRequest, res) => {
    if (req.isAuthenticated()) {
      res.redirect(ensureRelativeUrl(req.body.to) || '/');
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

        const allowedProps = ['email', 'networkIDs', 'password', 'superUser'];
        let userProperties = await getPropsToUpdate(
          allowedProps,
          body,
          params => injectOrganizationParams(req, params),
        );
        userProperties = await injectOrganizationParams(req, userProperties);
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
        const allowedProps = ['networkIDs', 'password', 'superUser'];

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

async function validateAndHashPassword(password) {
  if (
    typeof password !== 'string' ||
    password === '' ||
    password.length < MIN_PASSWORD_LENGTH
  ) {
    throw new Error(
      'Password must contain at least ' + MIN_PASSWORD_LENGTH + ' characters',
    );
  }

  const salt = await bcrypt.genSalt(SALT_GEN_ROUNDS);
  return await bcrypt.hash(password, salt);
}

function ensureRelativeUrl(url: ?string): ?string {
  if (url && (url.indexOf('/') !== 0 || url.indexOf('//') === 0)) {
    return null;
  }
  return url;
}

export default userMiddleware;
