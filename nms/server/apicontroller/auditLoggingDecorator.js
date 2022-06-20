/*
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
 * @flow
 * @format
 */

import type {ExpressResponse} from 'express';
import type {FBCNMSRequest} from '../auth/access';

const url = require('url');
import pathToRegexp from 'path-to-regexp';

import {AuditLogEntry} from '../../shared/sequelize_models';
// $FlowFixMe migrated to typescript
const logger = require('../../shared/logging.ts').getLogger(module);

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
    path: '/magma/networks/:networkId/gateways/:objectId/configs',
    type: 'gateway config',
  },
  {
    path: '/magma/networks/:networkId/subscribers/:objectId',
    type: 'Subscriber',
  },
  {
    path: '/magma/v1/lte/:networkId/subscribers/:objectId/(.*)',
    type: 'Subscriber',
  },
  {
    path: '/magma/networks/:networkId/subscribers',
    resolver: (req: FBCNMSRequest) => [req.body.id, 'Subscriber'],
  },
  {
    path: '/magma/networks/:networkId/tracing',
    type: 'Call Tracing',
  },
  {
    path: '/magma/networks/:networkId/tracing/(.*)',
    type: 'Call Tracing',
  },
  {
    path: '/magma/networks/:networkId/alerts/(.*)',
    type: 'Alerts',
  },
  {
    path: '/magma/networks/:networkId/tiers/:objectId',
    type: 'Network tier',
  },
  {
    path: '/magma/networks/:networkId/tiers',
    resolver: (req: FBCNMSRequest) => [req.body.id, 'Network tier'],
  },
  {
    path: '/magma/networks/:networkId/configs/cellular',
    resolver: (_, params) => [params[2], 'Network cellular configs'],
  },
  {
    path: '/magma/v1/networks/:networkId',
    resolver: (_, params) => [params[1], 'Network'],
  },
  {
    path: '/magma/v1/feg_lte/:networkId',
    resolver: (_, params) => [params[1], 'Federated LTE network'],
  },
  {
    path: '/magma/v1/feg_lte/:networkId/federation',
    resolver: (_, params) => [params[1], 'Federated LTE federation'],
  },
  {
    path: '/magma/v1/feg_lte/:networkId/subscriber_config',
    resolver: (_, params) => [params[1], 'Federated LTE subscriber config'],
  },
  {
    path: '/magma/v1/feg/:networkId',
    resolver: (_, params) => [params[1], 'Federation network'],
  },
  {
    path: '/magma/v1/lte/:networkId',
    resolver: (_, params) => [params[1], 'LTE network'],
  },
  {
    path: '/magma/v1/lte/:networkId/cellular/(.*)',
    resolver: (_, params) => [params[2], 'LTE network'],
  },
  {
    path: '/magma/v1/lte/:networkId/apns',
    resolver: (_, params) => [params[2], 'LTE network'],
  },
  {
    path: '/magma/v1/lte/:networkId/apns/(.*)',
    resolver: (_, params) => [params[2], 'LTE network'],
  },
  {
    path: '/magma/v1/lte/:networkId/gateway_pools/(.*)',
    resolver: (_, params) => [params[2], 'LTE network pools'],
  },
  {
    path: '/magma/v1/lte/:networkId/gateway_pools/(.*)',
    resolver: (_, params) => [params[2], 'LTE network pools'],
  },
  {
    path: '/magma/v1/lte/:networkId/cellular/(.*)',
    resolver: (_, params) => [params[2], 'LTE network'],
  },
  {
    path: '/magma/v1/lte/:networkId/dns/(.*)',
    resolver: (_, params) => [params[2], 'LTE network'],
  },
  {
    path: '/magma/v1/lte/:networkId/enodebs',
    resolver: (_, params) => [params[2], 'LTE managed eNodeBs'],
  },
  {
    path: '/magma/v1/lte/:networkId/enodebs/(.*)',
    resolver: (_, params) => [params[2], 'LTE network'],
  },
  {
    path: '/magma/v1/lte/:networkId/description',
    resolver: (_, params) => [params[2], 'LTE network'],
  },
  {
    path: '/magma/v1/lte/:networkId/gateways',
    resolver: req => [req.body.id, 'LTE gateway'],
  },
  {
    path: '/magma/v1/lte/:networkId/gateways/:objectId',
    resolver: (_, params) => [params[2], 'LTE gateway'],
  },
  {
    path: '/magma/v1/lte/:networkId/gateways/:objectId/(.*)',
    resolver: (_, params) => [params[2], 'LTE gateway'],
  },
  {
    path: '/magma/v1/lte/:networkId/subscribers/:objectId',
    type: 'subscriber',
  },
  {
    path: '/magma/v1/lte/:networkId/subscribers',
    resolver: (req: FBCNMSRequest) => [req.body.id, 'Subscriber'],
  },
  {
    path: '/magma/v1/lte/:networkId/policy_qos_profiles',
    resolver: (req: FBCNMSRequest) => [req.body.id, 'Policy qos profiles'],
  },
  {
    path: '/magma/v1/lte/:networkId/policy_qos_profiles/(.*)',
    resolver: (req: FBCNMSRequest) => [req.body.id, 'Policy qos profiles'],
  },
  {
    path: '/magma/v1/cwf/:networkId/gateways',
    resolver: req => [req.body.id, 'Carrier wifi gateway'],
  },
  {
    path: '/magma/v1/cwf/:networkId/gateways/:objectId',
    resolver: (_, params) => [params[2], 'Carrier wifi gateway'],
  },
  {
    path: '/magma/v1/feg/:networkId/gateways',
    resolver: req => [req.body.id, 'Federation gateway'],
  },
  {
    path: '/magma/v1/feg/:networkId/gateways/:objectId',
    resolver: (_, params) => [params[2], 'Federation gateway'],
  },
  {
    path: '/magma/v1/feg/:networkId/gateways/:objectId/federation',
    resolver: (_, params) => [params[2], 'Federation gateway config'],
  },
  {
    path: '/magma/v1/networks/:networkId/rules/policies',
    resolver: req => [req.body.id, 'Policy'],
  },
  {
    path: '/magma/v1/networks/:networkId/policies/rules/:objectId',
    resolver: (_, params) => [params[2], 'Policy'],
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
