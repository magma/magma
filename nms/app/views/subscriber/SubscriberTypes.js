/*
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

import type {subscriber} from '../../../generated/MagmaAPIBindings';

export type subscriberStaticIpsRowType = {
  apnName: string,
  staticIp: string,
};
export type subscriberForbiddenNetworkTypes = {
  nwTypes: string,
};
export type EditSubscriberProps = {
  subscriberState: subscriber,
  onSubscriberChange: (key: string, val: string | number | {}) => void,
  inputClass: string,
  onTrafficPolicyChange: (
    key: string,
    val: string | number | {},
    index: number,
  ) => void,
  onDeleteApn: (apn: {}) => void,
  onAddApnStaticIP: () => void,
  subProfiles: {},
  subscriberStaticIPRows: Array<subscriberStaticIpsRowType>,
  forbiddenNetworkTypes: Array<subscriberForbiddenNetworkTypes>,
  authKey: string,
  authOpc: string,
  setAuthKey: (key: string) => void,
  setAuthOpc: (key: string) => void,
};
