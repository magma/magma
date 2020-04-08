/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

const {
  DEV_MODE,
  LOG_FORMAT,
  LOG_LEVEL,
} = require('@fbcnms/platform-server/config');

// This must be done before any module imports to configure
// logging correctly
const logging = require('@fbcnms/logging');
logging.configure({
  LOG_FORMAT,
  LOG_LEVEL,
});

const {
  appMiddleware,
  csrfMiddleware,
  organizationMiddleware,
  sessionMiddleware,
  webpackSmartMiddleware,
} = require('@fbcnms/express-middleware');
const connectSession = require('connect-session-sequelize');
const express = require('express');
const passport = require('passport');
const path = require('path');
const paths = require('fbcnms-webpack-config/paths');
const fbcPassport = require('@fbcnms/auth/passport').default;
const session = require('express-session');
const {sequelize} = require('@fbcnms/sequelize-models');
const OrganizationLocalStrategy = require('@fbcnms/auth/strategies/OrganizationLocalStrategy')
  .default;

const {access, configureAccess} = require('@fbcnms/auth/access');
const {
  AccessRoles: {USER},
} = require('@fbcnms/auth/roles');

const devMode = process.env.NODE_ENV !== 'production';

// Create Sequelize Store
const SessionStore = connectSession(session.Store);
const sequelizeSessionStore = new SessionStore({db: sequelize});

// FBC express initialization
const app = express();
app.set('trust proxy', 1);
app.use(organizationMiddleware());
app.use(appMiddleware());
app.use(
  sessionMiddleware({
    devMode,
    sessionStore: sequelizeSessionStore,
    sessionToken:
      process.env.SESSION_TOKEN || 'fhcfvugnlkkgntihvlekctunhbbdbjiu',
  }),
);
app.use(passport.initialize());
app.use(passport.session()); // must be after sessionMiddleware

fbcPassport.use();
passport.use('local', OrganizationLocalStrategy());

// Views
app.set('views', path.join(__dirname, '..', 'views'));
app.set('view engine', 'pug');

// Routes
app.use(
  webpackSmartMiddleware({
    devMode: DEV_MODE,
    devWebpackConfig: require('../config/webpack.config.js'),
    distPath: paths.distPath,
  }),
);
app.use('/user', require('@fbcnms/auth/express').unprotectedUserRoutes());

app.use(configureAccess({loginUrl: '/user/login'}));

app.use('/', csrfMiddleware(), access(USER), require('./main/routes').default);

export default app;
