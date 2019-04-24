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
import {Strategy as LocalStrategy} from 'passport-local';
import {getUserFromRequest} from '../util';

export default function() {
  return new LocalStrategy(
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
  );
}
