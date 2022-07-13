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

import bcrypt from 'bcryptjs';

module.exports = {
  /**
   * @param {QueryInterface} queryInterface
   */
  up: async queryInterface => {
    const email = process.env.FB_TEST_USER;
    if (!email) {
      return Promise.resolve(null);
    }

    const salt = await bcrypt.genSalt(10);
    const passwordHash = await bcrypt.hash(email, salt);

    return queryInterface.bulkInsert(
      'Users',
      [
        {
          email: email,
          organization: 'master',
          password: passwordHash,
          role: 3,
          createdAt: new Date(),
          updatedAt: new Date(),
          networkIDs: '[]',
          tabs: '[]',
        },
        {
          email: email,
          organization: 'fb-test',
          password: passwordHash,
          role: 3,
          createdAt: new Date(),
          updatedAt: new Date(),
          networkIDs: '[]',
          tabs: '["inventory", "workorders", "automation"]',
        },
        {
          email: email,
          organization: 'magma-test',
          password: passwordHash,
          role: 3,
          createdAt: new Date(),
          updatedAt: new Date(),
          networkIDs: '["mpk_test"]',
          tabs: '["nms"]',
        },
      ],
      {},
    );
  },
  /**
   * @param {QueryInterface} queryInterface
   */
  down: queryInterface => {
    if (process.env.FB_TEST_USER) {
      // @ts-ignore
      return queryInterface.bulkDelete('Users', null, {});
    }
    return Promise.resolve(null);
  },
};
