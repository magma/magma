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
import passport from 'passport';
import {Strategy as LocalStrategy} from 'passport-local';

import {injectOrganizationParams} from './organization';

import type {StaticUserModel} from '@fbcnms/sequelize-models/models/user';
import type {FBCNMSMiddleWareRequest} from '@fbcnms/express-middleware';
type OutputRequest<T> = {
  logIn: (T, (err?: ?Error) => void) => void,
  logOut: () => void,
  logout: () => void,
  user: T,
  isAuthenticated: () => boolean,
  isUnauthenticated: () => boolean,
} & FBCNMSMiddleWareRequest;
// User is currently untyped, export as an object.
export type FBCNMSPassportRequest = OutputRequest<Object>;

type PassportConfig = {
  UserModel: StaticUserModel,
};

function use(config: PassportConfig) {
  const getUserFromRequest = async (
    req: FBCNMSMiddleWareRequest,
    email: string,
  ) => {
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
}

export default {
  use,
};
