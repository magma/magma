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

module.exports = {
  up: (queryInterface: QueryInterface, _Sequelize: DataTypes) => {
    return queryInterface.bulkInsert(
      'Organizations',
      [
        {
          id: '1',
          customDomains: '[]',
          name: 'master',
          networkIDs: '[]',
          tabs: '["admin"]',
          createdAt: '2019-02-11 20:05:05',
          updatedAt: '2019-02-11 20:05:05',
        },
        {
          id: '2',
          customDomains: '[]',
          name: 'fb-test',
          networkIDs: '["mpk_test"]',
          tabs: '["nms", "inventory", "workorders", "automation"]',
          createdAt: '2019-02-11 20:05:05',
          updatedAt: '2019-02-11 20:05:05',
        },
      ],
      {},
    );
  },

  down: (queryInterface: QueryInterface, _Sequelize: DataTypes) => {
    return queryInterface.bulkDelete('Organizations', null, {});
  },
};
