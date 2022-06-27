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
import {AssociateProp} from './AssociateTypes';
import {BuildOptions, DataTypes, Model} from 'sequelize';

export interface AuditLogEntryRawType {
  actingUserId: number;
  organization: string;
  mutationType: 'CREATE' | 'UPDATE' | 'DELETE';
  objectId: string;
  objectType: string;
  objectDisplayName: string;
  mutationData: {[key: string]: unknown};
  url: string;
  ipAddress: string;
  status: 'SUCCESS' | 'FAILURE';
  statusCode: string;
}

interface AuditLogEntryModel extends AuditLogEntryRawType, Model {
  readonly id: number;
}

type AuditLogEntryModelStatic = typeof Model & {
  new (values?: object, options?: BuildOptions): AuditLogEntryModel;
} & AssociateProp;

export default (sequelize: sequelize.Sequelize) => {
  return sequelize.define(
    'AuditLogEntry',
    {
      actingUserId: {
        type: DataTypes.INTEGER,
        allowNull: false,
      },
      organization: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      mutationType: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      objectId: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      objectDisplayName: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      objectType: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      mutationData: {
        type: DataTypes.JSON,
        allowNull: false,
      },
      url: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      ipAddress: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      status: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      statusCode: {
        type: DataTypes.STRING,
        allowNull: false,
      },
    },
    {},
  ) as AuditLogEntryModelStatic;
};
