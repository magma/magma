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
  jobTitle?: string,
|};

const generateString = length =>
  Math.random()
    .toString(36)
    .replace(/[^a-z]+/g, '')
    .substr(0, length || 5);
const randomNaturalNumber = (from, to) => {
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
  jobTitle: `${generateString(1).toUpperCase()}${generateString(
    randomNaturalNumber(3, 6),
  )} ${generateString(1).toUpperCase()}${generateString(
    randomNaturalNumber(4, 10),
  )}`,
}));
