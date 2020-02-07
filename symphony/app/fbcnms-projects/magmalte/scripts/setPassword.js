/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
'use strict';

import bcrypt from 'bcryptjs';
import {AccessRoles} from '@fbcnms/auth/roles';
import {User} from '@fbcnms/sequelize-models';

const SALT_GEN_ROUNDS = 10;

type UserObject = {
  email: string,
  password: string,
  superUser: boolean,
};

async function updateUser(user: User, userObject: UserObject) {
  const {password, superUser} = userObject;
  const salt = await bcrypt.genSalt(SALT_GEN_ROUNDS);
  const passwordHash = await bcrypt.hash(password, salt);
  await user.update({
    password: passwordHash,
    role: superUser ? AccessRoles.SUPERUSER : AccessRoles.USER,
  });
}

async function createUser(userObject: UserObject) {
  const {email, password, superUser} = userObject;
  const salt = await bcrypt.genSalt(SALT_GEN_ROUNDS);
  const passwordHash = await bcrypt.hash(password, salt);
  await User.create({
    email: email.toLowerCase(),
    password: passwordHash,
    role: superUser ? AccessRoles.SUPERUSER : AccessRoles.USER,
  });
}

async function createOrUpdateUser(userObject: UserObject) {
  const user = await User.findOne({
    where: {
      email: userObject.email.toLowerCase(),
    },
  });
  if (!user) {
    await createUser(userObject);
  } else {
    await updateUser(user, userObject);
  }
}

function main() {
  const args = process.argv.slice(2);
  if (args.length != 2) {
    console.log('Usage: setPassword.js <email> <password>');
    process.exit(1);
  }
  const userObject = {
    email: args[0],
    password: args[1],
    superUser: true,
  };
  createOrUpdateUser(userObject)
    .then(_res => {
      console.log('Success');
      process.exit();
    })
    .catch(err => {
      console.error(err);
      process.exit(1);
    });
}

main();
