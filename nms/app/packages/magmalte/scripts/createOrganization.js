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

import type {OrganizationType} from '@fbcnms/sequelize-models/models/organization';
import type {Tab} from '@fbcnms/types/tabs';

import Sequelize from 'sequelize';

import {Organization} from '@fbcnms/sequelize-models';
import {coerceToTab} from '@fbcnms/types/tabs';
import {difference} from 'lodash';
import {getProjectTabs} from '@fbcnms/projects/projects';
import {union} from 'lodash';

type OrganizationObject = {
  name: string,
  tabs: Array<Tab>,
  networkIDs: Array<string>,
  csvCharset: '',
};

async function updateOrganization(
  organization: OrganizationType,
  organizationObject: OrganizationObject,
) {
  console.log(
    `Updating organization ${organizationObject.name} to: ` +
      `tabs=${organizationObject.tabs.join(' ')}, ` +
      `networkIDs=[${organizationObject.networkIDs.join(' ')}]`,
  );
  await organization.update({
    tabs: organizationObject.tabs ?? ['inventory'],
    networkIDs: union(
      organization.networkIDs ?? [],
      organizationObject.networkIDs,
    ),
  });
}

async function createOrganization(organizationObject: OrganizationObject) {
  console.log(
    `Creating a new organization: name=${organizationObject.name}, ` +
      `tabs=${organizationObject.tabs.join(' ')}, ` +
      `networkIDs=[${organizationObject.networkIDs.join(' ')}]`,
  );
  await Organization.create({
    name: organizationObject.name,
    tabs: organizationObject.tabs,
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
  if (args.length < 2) {
    console.log(
      'Usage: createOrganization.js <name> <tab>,<tab>,... <networkID>,<networkID>, ...',
    );
    process.exit(1);
  }

  const validTabs = getProjectTabs();
  const tabs = args[1].split(',').map(tab => coerceToTab(tab));
  const invalidTabs = difference(
    tabs,
    validTabs.map(tab => tab.id),
  ).join(', ');
  if (invalidTabs) {
    console.log(
      `tab should be one of: ${validTabs
        .map(tab => tab.id)
        .join(', ')}. Got: ${invalidTabs}`,
    );
    process.exit(1);
  }

  const networkIDs = (args[2] || '').split(',');
  const organizationObject = {
    name: args[0],
    tabs: tabs,
    networkIDs: networkIDs,
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
