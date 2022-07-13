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
import {AccessRoles} from '../../../shared/roles';
import {BuildOptions, DataTypes, Model} from 'sequelize';
import {omit} from 'lodash';
import type {AssociateProp} from './AssociateTypes';

export interface UserRawType {
  id: number;
  networkIDs: Array<string>;
  isSuperUser: boolean;
  isReadOnlyUser: boolean;
  email: string;
  organization?: string;
  password: string;
  role: number;
}

export interface UserModel extends UserRawType, Model {
  readonly id: number;
}

export type UserModelStatic = typeof Model & {
  new (values?: object, options?: BuildOptions): UserModel;
} & AssociateProp;

export default (sequelize: sequelize.Sequelize) => {
  const attributes: sequelize.ModelAttributes<UserModel> = {
    email: DataTypes.STRING,
    organization: DataTypes.STRING,
    password: DataTypes.STRING,
    role: DataTypes.INTEGER,
    networkIDs: {
      type: DataTypes.JSON,
      allowNull: false,
      defaultValue: [],
      get() {
        return this.getDataValue('networkIDs') || [];
      },
    },
  };

  const options: sequelize.ModelOptions<UserModel> = {
    getterMethods: {
      isSuperUser() {
        return this.role === AccessRoles.SUPERUSER;
      },
      isReadOnlyUser() {
        return this.role === AccessRoles.READ_ONLY_USER;
      },
    },
  };
  const User = sequelize.define(
    'User',
    attributes,
    options as sequelize.ModelOptions,
  ) as UserModelStatic;

  // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
  User.prototype.toJSON = function () {
    // eslint-disable-next-line @typescript-eslint/no-unsafe-call, @typescript-eslint/no-unsafe-member-access
    return omit(this.get(), 'password');
  };
  return User;
};
