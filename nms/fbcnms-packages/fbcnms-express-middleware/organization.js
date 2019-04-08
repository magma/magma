/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Sequelize from 'sequelize';
import type {ExpressRequest, ExpressResponse, NextFunction} from 'express';

function getSubdomainList(host: ?string): Array<string> {
  if (!host) {
    return [];
  }
  const subdomainList = host.split('.');
  if (subdomainList) {
    subdomainList.splice(-1, 1);
  }
  return subdomainList;
}

export async function getOrganization(
  req: {get(field: string): string | void},
  OrganizationModel: any,
): Object {
  const host = req.get('host') || 'UNKNOWN_HOST';
  let org = await OrganizationModel.findOne({
    where: Sequelize.fn(
      'JSON_CONTAINS',
      Sequelize.col('customDomains'),
      `"${host}"`,
    ),
  });

  if (org) {
    return org;
  }

  const subDomains = getSubdomainList(host);
  if (subDomains.length != 1 && subDomains.length != 2) {
    throw new Error('Invalid organization!');
  }
  org = await OrganizationModel.findOne({
    where: {
      name: subDomains[0],
    },
  });
  if (!org) {
    throw new Error('Invalid organization!');
  }
  return org;
}

// We don't depend on organization to be there in other modules.
export type OrganizationRequestPart = {organization: () => Promise<Object>};
export type OrganizationMiddlewareRequest = ExpressRequest &
  $Shape<OrganizationRequestPart>;

export function organizationMiddleware({model}: {model: any}) {
  return (
    req: OrganizationMiddlewareRequest,
    res: ExpressResponse,
    next: NextFunction,
  ) => {
    try {
      req.organization = () => getOrganization(req, model);
      next();
    } catch (err) {
      res.status(404).send();
      next(err);
    }
  };
}
