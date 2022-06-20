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
 *
 * @flow strict-local
 * @format
 */
'use strict';

// $FlowFixMe migrated to typescript
import {AccessRoleLevel, AccessRoles} from '../../shared/roles';
// $FlowFixMe migrated to typescript
import {ErrorCodes} from '../../shared/errorCodes';

const path = require('path');
const addQueryParamsToUrl = require('./util').addQueryParamsToUrl;
// $FlowFixMe migrated to typescript
const logger = require('../../shared/logging.ts').getLogger(module);
const openRoutes = require('./openRoutes').default;

import type {ExpressResponse, NextFunction} from 'express';
import type {FBCNMSPassportRequest} from './passport';

type Options = {loginUrl: string};
// Final type, thus naming it as thus.
export type FBCNMSRequest = FBCNMSPassportRequest & {access: Options};

const validators = {
  [AccessRoles.USER]: (req: FBCNMSPassportRequest) => {
    return req.isAuthenticated();
  },
  [AccessRoles.SUPERUSER]: (req: FBCNMSPassportRequest) => {
    return req.user && req.user.role === AccessRoles.SUPERUSER;
  },
};

export const configureAccess = (options: Options) => {
  return function setup(
    req: FBCNMSPassportRequest & {access?: Options},
    _res: ExpressResponse,
    next: NextFunction,
  ) {
    req.access = options;
    next();
  };
};

export const access = (level: AccessRoleLevel) => {
  return async function access(
    req: FBCNMSRequest,
    res: ExpressResponse,
    next: NextFunction,
  ) {
    const normalizedURL = path.normalize(req.originalUrl);
    const isOpenRoute = openRoutes.some(route => normalizedURL.match(route));
    const hasPermission = validators[level](req);
    if (!isOpenRoute && req.user && req.organization) {
      const domainOrganization = await req.organization();
      const organization = req.user.organization;
      if (domainOrganization.name !== organization) {
        logger.error(
          'Strange bug, please fix! Organizations are Not Equal!! req.user.organization=' +
            (organization ?? ''),
          ', domainOrganization=' + domainOrganization.name,
        );
        req.logout();
        res.redirect(req.access.loginUrl);
        return;
      }
    }
    if (isOpenRoute || hasPermission) {
      // Continue to the next middleware if the user has permission
      next();
      return;
    }

    logger.info(
      `Client has no permission to view route: [%s]. They are ${
        req.user ? '' : 'not'
      } logged in`,
      req.hostname + req.originalUrl,
    );

    if (!req.user) {
      // No logged in user, attempt to redirect to login url
      const loginURL = addQueryParamsToUrl(req.access.loginUrl, {
        to: req.originalUrl,
      });
      res.format({
        // for axios requests, redirect does not work, so we return 403 error.
        json: () =>
          res.status(403).json({
            errorCode: ErrorCodes.USER_NOT_LOGGED_IN,
            description: 'You must login to see this',
          }),
        // for browser requests, simply redirect
        html: () => res.redirect(loginURL),
        default: () => res.redirect(loginURL),
      });
    } else {
      // if there is a logged in user, we shouldn't redirect to login page
      // because it would create an infinite loop since the login page redirects
      // back if the user is logged in.
      res.redirect('/');
    }
  };
};
