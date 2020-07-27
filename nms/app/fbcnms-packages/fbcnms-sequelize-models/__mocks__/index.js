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

jest.mock('../sequelizeConfig', () => {
  process.env.NODE_ENV = 'test';
  return {
    [process.env.NODE_ENV]: {
      username: null,
      password: null,
      database: 'db',
      dialect: 'sqlite',
      logging: false,
    },
  };
});

beforeAll(async () => {
  const {sequelize} = jest.requireActual('../');
  // running sync instead of migrations because of weird foreign key issues
  await sequelize.sync({force: true});
});

const realModels = jest.requireActual('../');
module.exports = realModels;
