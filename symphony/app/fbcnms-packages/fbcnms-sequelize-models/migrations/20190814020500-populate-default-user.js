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
import type {DataTypes, QueryInterface} from 'sequelize';

module.exports = {
  up: async (queryInterface: QueryInterface, _: DataTypes) => {
    const email = process.env.FB_TEST_USER;
    if (!email) {
      return Promise.resolve(null);
    }

    const salt = await bcrypt.genSalt(10);
    const passwordHash = await bcrypt.hash(email, salt);

    return queryInterface.bulkInsert(
      'Users',
      [
        {
          email: email,
          organization: 'master',
          password: passwordHash,
          role: 3,
          createdAt: new Date(),
          updatedAt: new Date(),
          networkIDs: '[]',
          tabs: '[]',
        },
        {
          email: email,
          organization: 'fb-test',
          password: passwordHash,
          role: 3,
          createdAt: new Date(),
          updatedAt: new Date(),
          networkIDs: '["mpk_test"]',
          tabs: '["inventory", "nms", "workorders", "automation"]',
        },
      ],
      {},
    );
  },

  down: (queryInterface: QueryInterface, _: DataTypes) => {
    if (process.env.FB_TEST_USER) {
      return queryInterface.bulkDelete('Users', null, {});
    }
    return Promise.resolve(null);
  },
};
