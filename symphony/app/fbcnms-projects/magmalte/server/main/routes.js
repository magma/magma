/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {AppContextAppData} from '@fbcnms/ui/context/AppContext';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

import express from 'express';
import staticDist from 'fbcnms-webpack-config/staticDist';
import userMiddleware from '@fbcnms/auth/express';
import {AccessRoles} from '@fbcnms/auth/roles';
import {MAPBOX_ACCESS_TOKEN} from '@fbcnms/platform-server/config';

import {access} from '@fbcnms/auth/access';
import {getEnabledFeatures} from '@fbcnms/platform-server/features';

const router = express.Router();

const handleReact = _tab =>
  async function (req: FBCNMSRequest, res) {
    const appData: AppContextAppData = {
      csrfToken: req.csrfToken(),
      tabs: ['nms'],
      user: req.user
        ? {
            tenant: '',
            email: req.user.email,
            isSuperUser: req.user.isSuperUser,
            isReadOnlyUser: req.user.isReadOnlyUser,
          }
        : {tenant: '', email: '', isSuperUser: false, isReadOnlyUser: false},
      enabledFeatures: await getEnabledFeatures(req, null),
      ssoEnabled: false,
      ssoSelectedType: 'none',
      csvCharset: null,
    };
    res.render('index', {
      staticDist,
      configJson: JSON.stringify({
        appData,
        MAPBOX_ACCESS_TOKEN: req.user && MAPBOX_ACCESS_TOKEN,
      }),
    });
  };

router.use('/healthz', (req: FBCNMSRequest, res) => res.send('OK'));
router.use(
  '/admin',
  access(AccessRoles.SUPERUSER),
  require('@fbcnms/platform-server/admin/routes').default,
);
router.get('/admin*', access(AccessRoles.SUPERUSER), handleReact('admin'));
router.use(
  '/nms/apicontroller',
  require('@fbcnms/platform-server/apicontroller/routes').default,
);
router.use(
  '/nms/network',
  require('@fbcnms/platform-server/network/routes').default,
);

router.use('/logger', require('@fbcnms/platform-server/logger/routes'));
router.use('/test', require('@fbcnms/platform-server/test/routes'));
router.use(
  '/user',
  userMiddleware({
    loginSuccessUrl: '/nms',
    loginFailureUrl: '/user/login?invalid=true',
  }),
);
router.get('/nms*', access(AccessRoles.USER), handleReact('nms'));

router.get('/', (req: FBCNMSRequest, res) => {
  res.redirect('/nms');
});

export default router;
