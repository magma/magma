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
 * @typedef { import("sequelize").WhereOptions } WhereOptions
 */

module.exports = {
  /**
   * @param {QueryInterface} queryInterface
   */
  up: queryInterface => {
    return queryInterface.bulkInsert(
      'Organizations',
      [
        {
          id: '1',
          customDomains: '[]',
          name: 'master',
          networkIDs: '[]',
          tabs: '["admin"]',
          createdAt: '2019-02-11 20:05:05',
          updatedAt: '2019-02-11 20:05:05',
        },
        {
          id: '2',
          customDomains: '[]',
          name: 'fb-test',
          networkIDs: '[]',
          tabs: '["inventory", "workorders", "automation"]',
          createdAt: '2019-02-11 20:05:05',
          updatedAt: '2019-02-11 20:05:05',
        },
        {
          id: '3',
          customDomains: '[]',
          name: 'magma-test',
          networkIDs: '["mpk_test"]',
          tabs: '["nms"]',
          createdAt: '2019-02-11 20:05:05',
          updatedAt: '2019-02-11 20:05:05',
        },
      ],
      {},
    );
  },

  /**
   * @param {QueryInterface} queryInterface
   */
  down: queryInterface => {
    // @ts-ignore
    return queryInterface.bulkDelete('Organizations', null, {});
  },
};
