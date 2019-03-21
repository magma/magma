/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {Request} from 'express';
import bcrypt from 'bcryptjs';
import passport from 'passport';
import {Strategy as LocalStrategy} from 'passport-local';
import {Strategy as FacebookStrategy} from 'passport-facebook';
import {UserVerificationTypes} from './types';

import {injectOrganizationParams} from './organization';

type PassportConfig = {
  UserModel: any,
  facebookLogin?: {
    appId: string,
    appSecret: string,
  },
};

function use(config: PassportConfig) {
  const getUserFromRequest = async (req: Request, email: string) => {
    const where = await injectOrganizationParams(req, {email});
    return await config.UserModel.findOne({where});
  };

  passport.serializeUser((user, done) => {
    done(null, user.id);
  });

  passport.deserializeUser(async (id, done) => {
    try {
      const user = await config.UserModel.findById(id);
      done(null, user);
    } catch (error) {
      done(error);
    }
  });

  passport.use(
    'local',
    new LocalStrategy(
      {
        usernameField: 'email',
        passwordField: 'password',
        passReqToCallback: true,
      },
      async (req, email, password, done) => {
        try {
          const user = await getUserFromRequest(req, email);
          if (!user) {
            return done(null, false, {
              message: 'Username or password invalid!',
            });
          }

          if (
            user.verificationType &&
            user.verificationType != UserVerificationTypes.PASSWORD
          ) {
            return done(null, false, {message: 'Wrong verification type'});
          }

          if (await bcrypt.compare(password, user.password)) {
            done(null, user);
          } else {
            done(null, false, {message: 'Invalid username or password!'});
          }
        } catch (error) {
          done(error);
        }
      },
    ),
  );

  if (config.facebookLogin && config.facebookLogin.appId) {
    passport.use(
      new FacebookStrategy(
        {
          clientID: config.facebookLogin.appId,
          clientSecret: config.facebookLogin.appSecret,
          callbackURL: '/user/login/facebook/callback',
          profileFields: ['id', 'emails', 'name'],
          passReqToCallback: true,
        },
        async (req, accessToken, refreshToken, profile, done) => {
          try {
            if (!profile.emails?.[0]?.value) {
              return done(null, false, {message: 'Failed to read user email'});
            }
            const email = profile.emails[0].value;
            const user = await getUserFromRequest(req, email);
            if (!user) {
              return done(null, false, {message: 'User not authorized'});
            }

            if (
              !user.verificationType ||
              user.verificationType != UserVerificationTypes.FACEBOOK
            ) {
              return done(null, false, {message: 'Wrong verification type'});
            }

            done(null, user, {message: 'User logged in'});
          } catch (e) {
            done(null, false, {message: 'Failed to login!'});
          }
        },
      ),
    );
  }
}

export default {
  use,
};
