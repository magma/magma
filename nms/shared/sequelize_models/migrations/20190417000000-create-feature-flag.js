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
 *
 * We are using JSDoc type annotations because renaming this file will cause
 * the migration to be re-executed.
 *
 * NEW MIGRATIONS SHOULD BE WRITTEN IN TYPESCRIPT!
 *
 * @typedef { import("sequelize").QueryInterface } QueryInterface
 */

import {DataTypes} from 'sequelize';

module.exports = {
  /**
   * @param {{ context: QueryInterface}} params
   */
  up: ({context: queryInterface}) => {
    return queryInterface.createTable('FeatureFlags', {
      id: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: DataTypes.INTEGER,
      },
      featureId: {
        allowNull: false,
        type: DataTypes.STRING,
      },
      organization: {
        allowNull: false,
        type: DataTypes.STRING,
      },
      enabled: {
        allowNull: false,
        type: DataTypes.BOOLEAN,
      },
      createdAt: {
        allowNull: false,
        type: DataTypes.DATE,
      },
      updatedAt: {
        allowNull: false,
        type: DataTypes.DATE,
      },
    });
  },
  /**
   * @param {{ context: QueryInterface}} params
   */
  down: ({context: queryInterface}) => {
    return queryInterface.dropTable('FeatureFlags');
  },
};
