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
import fs from 'fs';
import path from 'path';
import paths from './paths';

const DEV_MODE = process.env.NODE_ENV !== 'production';
const MANIFEST_FILE = path.join(paths.appSrc, '../static/dist/manifest.json');

let manifest: null | Record<string, string> = null;
if (fs.existsSync(MANIFEST_FILE)) {
  const manifestRaw = fs.readFileSync(MANIFEST_FILE).toString('utf8').trim();
  manifest = JSON.parse(manifestRaw) as Record<string, string>;
}
export default function staticDist(
  filename: string,
  projectName: string,
): string | null | undefined {
  if (DEV_MODE || !manifest) {
    const path = '/static/dist/' + filename;
    if (typeof projectName === 'string') {
      return '/' + projectName + path;
    }
    return path;
  }

  return manifest[filename] || '/dev/null/' + filename;
}
