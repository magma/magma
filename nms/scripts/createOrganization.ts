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

import OrchestartorAPI from '../server/api/OrchestratorAPI';
import axios from 'axios';
import {Organization, sequelize} from '../shared/sequelize_models';
import {OrganizationRawType} from '../shared/sequelize_models/models/organization';
import {
  createJoinedOrganization,
  emptyOrg,
  getOrganizationByName,
} from '../server/host/routes';
import {union} from 'lodash';

type OrganizationObject = {
  name: string;
  networkIDs: Array<string>;
  csvCharset: '';
};

async function updateOrganization(
  joinedOrganization: OrganizationRawType,
  organizationObject: OrganizationObject,
) {
  console.log(
    `Updating organization ${organizationObject.name} to: ` +
      `networkIDs=[${organizationObject.networkIDs.join(' ')}]`,
  );
  let nmsOrganization = await Organization.findByPk(joinedOrganization.id);
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
    await nmsOrganization.update(
      {
        networkIDs: union(
          joinedOrganization.networkIDs ?? [],
          organizationObject.networkIDs,
        ),
      },
      {
        transaction,
      },
    );
    const unionNetworkIDs = [
      ...new Set([
        ...(joinedOrganization.networkIDs ?? []),
        ...organizationObject.networkIDs,
      ]),
    ];
    if (organizationObject.networkIDs !== unionNetworkIDs) {
      await OrchestartorAPI.tenants.tenantsTenantIdPut({
        tenant: {
          id: joinedOrganization.id,
          name: joinedOrganization.name,
          networks: unionNetworkIDs,
        },
        tenantId: joinedOrganization.id,
      });
    }
  });
}

async function createOrganization(organizationObject: OrganizationObject) {
  console.log(
    `Creating a new organization: name=${organizationObject.name}, ` +
      `networkIDs=[${organizationObject.networkIDs.join(' ')}]`,
  );
  await createJoinedOrganization(
    organizationObject.name,
    organizationObject.networkIDs,
    [],
  );
}

async function createOrUpdateOrganization(
  organizationObject: OrganizationObject,
) {
  const joinedOrganization = await getOrganizationByName(
    organizationObject.name,
  );
  if (!joinedOrganization) {
    await Promise.all([createOrganization(organizationObject)]);
  } else {
    await updateOrganization(joinedOrganization, organizationObject);
  }
}

function main() {
  const args = process.argv.slice(2);
  if (args.length < 1) {
    console.log(
      'Usage: createOrganization.ts <name> <networkID>,<networkID>, ...',
    );
    process.exit(1);
  }

  const networkIDs = (args[1] || '').split(',');
  const organizationObject = {
    name: args[0],
    networkIDs,
    csvCharset: '',
  } as const;
  createOrUpdateOrganization(organizationObject)
    .then(() => {
      console.log('Success');
      process.exit();
    })
    .catch(err => {
      if (axios.isAxiosError(err)) {
        console.log(
          `Error: Status: ${
            err?.response?.status ?? 500
          }: ${(err as Error).toString()}`,
        );
      } else {
        console.log(err);
      }
      process.exit(1);
    });
}

main();
