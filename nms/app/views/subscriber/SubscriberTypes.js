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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {SubscriberActionType, SubscriberInfo} from './SubscriberUtils';
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
export type SubscribersDialogDetailProps = {
  // Subscribers to add, edit or delete
  setSubscribers: (Array<SubscriberInfo>) => void,
  subscribers: Array<SubscriberInfo>,
  // Formatting error (eg: field missing, wrong IMSI format)
  setAddError: (Array<string>) => void,
  addError: Array<string>,
  // Display dropzone if set to true
  setUpload: boolean => void,
  upload: boolean,
  onClose: () => void,
  // Add, edit or delete subscribers
  onSave: (Array<SubscriberInfo>, selectedSubscribers?: Array<string>) => void,
  error?: string,
  // Row added with the Add New Row button
  rowAdd: boolean,
  setRowAdd: boolean => void,
  // Delete, Edit or Add subscriber
  subscriberAction: SubscriberActionType,
};
