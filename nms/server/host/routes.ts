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

import OrchestartorAPI from '../api/OrchestratorAPI';
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
import {Tenant} from '../../generated';
import {User} from '../../shared/sequelize_models';
import {UserRawType} from '../../shared/sequelize_models/models/user';
import {getPropsToUpdate} from '../auth/util';
import type {FeatureID} from '../../shared/types/features';

import axios from 'axios';
import {
  OrganizationModel,
  OrganizationRawType,
} from '../../shared/sequelize_models/models/organization';

const logger = logging.getLogger(module);

const router = Router();

export const emptyOrg = {
  customDomains: [],
  csvCharset: '',
  ssoSelectedType: 'none' as const,
  ssoCert: '',
  ssoEntrypoint: '',
  ssoIssuer: '',
  ssoOidcClientID: '',
  ssoOidcClientSecret: '',
  ssoOidcConfigurationURL: '',
};

function joinTenantAndOrganization(
  organization: OrganizationRawType | undefined | null,
  tenant: Tenant,
): OrganizationRawType {
  if (!organization) {
    return {
      ...emptyOrg,
      id: tenant.id,
      name: tenant.name,
      networkIDs: tenant.networks,
    };
  } else {
    return {
      id: tenant.id,
      name: tenant.name,
      networkIDs: tenant.networks,
      customDomains: organization.customDomains,
      csvCharset: organization.csvCharset,
      ssoSelectedType: organization.ssoSelectedType,
      ssoCert: organization.ssoCert,
      ssoEntrypoint: organization.ssoEntrypoint,
      ssoIssuer: organization.ssoIssuer,
      ssoOidcClientID: organization.ssoOidcClientID,
      ssoOidcClientSecret: organization.ssoOidcClientSecret,
      ssoOidcConfigurationURL: organization.ssoOidcConfigurationURL,
    };
  }
}

function joinTenantsAndOrganizations(
  orc8rTenants: Array<Tenant>,
  nmsOrganizations: Record<number, OrganizationModel>,
): Array<OrganizationRawType> {
  return orc8rTenants.map(tenant => {
    return joinTenantAndOrganization(nmsOrganizations[tenant.id], tenant);
  });
}

async function getAllJoinedOrganizations(): Promise<
  Array<OrganizationRawType>
> {
  const nmsOrganizations = await Organization.findAll();
  const nmsOrganizationsMap: Record<number, OrganizationModel> = {};
  nmsOrganizations.forEach((org: OrganizationModel) => {
    nmsOrganizationsMap[org.id] = org;
  });
  const orc8rTenants = (await OrchestartorAPI.tenants.tenantsGet()).data;
  return joinTenantsAndOrganizations(orc8rTenants, nmsOrganizationsMap);
}

router.get(
  '/organization/async',
  asyncHandler(async (req: Request, res) => {
    try {
      const organizations = await getAllJoinedOrganizations();
      res.status(200).send({organizations});
    } catch (error) {
      res.status(500).send({error: (error as Error).toString()});
    }
  }),
);

export async function getOrganizationByName(
  name: string,
): Promise<OrganizationRawType | null> {
  const organizations = await getAllJoinedOrganizations();
  for (const org of organizations) {
    if (org.name === name) return org;
  }
  return null;
}

router.get(
  '/organization/async/:name',
  asyncHandler(async (req: Request<{name: string}>, res) => {
    try {
      const organization = await getOrganizationByName(req.params.name);
      if (organization) {
        res.status(200).send({organization});
      } else {
        res.status(404).send({
          error: new Error(
            `Organization with name ${req.params.name} not found`,
          ).toString(),
        });
      }
    } catch (error) {
      if (axios.isAxiosError(error)) {
        res
          .status(error?.response?.status ?? 500)
          .send({error: (error as Error).toString()});
      } else {
        res.status(500).send({error: (error as Error).toString()});
      }
    }
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

export async function createJoinedOrganization(
  name: string,
  networkIDs: Array<string>,
  customDomains: Array<string>,
): Promise<OrganizationRawType> {
  let createdOrganization = <OrganizationRawType>{};
  await sequelize.transaction(async transaction => {
    createdOrganization = await Organization.create(
      {
        name: name,
        networkIDs: networkIDs,
        customDomains: customDomains,
        csvCharset: '',
        ssoCert: '',
        ssoEntrypoint: '',
        ssoIssuer: '',
      },
      {transaction},
    );

    // not ideal: since the ID is generated in Organizations table (and not explicitly set) we have to used the ID chosen in NMS. If ID already in use in orc8r, we get an error.
    const tenantToCreate: Tenant = {
      id: createdOrganization.id,
      name: name,
      networks: networkIDs,
    };

    await OrchestartorAPI.tenants.tenantsPost({tenant: tenantToCreate});
  });
  return createdOrganization;
}

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
      const joinedOrganization = await getOrganizationByName(req.body.name);
      if (joinedOrganization) {
        return res.status(409).send({error: 'Organization already exists'});
      }
      try {
        const createdOrganization = await createJoinedOrganization(
          req.body.name,
          req.body.networkIDs,
          req.body.customDomains,
        );
        res.status(200).send({organization: createdOrganization});
      } catch (error) {
        if (axios.isAxiosError(error)) {
          res
            .status(error?.response?.status ?? 500)
            .send({error: (error as Error).toString()});
        } else {
          res.status(500).send({error: (error as Error).toString()});
        }
      }
    },
  ),
);

router.put(
  '/organization/async/:name',
  asyncHandler(async (req: Request<never, any, OrganizationRawType>, res) => {
    const joinedOrganization = await getOrganizationByName(req.body.name);
    if (!joinedOrganization) {
      return res.status(404).send({error: 'Organization does not exist'});
    }
    let nmsOrganization = await Organization.findByPk(joinedOrganization.id);
    let updatedOrganization: OrganizationRawType | null = null;
    try {
      await sequelize.transaction(async transaction => {
        if (!nmsOrganization) {
          nmsOrganization = await Organization.create(
            {
              ...emptyOrg,
              name: joinedOrganization.name,
              networkIDs: joinedOrganization.networkIDs,
            },
            {transaction},
          );
        }
        updatedOrganization = await nmsOrganization.update(req.body, {
          transaction,
        });
        // the comparison is sensitive to order, not sure it really matters but if so
        // we could also do: [...new Set([...req.body.networkIDs,...joinedOrganization.networkIDs,])]
        if (
          req.body?.networkIDs &&
          req.body.networkIDs !== joinedOrganization.networkIDs
        ) {
          await OrchestartorAPI.tenants.tenantsTenantIdPut({
            tenant: {
              id: joinedOrganization.id,
              name: req.body.name,
              networks: req.body.networkIDs,
            },
            tenantId: joinedOrganization.id,
          });
        }
      });
    } catch (error) {
      if (axios.isAxiosError(error)) {
        return res
          .status(error?.response?.status ?? 500)
          .send({error: (error as Error).toString()});
      } else {
        return res.status(500).send({error: (error as Error).toString()});
      }
    }
    res.status(200).send({organization: updatedOrganization});
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
      const joinedOrganization = await getOrganizationByName(req.params.name);
      if (!joinedOrganization) {
        return res.status(404).send({error: 'Organization does not exist'});
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
          let nmsOrganization = await Organization.findByPk(
            joinedOrganization.id,
          );
          if (!nmsOrganization) {
            nmsOrganization = await Organization.create({
              ...emptyOrg,
              name: joinedOrganization.name,
              networkIDs: joinedOrganization.networkIDs,
            });
          }
          if (nmsOrganization && nmsOrganization.ssoEntrypoint) {
            props.password = crypto.randomBytes(16).toString('hex');
          }
        }

        const user = await User.create(props);
        res.status(200).send({user});
      } catch (error) {
        res.status(400).send({error: (error as Error).toString()});
      }
    },
  ),
);

router.delete(
  '/organization/async/:id',
  asyncHandler(async (req: Request<{id: number}>, res) => {
    let orc8rTenant: Tenant | null = null;
    try {
      orc8rTenant = (
        await OrchestartorAPI.tenants.tenantsTenantIdGet({
          tenantId: req.params.id,
        })
      ).data;
    } catch (error) {
      if (axios.isAxiosError(error)) {
        return res
          .status(error?.response?.status ?? 500)
          .send({error: (error as Error).toString()});
      } else {
        return res.status(500).send({error: (error as Error).toString()});
      }
    }
    if (orc8rTenant !== null) {
      const nmsOrganization = await Organization.findByPk(req.params.id);

      try {
        await sequelize.transaction(async transaction => {
          if (nmsOrganization !== null) {
            await nmsOrganization.destroy({transaction});

            await User.destroy({
              where: {organization: nmsOrganization.name},
              individualHooks: true,
              transaction,
            });
          }

          await OrchestartorAPI.tenants.tenantsTenantIdDelete({
            tenantId: req.params.id,
          });
        });
      } catch (error) {
        if (axios.isAxiosError(error)) {
          return res
            .status(error?.response?.status ?? 500)
            .send({error: (error as Error).toString()});
        } else {
          return res.status(500).send({error: (error as Error).toString()});
        }
      }
      res.status(200).send({success: true});
    } else {
      res.status(409).send({
        error: new Error(
          `Organization with id ${req.params.id} not found in orc8r tenants table`,
        ).toString(),
      });
    }
  }),
);

export default router;
