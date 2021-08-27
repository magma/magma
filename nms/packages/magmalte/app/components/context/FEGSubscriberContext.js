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

import type {
  network_id,
  subscriber_id,
  subscriber_state,
} from '../../../generated/MagmaAPIBindings';

import React from 'react';

export type FEGSubscriberContextType = {
  sessionState: {
    [networkId: network_id]: {[subscriberId: subscriber_id]: subscriber_state},
  },
  setSessionState: (newSessionState: {
    [string]: {[string]: subscriber_state},
  }) => void,
};

export default React.createContext<FEGSubscriberContextType>({});
