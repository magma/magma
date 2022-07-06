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
 */
import openRoutes from './openRoutes';
import path from 'path';
import {AccessRoleLevel, AccessRoles} from '../../shared/roles';
import {ErrorCodes} from '../../shared/errorCodes';
import {addQueryParamsToUrl} from './util';
import {getLogger} from '../../shared/logging';

import type {NextFunction, Request, Response} from 'express';
const logger = getLogger(module);

type Options = {loginUrl: string};

const validators = {
  [AccessRoles.USER]: (req: Request) => {
    return req.isAuthenticated();
  },
  [AccessRoles.SUPERUSER]: (req: Request) => {
    return req.user && req.user.role === AccessRoles.SUPERUSER;
  },
  [AccessRoles.READ_ONLY_USER]: () => {
    return false;
  },
};

export const configureAccess = (options: Options) => {
  return function setup(req: Request, _res: Response, next: NextFunction) {
    req.access = options;
    next();
  };
};

export const access = (level: AccessRoleLevel) => {
  const handler = async function access(
    req: Request,
    res: Response,
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

  return (req: Request, res: Response, next: NextFunction) => {
    void handler(req, res, next);
  };
};
