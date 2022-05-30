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
        '<rootDir>/server/**/__tests__/*.[jt]s?(x)',
        '<rootDir>/shared/**/__tests__/*.[jt]s?(x)',
      ],
      transform: {
        '^.+\\.(js|ts|tsx)$': 'babel-jest',
      },
      resetMocks: true,
      restoreMocks: true,
    },
    {
      name: 'app',
      testEnvironment: 'jsdom',
      testMatch: ['<rootDir>/app/**/__tests__/*.[jt]s?(x)'],
      transform: {
        '^.+\\.(js|ts|tsx)$': 'babel-jest',
      },
      setupFilesAfterEnv: ['./jest.setup.app.ts'],
      resetMocks: true,
      restoreMocks: true,
    },
  ],
  testEnvironment: 'jsdom',
};
