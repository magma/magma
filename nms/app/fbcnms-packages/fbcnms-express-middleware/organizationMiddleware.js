/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import {Organization, jsonArrayContains} from '@fbcnms/sequelize-models';
import type {ExpressRequest, ExpressResponse, NextFunction} from 'express';
import type {OrganizationType} from '@fbcnms/sequelize-models/models/organization';

async function getOrganizationFromHost(
  host: string,
): Promise<?OrganizationType> {
  const org = await Organization.findOne({
    where: jsonArrayContains('customDomains', host),
  });
  if (org) {
    return org;
  }
  const subdomain = host.split('.')[0];
  return await Organization.findOne({
    where: {
      name: subdomain,
    },
  });
}

export async function getOrganization(req: {
  get(field: string): string | void,
}): Promise<OrganizationType> {
  for (const header of ['host', 'x-forwarded-host']) {
    const host = req.get(header);
    if (host != null && host !== '') {
      const org = await getOrganizationFromHost(host);
      if (org) {
        return org;
      }
    }
  }
  throw new Error('Invalid organization!');
}

// We don't depend on organization to be there in other modules.
export type OrganizationRequestPart = {
  organization: () => Promise<OrganizationType>,
};
export type OrganizationMiddlewareRequest = ExpressRequest &
  $Shape<OrganizationRequestPart>;

export default function organizationMiddleware() {
  return (
    req: OrganizationMiddlewareRequest,
    res: ExpressResponse,
    next: NextFunction,
  ) => {
    try {
      req.organization = () => getOrganization(req);
      next();
    } catch (err) {
      res.status(404).send();
      next(err);
    }
  };
}
