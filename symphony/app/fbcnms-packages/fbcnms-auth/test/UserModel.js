/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
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
    readOnly: false,
  },
  {
    id: '2',
    email: 'noorg@123.com',
    organization: 'nottakenintoconsideration',
    role: AccessRoles.USER,
    password: bcrypt.hashSync('password1234', bcrypt.genSaltSync(1)),
    readOnly: false,
  },
  {
    id: '3',
    email: 'superuser@123.com',
    organization: 'validorg',
    role: AccessRoles.SUPERUSER,
    password: bcrypt.hashSync('password1234', bcrypt.genSaltSync(1)),
    readOnly: false,
  },
];

export const USERS_EXPECTED = [
  {
    networkIDs: [],
    id: 1,
    email: 'valid@123.com',
    organization: 'validorg',
    role: 0,
    tabs: [],
    readOnly: false,
  },
  {
    networkIDs: [],
    id: 2,
    email: 'noorg@123.com',
    organization: 'nottakenintoconsideration',
    role: 0,
    tabs: [],
    readOnly: false,
  },
  {
    networkIDs: [],
    id: 3,
    email: 'superuser@123.com',
    organization: 'validorg',
    role: 3,
    tabs: [],
    readOnly: false,
  },
];
