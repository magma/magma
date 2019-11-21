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
import {User} from '@fbcnms/sequelize-models';

import type {FBCNMSMiddleWareRequest} from '@fbcnms/express-middleware';
import type {UserType} from '@fbcnms/sequelize-models/models/user';

type OutputRequest<T> = {
  logIn: (T, (err?: ?Error) => void) => void,
  logOut: () => void,
  logout: () => void,
  user: T,
  isAuthenticated: () => boolean,
  isUnauthenticated: () => boolean,
} & FBCNMSMiddleWareRequest;
export type FBCNMSPassportRequest = OutputRequest<UserType>;

function use() {
  passport.serializeUser((user, done) => {
    done(null, user.id);
  });

  passport.deserializeUser(async (id, done) => {
    try {
      const user = await User.findByPk(id);
      done(null, user);
    } catch (error) {
      done(error);
    }
  });
}

export default {
  use,
};
