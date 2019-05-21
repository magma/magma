/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import session from 'express-session';

type SessionMiddlewareOptions = {|
  devMode: boolean,
  sessionStore: session.Session,
  sessionToken: string,
|};

export default function sessionMiddleware(
  options: SessionMiddlewareOptions,
): Middleware {
  options.sessionStore.sync();
  return session({
    cookie: {
      secure: !options.devMode,
    },
    // Used to sign the session cookie
    secret: options.sessionToken,
    resave: false,
    saveUninitialized: true,
    store: options.sessionStore,
    unset: 'destroy',
  });
}
