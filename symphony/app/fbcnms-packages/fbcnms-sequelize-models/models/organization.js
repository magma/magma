/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Sequelize from 'sequelize';

import type {AssociateProp} from './AssociateTypes.flow';
import type {DataTypes, Model} from 'sequelize';
import type {SSOSelectedType} from '@fbcnms/types/auth';
import type {Tab} from '@fbcnms/types/tabs';

type OrganizationInitAttributes = {
  id?: number,
  name: string,
  tabs?: Array<Tab>,
  customDomains?: Array<string>,
  networkIDs: Array<string>,
  csvCharset: string,
  ssoSelectedType?: SSOSelectedType,
  ssoCert: string,
  ssoEntrypoint: string,
  ssoIssuer: string,
  ssoOidcClientID?: string,
  ssoOidcClientSecret?: string,
  ssoOidcConfigurationURL?: string,
};

export type OrganizationPlainAttributes = {
  id: number,
  name: string,
  tabs: Array<Tab>,
  customDomains: Array<string>,
  networkIDs: Array<string>,
  csvCharset: string,
  ssoSelectedType: SSOSelectedType,
  ssoCert: string,
  ssoEntrypoint: string,
  ssoIssuer: string,
  ssoOidcClientID: string,
  ssoOidcClientSecret: string,
  ssoOidcConfigurationURL: string,
};

type OrganizationGetters = {
  isMasterOrg: boolean,
};

type OrganizationAttributes = OrganizationPlainAttributes & OrganizationGetters;

export type OrganizationModel = Model<
  OrganizationAttributes,
  OrganizationInitAttributes,
  OrganizationPlainAttributes,
>;
export type OrganizationType = OrganizationModel & OrganizationAttributes;

export type StaticOrganizationModel = Class<OrganizationModel>;

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
      ssoSelectedType: {
        type: types.ENUM('none', 'saml', 'oidc'),
        allowNull: false,
        defaultValue: 'none',
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
      ssoOidcClientID: {
        type: types.STRING,
        allowNull: false,
        defaultValue: '',
      },
      ssoOidcClientSecret: {
        type: types.STRING,
        allowNull: false,
        defaultValue: '',
      },
      ssoOidcConfigurationURL: {
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
