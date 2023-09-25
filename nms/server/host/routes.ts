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

import OrchestratorAPI from '../api/OrchestratorAPI';
import Sequelize from 'sequelize';
import asyncHandler from '../util/asyncHandler';
import crypto from 'crypto';
import featureConfigs, {FeatureConfig} from '../features';
import logging from '../../shared/logging';
import {
  FeatureFlag,
  Organization,
  sequelize,
} from '../../shared/sequelize_models';
import {FeatureFlagModel} from '../../shared/sequelize_models/models/featureflag';
import {Request, Router} from 'express';
import {User} from '../../shared/sequelize_models';
import {UserRawType} from '../../shared/sequelize_models/models/user';
import {getPropsToUpdate} from '../auth/util';
import {
  rethrowUnlessNotFoundError,
  syncOrganizationWithOrc8rTenant,
} from '../util/tenantsSync';
import type {FeatureID} from '../../shared/types/features';

const logger = logging.getLogger(module);

const router = Router();

router.get(
  '/organization/async',
  asyncHandler(async (req: Request, res) => {
    const organizations = await Organization.findAll();
    res.status(200).send({organizations});
  }),
);

router.get(
  '/organization/async/:name',
  asyncHandler(async (req: Request<{name: string}>, res) => {
    const organization = await Organization.findOne({
      where: Sequelize.where(
        Sequelize.fn('lower', Sequelize.col('name')),
        Sequelize.fn('lower', req.params.name),
      ),
    });
    res.status(200).send({organization});
  }),
);

router.get(
  '/organization/async/:name/users',
  asyncHandler(async (req: Request<{name: string}>, res) => {
    const users = await User.findAll({
      where: {
        organization: req.params.name,
      },
    });
    res.status(200).send(users);
  }),
);

type FeatureFlagConfig = Record<string, {id: number; enabled: boolean}>;

const configFromFeatureFlag = (flag: FeatureFlagModel) => ({
  id: flag.id,
  enabled: flag.enabled,
});
router.get(
  '/feature/async',
  asyncHandler(async (req: Request, res) => {
    const results: Record<
      FeatureID,
      FeatureConfig & {
        config?: FeatureFlagConfig;
      }
    > = {...featureConfigs};
    (Object.keys(results) as Array<FeatureID>).forEach(
      id => (results[id].config = {}),
    );
    const featureFlags = await FeatureFlag.findAll();
    featureFlags.forEach(flag => {
      if (!results[flag.featureId]) {
        logger.error(
          'feature config is missing for featureId: ' + flag.featureId,
        );
      } else {
        results[flag.featureId].config![
          flag.organization
        ] = configFromFeatureFlag(flag);
      }
    });
    res.status(200).send(Object.values(results));
  }),
);

router.post(
  '/feature/async/:featureId',
  asyncHandler(
    async (
      req: Request<
        {featureId: FeatureID},
        any,
        {
          toUpdate: Record<number, {enabled: boolean}>;
          toDelete: Record<number, boolean>;
          toCreate: Array<{organization: string; enabled: boolean}>;
        }
      >,
      res,
    ) => {
      const featureId = req.params.featureId;
      const result: FeatureConfig & {
        config?: FeatureFlagConfig;
      } = featureConfigs[featureId];
      const {toUpdate, toDelete, toCreate} = req.body;
      const featureFlags = await FeatureFlag.findAll({where: {featureId}});
      await Promise.all(
        featureFlags.map(async flag => {
          if (toUpdate[flag.id]) {
            const newFlag = await flag.update({
              enabled: toUpdate[flag.id].enabled,
            });
            result.config![flag.organization] = configFromFeatureFlag(newFlag);
          } else if (toDelete[flag.id] !== undefined) {
            await FeatureFlag.destroy({where: {id: flag.id}});
            delete result.config![flag.organization];
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
          result.config![flag.organization] = configFromFeatureFlag(flag);
        }),
      );

      res.status(200).send(result);
    },
  ),
);

router.post(
  '/organization/async',
  asyncHandler(
    async (
      req: Request<
        never,
        any,
        {name: string; networkIDs: Array<string>; customDomains: Array<string>}
      >,
      res,
    ) => {
      let organization = await Organization.findOne({
        where: Sequelize.where(
          Sequelize.fn('lower', Sequelize.col('name')),
          Sequelize.fn('lower', req.body.name),
        ),
      });
      if (organization) {
        return res.status(409).send({message: 'Organization exists already'});
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
      await syncOrganizationWithOrc8rTenant(organization);
      res.status(200).send({organization});
    },
  ),
);

router.put(
  '/organization/async/:name',
  asyncHandler(async (req: Request<never, any, {name: string}>, res) => {
    const organization = await Organization.findOne({
      where: Sequelize.where(
        Sequelize.fn('lower', Sequelize.col('name')),
        Sequelize.fn('lower', req.body.name),
      ),
    });
    if (!organization) {
      return res.status(404).send({message: 'Organization does not exist'});
    }
    const updated = await organization.update(req.body);
    await syncOrganizationWithOrc8rTenant(updated);
    res.status(200).send({organization: updated});
  }),
);

const USER_PROPS = [
  'email',
  'networkIDs',
  'password',
  'role',
  'organization',
] as const;

router.post(
  '/organization/async/:name/add_user',
  asyncHandler(
    async (req: Request<{name: string}, any, Partial<UserRawType>>, res) => {
      const organization = await Organization.findOne({
        where: Sequelize.where(
          Sequelize.fn('lower', Sequelize.col('name')),
          Sequelize.fn('lower', req.params.name),
        ),
      });
      if (!organization) {
        return res.status(404).send({message: 'Organization does not exist'});
      }

      try {
        const props = await getPropsToUpdate(
          USER_PROPS,
          {
            organization: req.params.name,
            ...req.body,
          },
          params =>
            Promise.resolve({
              ...params,
              organization: req.params.name,
            }),
        );

        // this happens when the user is being added to an organization that
        // uses SSO for login, give it a random password
        if (props.password === undefined) {
          const organization = await Organization.findOne({
            where: Sequelize.where(
              Sequelize.fn('lower', Sequelize.col('name')),
              Sequelize.fn('lower', req.params.name),
            ),
          });
          if (organization && organization.ssoEntrypoint) {
            props.password = crypto.randomBytes(16).toString('hex');
          }
        }

        const user = await User.create(props);
        res.status(200).send({user});
      } catch (error) {
        res.status(400).send({message: (error as Error).toString()});
      }
    },
  ),
);

router.delete(
  '/organization/async/:id',
  asyncHandler(async (req: Request<{id: string}>, res) => {
    const organization = await Organization.findOne({
      where: {id: req.params.id},
    });

    if (!organization) {
      await deleteOrc8rTenant(+req.params.id);
      return res.status(200).send({success: true});
    }

    await sequelize.transaction(async transaction => {
      await organization.destroy({transaction});

      await User.destroy({
        where: {organization: organization.name},
        individualHooks: true,
        transaction,
      });
    });

    await deleteOrc8rTenant(+req.params.id);
    res.status(200).send({success: true});
  }),
);

async function deleteOrc8rTenant(organizationId: number) {
  try {
    await OrchestratorAPI.tenants.tenantsTenantIdDelete({
      tenantId: organizationId,
    });
  } catch (error) {
    // Ignore "not found" since there is no guarantee NMS and Orc8r are in sync
    rethrowUnlessNotFoundError(error);
  }
}

export default router;
