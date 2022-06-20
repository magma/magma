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

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '../auth/access';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {FeatureID} from '../../shared/types/features';

import Sequelize from 'sequelize';
import asyncHandler from '../util/asyncHandler';
import crypto from 'crypto';
import express from 'express';
import featureConfigs from '../features';
// $FlowFixMe migrated to typescript
import logging from '../../shared/logging';
import {FeatureFlag, Organization} from '../../shared/sequelize_models';
import {User} from '../../shared/sequelize_models';
import {getPropsToUpdate} from '../auth/util';

const logger = logging.getLogger(module);

const router: express.Router<FBCNMSRequest, ExpressResponse> = express.Router();

router.get(
  '/organization/async',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organizations = await Organization.findAll();
    res.status(200).send({organizations});
  }),
);

router.get(
  '/organization/async/:name',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organization = await Organization.findOne({
      where: {
        name: Sequelize.where(
          Sequelize.fn('lower', Sequelize.col('name')),
          Sequelize.fn('lower', req.params.name),
        ),
      },
    });
    res.status(200).send({organization});
  }),
);

router.get(
  '/organization/async/:name/users',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const users = await User.findAll({
      where: {
        organization: req.params.name,
      },
    });
    res.status(200).send(users);
  }),
);

const configFromFeatureFlag = flag => ({
  id: flag.id,
  enabled: flag.enabled,
});
router.get(
  '/feature/async',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    // $FlowFixMe: results needs to be typed correctly
    const results: {[string]: any} = {...featureConfigs};
    Object.keys(results).forEach(id => (results[id].config = {}));
    const featureFlags = await FeatureFlag.findAll();
    featureFlags.forEach(flag => {
      if (!results[flag.featureId]) {
        logger.error(
          'feature config is missing for featureId: ' + flag.featureId,
        );
      } else {
        results[flag.featureId].config[
          flag.organization
        ] = configFromFeatureFlag(flag);
      }
    });
    res.status(200).send(Object.values(results));
  }),
);

router.post(
  '/feature/async/:featureId',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    // $FlowFixMe: Ensure it's a FeatureID
    const featureId: FeatureID = req.params.featureId;
    // $FlowFixMe: results needs to be typed correctly
    const results: {[string]: any} = {...featureConfigs};
    results.config = {};
    const {toUpdate, toDelete, toCreate} = req.body;
    const featureFlags = await FeatureFlag.findAll({where: {featureId}});
    await Promise.all(
      featureFlags.map(async flag => {
        if (toUpdate[flag.id]) {
          const newFlag = await flag.update({
            enabled: toUpdate[flag.id].enabled,
          });
          results.config[flag.organization] = configFromFeatureFlag(newFlag);
        } else if (toDelete[flag.id] !== undefined) {
          await FeatureFlag.destroy({where: {id: flag.id}});
        }
      }),
    );

    await Promise.all(
      toCreate.map(async data => {
        const flag = await FeatureFlag.create({
          featureId,
          organization: data.organization,
          enabled: data.enabled,
        });

        results.config[flag.organization] = configFromFeatureFlag(flag);
      }),
    );

    res.status(200).send(results);
  }),
);

router.post(
  '/organization/async',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    let organization = await Organization.findOne({
      where: {
        name: Sequelize.where(
          Sequelize.fn('lower', Sequelize.col('name')),
          Sequelize.fn('lower', req.body.name),
        ),
      },
    });
    if (organization) {
      return res.status(404).send({error: 'Organization already exists'});
    }
    organization = await Organization.create({
      name: req.body.name,
      networkIDs: req.body.networkIDs,
      customDomains: req.body.customDomains,
      csvCharset: '',
      ssoCert: '',
      ssoEntrypoint: '',
      ssoIssuer: '',
    });
    res.status(200).send({organization});
  }),
);

router.put(
  '/organization/async/:name',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organization = await Organization.findOne({
      where: {
        name: Sequelize.where(
          Sequelize.fn('lower', Sequelize.col('name')),
          Sequelize.fn('lower', req.body.name),
        ),
      },
    });
    if (!organization) {
      return res.status(404).send({error: 'Organization does not exist'});
    }
    const updated = await organization.update(req.body);
    res.status(200).send({organization: updated});
  }),
);

const USER_PROPS = ['email', 'networkIDs', 'password', 'role', 'organization'];

router.post(
  '/organization/async/:name/add_user',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    const organization = await Organization.findOne({
      where: {
        name: Sequelize.where(
          Sequelize.fn('lower', Sequelize.col('name')),
          Sequelize.fn('lower', req.params.name),
        ),
      },
    });
    if (!organization) {
      return res.status(404).send({error: 'Organization does not exist'});
    }

    try {
      let props = {organization: req.params.name, ...req.body};
      props = await getPropsToUpdate(USER_PROPS, props, async params => ({
        ...params,
        organization: req.params.name,
      }));

      // this happens when the user is being added to an organization that
      // uses SSO for login, give it a random password
      if (props.password === undefined) {
        const organization = await Organization.findOne({
          where: {
            name: Sequelize.where(
              Sequelize.fn('lower', Sequelize.col('name')),
              Sequelize.fn('lower', req.params.name),
            ),
          },
        });
        if (organization && organization.ssoEntrypoint) {
          props.password = crypto.randomBytes(16).toString('hex');
        }
      }

      const user = await User.create(props);
      res.status(200).send({user});
    } catch (error) {
      res.status(400).send({error: error.toString()});
    }
  }),
);

router.delete(
  '/organization/async/:id',
  asyncHandler(async (req: FBCNMSRequest, res) => {
    await Organization.destroy({
      where: {id: req.params.id},
      individualHooks: true,
    });
    res.status(200).send({success: true});
  }),
);

export default router;
