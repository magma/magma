/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import fs from 'fs';

import {isEqual, sortBy} from 'lodash';

import MagmaV1API from '@fbcnms/platform-server/magma/index';
import {Organization} from '@fbcnms/sequelize-models';

import type {Datasource} from './GrafanaAPIType';
import type {FBCNMSRequest} from '@fbcnms/auth/access';
import type {GrafanaClient, GrafanaResponse} from './GrafanaAPI';
import type {OrganizationType} from '@fbcnms/sequelize-models/models/organization';
import type {UserType} from '@fbcnms/sequelize-models/models/user';
import type {tenant} from '../../fbcnms-magma-api';

export type GrafanaError = {
  response: GrafanaResponse<mixed>,
  message: string,
};

export const ORC8R_DATASOURCE_NAME = 'Orchestrator Datasource';

export async function HandleNewGrafanaUser(
  client: GrafanaClient,
  req: FBCNMSRequest,
): Promise<?GrafanaError> {
  const nmsOrg = await req.organization();
  const username = makeGrafanaUsername(req.user.id);

  // Check if user's organization already exists in Grafana
  let orgIDForUser: number;
  const orgResp = await client.getOrg(nmsOrg.name);
  switch (orgResp.status) {
    case 200:
      orgIDForUser = orgResp.data.id;
      break;
    case 404:
      const createOrgResp = await client.addOrg(nmsOrg.name);
      if (createOrgResp.status !== 200) {
        return {
          response: createOrgResp,
          message: 'Unexpected error creating organization',
        };
      }
      orgIDForUser = createOrgResp.data.orgId;
      break;
    default:
      return {
        response: orgResp,
        message: 'Unexpected error getting Grafana Organization',
      };
  }

  // Create new global user
  const createUserResp = await client.createUser({
    email: username,
    login: username,
    name: username,
    // Grafana uses AuthProxy so no password is required, but API still
    // requires the field. Login page is never accessible.
    password: '12345678',
  });
  if (createUserResp.status !== 200) {
    return {
      response: createUserResp,
      message: 'Unexpected error creating user',
    };
  }

  // Grafana will automatically create an org with name == user.email
  const newOrgResp = await client.getOrg(req.user.email);
  if (newOrgResp.status !== 200) {
    return {
      response: newOrgResp,
      message: 'Unexpected error getting organization',
    };
  }

  // Delete user's org
  const deleteOrgResp = await client.deleteOrg(newOrgResp.data.id);
  if (deleteOrgResp.status !== 200) {
    return {
      response: deleteOrgResp,
      message: 'Unexpected error deleting organization',
    };
  }

  // Add user to specified org
  const addToOrgResp = await client.addUserToOrg(orgIDForUser, {
    loginOrEmail: username,
    role: 'Editor',
  });
  if (addToOrgResp.status !== 200) {
    return {
      response: addToOrgResp,
      message: 'Unexpected error adding user to organization',
    };
  }
  return;
}

export async function HandleNewDatasource(
  client: GrafanaClient,
  req: FBCNMSRequest,
): Promise<?GrafanaError> {
  // Retrieve admin cert and key
  let cert, key;
  try {
    cert = fs.readFileSync(process.env.API_CERT_FILENAME || '');
    key = fs.readFileSync(process.env.API_PRIVATE_KEY_FILENAME || '');
  } catch (error) {
    return {
      response: {data: {}, status: 500},
      message: 'Could not retrieve cert for datasource ' + error,
    };
  }

  const nmsOrg = await req.organization();
  const grafanaOrgID = await getUserGrafanaOrgID(client, req.user);
  const nmsOrgID = nmsOrg.id;
  const apiHost = process.env.API_HOST;
  if (isNaN(grafanaOrgID) || apiHost === undefined || nmsOrgID === undefined) {
    return {
      response: {data: {}, status: 500},
      message: 'Could not get required information for datasource',
    };
  }

  // Create new datasource in Grafana
  const addDSResp = await client.createDatasource(
    makeDatasourceConfig({grafanaOrgID, nmsOrgID, apiHost, cert, key}),
    grafanaOrgID,
  );
  if (addDSResp.status !== 200) {
    return {
      response: {data: addDSResp.data, status: addDSResp.status},
      message: 'Could not create datasource ',
    };
  }
  return;
}

export async function HandleSyncOrganizations(): Promise<?GrafanaError> {
  const tenantMap = {};
  try {
    const orc8rTenants = await MagmaV1API.getTenants();
    orc8rTenants.forEach(tenant => {
      tenantMap[tenant.id] = tenant;
    });
  } catch (error) {
    return {
      response: {data: error.response.data, status: error.response.status},
      message: 'Error updating tenant',
    };
  }

  const nmsOrganizations = await Organization.findAll();
  for (const org of nmsOrganizations) {
    const orc8rTenant = tenantMap[org.id];
    try {
      // Update if tenant exists but is not equal to NMS Org
      if (orc8rTenant && !organizationsEqual(org, orc8rTenant)) {
        await MagmaV1API.putTenantsByTenantId({
          tenant: {id: org.id, name: org.name, networks: org.networkIDs},
          tenantId: org.id,
        });
      } else if (!orc8rTenant) {
        // Create new orc8r tenant if it didn't exist before
        await MagmaV1API.postTenants({
          tenant: {
            id: org.id,
            name: org.name,
            networks: org.networkIDs,
          },
        });
      }
    } catch (error) {
      return {
        response: {data: error.response.data, status: error.response.status},
        message: 'Error updating tenant',
      };
    }
  }
  return;
}

export function makeGrafanaUsername(userID: number): string {
  return `NMSUser_${userID}`;
}

async function getUserGrafanaOrgID(
  client: GrafanaClient,
  user: UserType,
): Promise<number> {
  if (user.organization === undefined) {
    return NaN;
  }
  const getOrgResp = await client.getOrg(user.organization);
  if (getOrgResp.data.id) {
    return getOrgResp.data.id;
  }
  return NaN;
}

function makeDatasourceConfig(params: {
  grafanaOrgID: number,
  nmsOrgID: number,
  apiHost: string,
  cert: Buffer,
  key: Buffer,
}): Datasource {
  return {
    name: ORC8R_DATASOURCE_NAME + '_' + params.grafanaOrgID,
    orgId: params.grafanaOrgID,
    type: 'orchestrator-grafana-datasource',
    access: 'proxy',
    url: `https://${params.apiHost}/magma/v1/tenants/${params.nmsOrgID}`,
    jsonData: {
      tlsAuth: true,
      tlsSkipVerify: true,
    },
    basicAuth: false,
    isDefault: true,
    readOnly: true,
    secureJsonData: {
      tlsClientCert: params.cert.toString(),
      tlsClientKey: params.key.toString(),
    },
  };
}

function organizationsEqual(
  nmsOrg: OrganizationType,
  orc8rTenant: tenant,
): boolean {
  return (
    nmsOrg.name == orc8rTenant.name &&
    isEqual(sortBy(nmsOrg.networkIDs), sortBy(orc8rTenant.networks))
  );
}
