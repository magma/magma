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
import passport from 'passport';
import staticDist from 'fbcnms-webpack-config/staticDist';
import {access} from './access';
import {AccessRoles} from './roles';
import {addQueryParamsToUrl} from './util';
import EmailValidator from 'email-validator';
import {User} from '@fbcnms/sequelize-models';

import type {FBCNMSRequest} from './access';

const SALT_GEN_ROUNDS = 10;
const MIN_PASSWORD_LENGTH = 10;

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
  verificationType: 'verificationType',
};

function userMiddleware(options: Options): express.Router {
  const router = express.Router();

  async function getPropsToUpdate(
    req: FBCNMSRequest,
    allowedProps: string[],
    body: {[string]: string},
  ) {
    allowedProps = allowedProps.filter(prop =>
      User.rawAttributes.hasOwnProperty(FIELD_MAP[prop]),
    );
    const userProperties = {};
    for (const prop of allowedProps) {
      if (body.hasOwnProperty(prop)) {
        switch (prop) {
          case 'email':
            const email = body[prop].toLowerCase();

            if (!EmailValidator.validate(body.email)) {
              throw new Error('Please enter a valid email');
            }

            // Check if user exists
            const where = await injectOrganizationParams(req, {email});
            if (await User.findOne({where})) {
              throw new Error(`${email} already exists`);
            }
            userProperties[prop] = email;
            break;
          case 'password':
            userProperties[prop] = await validateAndHashPassword(body[prop]);
            break;
          case 'superUser':
            userProperties.role = body[prop]
              ? AccessRoles.SUPERUSER
              : AccessRoles.USER;
            break;
          default:
            userProperties[prop] = body[prop];
            break;
        }
      }
    }
    return userProperties;
  }

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

  router.get('/login', (req: FBCNMSRequest, res) => {
    if (req.isAuthenticated()) {
      res.redirect(ensureRelativeUrl(req.body.to) || '/');
      return;
    }
    res.render('login', {
      staticDist,
      configJson: JSON.stringify({
        appData: {
          csrfToken: req.csrfToken(),
        },
      }),
    });
  });

  router.get('/logout', (req: FBCNMSRequest, res) => {
    if (req.isAuthenticated()) {
      req.logout();
    }
    res.redirect(options.loginFailureUrl);
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

        const allowedProps = [
          'email',
          'networkIDs',
          'password',
          'superUser',
          'verificationType',
        ];
        let userProperties = await getPropsToUpdate(req, allowedProps, body);
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
        const allowedProps = [
          'networkIDs',
          'password',
          'superUser',
          'verificationType',
        ];

        const userProperties = await getPropsToUpdate(
          req,
          allowedProps,
          req.body,
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
  if (!password || password.length < MIN_PASSWORD_LENGTH) {
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
