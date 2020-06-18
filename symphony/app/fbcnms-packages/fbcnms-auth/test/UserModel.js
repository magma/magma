/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import bcrypt from 'bcryptjs';
import {AccessRoles} from '../roles';

export const USERS = [
  {
    id: '1',
    email: 'valid@123.com',
    organization: 'validorg',
    role: AccessRoles.USER,
    password: bcrypt.hashSync('password1234', bcrypt.genSaltSync(1)),
    networkIDs: ['network1'],
  },
  {
    id: '2',
    email: 'noorg@123.com',
    organization: 'nottakenintoconsideration',
    role: AccessRoles.USER,
    password: bcrypt.hashSync('password1234', bcrypt.genSaltSync(1)),
  },
  {
    id: '3',
    email: 'superuser@123.com',
    organization: 'validorg',
    role: AccessRoles.SUPERUSER,
    password: bcrypt.hashSync('password1234', bcrypt.genSaltSync(1)),
  },
  {
    id: '4',
    email: 'readonlyuser@123.com',
    organization: 'readonlyorg',
    role: AccessRoles.READ_ONLY_USER,
    password: bcrypt.hashSync('password1234', bcrypt.genSaltSync(1)),
  },
];

export const USERS_EXPECTED = [
  {
    networkIDs: ['network1'],
    id: 1,
    email: 'valid@123.com',
    organization: 'validorg',
    role: AccessRoles.USER,
    tabs: [],
  },
  {
    networkIDs: [],
    id: 2,
    email: 'noorg@123.com',
    organization: 'nottakenintoconsideration',
    role: AccessRoles.USER,
    tabs: [],
  },
  {
    networkIDs: [],
    id: 3,
    email: 'superuser@123.com',
    organization: 'validorg',
    role: AccessRoles.SUPERUSER,
    tabs: [],
  },
  {
    networkIDs: [],
    id: 4,
    email: 'readonlyuser@123.com',
    organization: 'readonlyorg',
    role: AccessRoles.READ_ONLY_USER,
    tabs: [],
  },
];
