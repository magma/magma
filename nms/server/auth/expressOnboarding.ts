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

import staticDist from '../../config/staticDist';
import {AccessRoles} from '../../shared/roles';
import {NextFunction, Request, RequestHandler, Response, Router} from 'express';
import {Organization, User} from '../../shared/sequelize_models';
import {UserRawType} from '../../shared/sequelize_models/models/user';
import {getPropsToUpdate} from './util';
import {rateLimitMiddleware} from '../middleware';

export default function () {
  const asyncOnboardingMiddleware = async (
    req: Request,
    res: Response,
    next: NextFunction,
  ) => {
    if (req.isAuthenticated()) {
      res.redirect('/');
    } else if (await User.findOne()) {
      res.redirect('/user/login');
    }
    next();
  };

  const onboardingMiddleware: RequestHandler = (req, res, next) => {
    void asyncOnboardingMiddleware(req, res, next);
  };

  const router = Router();

  router.get('/onboarding', onboardingMiddleware, (req: Request, res) => {
    res.render('onboarding', {
      staticDist,
      configJson: JSON.stringify({
        appData: {
          csrfToken: req.csrfToken(),
        },
      }),
    });
  });

  router.post(
    '/onboarding',
    rateLimitMiddleware,
    onboardingMiddleware,
    async (req: Request<never, any, Partial<UserRawType>>, res) => {
      try {
        const allowedProps = ['email', 'password', 'organization'] as const;
        const userProps = await getPropsToUpdate(
          allowedProps,
          req.body,
          props =>
            Promise.resolve({
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
        res.status(400).send({message: (error as Error).toString()});
      }
    },
  );

  return router;
}
