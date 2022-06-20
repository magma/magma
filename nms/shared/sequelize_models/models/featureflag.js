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
// $FlowFixMe migrated to typescript
import type {AssociateProp} from './AssociateTypes';
import type {DataTypes, Model} from 'sequelize';

export type FeatureFlagRawType = {|
  featureId: string,
  organization: string,
  enabled: boolean,
|};

type FeatureFlagReadAttributes = {|
  ...FeatureFlagRawType,
  id: number,
|};

type FeatureFlagModel = Model<FeatureFlagReadAttributes, FeatureFlagRawType>;
export type StaticFeatureFlagModel = Class<FeatureFlagModel>;
export type FeatureFlagType = FeatureFlagModel & FeatureFlagRawType;

export default (
  sequelize: Sequelize,
  types: DataTypes,
): AssociateProp & StaticFeatureFlagModel => {
  return sequelize.define(
    'FeatureFlag',
    {
      featureId: {
        type: types.STRING,
        allowNull: false,
      },
      organization: {
        type: types.STRING,
        allowNull: false,
      },
      enabled: {
        type: types.BOOLEAN,
        allowNull: false,
      },
    },
    {},
  );
};
