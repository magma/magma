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
 * @flow
 * @format
 */

import bcrypt from 'bcryptjs';
import type {DataTypes, QueryInterface} from 'sequelize';

module.exports = {
  up: async (queryInterface: QueryInterface, _: DataTypes) => {
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

  down: (queryInterface: QueryInterface, _: DataTypes) => {
    if (process.env.FB_TEST_USER) {
      return queryInterface.bulkDelete('Users', null, {});
    }
    return Promise.resolve(null);
  },
};
