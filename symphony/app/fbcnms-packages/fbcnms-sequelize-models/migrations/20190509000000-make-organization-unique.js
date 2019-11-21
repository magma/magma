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
const CONSTRAINT_NAME = 'unique_organization_name';

module.exports = {
  up: (queryInterface: QueryInterface, _types: DataTypes) => {
    return queryInterface.addConstraint('Organizations', ['name'], {
      type: 'unique',
      name: CONSTRAINT_NAME,
    });
  },

  down: (queryInterface: QueryInterface, _types: DataTypes) => {
    return queryInterface.removeConstraint('Organizations', CONSTRAINT_NAME);
  },
};
