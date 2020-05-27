/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '@fbcnms/auth/access';

const url = require('url');
import pathToRegexp from 'path-to-regexp';

import {AuditLogEntry} from '@fbcnms/sequelize-models';
const logger = require('@fbcnms/logging').getLogger(module);

const defaultResolver = (req: FBCNMSRequest, type: string) => {
  const {search} = url.parse(req.originalUrl);
  const params = new URLSearchParams(search ?? '');
  return [params.get('requested_id'), type];
};

const PATHS: Array<{
  path: string,
  type?: string,
  resolver?: (FBCNMSRequest, string[]) => [?string, ?string],
}> = [
  {
    path: '/magma/networks/:networkId/gateways',
    resolver: (req: FBCNMSRequest) => defaultResolver(req, 'gateway'),
  },
  {
    path: '/magma/networks/:networkId/configs/devices',
    resolver: (req: FBCNMSRequest) => defaultResolver(req, 'device'),
  },
  {
    path: '/magma/networks/:networkId/gateways/:objectId',
    type: 'gateway',
  },
  {
    path: '/magma/networks/:networkId/gateways/:objectId/configs/cellular',
    type: 'gateway cellular config',
  },
  {
    path: '/magma/networks/:networkId/gateways/:objectId/configs/wifi',
    type: 'gateway wifi config',
  },
  {
    path: '/magma/networks/:networkId/gateways/:objectId/configs/tarazed',
    type: 'gateway tarazed config',
  },
  {
    path: '/magma/networks/:networkId/gateways/:objectId/configs/devmand',
    type: 'gateway devmand config',
  },
  {
    path: '/magma/networks/:networkId/gateways/:objectId/configs',
    type: 'gateway config',
  },
  {
    path: '/magma/networks/:networkId/subscribers/:objectId',
    type: 'subscriber',
  },
  {
    path: '/magma/networks/:networkId/subscribers',
    resolver: (req: FBCNMSRequest) => [req.body.id, 'subscriber'],
  },
  {
    path: '/magma/networks/:networkId/tiers/:objectId',
    type: 'network tier',
  },
  {
    path: '/magma/networks/:networkId/tiers',
    resolver: (req: FBCNMSRequest) => [req.body.id, 'network tier'],
  },
  {
    path: '/magma/networks/:networkId/configs/cellular',
    resolver: (_, params) => [params[2], 'network cellular configs'],
  },
  {
    path: '/magma/networks/:networkId/configs/wifi',
    resolver: (_, params) => [params[2], 'network wifi configs'],
  },
  {
    path: '/magma/v1/networks/:networkId',
    resolver: (_, params) => [params[1], 'network'],
  },
  {
    path: '/magma/v1/cwf/:networkId/gateways',
    resolver: req => [req.body.id, 'carrier wifi gateway'],
  },
  {
    path: '/magma/v1/cwf/:networkId/gateways/:objectId',
    resolver: (_, params) => [params[2], 'carrier wifi gateway'],
  },
  {
    path: '/magma/v1/feg/:networkId/gateways',
    resolver: req => [req.body.id, 'federation gateway'],
  },
  {
    path: '/magma/v1/feg/:networkId/gateways/:objectId',
    resolver: (_, params) => [params[2], 'federation gateway'],
  },
  {
    path: '/magma/v1/feg/:networkId/gateways/:objectId/federation',
    resolver: (_, params) => [params[2], 'federation gateway config'],
  },
  {
    path: '/magma/v1/networks/:networkId/rules/policies',
    resolver: req => [req.body.id, 'policy'],
  },
  {
    path: '/magma/v1/networks/:networkId/policies/rules/:objectId',
    resolver: (_, params) => [params[2], 'policy'],
  },
];

const MUTATION_TYPE_MAP = {
  POST: 'CREATE',
  PUT: 'UPDATE',
  DELETE: 'DELETE',
};

function getObjectIdAndType(req: FBCNMSRequest): [?string, ?string] {
  const parsed = url.parse(
    req.originalUrl.replace(/^\/nms\/apicontroller/, ''),
  );
  for (let i = 0; i < PATHS.length; i++) {
    const params = pathToRegexp(PATHS[i].path).exec(parsed.pathname);
    if (params) {
      return PATHS[i].resolver
        ? PATHS[i].resolver(req, params)
        : [params[2], PATHS[i].type];
    }
  }

  return [null, null];
}

export default async function auditLoggingDecorator(
  proxyRes: ExpressResponse,
  proxyResData: Buffer,
  userReq: FBCNMSRequest,
  _userRes: ExpressResponse,
) {
  if (!MUTATION_TYPE_MAP[userReq.method]) {
    return proxyResData;
  }

  const [objectId, objectType] = getObjectIdAndType(userReq);
  if (!objectId || !objectType) {
    return proxyResData;
  }

  let organizationName = '';
  if (userReq.organization) {
    const organization = await userReq.organization();
    organizationName = organization.name;
  }

  const data = {
    actingUserId: userReq.user.id,
    actingUserEmail: userReq.user.email,
    organization: organizationName,
    mutationType: MUTATION_TYPE_MAP[userReq.method],
    objectId,
    objectType,
    objectDisplayName: objectId,
    mutationData: userReq.body,
    url: userReq.originalUrl,
    ipAddress: userReq.ip,
    status: proxyRes.statusCode < 300 ? 'SUCCESS' : 'FAILURE',
    statusCode: `${proxyRes.statusCode}`,
  };

  try {
    await AuditLogEntry.create(data);
  } catch (error) {
    logger.error('Error creating AuditLogEntry', error);
  }
  return proxyResData;
}
