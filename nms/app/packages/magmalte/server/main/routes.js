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

import type {AppContextAppData} from '@fbcnms/ui/context/AppContext';
import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

import adminRoutes from '../admin/routes';
import apiControllerRoutes from '../apicontroller/routes';
import asyncHandler from '@fbcnms/util/asyncHandler';
import express from 'express';
import networkRoutes from '../network/routes';
import path from 'path';
import staticDist from '@fbcnms/webpack-config/staticDist';
import userMiddleware from '@fbcnms/auth/express';
import {AccessRoles} from '@fbcnms/auth/roles';
import {MAPBOX_ACCESS_TOKEN} from '@fbcnms/platform-server/config';

import {TABS} from '@fbcnms/types/tabs';
import {access} from '@fbcnms/auth/access';
import {getEnabledFeatures} from '@fbcnms/platform-server/features';
import {masterOrgMiddleware} from '@fbcnms/platform-server/master/middleware';

const router: express.Router<FBCNMSRequest, ExpressResponse> = express.Router();

const handleReact = tab =>
  async function (req: FBCNMSRequest, res) {
    const organization = req.organization ? await req.organization() : null;
    const orgTabs = organization?.tabs || [];
    if (TABS[tab] && orgTabs.indexOf(tab) === -1) {
      res.redirect(
        organization?.tabs.length ? `/${organization.tabs[0]}` : '/',
      );
      return;
    }
    const appData: AppContextAppData = {
      csrfToken: req.csrfToken(),
      tabs: organization?.tabs || [],
      user: req.user
        ? {
            tenant: organization?.name || '',
            email: req.user.email,
            isSuperUser: req.user.isSuperUser,
            isReadOnlyUser: req.user.isReadOnlyUser,
          }
        : {tenant: '', email: '', isSuperUser: false, isReadOnlyUser: false},
      enabledFeatures: await getEnabledFeatures(req, organization?.name),
      ssoEnabled: !!organization?.ssoEntrypoint,
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
router.use('/admin', access(AccessRoles.SUPERUSER), adminRoutes);
router.get('/admin*', access(AccessRoles.SUPERUSER), handleReact('admin'));
router.use('/nms/apicontroller', apiControllerRoutes);
router.use('/nms/network', networkRoutes);
router.use('/nms/static', express.static(path.join(__dirname, '../static')));

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

const masterRouter = require('@fbcnms/platform-server/master/routes');
router.use('/master', masterOrgMiddleware, masterRouter.default);

async function handleMaster(req: FBCNMSRequest, res) {
  const appData: AppContextAppData = {
    csrfToken: req.csrfToken(),
    user: {
      tenant: 'master',
      email: req.user.email,
      isSuperUser: req.user.isSuperUser,
      isReadOnlyUser: req.user.isReadOnlyUser,
    },
    enabledFeatures: await getEnabledFeatures(req, 'master'),
    tabs: [],
    ssoEnabled: false,
    ssoSelectedType: 'none',
    csvCharset: null,
  };
  res.render('master', {
    staticDist,
    configJson: JSON.stringify({appData}),
  });
}

router.get('/master*', masterOrgMiddleware, handleMaster);

router.get(
  '/*',
  access(AccessRoles.USER),
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organization = await req.organization();
    if (organization.isMasterOrg) {
      res.redirect('/master');
    } else if (req.user.tabs && req.user.tabs.length > 0) {
      res.redirect(req.user.tabs[0]);
    } else if (organization.tabs && organization.tabs.length > 0) {
      res.redirect(organization.tabs[0]);
    } else {
      console.warn(
        `no tabs for user ${req.user.email}, organization ${organization.name}`,
      );
      res.redirect('/nms');
    }
  }),
);

export default router;
