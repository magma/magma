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

import type {EmbeddedData} from '../../shared/types/embeddedData';
import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '../auth/access';

import MagmaV1API from '../magma/index';
import adminRoutes from '../admin/routes';
import apiControllerRoutes from '../apicontroller/routes';
import asyncHandler from '../util/asyncHandler';
import express from 'express';
import hostRoutes from '../host/routes';
import loggerRoutes from '../logger/routes';
import networkRoutes from '../network/routes';
import path from 'path';
import staticDist from '../../config/staticDist';
import testRoutes from '../test/routes';
import userMiddleware from '../auth/express';
import {AccessRoles} from '../../shared/roles';
import {TABS} from '../../fbc_js_core/types/tabs';
import {access} from '../auth/access';
import {getEnabledFeatures} from '../features';
import {hostOrgMiddleware} from '../host/middleware';

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
    const appData: EmbeddedData = {
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
      }),
    });
  };

router.use('/healthz', (req: FBCNMSRequest, res) => res.send('OK'));
router.use('/version', (req: FBCNMSRequest, res) =>
  res.send(process.env.VERSION_TAG),
);
router.use('/admin', access(AccessRoles.SUPERUSER), adminRoutes);
router.get('/admin*', access(AccessRoles.SUPERUSER), handleReact('admin'));
router.use('/nms/apicontroller', apiControllerRoutes);
router.use('/nms/network', networkRoutes);
router.use('/nms/static', express.static(path.join(__dirname, '../static')));

router.use('/logger', loggerRoutes);
router.use('/test', testRoutes);
router.use(
  '/user',
  userMiddleware({
    loginSuccessUrl: '/nms',
    loginFailureUrl: '/user/login?invalid=true',
  }),
);
router.get('/nms*', access(AccessRoles.USER), handleReact('nms'));

router.get(
  '/host/networks/async',
  asyncHandler(async (_: FBCNMSRequest, res) => {
    const networks = await MagmaV1API.getNetworks();
    res.status(200).send(networks);
  }),
);

router.use('/host', hostOrgMiddleware, hostRoutes);

async function handleHost(req: FBCNMSRequest, res) {
  const appData: EmbeddedData = {
    csrfToken: req.csrfToken(),
    user: {
      tenant: 'host',
      email: req.user.email,
      isSuperUser: req.user.isSuperUser,
      isReadOnlyUser: req.user.isReadOnlyUser,
    },
    enabledFeatures: await getEnabledFeatures(req, 'host'),
    tabs: [],
    ssoEnabled: false,
    ssoSelectedType: 'none',
    csvCharset: null,
  };
  res.render('host', {
    staticDist,
    configJson: JSON.stringify({appData}),
  });
}

router.get('/host*', hostOrgMiddleware, handleHost);

router.get(
  '/*',
  access(AccessRoles.USER),
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organization = await req.organization();
    if (organization.isHostOrg) {
      res.redirect('/host');
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
