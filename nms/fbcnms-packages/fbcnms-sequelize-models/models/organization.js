/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Sequelize from 'sequelize';

import type {DataTypes, Model} from 'sequelize';
import type {AssociateProp} from './AssociateTypes.flow';

export type OrganizationRawType = {
  name: string,
  tabs?: Array<string>,
  customDomains?: Array<string>,
  networkIDs: Array<string>,
};

type OrganizationModel = Model<OrganizationRawType>;
export type StaticOrganizationModel = Class<OrganizationModel>;
export type OrganizationType = OrganizationModel & OrganizationRawType;

export default (
  sequelize: Sequelize,
  types: DataTypes,
): StaticOrganizationModel & AssociateProp => {
  const Organization = sequelize.define(
    'Organization',
    {
      name: types.STRING,
      tabs: {
        type: types.JSON,
        allowNull: false,
        defaultValue: [],
      },
      customDomains: {
        type: types.JSON,
        allowNull: false,
        defaultValue: [],
        get() {
          return this.getDataValue('customDomains') || [];
        },
      },
      networkIDs: {
        type: types.JSON,
        allowNull: false,
        defaultValue: [],
      },
    },
    {},
  );
  Organization.associate = function(_models) {
    // associations can be defined here
  };
  return Organization;
};
