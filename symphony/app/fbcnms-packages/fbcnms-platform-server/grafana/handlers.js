/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import {isEqual, sortBy} from 'lodash';

import MagmaV1API from '@fbcnms/platform-server/magma/index';
import {Organization} from '@fbcnms/sequelize-models';
import {apiCredentials} from '../config';

import type {Datasource, PostDatasource} from './GrafanaAPIType';
import type {FBCNMSRequest} from '@fbcnms/auth/access';
import type {GrafanaClient} from './GrafanaAPI';
import type {OrganizationType} from '@fbcnms/sequelize-models/models/organization';
import type {UserType} from '@fbcnms/sequelize-models/models/user';
import type {tenant} from '../../fbcnms-magma-api';

export type Task = {name: string, status: number, message: string};

export const ORC8R_DATASOURCE_NAME = 'Orchestrator Datasource';

export async function syncGrafanaUser(
  client: GrafanaClient,
  req: FBCNMSRequest,
): Promise<{completedTasks: Array<Task>, errorTask?: Task}> {
  const completedTasks: Array<Task> = [];

  const nmsOrg = await req.organization();
  const username = makeGrafanaUsername(req.user.id);

  // Check if user's organization already exists in Grafana
  let orgIDForUser: number;
  const orgResp = await client.getOrg(nmsOrg.name);
  switch (orgResp.status) {
    case 200:
      orgIDForUser = orgResp.data.id;
      completedTasks.push({
        name: 'Retrieve Grafana organization',
        status: orgResp.status,
        message: orgResp.data,
      });
      break;
    case 404:
      const createOrgResp = await client.addOrg(nmsOrg.name);
      if (createOrgResp.status !== 200) {
        return {
          completedTasks,
          errorTask: {
            name: 'Add grafana organization',
            status: createOrgResp.status,
            message: createOrgResp.data,
          },
        };
      }
      orgIDForUser = createOrgResp.data.orgId;
      completedTasks.push({
        name: 'Add grafana organization',
        status: 200,
        message: createOrgResp.data,
      });
      break;
    default:
      return {
        completedTasks,
        errorTask: {
          name: 'Retrieve grafana organization',
          status: orgResp.status,
          message: orgResp.data,
        },
      };
  }

  const getUserResp = await client.getUser(username);
  if (getUserResp.status !== 200) {
    const createUserResp = await createNewUser(client, username);
    completedTasks.push(...createUserResp.completedTasks);
  }

  let userInCorrectOrg = false;
  try {
    userInCorrectOrg = await checkIfUserInOrg(client, username, orgIDForUser);
    completedTasks.push({
      name: 'Check if user in correct organization',
      status: 200,
      message: userInCorrectOrg
        ? 'User in organization'
        : 'User not in organization',
    });
  } catch (error) {
    return {
      completedTasks,
      errorTask: {
        name: 'Check if user in correct organization',
        status: error.status,
        message: error.data,
      },
    };
  }
  if (userInCorrectOrg) {
    return {completedTasks};
  }

  // Add user to specified org
  const addToOrgResp = await client.addUserToOrg(orgIDForUser, {
    loginOrEmail: username,
    role: 'Editor',
  });
  if (addToOrgResp.status !== 200) {
    return {
      completedTasks,
      errorTask: {
        name: 'Add User to Organization',
        status: addToOrgResp.status,
        message: addToOrgResp.data,
      },
    };
  }
  completedTasks.push({
    name: 'Add User to Organization',
    status: addToOrgResp.status,
    message: addToOrgResp.data,
  });
  return {completedTasks};
}

async function createNewUser(
  client: GrafanaClient,
  username: string,
): Promise<{completedTasks: Array<Task>, errorTask?: Task}> {
  const completedTasks: Array<Task> = [];
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
      completedTasks,
      errorTask: {
        name: 'Create Grafana User',
        status: createUserResp.status,
        message: createUserResp.data,
      },
    };
  }
  completedTasks.push({
    name: 'Create Grafana User',
    status: createUserResp.status,
    message: createUserResp.data,
  });

  // Grafana will automatically create an org with name == username
  const newOrgResp = await client.getOrg(username);
  if (newOrgResp.status !== 200) {
    return {
      completedTasks,
      errorTask: {
        name: "Retrieve User's Grafana Organization",
        status: newOrgResp.status,
        message: newOrgResp.data,
      },
    };
  }
  completedTasks.push({
    name: "Retrieve User's Grafana Organization",
    status: newOrgResp.status,
    message: newOrgResp.data,
  });

  // Delete user's org
  const deleteOrgResp = await client.deleteOrg(newOrgResp.data.id);
  if (deleteOrgResp.status !== 200) {
    return {
      completedTasks,
      errorTask: {
        name: "Delete User's Grafana Organization",
        status: deleteOrgResp.status,
        message: deleteOrgResp.data,
      },
    };
  }
  completedTasks.push({
    name: "Delete User's Grafana Organization",
    status: deleteOrgResp.status,
    message: deleteOrgResp.data,
  });
  return {completedTasks};
}

export async function syncDatasource(
  client: GrafanaClient,
  req: FBCNMSRequest,
): Promise<{completedTasks: Array<Task>, errorTask?: Task}> {
  const completedTasks: Array<Task> = [];
  // Retrieve admin cert and key
  const tryCreds = apiCredentials();
  if (tryCreds.cert === undefined || tryCreds.key === undefined) {
    return {
      completedTasks,
      errorTask: {
        name: 'Retrieve certs for datasource',
        status: 500,
        message: 'Could not retrieve certs for datasource',
      },
    };
  }
  const creds = {cert: tryCreds.cert, key: tryCreds.key};
  completedTasks.push({
    name: 'Retrieve certs for datasource',
    status: 200,
    message: 'success',
  });

  const nmsOrg = await req.organization();
  const grafanaOrgID = await getUserGrafanaOrgID(client, req.user);
  const nmsOrgID = nmsOrg.id;
  const apiHost = process.env.API_HOST;
  if (isNaN(grafanaOrgID) || apiHost === undefined || nmsOrgID === undefined) {
    return {
      completedTasks,
      errorTask: {
        name: 'Get required information for datasource',
        status: 500,
        message: `GrafanaOrgID: ${grafanaOrgID},
         apiHost: ${apiHost || ''},
         nmsOrgID: ${nmsOrgID}`,
      },
    };
  }
  completedTasks.push({
    name: 'Get required information for datasource',
    status: 200,
    message: 'success',
  });

  const getDSResp = await client.getDatasources(grafanaOrgID);
  if (getDSResp.status !== 200) {
    return {
      completedTasks,
      errorTask: {
        name: 'Retrieve datasources',
        status: getDSResp.status,
        message: getDSResp.data,
      },
    };
  }
  const newDSParams: DatasourceParams = {
    grafanaOrgID,
    nmsOrgID,
    apiHost,
    cert: creds.cert,
    key: creds.key,
  };

  const ds = getOrc8rDatasource(getDSResp.data);
  completedTasks.push({
    name: 'Checked datasource exists in org',
    status: 200,
    message: ds ? 'Datsource exists' : 'Datasource does not exist',
  });
  if (ds) {
    // Update Datasource if parameters have changed
    const updateDSResp = await updateDatasourceIfChanged({
      oldDS: ds,
      newDSParams,
      client,
    });
    completedTasks.push(...updateDSResp.completedTasks);
    if (updateDSResp.errorTask) {
      return {
        completedTasks,
        errorTask: updateDSResp.errorTask,
      };
    }
    return {completedTasks};
  }

  // Create new datasource in Grafana
  const addDSResp = await client.createDatasource(
    makeDatasourceConfig(newDSParams),
    grafanaOrgID,
  );
  if (addDSResp.status !== 200) {
    return {
      completedTasks,
      errorTask: {
        name: 'Create datasource',
        status: addDSResp.status,
        message: addDSResp.data,
      },
    };
  }
  completedTasks.push({
    name: 'Create datasource',
    status: addDSResp.status,
    message: addDSResp.data,
  });
  return {completedTasks};
}

type updateDatasourceArgs = {
  oldDS: Datasource,
  newDSParams: DatasourceParams,
  client: GrafanaClient,
};

async function updateDatasourceIfChanged({
  oldDS,
  newDSParams,
  client,
}: updateDatasourceArgs): Promise<{
  completedTasks: Array<Task>,
  errorTask?: Task,
}> {
  const completedTasks: Array<Task> = [];
  // Make sure API Endpoint matches
  if (oldDS.url === makeAPIUrl(newDSParams.apiHost, newDSParams.nmsOrgID)) {
    return {completedTasks};
  }
  const updatedDS = makeDatasourceConfig(newDSParams);
  const updateDSResp = await client.updateDatasource(
    oldDS.id,
    newDSParams.grafanaOrgID,
    updatedDS,
  );
  if (updateDSResp.status !== 200) {
    return {
      completedTasks,
      errorTask: {
        name: 'Update datasource',
        status: updateDSResp.status,
        message: updateDSResp.data,
      },
    };
  }
  completedTasks.push({
    name: 'Update datasource',
    status: updateDSResp.status,
    message: updateDSResp.data,
  });
  return {completedTasks};
}

export async function syncTenants(): Promise<{
  completedTasks: Array<Task>,
  errorTask?: Task,
}> {
  const completedTasks: Array<Task> = [];
  const tenantMap = {};
  try {
    const orc8rTenants = await MagmaV1API.getTenants();
    orc8rTenants.forEach(tenant => {
      tenantMap[tenant.id] = tenant;
    });
    completedTasks.push({
      name: 'Retrieve Magma Tenants',
      status: 200,
      message: 'success',
    });
  } catch (error) {
    return {
      completedTasks,
      errorTask: {
        name: 'Retrieve Magma Tenants',
        status: error.response.status,
        message: error.response.data,
      },
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
        completedTasks.push({
          name: 'Update Magma Tenant',
          status: 200,
          message: 'success',
        });
      } else if (!orc8rTenant) {
        // Create new orc8r tenant if it didn't exist before
        await MagmaV1API.postTenants({
          tenant: {id: org.id, name: org.name, networks: org.networkIDs},
        });
        completedTasks.push({
          name: 'Create Magma Tenant',
          status: 200,
          message: 'success',
        });
      }
    } catch (error) {
      return {
        completedTasks,
        errorTask: {
          name: 'Update Magma Tenants',
          status: error.response.status,
          message: error.response.data,
        },
      };
    }
  }
  return {completedTasks};
}

export function makeGrafanaUsername(userID: number): string {
  return `NMSUser_${userID}`;
}

async function checkIfUserInOrg(
  client: GrafanaClient,
  username: string,
  orgID: number,
): Promise<boolean> {
  const getUsersResp = await client.getUsersInOrg(orgID);
  return getUsersResp.data.some(user => user.login === username);
}

function getOrc8rDatasource(datasources: Array<Datasource>): ?Datasource {
  return datasources.find(ds => ds.name.startsWith(ORC8R_DATASOURCE_NAME));
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

type DatasourceParams = {
  grafanaOrgID: number,
  nmsOrgID: number,
  apiHost: string,
  cert: string | Buffer,
  key: string | Buffer,
};

function makeDatasourceConfig(params: DatasourceParams): PostDatasource {
  return {
    name: ORC8R_DATASOURCE_NAME + '_' + params.grafanaOrgID,
    orgId: params.grafanaOrgID,
    type: 'prometheus',
    access: 'proxy',
    url: makeAPIUrl(params.apiHost, params.nmsOrgID),
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

function makeAPIUrl(apiHost: string, nmsOrgID: number): string {
  return `https://${apiHost}/magma/v1/tenants/${nmsOrgID}/metrics`;
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
