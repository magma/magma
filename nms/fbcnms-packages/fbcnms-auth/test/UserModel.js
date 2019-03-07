/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {AccessRoles} from '../roles';
import bcrypt from 'bcryptjs';
import {UserVerificationTypes} from '../types';

export const USERS = [
  {
    id: '1',
    email: 'valid@123.com',
    organization: 'validorg',
    role: AccessRoles.USER,
    password: bcrypt.hashSync('password1234', bcrypt.genSaltSync(1)),
    verificationType: UserVerificationTypes.PASSWORD,
  },
  {
    id: '2',
    email: 'noorg@123.com',
    organization: 'nottakenintoconsideration',
    role: AccessRoles.USER,
    password: bcrypt.hashSync('password1234', bcrypt.genSaltSync(1)),
    verificationType: UserVerificationTypes.PASSWORD,
  },
  {
    id: '3',
    email: 'superuser@123.com',
    organization: 'validorg',
    role: AccessRoles.SUPERUSER,
    password: bcrypt.hashSync('password1234', bcrypt.genSaltSync(1)),
    verificationType: UserVerificationTypes.PASSWORD,
  },
];

export const USERS_EXPECTED = [
  {
    isSuperUser: false,
    networkIDs: [],
    id: 1,
    email: 'valid@123.com',
    organization: 'validorg',
    role: 0,
  },
  {
    isSuperUser: false,
    networkIDs: [],
    id: 2,
    email: 'noorg@123.com',
    organization: 'nottakenintoconsideration',
    role: 0,
  },
  {
    isSuperUser: true,
    networkIDs: [],
    id: 3,
    email: 'superuser@123.com',
    organization: 'validorg',
    role: 3,
  },
];
