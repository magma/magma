/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {FBCNMSRequest} from '@fbcnms/auth/access';
import type {GrafanaClient, GrafanaResponse} from './GrafanaAPI';

export type GrafanaError = {
  // status: number,
  // data: mixed,
  response: GrafanaResponse<mixed>,
  message: string,
};

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

export function makeGrafanaUsername(userID: number): string {
  return `NMSUser_${userID}`;
}
