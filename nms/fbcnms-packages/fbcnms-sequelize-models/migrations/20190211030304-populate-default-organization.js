/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

module.exports = {
  up: (queryInterface, _Sequelize) => {
    return queryInterface.bulkInsert(
      'Organizations',
      [
        {
          id: '1',
          customDomains: '[]',
          name: 'fb-test',
          networkIDs: '["mpk_test"]',
          createdAt: '2019-02-11 20:05:05',
          updatedAt: '2019-02-11 20:05:05',
        },
      ],
      {},
    );
  },

  down: (queryInterface, _Sequelize) => {
    return queryInterface.bulkDelete('Organizations', null, {});
  },
};
