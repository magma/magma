/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import axios from 'axios';

import type {
  AddOrgUserResponse,
  CreateDatasourceResponse,
  CreateOrgResponse,
  CreateUserResponse,
  DeleteOrgResponse,
  GetDatasourcesResponse,
  GetHealthResponse,
  GetOrgUsersResponse,
  GetUserResponse,
  OrgUser,
  Organization,
  PostDatasource,
  User,
} from './GrafanaAPIType';

type GrafanaPromise<T> = GrafanaPromise<T>;

export type GrafanaResponse<T> = {
  status: number,
  data: T,
};

export type GrafanaClient = {
  getUser: string => GrafanaPromise<GetUserResponse>,
  createUser: User => GrafanaPromise<CreateUserResponse>,

  getOrg: string => GrafanaPromise<Organization>,
  addOrg: string => GrafanaPromise<CreateOrgResponse>,
  deleteOrg: number => GrafanaPromise<DeleteOrgResponse>,
  addUserToOrg: (
    orgID: number,
    user: OrgUser,
  ) => GrafanaPromise<AddOrgUserResponse>,
  getUsersInOrg: (orgID: number) => GrafanaPromise<GetOrgUsersResponse>,

  createDatasource: (
    ds: PostDatasource,
    orgID: number,
  ) => GrafanaPromise<CreateDatasourceResponse>,
  updateDatasource: (
    dsID: number,
    orgID: number,
    ds: PostDatasource,
  ) => GrafanaPromise<CreateDatasourceResponse>,
  getDatasources: (orgID: number) => GrafanaPromise<GetDatasourcesResponse>,

  getHealth: () => GrafanaPromise<GetHealthResponse>,
};

type axiosRequest = {
  url: string,
  method: string,
  query?: {[string]: string},
  body?: mixed,
  headers?: {[string]: string},
};

async function request<T>(req: axiosRequest): GrafanaPromise<T> {
  try {
    const res = await axios(req);
    return {status: res.status, data: res.data};
  } catch (error) {
    return {status: error.response.status, data: error.response.data};
  }
}

const client = (
  apiURL: string,
  constHeaders: {[string]: string},
): GrafanaClient => ({
  async getUser(loginOrEmail: string): GrafanaPromise<GetUserResponse> {
    return request({
      url: apiURL + `/api/users/lookup`,
      params: {loginOrEmail: loginOrEmail},
      method: 'GET',
      headers: constHeaders,
    });
  },

  async createUser(user: User): GrafanaPromise<CreateUserResponse> {
    return request({
      url: apiURL + `/api/admin/users`,
      method: 'POST',
      data: user,
      headers: constHeaders,
    });
  },

  async getOrg(orgName: string): GrafanaPromise<Organization> {
    return request({
      url: apiURL + `/api/orgs/name/${orgName}`,
      method: 'GET',
      headers: constHeaders,
    });
  },

  async addOrg(orgName: string): GrafanaPromise<CreateOrgResponse> {
    return request({
      url: apiURL + '/api/orgs',
      method: 'POST',
      data: {name: orgName},
      headers: constHeaders,
    });
  },

  async deleteOrg(orgID: number): GrafanaPromise<DeleteOrgResponse> {
    return request({
      url: apiURL + `/api/orgs/${orgID}`,
      method: 'DELETE',
      headers: constHeaders,
    });
  },

  async addUserToOrg(
    orgID: number,
    user: OrgUser,
  ): GrafanaPromise<AddOrgUserResponse> {
    return request({
      url: apiURL + `/api/orgs/${orgID}/users`,
      method: 'POST',
      data: user,
      headers: constHeaders,
    });
  },

  async getUsersInOrg(orgID: number): GrafanaPromise<GetOrgUsersResponse> {
    return request({
      url: apiURL + `/api/orgs/${orgID}/users`,
      method: 'GET',
      headers: {...constHeaders, 'X-Grafana-Org-Id': orgID.toString()},
    });
  },

  async createDatasource(
    ds: PostDatasource,
    orgId: number,
  ): GrafanaPromise<CreateDatasourceResponse> {
    return request({
      url: apiURL + `/api/datasources`,
      method: 'POST',
      data: ds,
      headers: {...constHeaders, 'X-Grafana-Org-Id': orgId.toString()},
    });
  },

  async updateDatasource(
    dsID: number,
    orgID: number,
    ds: PostDatasource,
  ): GrafanaPromise<CreateDatasourceResponse> {
    return request({
      url: apiURL + `/api/datasources/${dsID}`,
      method: 'PUT',
      data: ds,
      headers: {...constHeaders, 'X-Grafana-Org-Id': orgID.toString()},
    });
  },

  async getDatasources(orgID: number): GrafanaPromise<GetDatasourcesResponse> {
    return request({
      url: apiURL + `/api/datasources`,
      method: 'GET',
      headers: {...constHeaders, 'X-Grafana-Org-Id': orgID.toString()},
    });
  },

  async getHealth(): GrafanaPromise<GetHealthResponse> {
    return request({
      url: apiURL + `/api/health`,
      method: 'GET',
      headers: constHeaders,
    });
  },
});

export default client;
