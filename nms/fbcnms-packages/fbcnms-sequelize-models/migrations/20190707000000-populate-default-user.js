/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {DataTypes, QueryInterface} from 'sequelize';
import nullthrows from '@fbcnms/util/nullthrows';
import bcrypt from 'bcryptjs';

module.exports = {
  up: async (queryInterface: QueryInterface, _Sequelize: DataTypes) => {
    if (!process.env.FB_TEST_USER) {
      return Promise.resolve(null);
    }

    const testuser = nullthrows(process.env.FB_TEST_USER);

    const salt = await bcrypt.genSalt(10);
    const passwordHash = await bcrypt.hash(testuser, salt);

    return queryInterface.bulkInsert(
      'Users',
      [
        {
          email: testuser,
          organization: 'master',
          password: passwordHash,
          role: 3,
          createdAt: '2019-07-07 15:19:42',
          updatedAt: '2019-07-07 15:19:42',
          networkIDs: '[]',
        },
        {
          email: testuser,
          organization: 'fb-test',
          password: passwordHash,
          role: 3,
          createdAt: '2019-07-07 15:19:42',
          updatedAt: '2019-07-07 15:19:42',
          networkIDs: '["mpk_test"]',
        },
      ],
      {},
    );
  },

  down: (queryInterface: QueryInterface, _Sequelize: DataTypes) => {
    if (!process.env.FB_TEST_USER) {
      return Promise.resolve(null);
    }
    return queryInterface.bulkDelete('Users', null, {});
  },
};
