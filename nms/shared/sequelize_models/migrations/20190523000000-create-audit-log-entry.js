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
 * @typedef { import("sequelize").DataTypes } DataTypes
 */

module.exports = {
  /**
   * @param {QueryInterface} queryInterface
   * @param {DataTypes} types
   */
  up: (queryInterface, types) => {
    return queryInterface.createTable('AuditLogEntries', {
      id: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: types.INTEGER,
      },
      actingUserId: {
        type: types.INTEGER,
        allowNull: false,
      },
      organization: {
        type: types.STRING,
        allowNull: false,
      },
      mutationType: {
        type: types.STRING,
        allowNull: false,
      },
      objectId: {
        type: types.STRING,
        allowNull: false,
      },
      objectDisplayName: {
        type: types.STRING,
        allowNull: false,
      },
      objectType: {
        type: types.STRING,
        allowNull: false,
      },
      mutationData: {
        type: types.JSON,
        allowNull: false,
      },
      url: {
        type: types.STRING,
        allowNull: false,
      },
      ipAddress: {
        type: types.STRING,
        allowNull: false,
      },
      status: {
        type: types.STRING,
        allowNull: false,
      },
      statusCode: {
        type: types.STRING,
        allowNull: false,
      },
      createdAt: {
        allowNull: false,
        type: types.DATE,
      },
      updatedAt: {
        allowNull: false,
        type: types.DATE,
      },
    });
  },

  /**
   * @param {QueryInterface} queryInterface
   */
  down: queryInterface => {
    return queryInterface.dropTable('AuditLogEntries');
  },
};
