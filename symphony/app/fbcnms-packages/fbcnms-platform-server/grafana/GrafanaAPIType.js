/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export type User = {
  name: string,
  email: string,
  login: string,
  password: string,
};

export type GetUserResponse = {
  id: string,
  email: string,
  name: string,
  login: string,
  theme: string,
  orgId: number,
  isGrafanaAdmin: boolean,
  isDisabled: boolean,
  isExternal: boolean,
  authLabels: Array<string>,
  updatedAt: string,
  createdAt: string,
};

export type OrgUser = {
  loginOrEmail: string,
  role: 'Admin' | 'Editor' | 'Viewer',
};

export type Organization = {
  id: number,
  name: string,
};

export type CreateOrgResponse = {
  orgId: number,
  message: string,
};

export type DeleteOrgResponse = {
  message: string,
};

export type AddOrgUserResponse = {
  message: string,
};

export type GetOrgUsersResponse = Array<OrgUserResponse>;

type OrgUserResponse = {
  orgId: number,
  userId: number,
  email: string,
  avatarUrl: string,
  login: string,
  role: string,
  lastSeenAt: string,
  lastSeenAtAge: string,
};

export type CreateUserResponse = {
  id: number,
  message: string,
};

export type PostDatasource = {
  orgId: number,
  name: string,
  type: string,
  typeLogoUrl?: string,
  access: string,
  url: string,
  password?: string,
  user?: string,
  database?: string,
  basicAuth: boolean,
  basicAuthUser?: string,
  basicAuthPassword?: string,
  withCredentials?: boolean,
  isDefault: boolean,
  jsonData: {tlsAuth: boolean, tlsSkipVerify: boolean},
  secureJsonData: {tlsClientCert: string, tlsClientKey: string},
  version?: number,
  readOnly: boolean,
};

export type Datasource = {
  id: number,
  orgId: number,
  name: string,
  type: string,
  typeLogoUrl?: string,
  access: string,
  url: string,
  password: string,
  user: string,
  database: string,
  basicAuth: boolean,
  basicAuthUser: string,
  withCredentials?: boolean,
  isDefault: boolean,
  jsonData: {tlsAuth?: boolean, tlsSkipVerify?: boolean},
  secureJsonData: {
    tlsClientCert?: string,
    tlsClientKey?: string,
    basicAuthPassword?: string,
  },
  version: number,
  readOnly: boolean,
};

export type CreateDatasourceResponse = {
  datasource: Datasource,
  id: number,
  message: string,
  name: string,
};

export type GetDatasourcesResponse = Array<Datasource>;

export type GetHealthResponse = {
  commit: string,
  database: string,
  version: string,
};
