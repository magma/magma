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
 */

import sequelize from 'sequelize';
import {BuildOptions, DataTypes, Model} from 'sequelize';
import type {AssociateProp} from './AssociateTypes';
import type {SSOSelectedType} from '../../types/auth';

export interface OrganizationPlainAttributes {
  id: number;
  name: string;
  customDomains: Array<string>;
  networkIDs: Array<string>;
  csvCharset: string;
  ssoSelectedType: SSOSelectedType;
  ssoCert: string;
  ssoEntrypoint: string;
  ssoIssuer: string;
  ssoOidcClientID: string;
  ssoOidcClientSecret: string;
  ssoOidcConfigurationURL: string;
}

export interface OrganizationModel extends OrganizationPlainAttributes, Model {
  isHostOrg: boolean;
}

type OrganizationModelStatic = typeof Model & {
  new (values?: object, options?: BuildOptions): OrganizationModel;
} & AssociateProp;

const HOST_ORG = 'host';

export default (sequelize: sequelize.Sequelize) => {
  const attributes: sequelize.ModelAttributes<OrganizationModel> = {
    name: DataTypes.STRING,
    csvCharset: DataTypes.STRING,
    customDomains: {
      type: DataTypes.JSON,
      allowNull: false,
      defaultValue: [],
      get() {
        return this.getDataValue('customDomains') || [];
      },
    },
    networkIDs: {
      type: DataTypes.JSON,
      allowNull: false,
      defaultValue: [],
    },
    ssoSelectedType: {
      type: DataTypes.ENUM('none', 'saml', 'oidc'),
      allowNull: false,
      defaultValue: 'none',
    },
    ssoCert: {
      type: DataTypes.TEXT,
      allowNull: false,
      defaultValue: '',
    },
    ssoEntrypoint: {
      type: DataTypes.STRING,
      allowNull: false,
      defaultValue: '',
    },
    ssoIssuer: {
      type: DataTypes.STRING,
      allowNull: false,
      defaultValue: '',
    },
    ssoOidcClientID: {
      type: DataTypes.STRING,
      allowNull: false,
      defaultValue: '',
    },
    ssoOidcClientSecret: {
      type: DataTypes.STRING,
      allowNull: false,
      defaultValue: '',
    },
    ssoOidcConfigurationURL: {
      type: DataTypes.STRING,
      allowNull: false,
      defaultValue: '',
    },
  };
  const options: sequelize.ModelOptions<OrganizationModel> = {
    getterMethods: {
      isHostOrg() {
        return this.name === HOST_ORG;
      },
    },
  };
  const Organization = sequelize.define(
    'Organization',
    attributes,
    options as sequelize.ModelOptions,
  ) as OrganizationModelStatic;

  Organization.addHook(
    'beforeCreate',
    'nameToLowerCase',
    (organization: OrganizationModel) => {
      organization.name = organization.name.toLowerCase();
    },
  );

  return Organization;
};
