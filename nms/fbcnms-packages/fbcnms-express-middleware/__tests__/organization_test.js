/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {getOrganization} from '@fbcnms/express-middleware/organizationMiddleware';
import {find} from 'lodash-es';

import type {StaticOrganizationModel} from '@fbcnms/sequelize-models/models/organization';

const ORGS = [
  {
    id: '1',
    name: 'custom_domain_org',
    customDomains: ['subdomain.localtest.me'],
  },
  {
    id: '2',
    name: 'subdomain',
    customDomains: [],
  },
];

// $FlowIgnore We know this is wrong, but don't want to deal with it.
const MockOrganization: StaticOrganizationModel = {
  findOne: ({where}) => {
    if (where.name) {
      return find(ORGS, where);
    }

    return find(ORGS, {customDomains: [JSON.parse(where.args[1])]});
  },
};

describe('organization tests', () => {
  it('should allow custom domain', async () => {
    const request = {
      get: () => 'subdomain.localtest.me',
    };

    const org = await getOrganization(request, MockOrganization);
    expect(org.name).toBe(ORGS[0].name);
  });

  it('should allow org by subdomain', async () => {
    const request = {
      get: () => 'subdomain.phbcloud.io',
    };

    const org = await getOrganization(request, MockOrganization);
    expect(org.name).toBe(ORGS[1].name);
  });

  it('should throw an exception when no org is found', async () => {
    const request = {
      get: () => 'unknowndomain.com',
    };

    await expect(getOrganization(request, MockOrganization)).rejects.toThrow();
  });
});
