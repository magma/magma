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
import {FeatureID} from '../../types/features';
import type {AssociateProp} from './AssociateTypes';

export interface FeatureFlagRawType {
  featureId: FeatureID;
  organization: string;
  enabled: boolean;
  readonly id: number;
}

export interface FeatureFlagModel
  extends FeatureFlagRawType,
    Model<FeatureFlagRawType, Partial<FeatureFlagRawType>> {}

type FeatureFlagModelStatic = typeof Model & {
  new (values?: object, options?: BuildOptions): FeatureFlagModel;
} & AssociateProp;

export default (sequelize: sequelize.Sequelize) => {
  return sequelize.define(
    'FeatureFlag',
    {
      featureId: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      organization: {
        type: DataTypes.STRING,
        allowNull: false,
      },
      enabled: {
        type: DataTypes.BOOLEAN,
        allowNull: false,
      },
    },
    {},
  ) as FeatureFlagModelStatic;
};
