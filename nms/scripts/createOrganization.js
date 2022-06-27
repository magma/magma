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
'use strict';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {OrganizationType} from '../shared/sequelize_models/models/organization';

import Sequelize from 'sequelize';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {Organization} from '../shared/sequelize_models';
import {union} from 'lodash';

type OrganizationObject = {
  name: string,
  networkIDs: Array<string>,
  csvCharset: '',
};

async function updateOrganization(
  organization: OrganizationType,
  organizationObject: OrganizationObject,
) {
  console.log(
    `Updating organization ${organizationObject.name} to: ` +
      `networkIDs=[${organizationObject.networkIDs.join(' ')}]`,
  );
  await organization.update({
    networkIDs: union(
      organization.networkIDs ?? [],
      organizationObject.networkIDs,
    ),
  });
}

async function createOrganization(organizationObject: OrganizationObject) {
  console.log(
    `Creating a new organization: name=${organizationObject.name}, ` +
      `networkIDs=[${organizationObject.networkIDs.join(' ')}]`,
  );
  await Organization.create({
    name: organizationObject.name,
    networkIDs: organizationObject.networkIDs,
    csvCharset: '',
    ssoCert: '',
    ssoEntrypoint: '',
    ssoIssuer: '',
  });
}

async function createOrUpdateOrganization(
  organizationObject: OrganizationObject,
) {
  const organization = await Organization.findOne({
    where: {
      name: Sequelize.where(
        Sequelize.fn('lower', Sequelize.col('name')),
        Sequelize.fn('lower', organizationObject.name),
      ),
    },
  });
  if (!organization) {
    await Promise.all([createOrganization(organizationObject)]);
  } else {
    await updateOrganization(organization, organizationObject);
  }
}

function main() {
  const args = process.argv.slice(2);
  if (args.length < 1) {
    console.log(
      'Usage: createOrganization.js <name> <networkID>,<networkID>, ...',
    );
    process.exit(1);
  }

  const networkIDs = (args[1] || '').split(',');
  const organizationObject = {
    name: args[0],
    networkIDs,
    csvCharset: '',
  };
  createOrUpdateOrganization(organizationObject)
    .then(_res => {
      console.log('Success');
      process.exit();
    })
    .catch(err => {
      console.error(err);
      process.exit(1);
    });
}

main();
