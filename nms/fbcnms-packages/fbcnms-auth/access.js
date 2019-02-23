/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

import type {AccessRoleLevel} from './roles';

const path = require('path');

const {AccessRoles} = require('./roles');
const express = require('express');
const logger = require('@fbcnms/logging').getLogger(module);
const openRoutes = require('./openRoutes').default;

const validators = {
  [AccessRoles.USER]: req => {
    return req.isAuthenticated();
  },
  [AccessRoles.SUPERUSER]: req => {
    return req.user && req.user.role === AccessRoles.SUPERUSER;
  },
};

export const configureAccess = (options: {loginUrl: string}) => {
  return function setup(
    req: express.Request,
    _res: express.Response,
    next: express.Next,
  ) {
    req.access = options;
    next();
  };
};

export const access = (level: AccessRoleLevel) => {
  return async function access(
    req: express.Request,
    res: express.Response,
    next: express.Next,
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
            organization,
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

    logger.debug(
      'Client has no permission to view route: [%s], redirecting to /',
      req.originalUrl,
    );
    res.redirect(req.access.loginUrl);
  };
};
