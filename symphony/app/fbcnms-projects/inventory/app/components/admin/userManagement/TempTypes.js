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

export type UserStatus = $Keys<typeof USER_STATUSES>;
export type UserRole = $Keys<typeof USER_ROLES>;

export type User = {
  authId: string,
  firstName: string,
  lastName: string,
  role: UserRole,
  status: UserStatus,
  photoId?: string,
  employmentTypeId?: string,
};
