/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {
  UserRole,
  UserStatus,
} from './__generated__/UsersView_UsersQuery.graphql';

import fbt from 'fbt';

type KeyValueEnum<TValues> = {
  [key: TValues]: {
    key: TValues,
    value: string,
  },
};

export const USER_ROLES: KeyValueEnum<UserRole> = {
  USER: {
    key: 'USER',
    value: `${fbt('User', '')}`,
  },
  ADMIN: {
    key: 'ADMIN',
    value: `${fbt('Admin', '')}`,
  },
  OWNER: {
    key: 'OWNER',
    value: `${fbt('Owner', '')}`,
  },
};
export const USER_STATUSES: KeyValueEnum<UserStatus> = {
  ACTIVE: {
    key: 'ACTIVE',
    value: `${fbt('Active', '')}`,
  },
  DEACTIVATED: {
    key: 'DEACTIVATED',
    value: `${fbt('Deactivated', '')}`,
  },
};
export const EMPLOYMENT_TYPES: KeyValueEnum<string> = {
  FULL_TIME: {
    key: 'FULL_TIME',
    value: `${fbt('Full Time', '')}`,
  },
  CONTRUCTOR: {
    key: 'CONTRACTOR',
    value: `${fbt('Contractor', '')}`,
  },
};

export type EmploymentType = $Keys<typeof EMPLOYMENT_TYPES>;

export type User = {|
  id: string,
  authID: string,
  firstName: string,
  lastName: string,
  role: UserRole,
  status: UserStatus,
  photoId?: string,
  employmentType?: EmploymentType,
  employeeID?: string,
  jobTitle?: string,
  phoneNumber?: string,
|};

const generateString = length =>
  Math.random()
    .toString(36)
    .replace(/[^a-z]+/g, '')
    .substr(0, length || 5);
const randomNaturalNumber = (from: number, to: number) => {
  from = from ?? 1;
  to = to > from ? to : from + 1;
  const range = to - from;
  return Math.round(Math.random() * range) + from;
};
const itemFromArray = arr => arr[randomNaturalNumber(0, arr.length - 1)];

export const GROUP_STATUSES = {
  Active: 'Active',
  Inactive: 'Inactive',
};

export type GroupStatus = $Keys<typeof GROUP_STATUSES>;

export type UserPermissionsGroup = {|
  id: string,
  name: string,
  description: string,
  status: GroupStatus,
  members: Array<string>,
|};

const generateGroupStatus = () => itemFromArray(Object.keys(GROUP_STATUSES));

export const TEMP_GROUPS: Array<UserPermissionsGroup> = [...new Array(5)].map(
  _ => ({
    id: `${generateString(10)}`,
    name: `${generateString(1).toUpperCase()}${generateString(
      randomNaturalNumber(4, 8),
    )} ${generateString(1).toUpperCase()}${generateString(
      randomNaturalNumber(4, 10),
    )}`,
    description: `${generateString(1).toUpperCase()}${generateString(
      randomNaturalNumber(4, 8),
    )} ${[...new Array(randomNaturalNumber(1, 7))]
      .map(_ => generateString(randomNaturalNumber(4, 10)))
      .join(' ')}`,
    // eslint-disable-next-line no-warning-comments
    // $FlowFixMe: it is temporary
    status: generateGroupStatus(),
    members: [],
  }),
);
