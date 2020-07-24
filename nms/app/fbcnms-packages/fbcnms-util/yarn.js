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

const _ = require('lodash');
const fs = require('fs');
const glob = require('glob');
const path = require('path');

type Dependencies = {
  [key: string]: string,
};

type WorkspaceConfig = {
  packages: Array<string>,
  nohoist: Array<string>,
};

type Manifest = {
  dependencies?: Dependencies,
  devDependencies?: Dependencies,
  peerDependencies?: Dependencies,
  optionalDependencies?: Dependencies,
  workspaces?: WorkspaceConfig,
};

export function resolveWorkspaces(
  root: string,
  rootManifest: Manifest,
): Manifest[] {
  if (!rootManifest.workspaces) {
    return [];
  }

  const files = rootManifest.workspaces.packages.map(pattern =>
    glob.sync(pattern.replace(/\/?$/, '/+(package.json)'), {
      cwd: root,
      ignore: pattern.replace(/\/?$/, '/node_modules/**/+(package.json)'),
    }),
  );

  return _.flatten(files).map(file => readManifest(path.join(root, file)));
}

export function readManifest(file: string): Manifest {
  return JSON.parse(fs.readFileSync(file, 'utf8'));
}
