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

const fs = require('fs');
const path = require('path');

const appDirectory = fs.realpathSync(process.cwd());

const resolveApp = (relativePath: string) =>
  path.resolve(appDirectory, relativePath);

module.exports = {
  appIndexJs: resolveApp('app/main.js'),
  loginJs: resolveApp('app/login.js'),
  masterJs: resolveApp('app/master.js'),
  appSrc: resolveApp('app'),
  distPath: resolveApp('static/dist'),
  packagesDir: resolveApp('../../packages'),
  fbcnmsDir: path.dirname(
    path.dirname(require.resolve('@fbcnms/babel-register')),
  ),
  resolveApp,
};
