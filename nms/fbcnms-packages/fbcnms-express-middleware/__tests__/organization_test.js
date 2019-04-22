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
import {Organization} from '@fbcnms/sequelize-models';

const ORGS = [
  {
    id: '1',
    name: 'custom_domain_org',
    customDomains: ['subdomain.localtest.me'],
    networkIDs: [],
  },
  {
    id: '2',
    name: 'subdomain',
    customDomains: [],
    networkIDs: [],
  },
];

describe('organization tests', () => {
  beforeEach(async () => {
    ORGS.forEach(async organization => await Organization.create(organization));
  });

  it('should allow custom domain', async () => {
    const request = {
      get: () => 'subdomain.localtest.me',
    };

    const org = await getOrganization(request);
    expect(org.name).toBe(ORGS[0].name);
  });

  it('should allow org by subdomain', async () => {
    const request = {
      get: () => 'subdomain.phbcloud.io',
    };

    const org = await getOrganization(request);
    expect(org.name).toBe(ORGS[1].name);
  });

  it('should throw an exception when no org is found', async () => {
    const request = {
      get: () => 'unknowndomain.com',
    };

    await expect(getOrganization(request)).rejects.toThrow();
  });
});
