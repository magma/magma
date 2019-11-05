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

import type {AssociateProp} from './AssociateTypes.flow';
import type {DataTypes, Model} from 'sequelize';

export type OrganizationRawType = {
  name: string,
  tabs?: Array<string>,
  customDomains?: Array<string>,
  networkIDs: Array<string>,
  csvCharset: string,
  ssoCert: string,
  ssoEntrypoint: string,
  ssoIssuer: string,
};

type OrganizationGetters = {
  isMasterOrg: boolean,
};

type OrganizationModel = Model<
  OrganizationRawType & OrganizationGetters,
  OrganizationRawType,
>;
export type StaticOrganizationModel = Class<OrganizationModel>;
export type OrganizationType = OrganizationModel &
  OrganizationRawType &
  OrganizationGetters;

const MASTER_ORG = 'master';

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
      csvCharset: types.STRING,
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
      ssoCert: {
        type: types.TEXT,
        allowNull: false,
        defaultValue: '',
      },
      ssoEntrypoint: {
        type: types.STRING,
        allowNull: false,
        defaultValue: '',
      },
      ssoIssuer: {
        type: types.STRING,
        allowNull: false,
        defaultValue: '',
      },
    },
    {
      getterMethods: {
        isMasterOrg() {
          return this.name === MASTER_ORG;
        },
      },
    },
  );
  Organization.associate = function(_models) {
    // associations can be defined here
  };
  return Organization;
};
