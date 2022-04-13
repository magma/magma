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

const TEST_SUBDOMAIN = '-test';
const LOCALHOST = 'localhost';
const PHB_SUBDOMAIN = 'purpleheadband.cloud';

export function isTestEnv(): boolean {
  return (
    window.location.hostname.includes(TEST_SUBDOMAIN) ||
    window.location.hostname.includes(LOCALHOST)
  );
}

export function isPhbProdEnv(): boolean {
  return window.location.hostname.includes(PHB_SUBDOMAIN);
}
