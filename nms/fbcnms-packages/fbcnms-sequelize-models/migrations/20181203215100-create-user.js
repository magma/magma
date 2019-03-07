/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */
'use strict';

module.exports = {
  up: (queryInterface, Sequelize) => {
    return queryInterface.addColumn('Users', 'organization', {
      allowNull: false,
      defaultValue: 'fb-test',
      type: Sequelize.STRING,
    });
  },

  down: (queryInterface, _Sequelize) => {
    return queryInterface.removeColumn('Users', 'organization');
  },
};
