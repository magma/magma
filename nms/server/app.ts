/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// This must be done before any module imports to configure
// logging correctly
import logging from '../shared/logging';
logging.configure({
  LOG_FORMAT,
  LOG_LEVEL,
});

import OrganizationLocalStrategy from './auth/strategies/OrganizationLocalStrategy';
import OrganizationSamlStrategy from './auth/strategies/OrganizationSamlStrategy';
import alertRoutes from './alerts/routes';
import connectSession from 'connect-session-sequelize';
import crypto from 'crypto';
import devWebpackConfig from '../config/webpack.config';
import express from 'express';
import fbcPassport from './auth/passport';
import fs from 'fs';
import grafanaRoutes from './grafana/routes';
import mainRoutes from './main/routes';
import passport from 'passport';
import path from 'path';
import paths from '../config/paths';
import session from 'express-session';
import {AccessRoles} from '../shared/roles';
import {DEV_MODE, LOG_FORMAT, LOG_LEVEL} from '../config/config';
import {access, configureAccess} from './auth/access';
import {
  appMiddleware,
  csrfMiddleware,
  organizationMiddleware,
  sessionMiddleware,
  webpackSmartMiddleware,
} from './middleware';
import {sequelize} from '../shared/sequelize_models';
import {unprotectedUserRoutes} from './auth/express';

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
    devMode: DEV_MODE,
    sessionStore: sequelizeSessionStore,
    sessionToken:
      process.env.SESSION_TOKEN || crypto.randomBytes(32).toString('hex'),
  }),
);
app.use(passport.initialize());
app.use(passport.session()); // must be after sessionMiddleware

fbcPassport.use();
passport.use('local', OrganizationLocalStrategy());
passport.use(
  'saml',
  OrganizationSamlStrategy({
    urlPrefix: '/user',
  }),
);

// Views
app.set('views', path.join(__dirname, '../server/', 'views'));
app.set('view engine', 'pug');

const distPath = paths.distPath;
// Routes
// TO DO - fix this in webpack-dev-middleware code in fbc-js-core
app.use(
  webpackSmartMiddleware({
    devMode: DEV_MODE,
    devWebpackConfig,
    distPath,
  }),
);
app.use('/user', unprotectedUserRoutes());

app.use(configureAccess({loginUrl: '/user/login'}));

// Grafana uses its own CSRF, so we don't need to handle it on our side.
// Grafana can access all metrics of an org, so it must be restricted
// to superusers
app.use('/grafana', access(AccessRoles.SUPERUSER), grafanaRoutes);

// add lte metrics json file handler
const lteMetricsJsonData = fs.readFileSync(
  path.join(__dirname, '..', 'api/data/LteMetrics.json'),
  'utf-8',
);
const alertLinksJsonData = fs.readFileSync(
  path.join(__dirname, '..', 'api/data/AlertLinks.json'),
  'utf-8',
);
app.get('/api/data/LteMetrics', (req, res) => res.send(lteMetricsJsonData));
app.get('/api/data/AlertLinks', (req, res) => res.send(alertLinksJsonData));

// Trigger syncing of automatically generated alerts
app.use('/sync_alerts', access(AccessRoles.USER), alertRoutes);

app.use('/', csrfMiddleware(), access(AccessRoles.USER), mainRoutes);

export default app;
