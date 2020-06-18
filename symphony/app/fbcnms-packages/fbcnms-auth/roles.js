/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export type AccessRoleLevel = 0 | 1 | 3;

export const AccessRoles = {
  USER: 0,
  READ_ONLY_USER: 1,
  SUPERUSER: 3,
};

export function accessRoleToString(role: number): string {
  if (role === AccessRoles.SUPERUSER) {
    return 'OWNER';
  }
  return 'USER';
}
