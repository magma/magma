/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import passport from 'passport';
import OrganizationLocalStrategy from './strategies/OrganizationLocalStrategy';
import {User} from '@fbcnms/sequelize-models';

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

function use() {
  passport.serializeUser((user, done) => {
    done(null, user.id);
  });

  passport.deserializeUser(async (id, done) => {
    try {
      const user = await User.findById(id);
      done(null, user);
    } catch (error) {
      done(error);
    }
  });

  passport.use('local', OrganizationLocalStrategy());
}

export default {
  use,
};
