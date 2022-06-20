/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */

import Sequelize from 'sequelize';

import type {AssociateProp} from './AssociateTypes.flow';
import type {DataTypes, Model} from 'sequelize';
// $FlowFixMe migrated to typescript
import type {SSOSelectedType} from '../../types/auth';

type OrganizationInitAttributes = {
  id?: number,
  name: string,
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
  isHostOrg: boolean,
};

type OrganizationAttributes = OrganizationPlainAttributes & OrganizationGetters;

export type OrganizationModel = Model<
  OrganizationAttributes,
  OrganizationInitAttributes,
  OrganizationPlainAttributes,
>;
export type OrganizationType = OrganizationModel & OrganizationAttributes;

export type StaticOrganizationModel = Class<OrganizationModel>;

const HOST_ORG = 'host';

export default (
  sequelize: Sequelize,
  types: DataTypes,
): StaticOrganizationModel & AssociateProp => {
  const Organization = sequelize.define(
    'Organization',
    {
      name: types.STRING,
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
        isHostOrg() {
          return this.name === HOST_ORG;
        },
      },
    },
  );
  Organization.addHook('beforeCreate', 'nameToLowerCase', organization => {
    organization.name = organization.name.toLowerCase();
    return organization;
  });
  Organization.associate = function (_models) {
    // associations can be defined here
  };
  return Organization;
};
