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

export async function getOrganization(req: {
  get(field: string): string | void,
}): Promise<OrganizationType> {
  const host = req.get('host') || 'UNKNOWN_HOST';
  let org = await Organization.findOne({
    where: jsonArrayContains('customDomains', host),
  });
  if (org) {
    return org;
  }

  const subDomain = host.split('.')[0];
  org = await Organization.findOne({
    where: {
      name: subDomain,
    },
  });
  if (!org) {
    throw new Error('Invalid organization!');
  }
  return org;
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
