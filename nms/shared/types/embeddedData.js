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

import type {FeatureID} from './features';
import type {SSOSelectedType} from './auth';
import type {Tab} from './tabs';

export type User = {
  tenant: string,
  email: string,
  isSuperUser: boolean,
  isReadOnlyUser: boolean,
};

export type EmbeddedData = {|
  csrfToken: string,
  tabs: $ReadOnlyArray<Tab>,
  user: User,
  enabledFeatures: FeatureID[],
  ssoEnabled: boolean,
  ssoSelectedType: SSOSelectedType,
  csvCharset: ?string,
|};
