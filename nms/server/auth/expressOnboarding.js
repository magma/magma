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

import type {ExpressResponse, NextFunction} from 'express';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {FBCNMSRequest} from './access';

import express from 'express';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import staticDist from '../../config/staticDist';
// $FlowFixMe migrated to typescript
import {AccessRoles} from '../../shared/roles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {Organization, User} from '../../shared/sequelize_models';
// $FlowFixMe migrated to typescript
import {getPropsToUpdate} from './util';

export default function () {
  const onboardingMiddleware = async (
    req: FBCNMSRequest,
    res: ExpressResponse,
    next: NextFunction,
  ) => {
    if (req.isAuthenticated()) {
      res.redirect('/');
    } else if (await User.findOne()) {
      res.redirect('/user/login');
    }
    next();
  };

  const router: express.Router<
    FBCNMSRequest,
    ExpressResponse,
  > = express.Router();

  router.get(
    '/onboarding',
    onboardingMiddleware,
    async (req: FBCNMSRequest, res) => {
      res.render('onboarding', {
        staticDist,
        configJson: JSON.stringify({
          appData: {
            csrfToken: req.csrfToken(),
          },
        }),
      });
    },
  );

  router.post(
    '/onboarding',
    onboardingMiddleware,
    async (req: FBCNMSRequest, res) => {
      try {
        const allowedProps = ['email', 'password', 'organization'];
        const userProps = await getPropsToUpdate(
          allowedProps,
          req.body,
          async props => ({
            ...props,
            organization: req.body.organization,
          }),
        );
        userProps.role = AccessRoles.SUPERUSER;
        await User.create(userProps);

        await Organization.create({
          name: req.body.organization,
          networkIDs: [],
          csvCharset: '',
          customDomains: [],
          ssoCert: '',
          ssoEntrypoint: '',
          ssoIssuer: '',
        });

        res.status(200).send({success: true});
      } catch (error) {
        res.status(400).send({error: error.toString()});
      }
    },
  );

  return router;
}
