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

export type CreateUserResponse = {
  id: number,
  message: string,
};
