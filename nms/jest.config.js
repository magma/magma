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
 * @format
 */

module.exports = {
  collectCoverageFrom: [
    '**/*.js',
    '!**/__mocks__/**',
    '!**/__tests__/**',
    '!**/node_modules/**',
  ],
  coverageReporters: ['json', 'html'],
  projects: [
    {
      name: 'server',
      testEnvironment: 'node',
      testMatch: [
        '<rootDir>/server/**/__tests__/*.js',
        '<rootDir>/shared/**/__tests__/*.js',
      ],
      transform: {
        '^.+\\.js$': 'babel-jest',
      },
    },
    {
      name: 'app',
      testEnvironment: 'jsdom',
      testMatch: ['<rootDir>/app/**/__tests__/*.js'],
      transform: {
        '^.+\\.js$': 'babel-jest',
      },
    },
  ],
  testEnvironment: 'jsdom',
};
