/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export const USER_ROLES = {
  User: 'User',
  Admin: 'Admin',
  Owner: 'Owner',
};

export const USER_STATUSES = {
  Active: 'Active',
  Deactivated: 'Deactivated',
  Deleted: 'Deleted',
};

export const EMPLOYMENT_TYPES = {
  FullTime: 'Full Time',
  Contructor: 'Contructor',
};

export type UserStatus = $Keys<typeof USER_STATUSES>;
export type UserRole = $Keys<typeof USER_ROLES>;
export type EmploymentType = $Keys<typeof EMPLOYMENT_TYPES>;

export type User = {|
  authId: string,
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
const generateRole = () => itemFromArray(Object.keys(USER_ROLES));
const generateStatus = () => itemFromArray(Object.keys(USER_STATUSES));
const generateEmployment = () => itemFromArray(Object.keys(EMPLOYMENT_TYPES));

export const TEMP_USERS: Array<User> = [...new Array(50)].map(_ => ({
  authId: `${generateString(randomNaturalNumber(8, 20))}@${generateString(
    randomNaturalNumber(2, 4),
  )}.${generateString(randomNaturalNumber(2, 3))}`,
  firstName: `${generateString(1).toUpperCase()}${generateString(
    randomNaturalNumber(4, 8),
  )}`,
  lastName: `${generateString(1).toUpperCase()}${generateString(
    randomNaturalNumber(4, 10),
  )}`,
  // eslint-disable-next-line no-warning-comments
  // $FlowFixMe: it is temporary
  role: generateRole(),
  // eslint-disable-next-line no-warning-comments
  // $FlowFixMe: it is temporary
  status: generateStatus(),
  // eslint-disable-next-line no-warning-comments
  // $FlowFixMe: it is temporary
  employmentType: generateEmployment(),
  employeeID: `${randomNaturalNumber(1000, 9999)}`,
  jobTitle: `${generateString(1).toUpperCase()}${generateString(
    randomNaturalNumber(3, 6),
  )} ${generateString(1).toUpperCase()}${generateString(
    randomNaturalNumber(4, 10),
  )}`,
}));

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
