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
 * @flow strict-local
 * @format
 */

import '@fbcnms/babel-register/polyfill';
import path from 'path';
import {merge, omit} from 'lodash';
import {readManifest, resolveWorkspaces} from '@fbcnms/util/yarn';

// Packages that can have duplicate versions (keep this to a minimum)
const PACKAGE_BLACKLIST = ['core-js'];

it('ensures no mismatched versions in workspaces', () => {
  const root = path.resolve(__dirname, '..');
  const rootManifest = readManifest(path.resolve(root, 'package.json'));
  const workspaces = resolveWorkspaces(root, rootManifest);

  const allManifests = [rootManifest, ...workspaces];

  const allDepsMap = merge(
    {},
    ...allManifests.map(manifest =>
      merge(
        {},
        manifest.dependencies,
        manifest.devDependencies,
        manifest.optionalDependencies,
        manifest.peerDependencies,
      ),
    ),
  );

  const o = packages => omit(packages, PACKAGE_BLACKLIST);
  const filteredDepsMap = o(allDepsMap);

  for (const manifest of allManifests) {
    expect(filteredDepsMap).toMatchObject(o(manifest.dependencies) || {});
    expect(filteredDepsMap).toMatchObject(o(manifest.devDependencies) || {});
    expect(filteredDepsMap).toMatchObject(o(manifest.peerDependencies) || {});
    expect(filteredDepsMap).toMatchObject(
      o(manifest.optionalDependencies) || {},
    );
  }
});
